package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
	CREATE TABLE IF NOT EXISTS scheduler (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date CHAR(8) NOT NULL DEFAULT "",
    	title VARCHAR(255) NOT NULL DEFAULT "",
    	comment TEXT NOT NULL DEFAULT "",
    	repeat VARCHAR(128) NOT NULL DEFAULT ""
	);

	CREATE INDEX IF NOT EXISTS index_date ON scheduler(date);
`

var db1 *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db1, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("Error with open %s: %v", dbFile, err)
	}
	if install {
		_, err = db1.Exec(schema)
	}

	return err
}

func Close() error {
	if db1 != nil {
		return db1.Close()
	}
	return nil
}
