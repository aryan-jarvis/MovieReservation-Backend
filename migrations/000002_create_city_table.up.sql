CREATE TABLE city (
    city_id SERIAL PRIMARY KEY,
    state_id INTEGER NOT NULL REFERENCES state(state_id) ON DELETE CASCADE,
    city_name VARCHAR(100) NOT NULL
);
