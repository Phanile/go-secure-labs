CREATE TABLE appleshop.sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_sessions_user
        FOREIGN KEY (user_id)
        REFERENCES appleshop.users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_sessions_token_hash 
    ON appleshop.sessions(token_hash);

CREATE INDEX idx_sessions_user_id 
    ON appleshop.sessions(user_id);