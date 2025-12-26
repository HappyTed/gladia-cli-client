package database

type Task struct {
	ID        int64  `db:"id"`
	FilePath  string `db:"file_path"` // file path
	FileURL   string `db:"file_url"`  // gladia upload
	ResultURL string `db:"result_url"`
	TaskID    string `db:"task_id"`
	Status    string `db:"status"`
	Results   []byte `db:"results"` // json data
}
