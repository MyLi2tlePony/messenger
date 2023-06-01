CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY UNIQUE NOT NULL,
    public_id TEXT UNIQUE,

    login    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,

    first_name  TEXT NOT NULL ,
    second_name TEXT NOT NULL

   --created TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tockens (
    id SERIAL PRIMARY KEY UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (Id) ON DELETE CASCADE,
    tocken TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY UNIQUE NOT NULL,
    Type    INTEGER NOT NULL,
    Created TIMESTAMP NOT NULL,
    Name    TEXT NOT NULL
);