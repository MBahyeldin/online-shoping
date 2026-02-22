DROP TRIGGER IF EXISTS set_updated_at_orders     ON orders;
DROP TRIGGER IF EXISTS set_updated_at_cart_items ON cart_items;
DROP TRIGGER IF EXISTS set_updated_at_carts      ON carts;
DROP TRIGGER IF EXISTS set_updated_at_products   ON products;
DROP TRIGGER IF EXISTS set_updated_at_categories ON categories;
DROP TRIGGER IF EXISTS set_updated_at_users      ON users;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS email_otps;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "pgcrypto";
