CREATE TABLE clothes_info (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    cloth_name VARCHAR(255) NOT NULL,
    cloth_cost INTEGER NOT NULL,
    cloth_size VARCHAR(4) NOT NULL CHECK (cloth_size IN ('s', 'l', 'xl', 'xxl')),
    version INTEGER NOT NULL DEFAULT 1
);
