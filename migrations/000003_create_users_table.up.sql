CREATE TABLE IF NOT EXISTS users(
    user_id INT GENERATED ALWAYS AS IDENTITY,
    username VARCHAR(15) UNIQUE,
    password_hash TEXT,
    role_id INT,
    PRIMARY KEY(user_id),
    FOREIGN KEY(role_id)
        REFERENCES roles(role_id)
);

-- pass same as username
INSERT INTO users(username, password_hash, role_id)
VALUES
    ('user1', '$2a$14$ymJHFkT1IO2PxAovxD83j.WNGpf5SqCP2zV9x/UoVzCMO6mvxDr4W', 2),
    ('manager1', '$2a$14$9eT25DD1a2lzrV0BjMZbleVVVlVqDvPelwEW00D375/Ho8C1QgVYG', 1);
