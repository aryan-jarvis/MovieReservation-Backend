CREATE TABLE ticket (
    ticket_id SERIAL PRIMARY KEY,
    amount NUMERIC(10,2) NOT NULL,
    transaction_id INTEGER NOT NULL REFERENCES transaction(transaction_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
