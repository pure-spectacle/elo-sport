package repositories

import "github.com/jmoiron/sqlx"

type StyleRepository struct {
	DB *sqlx.DB
}

func NewStyleRepository(db *sqlx.DB) *StyleRepository {
	return &StyleRepository{
		DB: db,
	}
}
