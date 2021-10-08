CREATE TABLE IF NOT EXISTS users(
    user_id INT GENERATED ALWAYS AS IDENTITY,
    username VARCHAR(15) NOT NULL UNIQUE,
    password_hash bytea NOT NULL,
    role_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    version INT DEFAULT 1,
    PRIMARY KEY(user_id),
    FOREIGN KEY(role_id)
        REFERENCES roles(role_id)
);

-- pass same as username
INSERT INTO users(username, password_hash, role_id, created_at)
VALUES
    ('user1', '$2a$14$ymJHFkT1IO2PxAovxD83j.WNGpf5SqCP2zV9x/UoVzCMO6mvxDr4W', 2, now()),
    ('user2', '$2a$14$aOjBMBnyUHjXyF4QjUrPVu4Njt3lC1YYHVYnstVjT1xBXRQ88Ll26', 2, now()),
    ('user3', '$2a$14$lnFvP.09dB35gURA2fGPsuGqp57yERh8ZR0bhH15TN.Tb08my.lyy', 2, now()),
    ('manager1', '$2a$14$9eT25DD1a2lzrV0BjMZbleVVVlVqDvPelwEW00D375/Ho8C1QgVYG', 1, now());
