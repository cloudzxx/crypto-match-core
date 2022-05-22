CREATE TABLE IF NOT EXISTS slice_history (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW(),
    slice_id VARCHAR(64) UNIQUE,
    begin_time TIMESTAMP,
    end_time TIMESTAMP,
    orders_count INT,
    balances_count INT
);

CREATE TABLE IF NOT EXISTS operlog (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW(),
    opt_type VARCHAR(32),
    opt_data JSONB,
    seq_id BIGINT
);

CREATE INDEX IF NOT EXISTS idx_operlog_seq ON operlog(seq_id);
CREATE INDEX IF NOT EXISTS idx_operlog_created ON operlog(created_at);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_id VARCHAR(64) UNIQUE,
    user_id INT,
    market VARCHAR(32),
    side VARCHAR(8),
    order_type VARCHAR(16),
    price VARCHAR(64),
    amount VARCHAR(64),
    left_amount VARCHAR(64),
    taker_fee VARCHAR(64),
    maker_fee VARCHAR(64),
    deal_stock VARCHAR(64),
    deal_money VARCHAR(64),
    deal_fee VARCHAR(64),
    source VARCHAR(64),
    status VARCHAR(16),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_user_market ON orders(user_id, market);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

CREATE TABLE IF NOT EXISTS balances (
    id BIGSERIAL PRIMARY KEY,
    user_id INT,
    asset VARCHAR(32),
    available VARCHAR(64),
    freeze VARCHAR(64),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, asset)
);

CREATE INDEX IF NOT EXISTS idx_balances_user ON balances(user_id);