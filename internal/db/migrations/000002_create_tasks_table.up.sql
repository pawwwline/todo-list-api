CREATE TABLE tasks(
    id SERIAL PRIMARY KEY,
    user_id INT,
    title VARCHAR(255),
    description VARCHAR(255),
    CONSTRAINT user_tasks
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);