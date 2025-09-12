package db

import (
	"database/sql"

	_ "github.com/lib/pq" 
)

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// 立即 ping（early ping）
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
