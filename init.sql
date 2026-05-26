CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL,
    password TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT NOT NULL,
    login TEXT PRIMARY KEY
);