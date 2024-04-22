CREATE TABLE IF NOT EXISTS users (
       id uuid PRIMARY KEY,
       username TEXT NOT NULL,
       email TEXT NOT NULL,
       password TEXT NOT NULL,
       first_name TEXT NOT NULL,
       last_name TEXT NOT NULL,
       bio TEXT NOT NULL,
       website TEXT NOT NULL,
       is_active BOOLEAN DEFAULT true,
       refresh_token TEXT NOT NULL,
       created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMP WITHOUT TIME ZONE,
       deleted_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX idx_unique_email ON users(email) WHERE deleted_at IS NULL;