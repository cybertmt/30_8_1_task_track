/*
    Схема БД для информационной системы
    отслеживания выполнения задач.
*/

DROP TABLE IF EXISTS tasks_labels, tasks, labels, users;

-- пользователи системы
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- метки задач
CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- задачи
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    opened BIGINT NOT NULL, -- время создания задачи
    closed BIGINT DEFAULT 0, -- время выполнения задачи
    author_id INTEGER REFERENCES users(id) ON DELETE CASCADE , -- автор задачи
    assigned_id INTEGER REFERENCES users(id) ON DELETE CASCADE, -- ответственный
    title TEXT, -- название задачи
    content TEXT -- текст задачи
);

-- связь многие-ко-многим между задачами и метками
CREATE TABLE tasks_labels (
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    label_id INTEGER REFERENCES labels(id) ON DELETE CASCADE
);

-- наполнение таблицы users
INSERT INTO users(name)
VALUES ('John'), ('Doe'), ('Richard'), ('Roe'), ('Cyber');

-- наполнение таблицы labels
INSERT INTO labels(name)
VALUES ('Go'), ('Java'), ('Python'), ('HTML'), ('CSS');