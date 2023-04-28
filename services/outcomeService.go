package services

import (
	"database/sql"
	"encoding/json"
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
		sqlStmt := `INSERT INTO outcome (bout_id, winner_id, loser_id) VALUES ($1, $2, $3) RETURNING outcome_id`
		err = dbconn.QueryRowx(sqlStmt, outcome.BoutId, outcome.WinnerId, outcome.LoserId).StructScan(&outcome)

		// Get winner id and loser id, and get their athlete scores
		// Create or update the athlete score
		loserScore, loserErr := o.athleteScoreService.GetAthleteScoreById(outcome.LoserId)
		winnerScore, winnerErr := o.athleteScoreService.GetAthleteScoreById(outcome.WinnerId)

		UpdateOrCreateAthleteScore(winnerScore[0], loserScore[0], outcome.IsDraw)

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
	var boutIdInOutcomeTableCount int
	var boutIdInBoutIdTableCount int
	err := dbconn.QueryRowx("SELECT COUNT(*) FROM outcome WHERE bout_id = $1", boutId).Scan(&boutIdInOutcomeTableCount)
	err2 := dbconn.QueryRowx("SELECT COUNT(*) FROM bout WHERE bout_id = $1", boutId).Scan(&boutIdInBoutIdTableCount)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	} else if err2 != nil {
		http.Error(w, err2.Error(), 400)
		return
	}

	if boutIdInOutcomeTableCount == 0 && boutIdInBoutIdTableCount == 1 {
		sqlStmt := `INSERT INTO outcome (bout_id, winner_id, loser_id) VALUES ($1, $2, $3) RETURNING outcome_id`
		loserScore, loserErr := o.athleteScoreService.GetAthleteScoreById(outcome.LoserId)
		winnerScore, winnerErr := o.athleteScoreService.GetAthleteScoreById(outcome.WinnerId)

		UpdateOrCreateAthleteScore(winnerScore[0], loserScore[0], outcome.IsDraw)
		err = dbconn.QueryRowx(sqlStmt, boutId, outcome.WinnerId, outcome.LoserId).StructScan(&outcome)

		if loserErr == nil && winnerErr == nil {
			json.NewEncoder(w).Encode(&outcome)
		} else {
			http.Error(w, err.Error(), 400)
			return
		}
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
