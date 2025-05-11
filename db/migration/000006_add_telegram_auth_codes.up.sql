CREATE TABLE telegram_auth_codes (
    auth_code VARCHAR(36) PRIMARY KEY,
    telegram_id BIGINT NULL,
    status VARCHAR NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'expired')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    
    CONSTRAINT fk_telegram_id 
      FOREIGN KEY(telegram_id) REFERENCES users(telegram_id)
);

CREATE INDEX idx_telegram_auth_codes_expires ON telegram_auth_codes(expires_at);