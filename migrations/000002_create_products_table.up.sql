CREATE TABLE IF NOT EXISTS products(
    product_id INT GENERATED ALWAYS AS IDENTITY,
    name TEXT,
    description TEXT,
    PRIMARY KEY(product_id)
);

INSERT INTO products(name, description)
VALUES
    ('Sparkling Water', 'Cheap sparkling water'),
    ('Water', 'Just water'),
    ('Bread', 'Just bread');