package database

type IDatabase interface {
	CreateTable() error
	Insert(Task) (int64, error)
	GetByID(int64) (Task, error)
	GetByTaskID(string) (Task, error)
	GetAll() ([]Task, error)
	Update(id int64, fileUrl string, resultUrl string, taskID string, status string, results []byte) error
	RemoveByID(id int64) error
	RemoveByTaskID(taskID string) error
}
