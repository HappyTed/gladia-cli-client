package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
}

func New(driver string, dsn string) (*Database, error) {

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	var db = &Database{conn: conn}

	if err = db.CreateTable(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT,
		file_url TEXT,
		result_url TEXT,
		task_id TEXT,
		status TEXT,
		results BLOB
	);`
	_, err := db.conn.Exec(query)
	return err
}

func (db *Database) Insert(row Task) (int64, error) {
	results, err := db.conn.Exec(
		"INSERT INTO requests(file_path, file_url, result_url, task_id, status, results) VALUES (?, ?, ?, ?, ?, ?)",
		row.FilePath, row.FileURL, row.ResultURL, row.TaskID, row.Status, row.Results,
	)
	if err != nil {
		log.Fatal(err)
	}

	id, err := results.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *Database) GetByID(id int64) (Task, error) {

	query := `SELECT * FROM requests
	WHERE id = ?`

	row := db.conn.QueryRow(query, id)

	var r Task
	err := row.Scan(&r.ID, &r.FilePath, &r.FileURL, &r.ResultURL, &r.TaskID, &r.Status, &r.Results)
	return r, err
}

func (db *Database) GetByTaskID(taskID string) (Task, error) {

	query := `SELECT * FROM requests
	WHERE task_id = ?`

	row := db.conn.QueryRow(query, taskID)

	var r Task
	err := row.Scan(&r.ID, &r.FilePath, &r.FileURL, &r.ResultURL, &r.TaskID, &r.Status, &r.Results)
	return r, err
}

func (db *Database) GetAll() ([]Task, error) {
	rows, err := db.conn.Query("SELECT * FROM requests ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Task
	for rows.Next() {
		var r Task
		rows.Scan(&r.ID, &r.FilePath, &r.FileURL, &r.ResultURL, &r.TaskID, &r.Status, &r.Results)
		requests = append(requests, r)
	}

	return requests, nil
}

func (db *Database) Update(
	id int64,
	fileUrl string,
	resultUrl string,
	taskID string,
	status string,
	results []byte,
) error {

	query := `UPDATE requests 
SET file_url = ?, result_url = ?, task_id = ?, status = ?, results = ? 
WHERE id = ?`

	_, err := db.conn.Exec(
		query,
		fileUrl, resultUrl, taskID, status, results, id,
	)

	return err
}

func (db *Database) RemoveByID(id int64) error {
	_, err := db.conn.Exec("DELETE FROM requests WHERE id = ?", id)
	return err
}

func (db *Database) RemoveByTaskID(taskID string) error {
	_, err := db.conn.Exec("DELETE FROM requests WHERE task_id = ?", taskID)
	return err
}
