package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ronin/models"

	"github.com/gorilla/mux"
	// "github.com/jmoiron/sqlx"
)

func GetAllBouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var bouts = models.GetBouts()

	sqlStmt := `SELECT * FROM bout`
	rows, err := dbconn.Queryx(sqlStmt)

	if err == nil {
		var tempBout = models.GetBout()

		for rows.Next() {
			err = rows.StructScan(&tempBout)
			bouts = append(bouts, tempBout)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&bouts)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func GetBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bouts = models.GetBouts()
	vars := mux.Vars(r)
	id := vars["bout_id"]
	var tempBout = models.GetBout()
	sqlStmt := `SELECT * FROM bout where bout_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)

	if err == nil {
		for rows.Next() {
			err = rows.StructScan(&tempBout)
			bouts = append(bouts, tempBout)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&bouts)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func CreateBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bout = models.GetBout()
	_ = json.NewDecoder(r.Body).Decode(&bout)
	var nextBoutId int = getNextBoudId()
	sqlStmt := `INSERT INTO bout (bout_id, challenger_id, acceptor_id, accepted, points) VALUES ($1, $2, $3, $4, $5)`
	_, err := dbconn.Exec(sqlStmt, nextBoutId, bout.ChallengerId, bout.AcceptorId, bout.Accepted, bout.Points)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&bout)
}

func UpdateBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bout = models.GetBout()
	_ = json.NewDecoder(r.Body).Decode(&bout)
	vars := mux.Vars(r)
	id := vars["bout_id"]
	sqlStmt := `UPDATE bout SET challenger_id = $1, acceptor_id = $2, accepted = $3, points = $4 WHERE bout_id = $5`
	_, err := dbconn.Exec(sqlStmt, bout.ChallengerId, bout.AcceptorId, bout.Accepted, bout.Points, id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&bout)
}

func DeleteBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["bout_id"]
	sqlStmt := `DELETE FROM bout WHERE bout_id = $1`
	_, err := dbconn.Exec(sqlStmt, id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&id)
}

func getNextBoudId() int {
	sqlStmt := `SELECT MAX(bout_id) FROM bout`
	var maxId int
	err := dbconn.Get(&maxId, sqlStmt)
	if err != nil {
		fmt.Println(err)
	}
	if maxId == 0 {
		return 1
	}
	return maxId + 1
}
