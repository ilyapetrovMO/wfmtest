CREATE TABLE IF NOT EXISTS storage(
    storage_id INT GENERATED ALWAYS AS IDENTITY,
    product_id INT NOT NULL UNIQUE,
    stored_amount INT DEFAULT 0,
    location TEXT NOT NULL,
    CHECK (stored_amount >= 0),
    PRIMARY KEY(storage_id),
    FOREIGN KEY(product_id)
        REFERENCES products(product_id)
);