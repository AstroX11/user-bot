package sql

import (
	"database/sql"
	"errors"
)

const createPrefixTableQuery = `
CREATE TABLE IF NOT EXISTS prefix (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	value TEXT NOT NULL
);
`

func InitPrefixTable() error {
	_, err := Conn.Exec(createPrefixTableQuery)
	if err != nil {
		return err
	}

	_, err = Conn.Exec(`INSERT OR IGNORE INTO prefix (id, value) VALUES (1, "/")`)
	return err
}

func GetPrefix() (string, error) {
	var prefix string
	err := Conn.QueryRow(`SELECT value FROM prefix WHERE id = 1`).Scan(&prefix)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return prefix, nil
}

func SetPrefix(newPrefix string) error {
	_, err := Conn.Exec(`UPDATE prefix SET value = ? WHERE id = 1`, newPrefix)
	return err
}
