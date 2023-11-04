package repositories

import (
	"ronin/models"

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

func (repo *OutcomeRepository) GetAllOutcomes() ([]models.Outcome, error) {
	var outcomes []models.Outcome
	var tempOutcome models.Outcome

	sqlStmt := `SELECT * FROM outcome`
	rows, err := repo.DB.Queryx(sqlStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(&tempOutcome)
		if err != nil {
			return nil, err
		}
		outcomes = append(outcomes, tempOutcome)
	}

	return outcomes, nil
}

func (repo *OutcomeRepository) GetOutcomeById(id string) (models.Outcome, error) {
	var outcome models.Outcome
	sqlStmt := `SELECT * FROM outcome WHERE outcome_id = $1`
	err := repo.DB.QueryRowx(sqlStmt, id).StructScan(&outcome)
	if err != nil {
		return models.Outcome{}, err
	}
	return outcome, nil
}

func (repo *OutcomeRepository) CreateOutcome(outcome models.Outcome) (models.Outcome, error) {
	sqlStmt := `INSERT INTO outcome (bout_id, winner_id, loser_id, style_id) VALUES ($1, $2, $3, $4) RETURNING outcome_id`
	_, err := repo.DB.Exec(sqlStmt, outcome.BoutId, outcome.WinnerId, outcome.LoserId, outcome.StyleId)
	if err != nil {
		return models.Outcome{}, err
	}
	return outcome, nil
}

func (repo *OutcomeRepository) CreateOutcomeByBoutIdNotDraw(outcome *models.Outcome, boutId string) error {
	sqlStmt := `INSERT INTO outcome (bout_id, winner_id, loser_id, is_draw, style_id) VALUES ($1, $2, $3, $4, $5) RETURNING outcome_id`
	err := repo.DB.QueryRowx(sqlStmt, boutId, outcome.WinnerId, outcome.LoserId, outcome.IsDraw, outcome.StyleId).StructScan(outcome)
	if err != nil {
		return err
	}
	sqlStmt = `UPDATE athlete_record SET wins = wins + 1 WHERE athlete_id = $1`
	_, err = repo.DB.Exec(sqlStmt, outcome.WinnerId)
	if err != nil {
		return err
	}

	sqlStmt = `UPDATE athlete_record SET losses = losses + 1 WHERE athlete_id = $1`
	_, err = repo.DB.Exec(sqlStmt, outcome.LoserId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *OutcomeRepository) CreateOutcomeByBoutIdDraw(outcome *models.Outcome, boutId string) error {
	sqlStmt := `INSERT INTO outcome (bout_id, winner_id, loser_id, is_draw, style_id) VALUES ($1, null, null, true, $2) RETURNING outcome_id`
	err := repo.DB.QueryRowx(sqlStmt, boutId, outcome.StyleId).StructScan(outcome)
	if err != nil {
		return err
	}

	sqlStmt = `UPDATE athlete_record SET draws = draws + 1 WHERE athlete_id = $1`
	_, err = repo.DB.Exec(sqlStmt, outcome.WinnerId)
	if err != nil {
		return err
	}

	sqlStmt = `UPDATE athlete_record SET draws = draws + 1 WHERE athlete_id = $1`
	_, err = repo.DB.Exec(sqlStmt, outcome.LoserId)
	if err != nil {
		return err
	}
	return nil
}
