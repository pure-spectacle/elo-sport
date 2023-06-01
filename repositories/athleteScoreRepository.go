package repositories

import (
	"github.com/jmoiron/sqlx"
)

type AthleteScoreRepository struct {
	DB *sqlx.DB
}

func NewAthleteScoreRepository(db *sqlx.DB) *AthleteScoreRepository {
	return &AthleteScoreRepository{
		DB: db,
	}
}
