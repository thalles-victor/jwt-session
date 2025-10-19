CREATE TABLE IF NOT EXISTS "sessions" (
    id VARCHAR(40) PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    browser TEXT,
    ip TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);