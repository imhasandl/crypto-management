package database

import (
	"database/sql"
)

type DB struct {
	*sql.DB
}

func InitDatabase(dbURL string) (*DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{
		db,
	}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
