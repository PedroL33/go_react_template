CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    thumbnail_url VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),  -- e.g., work, personal, etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);