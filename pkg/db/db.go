package db

import (
	"os"
	"database/sql"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func Init(dbfile string) error {
	_, err := os.Stat(dbfile)
	var install bool
	if err != nil {
		install = true
	}
	db, err = sql.Open("sqlite", dbfile)
	if err != nil {
		return err
	}
	if install {
		schema := `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(64) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(64) NOT NULL DEFAULT ""
);
CREATE INDEX date_idx ON scheduler (date);
`
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
