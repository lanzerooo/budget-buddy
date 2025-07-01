CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    amount NUMERIC NOT NULL,
    category TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
