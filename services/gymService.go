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

func GetAllGyms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var gyms = models.GetGyms()

	sqlStmt := `SELECT * FROM gym`
	rows, err := dbconn.Queryx(sqlStmt)

	if err == nil {
		var tempGym = models.GetGym()

		for rows.Next() {
			err = rows.StructScan(&tempGym)
			gyms = append(gyms, tempGym)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&gyms)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func CreateGym(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var gym = models.GetGym()
	_ = json.NewDecoder(r.Body).Decode(&gym)
	sqlStmt := `INSERT INTO gym (gym_name, gym_address, gym_city, gym_state, gym_zip, gym_phone, gym_email, gym_website, gym_description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING gym_id`
	err := dbconn.QueryRowx(sqlStmt, gym.Name, gym.Address, gym.City, gym.State, gym.Zip, gym.Phone, gym.Email, gym.Website, gym.Description).StructScan(&gym)

	if err == nil {
		json.NewEncoder(w).Encode(&gym)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func GetGym(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var gyms = models.GetGyms()
	vars := mux.Vars(r)
	id := vars["gym_id"]
	var tempGym = models.GetGym()
	sqlStmt := `SELECT * FROM gym where gym_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)

	if err == nil {
		for rows.Next() {
			err = rows.StructScan(&tempGym)
			gyms = append(gyms, tempGym)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&gyms)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}