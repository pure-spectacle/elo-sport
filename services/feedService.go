package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"ronin/repositories"

	"github.com/gorilla/mux"
	// "github.com/jmoiron/sqlx"
)

func GetFeedByAthleteId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["athlete_id"]

	repo := repositories.NewFeedRepository(dbconn)
	feed, err := repo.GetFeedByAthleteId(id)

	switch err {
	case sql.ErrNoRows:
		log.Println("no rows returns.")
		http.Error(w, err.Error(), 204)
	case nil:
		json.NewEncoder(w).Encode(&feed)
	default:
		http.Error(w, err.Error(), 400)
		return
	}
}
