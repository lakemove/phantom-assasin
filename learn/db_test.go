package learn

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type testTable struct {
	Id        string
	Log       string
	Text      []byte
	CreatedAt time.Time
}

func createDB() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", ":memory:")
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(8)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
	//defer db.Close()
	if err != nil {
		return
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS TEST (
			ID VARCHAR(100) PRIMARY KEY,
			TEXT BLOB,
			LOG VARCHAR(100),
			CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	return
}

func updateRows(db *sql.DB, table string, values map[string]interface{}, where map[string]interface{}) (err error) {
	var (
		params       []string
		vals         []interface{}
		where_params []string
		where_vals   []interface{}
	)
	for k, v := range values {
		params = append(params, fmt.Sprintf("%s=:%s", k, k))
		vals = append(vals, v)
	}
	for k, v := range where {
		where_params = append(where_params, fmt.Sprintf("%s=:%s", k, k))
		where_vals = append(where_vals, v)
	}
	q := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(params, ","), strings.Join(where_params, ","))
	_, err = db.Exec(q, append(vals, where_vals...)...) //this reads ugly, only to concat 2 slices
	return
}

func insertRow(db *sql.DB, table string, values map[string]interface{}) (err error) {
	var (
		names  []string
		params []string
		vals   []interface{}
	)
	for k, v := range values {
		names = append(names, k)
		params = append(params, fmt.Sprintf(":%s", k))
		vals = append(vals, v)
	}
	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(names, ","), strings.Join(params, ","))
	_, err = db.Exec(q, vals...)
	return
}

func TestUpdateRow(t *testing.T) {
	db, err := createDB()
	assert.Nil(t, err)
	err = insertRow(db, "TEST", map[string]interface{}{
		"ID":   "id001",
		"LOG":  "log start",
		"TEXT": []byte("text content..."),
	})
	assert.Nil(t, err)
	var rowlog string
	row := db.QueryRow("SELECT LOG FROM TEST WHERE ID='id001'")
	row.Scan(&rowlog)
	assert.Equal(t, "log start", rowlog)
	err = updateRows(db, "TEST", map[string]interface{}{"LOG": "log end"}, map[string]interface{}{"ID": "id001"})
	assert.Nil(t, err)
	row = db.QueryRow("SELECT LOG FROM TEST")
	row.Scan(&rowlog)
	assert.Equal(t, "log end", rowlog)
}
func TestInsertRow(t *testing.T) {
	db, err := createDB()
	assert.Nil(t, err)
	row0 := testTable{
		Id:        "id002",
		Log:       "first record",
		Text:      []byte("text content..."),
		CreatedAt: time.Now().UTC(),
	}
	err = insertRow(db, "TEST", map[string]interface{}{
		"ID":         row0.Id,
		"LOG":        row0.Log,
		"TEXT":       row0.Text,
		"CREATED_AT": row0.CreatedAt,
	})
	assert.Nil(t, err)
	row := db.QueryRow("SELECT ID, LOG, TEXT, CREATED_AT FROM TEST")
	row00 := testTable{}
	err = row.Scan(&row00.Id, &row00.Log, &row00.Text, &row00.CreatedAt)
	assert.Nil(t, err)
	assert.Equal(t, row0, row00)
}
