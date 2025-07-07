CREATE TABLE movie (
    movie_id SERIAL PRIMARY KEY,
    movie_name VARCHAR(255) NOT NULL,
    movie_description TEXT,
    duration INTEGER, -- duration in minutes
    language VARCHAR(50),
    genre VARCHAR(100),
    poster_url TEXT,
    rating NUMERIC(2,1), -- e.g., 8.5
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
