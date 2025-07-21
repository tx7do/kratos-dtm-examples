TRUNCATE TABLE stocks RESTART IDENTITY;
INSERT INTO stocks (id, created_at, product_id, quantity)
VALUES (1, now(), 1, 100);
SELECT setval('stocks_id_seq', (SELECT MAX(id) FROM stocks));

TRUNCATE TABLE products RESTART IDENTITY;
INSERT INTO products (id, created_at, name, price)
VALUES (1, now(), 'Apple iPhone 13', 799.99);
SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));

TRUNCATE TABLE users RESTART IDENTITY;
INSERT INTO users (id, created_at, username, nickname)
VALUES (1, now(), 'john_doe', 'John Doe');
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));