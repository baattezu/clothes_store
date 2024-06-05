-- Create table for orders
CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        user_id VARCHAR(255) NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create table for order items
CREATE TABLE order_items (
                             id SERIAL PRIMARY KEY,
                             order_id INT NOT NULL REFERENCES orders(id),
                             product_id VARCHAR(255) NOT NULL,
                             quantity INT NOT NULL,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create table for cart items
CREATE TABLE cart_items (
                            id SERIAL PRIMARY KEY,
                            user_id VARCHAR(255) NOT NULL,
                            product_id VARCHAR(255) NOT NULL,
                            quantity INT NOT NULL,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            UNIQUE (user_id, product_id) -- добавление уникального ограничения
);
