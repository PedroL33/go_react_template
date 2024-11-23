CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    due_date TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending',  -- 'pending', 'in_progress', 'completed'
    priority INT DEFAULT 1,  -- 1 = low, 2 = medium, 3 = high
    category VARCHAR(100),  -- e.g., work, personal, etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);