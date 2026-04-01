CREATE TABLE IF NOT EXISTS payments (
                                        id TEXT PRIMARY KEY,
                                        order_id TEXT NOT NULL,
                                        transaction_id TEXT,
                                        amount BIGINT NOT NULL,
                                        status TEXT NOT NULL,
                                        created_at TIMESTAMP NOT NULL
);