CREATE TABLE IF NOT EXISTS orders(
    order_id INT GENERATED ALWAYS AS IDENTITY,
    user_id INT,
    product_id INT,
    amount INT,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    CHECK (amount > 0),
    PRIMARY KEY(order_id),
    FOREIGN KEY(user_id)
        REFERENCES users(user_id),
    FOREIGN KEY(product_id)
        REFERENCES products(product_id)
);