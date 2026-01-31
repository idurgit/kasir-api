-- Add category_id and category_name columns to products table
ALTER TABLE products 
ADD COLUMN category_id INTEGER,
ADD COLUMN category_name VARCHAR(255);

-- Optional: Add index on category_id for better query performance
CREATE INDEX idx_products_category_id ON products(category_id);
