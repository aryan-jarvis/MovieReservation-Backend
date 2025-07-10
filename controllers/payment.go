package controllers

import (
	"backend/models"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	Amount float64 `json:"amount"`
	ShowID int     `json:"show_id"`
	Seats  string  `json:"seats"`
}

func GenerateTransactionID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("TXN%v", rand.Intn(100000000))
}

func GenerateHash(data string) string {
	hash := sha512.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func InitiatePayment(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userMap := userRaw.(map[string]interface{})
	var userID int
	switch v := userMap["user_id"].(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	user := models.User{
		UserID: userID,
		Name:   userMap["name"].(string),
		Email:  userMap["email"].(string),
	}

	// bind json payload to payment struct and validate
	var request PaymentRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// payu config from .env
	merchantKey := os.Getenv("PAYU_MERCHANT_KEY")
	merchantSalt := os.Getenv("PAYU_MERCHANT_SALT")
	payuBaseURL := os.Getenv("PAYU_BASE_URL")

	if merchantKey == "" || merchantSalt == "" || payuBaseURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PayU configuration missing"})
		return
	}

	transactionID := GenerateTransactionID()
	amountStr := fmt.Sprintf("%.2f", request.Amount)
	productInfo := "MovieTickets"
	firstName := user.Name
	email := user.Email
	phone := "9999999999"
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "BASE_URL is not configured"})
		return
	}

	successURL := fmt.Sprintf("%s/api/payment/success", baseURL)
	failureURL := fmt.Sprintf("%s/api/payment/failure", baseURL)

	// payu compatible hash
	hashString := fmt.Sprintf("%s|%s|%s|%s|%s|%s|||||||||||%s",
		merchantKey,
		transactionID,
		amountStr,
		productInfo,
		firstName,
		email,
		merchantSalt,
	)

	// compute SHA-512 hash of the constructed string
	hash := GenerateHash(hashString)

	booking := models.Booking{
		UserID: user.UserID,
		TxnID:  transactionID,
		Amount: int(request.Amount),
		Status: "pending",
		Seats:  request.Seats,
		ShowID: uint(request.ShowID),
	}

	// save booking with status pending
	if err := models.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// autosubmit payu form
	payuForm := fmt.Sprintf(`
		<form id="payuForm" method="post" action="%s/_payment">
			<input type="hidden" name="key" value="%s" />
			<input type="hidden" name="txnid" value="%s" />
			<input type="hidden" name="amount" value="%s" />
			<input type="hidden" name="productinfo" value="%s" />
			<input type="hidden" name="firstname" value="%s" />
			<input type="hidden" name="email" value="%s" />
			<input type="hidden" name="phone" value="%s" />
			<input type="hidden" name="surl" value="%s" />
			<input type="hidden" name="furl" value="%s" />
			<input type="hidden" name="hash" value="%s" />
		</form>
		<script type="text/javascript">
			document.getElementById("payuForm").submit();
		</script>
	`, payuBaseURL, merchantKey, transactionID, amountStr, productInfo, firstName, email, phone, successURL, failureURL, hash)

	// send auto submitting form back to the browser
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(payuForm))
}

func PaymentSuccessHandler(c *gin.Context) {
	// parsing the incoming form
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// flatten the form data into a single map
	params := make(map[string]string)
	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// the params map holds all the response fields from PayU
	log.Println("Payment Success Params:", params)

	// extract the fields
	txnID := params["txnid"]
	postedHash := params["hash"]
	status := params["status"]
	email := params["email"]
	amount := params["amount"]
	productInfo := params["productinfo"]
	firstName := params["firstname"]

	merchantSalt := os.Getenv("PAYU_MERCHANT_SALT")
	merchantKey := os.Getenv("PAYU_MERCHANT_KEY")

	// hash sequence as per PayU docs
	hashParts := []string{
		merchantSalt,
		status,
		"", "", "", "", "", "", "", "", "", "",
		email,
		firstName,
		productInfo,
		amount,
		txnID,
		merchantKey,
	}
	hashString := strings.Join(hashParts, "|")

	// compute hash
	hash := sha512.New()
	hash.Write([]byte(hashString))
	computedHash := hex.EncodeToString(hash.Sum(nil))

	// compare the hash
	if computedHash != postedHash {
		log.Printf("Hash mismatch:\nComputed: %s\nPosted: %s", computedHash, postedHash)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hash"})
		return
	}

	// updating booking status to 'success'
	if err := models.DB.Model(&models.Booking{}).Where("txn_id = ?", txnID).Update("status", "success").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
		return
	}
	frontendBaseURL := os.Getenv("FRONTEND_BASE_URL")
	if frontendBaseURL == "" {
		log.Fatal("FRONTEND_BASE_URL is not set")
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s/payment-success?txnid=%s", frontendBaseURL, txnID))
}

func PaymentFailureHandler(c *gin.Context) {
	// parse the input data
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// flatten the input data into a single map
	params := make(map[string]string)
	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	log.Println("Payment Failure Params:", params)

	// transaction id
	txnID := params["txnid"]

	// update booking status to failed
	if err := models.DB.Model(&models.Booking{}).Where("txn_id = ?", txnID).Update("status", "failed").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
		return
	}

	frontendBaseURL := os.Getenv("FRONTEND_BASE_URL")
	if frontendBaseURL == "" {
		log.Fatal("FRONTEND_BASE_URL is not set")
	}

	c.Redirect(http.StatusFound, frontendBaseURL+"/payment-failure")
}
