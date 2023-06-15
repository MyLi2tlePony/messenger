CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY UNIQUE NOT NULL,
    public_id TEXT UNIQUE,

    login    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,

    first_name  TEXT NOT NULL,
    second_name TEXT NOT NULL,

    created TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (Id) ON DELETE CASCADE,
    token TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY UNIQUE NOT NULL,

    name TEXT NOT NULL,
    description TEXT NOT NULL,
    open BOOLEAN NOT NULL,

    created TIMESTAMP NOT NULL
);
