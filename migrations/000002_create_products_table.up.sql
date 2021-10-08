CREATE TABLE IF NOT EXISTS products(
    product_id INT GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    in_storage INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    version int DEFAULT 1,
    PRIMARY KEY(product_id),
    CHECK (in_storage >= 0)
);

INSERT INTO products(name, description, created_at, in_storage)
VALUES
    ('Sparkling Water', 'Cheap sparkling water', '2020-01-01', 1),
    ('Water', 'Just water', '2021-05-12', 20),
    ('Bread', 'Just bread', '2018-12-20',100);