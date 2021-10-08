CREATE TABLE IF NOT EXISTS cart_items (
    cart_item_id int GENERATED ALWAYS AS IDENTITY,
    cart_id int NOT NULL,
    product_id int NOT NULL,
    amount int default 1 NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    version int DEFAULT 1,
    CHECK (amount > 0),
    PRIMARY KEY(cart_item_id),
    FOREIGN KEY(cart_id)
        REFERENCES carts(cart_id),
    FOREIGN KEY(product_id)
        REFERENCES products(product_id)
);