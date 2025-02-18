DO $$
BEGIN
    CREATE TABLE IF NOT EXISTS products (
        product_id varchar(50) NOT NULL,
        product_name VARCHAR(100) NOT NULL,
        description TEXT,
        price NUMERIC(10, 2) NOT NULL,
        store_id varchar(50) NOT NULL,
        created_at varchar(50) null,
        updated_at VARCHAR(50) NULL,
        FOREIGN KEY (store_id) REFERENCES stores (store_id) ON DELETE CASCADE
    );
END $$;