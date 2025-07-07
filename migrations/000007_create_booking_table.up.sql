CREATE TABLE booking (
    booking_id SERIAL PRIMARY KEY,
    booking_status VARCHAR(50) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES user(user_id) ON DELETE CASCADE,
    total_tickets INTEGER NOT NULL,
    show_id INTEGER NOT NULL REFERENCES show(show_id) ON DELETE CASCADE,
    booking_time TIMESTAMP WITH TIME ZONE DEFAULT now(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
