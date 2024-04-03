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
                         user_id INT NOT NULL,
                         ingredients JSONB,
                         steps JSONB
);

-- Создание таблицы КБЖУ рецепта
CREATE TABLE recipe_nutrition (
                                  nutrition_id SERIAL PRIMARY KEY,
                                  recipe_id INT,
                                  weight_total DECIMAL(10,2) DEFAULT 0,
                                  calories_total DECIMAL(10,2) DEFAULT 0,
                                  proteins_total DECIMAL(10,2) DEFAULT 0,
                                  fats_total DECIMAL(10,2) DEFAULT 0,
                                  carbohydrates_total DECIMAL(10,2) DEFAULT 0,
                                  reason_why_not varchar(255) DEFAULT 'Подсчет КБЖУ в процессе',
                                  FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id) ON DELETE CASCADE
);


CREATE TABLE ingredients (
                             ingredient_id SERIAL PRIMARY KEY,
                             name VARCHAR(255) NOT NULL UNIQUE,
                             calories_per_100g DECIMAL(10,2),
                             proteins_per_100g DECIMAL(10,2),
                             fats_per_100g DECIMAL(10,2),
                             carbohydrates_per_100g DECIMAL(10,2)
);

INSERT INTO ingredients (name, calories_per_100g, proteins_per_100g, fats_per_100g, carbohydrates_per_100g)
VALUES
    ('морковь', 41, 0.9, 0.2, 9.6),
    ('картофель', 77, 2, 0.1, 17),
    ('лук', 40, 1.1, 0.1, 9.3),
    ('помидор', 18, 0.9, 0.2, 3.9),
    ('огурец', 16, 0.7, 0.1, 3.6),
    ('говядина', 250, 26, 15, 0),
    ('свинина', 242, 26, 16, 0),
    ('рис', 130, 2.7, 0.3, 28),
    ('макароны', 130, 5, 1, 25),
    ('яйцо куриное', 155, 12.6, 10.6, 0.8),
    ('молоко', 42, 3.2, 1.8, 4.7),
    ('сыр', 403, 25, 33, 2.2),
    ('сливочное масло', 717, 0.8, 81.1, 0.6),
    ('оливковое масло', 884, 0, 100, 0),
    ('сахар', 387, 0, 0, 99.8),
    ('мед', 304, 0.3, 0, 82.4),
    ('соль', 0, 0, 0, 0),
    ('перец', 251, 10.4, 3.3, 44.2),
    ('куркума', 312, 9.68, 3.25, 64.93),
    ('куриное филе', 165, 31, 3.6, 0),
    ('курица', 165, 31, 3.6, 0),
    ('салат', 15, 1.4, 0.2, 2.9),
    ('хлебные гренки', 416, 9, 14, 63),
    ('пармезан', 420, 32, 29, 4),
    ('соус цезарь', 300, 2, 31, 3);




