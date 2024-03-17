CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    password_changed_at timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    role VARCHAR(20) NOT NULL DEFAULT 'client' CHECK (role IN ('administrator', 'client')),
    created_at timestamp NOT NULL DEFAULT (now())
);
