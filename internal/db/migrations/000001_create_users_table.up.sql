CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password CHAR(64) NOT NULL
);