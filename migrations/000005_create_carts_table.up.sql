CREATE TABLE IF NOT EXISTS carts (
    cart_id INT GENERATED ALWAYS AS IDENTITY,
    user_id INT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    version int DEFAULT 1,
    PRIMARY KEY(cart_id),
    FOREIGN KEY(user_id)
        REFERENCES users(user_id)
);

INSERT INTO carts (user_id, created_at)
VALUES
    (1, now()),
    (2, now()),
    (3, now());