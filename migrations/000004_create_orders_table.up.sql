CREATE TABLE IF NOT EXISTS orders(
    order_id INT GENERATED ALWAYS AS IDENTITY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    version int DEFAULT 1,
    CHECK (amount > 0),
    PRIMARY KEY(order_id),
    FOREIGN KEY(user_id)
        REFERENCES users(user_id),
    FOREIGN KEY(product_id)
        REFERENCES products(product_id)
);