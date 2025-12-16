package database

type Request struct {
	ID       int    `db:"id"`
	FilePath string `db:"file_path"` // file path
	FileUrl  string `db:"file_url"`  // gladia upload
	TaskID   string `db:"task_id"`
	Status   string `db:"status"`
	Results  []byte `db:"results"` // json data
}
