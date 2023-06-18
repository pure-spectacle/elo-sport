package repositories

import (
	"ronin/models"

	"github.com/jmoiron/sqlx"
)

type StyleRepository struct {
	DB *sqlx.DB
}

func NewStyleRepository(db *sqlx.DB) *StyleRepository {
	return &StyleRepository{
		DB: db,
	}
}

func (repo *StyleRepository) GetAllStyles() ([]models.Style, error) {
	var styles []models.Style
	var tempStyle models.Style

	sqlStmt := `SELECT * FROM style`
	rows, err := repo.DB.Queryx(sqlStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(&tempStyle)
		if err != nil {
			return nil, err
		}
		styles = append(styles, tempStyle)
	}

	return styles, nil
}

func (repo *StyleRepository) CreateStyle(style models.Style) error {
	sqlStmt := `INSERT INTO style (style_name) VALUES ($1) RETURNING style_id`
	err := repo.DB.QueryRowx(sqlStmt, style.StyleName).Scan(&style)
	if err != nil {
		return err
	}
	return nil
}

func (repo *StyleRepository) RegisterAthleteToStyle(athleteId int, styleId int) error {
	sqlStmt := `INSERT INTO athlete_style (athlete_id, style_id) VALUES ($1, $2)`
	_, err := repo.DB.Exec(sqlStmt, athleteId, styleId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *StyleRepository) GetCommonStyles(acceptorId, challengerId string) ([]models.Style, error) {
	var styles []models.Style
	sqlStmt := `SELECT s.style_id, s.style_name
	FROM style AS s
	JOIN athlete_style AS as1 ON s.style_id = as1.style_id
	JOIN athlete_style AS as2 ON s.style_id = as2.style_id
	WHERE as1.athlete_id=$1 and as2.athlete_id=$2`
	rows, err := repo.DB.Queryx(sqlStmt, acceptorId, challengerId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var style models.Style
		err = rows.StructScan(&style)
		if err != nil {
			return nil, err
		}
		styles = append(styles, style)
	}
	return styles, nil
}
