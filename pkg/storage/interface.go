package storage

import "task_tracker/pkg/storage/postgres"

type Interface interface {
	Tasks(int, int) ([]postgres.Task, error)
	NewTask(postgres.Task) (int, error)
	AuthorTasks(string) ([]postgres.Task, error)
	LabelTasks(string) ([]postgres.Task, error)
	UpdateTasks(int, string) ([]postgres.Task, error)
	DeleteTasks(int) ([]postgres.Task, error)
	// и еще много методов
}

