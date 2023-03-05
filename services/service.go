package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ronin/models"

	// "strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	// "github.com/google/uuid"
)

var dbconn *sqlx.DB

func GetAllAthletes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var athletes = models.GetAthletes()

	sqlStmt := `SELECT * FROM athlete`
	rows, err := dbconn.Queryx(sqlStmt)

	if err == nil {
		var tempAthlete = models.GetAthlete()

		for rows.Next() {
			err = rows.StructScan(&tempAthlete)
			athletes = append(athletes, tempAthlete)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&athletes)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func GetAthlete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var athletes = models.GetAthletes()
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	var tempAthlete = models.GetAthlete()
	sqlStmt := `SELECT * FROM athlete where athlete_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)
	for rows.Next() {
		err2 := rows.StructScan(&tempAthlete)
		athletes = append(athletes, tempAthlete)
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		json.NewEncoder(w).Encode(&athletes)
	}

}

func CreateAthlete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var athlete models.Athlete
	err := json.NewDecoder(r.Body).Decode(&athlete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var maxAthleteId int
	err = dbconn.QueryRowx("SELECT MAX(athlete_id) FROM athlete").Scan(&maxAthleteId)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	athlete.AthleteId = maxAthleteId + 1

	sqlStatement := `INSERT INTO athlete (athlete_id, gym_id, first_name, last_name, username, birth_date, wins, losses)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = dbconn.Queryx(sqlStatement, athlete.AthleteId, athlete.GymId, athlete.FirstName, athlete.LastName, athlete.Username, athlete.BirthDate, athlete.Wins, athlete.Losses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Athlete created successfully")

}

func SetDB(db *sqlx.DB) {
	dbconn = db
}
