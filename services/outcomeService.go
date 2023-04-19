package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"ronin/models"

	"github.com/gorilla/mux"
)

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

func CreateOutcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcome = models.GetOutcome()
	_ = json.NewDecoder(r.Body).Decode(&outcome)
	sqlStmt := `INSERT INTO outcome (bout_id, winner_id, loser_id, disputed) VALUES ($1, $2, $3, $4) RETURNING outcome_id`
	err := dbconn.QueryRowx(sqlStmt, outcome.BoutId, outcome.WinnerId, outcome.LoserId, outcome.Disputed).StructScan(&outcome)

	if err == nil {
		json.NewEncoder(w).Encode(&outcome)
	} else {
		http.Error(w, err.Error(), 400)
		return
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

func CreateOutcomeByBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcome = models.GetOutcome()
	vars := mux.Vars(r)
	boutId := vars["bout_id"]
	_ = json.NewDecoder(r.Body).Decode(&outcome)
	sqlStmt := `INSERT INTO outcome (bout_id, winner_id, loser_id, disputed) VALUES ($1, $2, $3, $4) RETURNING outcome_id`
	err := dbconn.QueryRowx(sqlStmt, boutId, outcome.WinnerId, outcome.LoserId, outcome.Disputed).StructScan(&outcome)

	if err == nil {
		json.NewEncoder(w).Encode(&outcome)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

