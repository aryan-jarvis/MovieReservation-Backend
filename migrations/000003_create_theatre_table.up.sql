CREATE TABLE theatre (
    theatre_id SERIAL PRIMARY KEY,
    theatre_name VARCHAR(255) NOT NULL,
    theatre_location VARCHAR(255),
    city_id INTEGER NOT NULL REFERENCES city(city_id) ON DELETE CASCADE,
    total_seats INTEGER NOT NULL,
    theatre_timing VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
