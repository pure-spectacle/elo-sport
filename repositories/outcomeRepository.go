package repositories

import (
	"github.com/jmoiron/sqlx"
)

type OutcomeRepository struct {
	DB *sqlx.DB
}

func NewOutcomeRepository(db *sqlx.DB) *OutcomeRepository {
	return &OutcomeRepository{
		DB: db,
	}
}
