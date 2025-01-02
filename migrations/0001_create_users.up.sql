CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    two_factor_secret TEXT,
    is_two_factor_enabled BOOLEAN DEFAULT false,
    CHECK 
    (
        (
        two_factor_secret != NULL AND is_two_factor_enabled = true
        ) 
        OR 
        (
        two_factor_secret = NULL AND is_two_factor_enabled = false
        )
    ),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);