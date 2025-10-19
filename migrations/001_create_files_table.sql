CREATE TABLE IF NOT EXISTS files(
    id VARCHAR (36) PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    size BIGINT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_files_filename ON files(filename);