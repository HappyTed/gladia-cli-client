package database

type IDatabase interface {
	CreateTable() error
	Insert(Request) (int64, error)
	GetByID(int64) (Request, error)
	GetByTaskID(string) (Request, error)
	GetAll() ([]Request, error)
	Update(id int64, fileUrl string, taskID string, status string, results []byte) error
	DeleteByID(id int64) error
	DeleteByTaskID(taskID string) error
}
