CREATE TABLE IF NOT EXISTS "recoveries" (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(40) UNIQUE NOT NULL,
    email VARCHAR(140) UNIQUE NOT NULL,
    code VARCHAR(255),
    attempts INT,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expired BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_email FOREIGN KEY (email) REFERENCES users(email)
);
