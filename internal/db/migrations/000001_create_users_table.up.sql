CREATE TABLE users(
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL, 
    hashed_password varchar(255) NOT NULL,
    created_at timestamp NOT NULL  DEFAULT NOW()
);