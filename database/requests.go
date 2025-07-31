package database

import (
	"time"

	"github.com/google/uuid"
)

func (db *DB) SaveCoinPrice(coin string, price int, ts time.Time) error {
	query := `INSERT INTO coins (id, coin, price, timestamp) VALUES ($1, $2, $3, $4)`
	
	_, err := db.Exec(query, uuid.New(), coin, price, ts)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetNearestPrice(coin string, ts time.Time) (int, time.Time, error) {
	query := `SELECT price, timestamp FROM coins WHERE coin=$1 ORDER BY ABS(EXTRACT(EPOCH FROM timestamp) - $2) LIMIT 1`
	row := db.QueryRow(query, coin, ts.Unix())

	var price int
	var timestamp time.Time

	err := row.Scan(&price, &timestamp)
	if err != nil {
		return 0, time.Time{}, err
	}
	
	return price, timestamp, nil
}
