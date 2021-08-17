package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"time"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id, labelId, authorId, assignedId int
	err := s.db.QueryRow(context.Background(), `
		SELECT id FROM users
		ORDER BY random()
		LIMIT 1;
		`,
	).Scan(&authorId)
	err = s.db.QueryRow(context.Background(), `
		SELECT id FROM users
		ORDER BY random()
		LIMIT 1;
		`,
	).Scan(&assignedId)
	err = s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content, opened, author_id, assigned_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
		`,
		t.Title,
		t.Content,
		strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		authorId,
		assignedId,
	).Scan(&id)
	err = s.db.QueryRow(context.Background(), `
		SELECT id FROM labels
		ORDER BY random()
		LIMIT 1;
		`,
	).Scan(&labelId)
	_, err = s.db.Exec(context.Background(), `
		INSERT INTO tasks_labels (task_id, label_id)
		VALUES ($1, $2);
		`,
		id,
		labelId,
	)
	return id, err
}

// AuthorTasks возвращает список задач из БД по автору.
func (s *Storage) AuthorTasks(authorName string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			tasks.id,
			tasks.opened,
			tasks.closed,
			tasks.author_id,
			tasks.assigned_id,
			tasks.title,
			tasks.content
		FROM tasks, users
		WHERE
			tasks.author_id = users.id AND users.name = $1
		ORDER BY tasks.id;
	`,
		authorName,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// LabelTasks возвращает список задач из БД по тегу.
func (s *Storage) LabelTasks(labelName string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			tasks.id,
			tasks.opened,
			tasks.closed,
			tasks.author_id,
			tasks.assigned_id,
			tasks.title,
			tasks.content
		FROM labels
		INNER JOIN tasks_labels 
			ON tasks_labels.label_id = labels.id
		INNER JOIN tasks 
			ON tasks.id = tasks_labels.task_id
		WHERE
			labels.name = $1
		ORDER BY tasks.id;
	`,
		labelName,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// UpdateTasks обновляет Content задачи по id.
func (s *Storage) UpdateTasks(taskId int, content string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		UPDATE tasks
		SET content = $2 
		WHERE tasks.id = $1
		RETURNING *;
	`,
		taskId,
		content,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// DeleteTasks удаляет задачу по id.
func (s *Storage) DeleteTasks(taskId int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		DELETE FROM tasks
		WHERE tasks.id = $1
		RETURNING *;
	`,
		taskId,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}