create database recipegram
    with owner postgres;

\c recipegram;

CREATE TABLE users (
                       user_id SERIAL PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL
);

-- Создание таблицы рецептов
CREATE TABLE recipes (
                         recipe_id SERIAL PRIMARY KEY,
                         title VARCHAR(255) NOT NULL,
                         description TEXT,
                         user_id INTEGER REFERENCES users(user_id) NOT NULL,
                         ingredients TEXT[] NOT NULL,
                         steps JSONB
);






