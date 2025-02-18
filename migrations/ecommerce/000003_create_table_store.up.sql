DO $$
BEGIN
    CREATE TABLE IF NOT EXISTS stores (
        store_id VARCHAR(50) PRIMARY KEY,
        store_name VARCHAR(100) NOT NULL,
        description TEXT,
        user_id VARCHAR(50) NOT NULL,
        created_at VARCHAR(50) NULL,
        updated_at VARCHAR(50) NULL,
        FOREIGN KEY (user_id) REFERENCES users (uuid) ON DELETE CASCADE
    );
END $$;
