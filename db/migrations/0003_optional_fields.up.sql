ALTER TABLE products
    MODIFY sku VARCHAR(64) NULL;

ALTER TABLE visit_items
    MODIFY present_quantity INT NULL,
    MODIFY store_quantity INT NULL,
    MODIFY price DECIMAL(10,2) NULL;
