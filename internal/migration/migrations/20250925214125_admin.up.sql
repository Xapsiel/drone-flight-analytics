
CREATE TYPE action_type AS ENUM ('upload_file', 'delete_file');

CREATE TABLE IF NOT EXISTS files(
    id SERIAL PRIMARY KEY ,
    user_id TEXT NOT NULL ,
    filename TEXT NOT NULL ,
    size BIGINT,
    uploaded_at TIMESTAMP,
    metadata JSONB
);

CREATE TABLE IF NOT EXISTS admin_logs(
    id SERIAL PRIMARY KEY ,
    user_id TEXT NOT NULL ,
    action action_type NOT NULL ,
    file_id INTEGER REFERENCES files(id),
    details JSONB,
    created_at TIMESTAMP DEFAULT  CURRENT_TIMESTAMP
)
