CREATE TABLE IF NOT EXISTS products(
    product_id INT GENERATED ALWAYS AS IDENTITY,
    name TEXT,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY(product_id)
);

INSERT INTO products(name, description, created_at)
VALUES
    ('Sparkling Water', 'Cheap sparkling water', '2020-01-01'),
    ('Water', 'Just water', '2021-05-12'),
    ('Bread', 'Just bread', '2018-12-20');