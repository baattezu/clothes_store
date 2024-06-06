CREATE TABLE user_info (
                           id SERIAL PRIMARY KEY,
                           created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                           updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                           fname VARCHAR(500) NOT NULL,
                           sname VARCHAR(500) NOT NULL,
                           email VARCHAR(255) NOT NULL UNIQUE,
                           password_hash BYTEA NOT NULL,
                           activated BOOLEAN NOT NULL DEFAULT false,
                           permission_id INT NOT NULL,
                           version INT NOT NULL DEFAULT 1
);
CREATE TABLE tokens (
                        id SERIAL PRIMARY KEY,
                        user_id INT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE,
                        hash BYTEA NOT NULL,
                        scope VARCHAR(255) NOT NULL,
                        expiry TIMESTAMPTZ NOT NULL,
                        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

ALTER TABLE user_info ADD CONSTRAINT fk_permissions
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE;

INSERT INTO permissions (code) VALUES
                                  ('read'),    ('write');