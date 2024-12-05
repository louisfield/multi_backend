CREATE TABLE words (
    id SERIAL PRIMARY KEY,
    word VARCHAR(255) NOT NULL,
    points INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);