package repositories

import (
	"ronin/models"

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

type AthleteScoreService struct{}

func (repo *AthleteScoreRepository) GetAllAthleteScores() ([]models.AthleteScore, error) {
	var athleteScores []models.AthleteScore
	var tempAthleteScore models.AthleteScore

	sqlStmt := `SELECT * FROM athlete_score`
	rows, err := repo.DB.Queryx(sqlStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(&tempAthleteScore)
		if err != nil {
			return nil, err
		}
		athleteScores = append(athleteScores, tempAthleteScore)
	}

	return athleteScores, nil
}

func (repo *AthleteScoreRepository) GetAthleteStyleScoresById(id string) ([]models.AthleteStyleScore, error) {
	var athleteScores []models.AthleteStyleScore
	sqlStmt := `SELECT a.score, s.style_name
	FROM athlete_score AS a
	JOIN style AS s ON s.style_id = a.style_id
	JOIN (
		SELECT athlete_id, style_id, MAX(updated_dt) AS max_updated_dt
		FROM athlete_score
		GROUP BY athlete_id, style_id
	) AS max_dt ON max_dt.athlete_id = a.athlete_id AND max_dt.style_id = a.style_id AND max_dt.max_updated_dt = a.updated_dt
	WHERE a.athlete_id = $1`
	err := repo.DB.Select(&athleteScores, sqlStmt, id)
	if err != nil {
		return nil, err
	}
	return athleteScores, nil

}

func (repo *AthleteScoreRepository) GetAllAthleteScoresByAthleteId(id string) ([]models.AthleteScore, error) {
	var athleteScores []models.AthleteScore

	sqlStmt := `SELECT * FROM athlete_score where athlete_id = $1`
	err := repo.DB.QueryRowx(sqlStmt, id).StructScan(&athleteScores)
	if err != nil {
		return nil, err
	}

	return athleteScores, nil
}

func (repo *AthleteScoreRepository) GetAthleteScoreByStyle(id int, style int) (models.AthleteScore, error) {
	var athleteScore models.AthleteScore

	sqlStmt := `WITH ranked_scores AS (
		SELECT
			athlete_id,
			style_id,
			score,
			updated_dt,
			outcome_id,
			ROW_NUMBER() OVER (PARTITION BY athlete_id, style_id ORDER BY updated_dt DESC) AS rank
		FROM
			athlete_score
		WHERE
			athlete_id = $1 AND style_id = $2
	)
	SELECT
		athlete_id,
		style_id,
		score,
		updated_dt
	FROM
		ranked_scores
	WHERE
		rank = 1`
	err := repo.DB.QueryRowx(sqlStmt, id, style).StructScan(&athleteScore)
	if err != nil {
		return models.AthleteScore{}, err
	}

	return athleteScore, nil
}

func (repo *AthleteScoreRepository) UpdateAthleteScore(score int, athleteId int, styleId int, outcomeId int) error {
	sqlStmt := `INSERT INTO athlete_score (score, athlete_id, style_id, outcome_id) VALUES ($1, $2, $3, $4)`
	_, err := repo.DB.Exec(sqlStmt, score, athleteId, styleId, outcomeId)
	if err != nil {
		return err
	}

	return nil
}

func (repo *AthleteScoreRepository) CreateAthleteScoreUponRegistration(athleteId int, styleId int) error {
	sqlStmt := `INSERT INTO athlete_score (athlete_id, style_id, score) VALUES ($1, $2, $3)`
	_, err := repo.DB.Exec(sqlStmt, athleteId, styleId, 400)
	if err != nil {
		return err
	}

	return nil
}
