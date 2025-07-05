CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(30) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL, 
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
