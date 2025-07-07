CREATE TABLE transaction (
    transaction_id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES booking(booking_id) ON DELETE CASCADE,
    total_amount NUMERIC(10,2) NOT NULL,
    payment_method VARCHAR(50),
    transaction_time TIMESTAMP WITH TIME ZONE DEFAULT now(),
    payment_status VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
