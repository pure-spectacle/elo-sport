package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"ronin/models"

	"github.com/gorilla/mux"
	// "github.com/jmoiron/sqlx"
)

func GetFeedByAthleteId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var feed []models.Feed
	vars := mux.Vars(r)
	id := vars["athlete_id"]

	sqlStmt := `SELECT * FROM feed where athlete_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)
	if err == nil {
		defer rows.Close()
		var tempFeed = models.GetFeed()
		for rows.Next() {
			err = rows.StructScan(&tempFeed)
			feed = append(feed, tempFeed)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&feed)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}