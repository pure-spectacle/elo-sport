package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"ronin/models"

	"github.com/gorilla/mux"
)

type OutcomeService struct {
	athleteScoreService *AthleteScoreService
}

func NewOutcomeService(athleteScoreService *AthleteScoreService) *OutcomeService {
	return &OutcomeService{athleteScoreService: athleteScoreService}
}

func GetAllOutcomes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var outcomes = models.GetOutcomes()

	sqlStmt := `SELECT * FROM outcome`
	rows, err := dbconn.Queryx(sqlStmt)

	if err == nil {
		var tempOutcome = models.GetOutcome()

		for rows.Next() {
			err = rows.StructScan(&tempOutcome)
			outcomes = append(outcomes, tempOutcome)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&outcomes)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func GetOutcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcomes = models.GetOutcomes()
	vars := mux.Vars(r)
	id := vars["outcome_id"]
	var tempOutcome = models.GetOutcome()
	sqlStmt := `SELECT * FROM outcome where outcome_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)

	if err == nil {
		for rows.Next() {
			err = rows.StructScan(&tempOutcome)
			outcomes = append(outcomes, tempOutcome)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&outcomes)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (o *OutcomeService) CreateOutcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcome = models.GetOutcome()
	_ = json.NewDecoder(r.Body).Decode(&outcome)

	// Check if bout_id already exists in the outcome table
	var count int
	err := dbconn.QueryRowx("SELECT COUNT(*) FROM outcome WHERE bout_id = $1", outcome.BoutId).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if count < 1 {
		sqlStmt := `INSERT INTO outcome (bout_id, winner_id, loser_id, style_id) VALUES ($1, $2, $3, $4) RETURNING outcome_id`
		err = dbconn.QueryRowx(sqlStmt, outcome.BoutId, outcome.WinnerId, outcome.LoserId, outcome.StyleId).StructScan(&outcome)

		// Get winner id and loser id, and get their athlete scores
		// Create or update the athlete score
		loserScore, loserErr := o.athleteScoreService.GetAthleteScoreById(outcome.LoserId, outcome.StyleId)
		winnerScore, winnerErr := o.athleteScoreService.GetAthleteScoreById(outcome.WinnerId, outcome.StyleId)

		CreateAthleteScore(winnerScore, loserScore, outcome.IsDraw)

		if winnerErr == nil && loserErr == nil {
			json.NewEncoder(w).Encode(&outcome)
		} else {
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		// Send a notification to the user when an outcome already exists for the bout
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := map[string]string{
			"error": "An outcome already exists for this bout. Please create another bout if you would like to challenge your opponent again.",
		}
		json.NewEncoder(w).Encode(errorMessage)
	}
}

func GetOutcomeByBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcomes = models.GetOutcomes()
	vars := mux.Vars(r)
	id := vars["bout_id"]
	var tempOutcome = models.GetOutcome()
	sqlStmt := `SELECT * FROM outcome where bout_id = $1`

	rows, err := dbconn.Queryx(sqlStmt, id)

	if err == nil {
		for rows.Next() {
			err = rows.StructScan(&tempOutcome)
			outcomes = append(outcomes, tempOutcome)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&outcomes)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (o *OutcomeService) CreateOutcomeByBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcome = models.GetOutcome()
	vars := mux.Vars(r)
	boutId := vars["bout_id"]
	_ = json.NewDecoder(r.Body).Decode(&outcome)

	// Check if bout_id already exists in the outcome table
	var boutIdInOutcomeTableCount, boutIdInBoutIdTableCount int
	err := dbconn.QueryRowx("SELECT COUNT(*) FROM outcome WHERE bout_id = $1", boutId).Scan(&boutIdInOutcomeTableCount)
	err2 := dbconn.QueryRowx("SELECT COUNT(*) FROM bout WHERE bout_id = $1", boutId).Scan(&boutIdInBoutIdTableCount)
	if err != nil || err2 != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if boutIdInOutcomeTableCount == 0 && boutIdInBoutIdTableCount == 1 {
		err := o.insertOutcomeAndUpdateAthleteScores(&outcome, boutId)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		json.NewEncoder(w).Encode(&outcome)
	} else if boutIdInOutcomeTableCount > 0 {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := map[string]string{
			"error": "An outcome already exists for this bout. Please create another bout if you would like to challenge your opponent again.",
		}
		json.NewEncoder(w).Encode(errorMessage)
	} else if boutIdInBoutIdTableCount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := map[string]string{
			"error": "No bout was found.",
		}
		json.NewEncoder(w).Encode(errorMessage)
	}
}

func (o *OutcomeService) insertOutcomeAndUpdateAthleteScores(outcome *models.Outcome, boutId string) error {
	var sqlStmt string
	if !outcome.IsDraw {
		sqlStmt = `INSERT INTO outcome (bout_id, winner_id, loser_id, is_draw, style_id) VALUES ($1, $2, $3, $4, $5) RETURNING outcome_id`
		err := dbconn.QueryRowx(sqlStmt, boutId, outcome.WinnerId, outcome.LoserId, outcome.IsDraw, outcome.StyleId).StructScan(outcome)
		if err != nil {
			return err
		}

		sqlStmt = `UPDATE athlete_record SET wins = wins + 1 WHERE athlete_id = $1`
		_, err = dbconn.Exec(sqlStmt, outcome.WinnerId)
		if err != nil {
			return err
		}

		sqlStmt = `UPDATE athlete_record SET losses = losses + 1 WHERE athlete_id = $1`
		_, err = dbconn.Exec(sqlStmt, outcome.LoserId)
		if err != nil {
			return err
		}
	} else {
		sqlStmt = `INSERT INTO outcome (bout_id, winner_id, loser_id, is_draw, style_id) VALUES ($1, null, null, true, $2) RETURNING outcome_id`
		err := dbconn.QueryRowx(sqlStmt, boutId, outcome.StyleId).StructScan(outcome)
		if err != nil {
			return err
		}

		sqlStmt = `UPDATE athlete_record SET draws = draws + 1 WHERE athlete_id = $1`
		_, err = dbconn.Exec(sqlStmt, outcome.WinnerId)
		if err != nil {
			return err
		}

		sqlStmt = `UPDATE athlete_record SET draws = draws + 1 WHERE athlete_id = $1`
		_, err = dbconn.Exec(sqlStmt, outcome.LoserId)
		if err != nil {
			return err
		}
	}

	updateStatement := `UPDATE bout SET completed = true WHERE bout_id = $1`
	_, err := dbconn.Exec(updateStatement, boutId)
	if err != nil {
		return err
	}

	loserScore, loserErr := o.athleteScoreService.GetAthleteScoreById(outcome.LoserId, outcome.StyleId)
	winnerScore, winnerErr := o.athleteScoreService.GetAthleteScoreById(outcome.WinnerId, outcome.StyleId)
	if loserErr != nil || winnerErr != nil {
		return errors.New("Error fetching athlete scores")
	}

	CreateAthleteScore(winnerScore, loserScore, outcome.IsDraw)
	return nil
}
