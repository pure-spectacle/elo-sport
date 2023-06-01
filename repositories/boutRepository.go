package repositories

import (
	"github.com/jmoiron/sqlx"
)

type BoutRepository struct {
	DB *sqlx.DB
}

func NewBoutRepository(db *sqlx.DB) *BoutRepository {
	return &BoutRepository{
		DB: db,
	}
}
