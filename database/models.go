package database

import (
	"time"

	"github.com/google/uuid"
)

type Coin struct {
	ID        uuid.UUID
	Coin      string
	Price     int
	Timestamp time.Time
}
