package sql

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

var Conn = Sqlite()

func Sqlite() *sql.DB {
	dsn := "file:whatsmeow.db" +
		"?_pragma=foreign_keys=on" +
		"&_pragma=journal_mode=WAL" +
		"&_pragma=synchronous=NORMAL" +
		"&_pragma=cache_size=10000" +
		"&_pragma=locking_mode=EXCLUSIVE" +
		"&_pragma=busy_timeout=5000" +
		"&cache=shared"

	db, _ := sql.Open("sqlite", dsn)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db
}
