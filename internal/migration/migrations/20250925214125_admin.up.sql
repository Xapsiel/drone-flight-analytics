
CREATE TYPE action_type AS ENUM ('upload_file', 'delete_file');
CREATE TYPE file_status AS ENUM ('error', 'processing','parsed');
CREATE TABLE IF NOT EXISTS files(
    id SERIAL PRIMARY KEY ,
    user_id TEXT NOT NULL ,
    filename TEXT NOT NULL ,
    size BIGINT,
    uploaded_at TIMESTAMP DEFAULT  CURRENT_TIMESTAMP,
    valid_count INT ,
    error_count INT,
    metadata JSONB,
    status file_status,
    UNIQUE (metadata)
);

