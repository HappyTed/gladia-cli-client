package database

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	PATH   = "/tmp/test.db"
	DRIVER = "sqlite3"
)

func TestInit(t *testing.T) {

	err := os.Remove(PATH)
	if err != nil {
		require.ErrorIs(t, err, os.ErrNotExist)
	}

	db, err := New(DRIVER, PATH)
	require.NoError(t, err)

	// проверим, что База создалась
	_, err = os.Stat(PATH)
	require.NoError(t, err)

	query := `
		SELECT 1
		FROM sqlite_master
		WHERE type='table' AND name=?;
	`

	var exists int
	err = db.conn.QueryRow(query, "requests").Scan(&exists)
	require.NoError(t, err)
	require.Equal(t, 1, exists)

	err = db.conn.Close()
	if err != nil {
		require.NoError(t, err)
	}
}

func TestInsert(t *testing.T) {

	err := os.Remove(PATH)
	if err != nil {
		require.ErrorIs(t, err, os.ErrNotExist)
	}

	db, err := New(DRIVER, PATH)
	require.NoError(t, err)

	data := Task{
		FilePath:  "fp",
		FileURL:   "fu",
		ResultURL: "ru",
		TaskID:    "123",
		Status:    "ok",
		Results:   []byte{0x00, 0x01, 0xFF, 0x10},
	}

	id, err := db.Insert(data)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)
	data.ID = id

	// проверяем что данные вставлены правильно
	query := `
		SELECT id, file_path, file_url, result_url, task_id, status, results
		FROM requests
		WHERE id = ?
	`

	var result Task
	err = db.conn.QueryRow(query, id).Scan(
		&result.ID,
		&result.FilePath,
		&result.FileURL,
		&result.ResultURL,
		&result.TaskID,
		&result.Status,
		&result.Results,
	)

	assert.Equal(t, data, result)
}

func TestGetByID(t *testing.T) {

	err := os.Remove(PATH)
	if err != nil {
		require.ErrorIs(t, err, os.ErrNotExist)
	}

	db, err := New(DRIVER, PATH)
	require.NoError(t, err)

	data := Task{
		FilePath:  "fp",
		FileURL:   "fu",
		ResultURL: "ru",
		TaskID:    "123",
		Status:    "ok",
		Results:   []byte{0x00, 0x01, 0xFF, 0x10},
	}
	other_1 := Task{
		FilePath:  "1",
		FileURL:   "1",
		ResultURL: "1",
		TaskID:    "1",
		Status:    "1",
		Results:   []byte{0x00},
	}
	other_2 := Task{
		FilePath:  "2",
		FileURL:   "2",
		ResultURL: "2",
		TaskID:    "2",
		Status:    "2",
		Results:   []byte{0x001},
	}

	id, err := db.Insert(data)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)
	data.ID = id

	_, err = db.Insert(other_1)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)

	_, err = db.Insert(other_2)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)

	result, err := db.GetByID(id)

	assert.Equal(t, data, result)
}

func TestGetByTaskID(t *testing.T) {

	err := os.Remove(PATH)
	if err != nil {
		require.ErrorIs(t, err, os.ErrNotExist)
	}

	db, err := New(DRIVER, PATH)
	require.NoError(t, err)

	data := Task{
		FilePath:  "fp",
		FileURL:   "fu",
		ResultURL: "ru",
		TaskID:    "123",
		Status:    "ok",
		Results:   []byte{0x00, 0x01, 0xFF, 0x10},
	}
	other_1 := Task{
		FilePath:  "1",
		FileURL:   "1",
		ResultURL: "1",
		TaskID:    "1",
		Status:    "1",
		Results:   []byte{0x00},
	}
	other_2 := Task{
		FilePath:  "2",
		FileURL:   "2",
		ResultURL: "2",
		TaskID:    "2",
		Status:    "2",
		Results:   []byte{0x001},
	}

	id, err := db.Insert(data)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)
	data.ID = id

	_, err = db.Insert(other_1)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)

	_, err = db.Insert(other_2)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)

	result, err := db.GetByTaskID(data.TaskID)

	assert.Equal(t, data, result)
}

func TestRemoveByID(t *testing.T) {

	err := os.Remove(PATH)
	if err != nil {
		require.ErrorIs(t, err, os.ErrNotExist)
	}

	db, err := New(DRIVER, PATH)
	require.NoError(t, err)

	data := Task{
		FilePath: "fp",
		FileURL:  "fu",
		TaskID:   "123",
		Status:   "ok",
		Results:  []byte{0x00, 0x01, 0xFF, 0x10},
	}
	other_1 := Task{
		FilePath: "1",
		FileURL:  "1",
		TaskID:   "1",
		Status:   "1",
		Results:  []byte{0x00},
	}
	other_2 := Task{
		FilePath: "2",
		FileURL:  "2",
		TaskID:   "2",
		Status:   "2",
		Results:  []byte{0x001},
	}

	id, err := db.Insert(data)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)
	data.ID = id

	_, err = db.Insert(other_1)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)

	_, err = db.Insert(other_2)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)

	err = db.RemoveByID(id)
	assert.NoError(t, err)

	_, err = db.GetByID(id)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestRemoveByTaskID(t *testing.T) {

	err := os.Remove(PATH)
	if err != nil {
		require.ErrorIs(t, err, os.ErrNotExist)
	}

	db, err := New(DRIVER, PATH)
	require.NoError(t, err)

	var taskID = "TestRemoveByTaskID"

	data := Task{
		FilePath: "fp",
		FileURL:  "fu",
		TaskID:   taskID,
		Status:   "ok",
		Results:  []byte{0x00, 0x01, 0xFF, 0x10},
	}
	other_1 := Task{
		FilePath: "1",
		FileURL:  "1",
		TaskID:   "1",
		Status:   "1",
		Results:  []byte{0x00},
	}
	other_2 := Task{
		FilePath: "2",
		FileURL:  "2",
		TaskID:   "2",
		Status:   "2",
		Results:  []byte{0x001},
	}

	id, err := db.Insert(data)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)
	data.ID = id

	_, err = db.Insert(other_1)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)

	_, err = db.Insert(other_2)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)

	err = db.RemoveByTaskID(taskID)
	assert.NoError(t, err)

	_, err = db.GetByTaskID(taskID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestUpdate(t *testing.T) {

	err := os.Remove(PATH)
	if err != nil {
		require.ErrorIs(t, err, os.ErrNotExist)
	}

	db, err := New(DRIVER, PATH)
	require.NoError(t, err)

	data := Task{
		FilePath:  "fp",
		FileURL:   "fu",
		ResultURL: "ru",
		TaskID:    "123",
		Status:    "ok",
		Results:   []byte{0x00, 0x01, 0xFF, 0x10},
	}

	id, err := db.Insert(data)
	require.NoError(t, err)
	assert.NotEqual(t, id, 0)

	expected := Task{
		ID:        id,
		FilePath:  data.FilePath,
		FileURL:   "fu_update",
		ResultURL: "ru_update",
		TaskID:    "task_update",
		Status:    "status_update",
		Results:   []byte{0xFF},
	}

	err = db.Update(id, expected.FileURL, expected.ResultURL, expected.TaskID, expected.Status, expected.Results)
	require.NoError(t, err)

	// проверяем что данные обновлены правильно
	result, err := db.GetByID(id)

	assert.Equal(t, expected, result)
}
