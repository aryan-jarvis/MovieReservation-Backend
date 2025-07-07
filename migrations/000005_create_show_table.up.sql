CREATE TABLE show (
    show_id SERIAL PRIMARY KEY,
    theatre_id INTEGER NOT NULL REFERENCES theatre(theatre_id) ON DELETE CASCADE,
    movie_id INTEGER NOT NULL REFERENCES movie(movie_id) ON DELETE CASCADE,
    show_time TIMESTAMP WITH TIME ZONE NOT NULL,
    total_seats INTEGER NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    show_duration INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
