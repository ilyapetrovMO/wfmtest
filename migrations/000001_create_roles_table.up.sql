CREATE TABLE IF NOT EXISTS roles (
    role_id INT,
    name TEXT NOT NULL UNIQUE,
    PRIMARY KEY(role_id)
);

INSERT INTO roles(role_id, name)
VALUES
    (1, 'manager'),
    (2, 'user');