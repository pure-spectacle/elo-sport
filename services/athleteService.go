package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ronin/models"

	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	// "github.com/google/uuid"
)

var dbconn *sqlx.DB

type AthleteUsername struct {
	Username string `json:"username" db:"username"`
}

type AthleteId struct {
	AthleteId int `json:"athleteId" db:"athlete_id"`
}

func GetAllAthleteUsernames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sqlStmt := `SELECT username FROM athlete`
	rows, err := dbconn.Queryx(sqlStmt)

	if err == nil {
		var usernames []string
		var tempUsername AthleteUsername

		for rows.Next() {
			err = rows.StructScan(&tempUsername)
			usernames = append(usernames, tempUsername.Username)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), http.StatusNoContent)
			}
		case nil:
			json.NewEncoder(w).Encode(&usernames)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}
}

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
				http.Error(w, err.Error(), http.StatusNoContent)
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

func GetAthleteByUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var athletes = models.GetAthletes()
	vars := mux.Vars(r)
	username := vars["username"]
	var tempAthlete = models.GetAthlete()
	sqlStmt := `SELECT * FROM athlete where username = $1`
	rows, err := dbconn.Queryx(sqlStmt, username)
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

func IsAuthorizedUser(w http.ResponseWriter, r *http.Request) {
	var athlete models.Athlete
	err := json.NewDecoder(r.Body).Decode(&athlete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var athleteId int
	sqlStmt := `SELECT count(*) FROM athlete where username = $1 and password = $2`
	err = dbconn.QueryRow(sqlStmt, athlete.Username, athlete.Password).Scan(&athleteId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if athleteId == 1 {
		var tempAthlete = models.GetAthlete()
		sqlStmt := `SELECT * FROM athlete where username = $1`
		rows, err := dbconn.Queryx(sqlStmt, athlete.Username)
		for rows.Next() {
			err2 := rows.StructScan(&tempAthlete)
			if err2 != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return

			}
			idObj := AthleteId{AthleteId: tempAthlete.AthleteId}
			json.NewEncoder(w).Encode(&idObj)
		}
	} else {
		json.NewEncoder(w).Encode(false)
	}
}

func CreateAthlete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var athlete models.Athlete
	err := json.NewDecoder(r.Body).Decode(&athlete)
	var athleteId int
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO athlete (first_name, last_name, username, birth_date, email, password)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING athlete_id`
	rows, err := dbconn.Queryx(sqlStatement, athlete.FirstName, athlete.LastName, athlete.Username, athlete.BirthDate, athlete.Email, athlete.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		for rows.Next() {
			err = rows.Scan(&athleteId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sqlStmt := `INSERT INTO athlete_record (athlete_id, wins, losses, draws) VALUES ($1, 0, 0, 0)`
			_, err = dbconn.Queryx(sqlStmt, athleteId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Create AthleteId instance and encode it as JSON
			idObj := AthleteId{AthleteId: athleteId}
			json.NewEncoder(w).Encode(idObj)
		}
	}
}

func UpdateAthlete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var athlete models.Athlete
	err := json.NewDecoder(r.Body).Decode(&athlete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE athlete SET first_name = $2, last_name = $3, username = $4, birth_date = $5, email = $6, password = $7 WHERE athlete_id = $8`
	_, err = dbconn.Queryx(sqlStatement, athlete.FirstName, athlete.LastName, athlete.Username, athlete.BirthDate, athlete.Email, athlete.Password, athlete.AthleteId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Athlete updated successfully")

}

func DeleteAthlete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["athlete_id"]

	sqlStatement := `DELETE FROM athlete WHERE athlete_id = $1`
	_, err := dbconn.Queryx(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Athlete deleted successfully")

}

func GetAthleteRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var record models.AthleteRecord
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	sqlStmt := `SELECT * FROM athlete_record where athlete_id = $1`
	row, err := dbconn.Queryx(sqlStmt, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		for row.Next() {
			err2 := row.StructScan(&record)
			if err2 != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		json.NewEncoder(w).Encode(&record)
	}
}

func FollowAthlete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var follow models.Follow
	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sqlStatement := `INSERT INTO following (follower_id, followed_id) VALUES ($1, $2)`
	_, err = dbconn.Queryx(sqlStatement, follow.FollowerId, follow.FollowedId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Athlete followed successfully")
}

func UnfollowAthlete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	followerId, err := strconv.Atoi(vars["followerId"])
	if err != nil {
		http.Error(w, "Invalid followerId", http.StatusBadRequest)
		return
	}
	followedId, err := strconv.Atoi(vars["followedId"])
	if err != nil {
		http.Error(w, "Invalid followedId", http.StatusBadRequest)
		return
	}
	sqlStatement := `DELETE FROM following WHERE follower_id = $1 AND followed_id = $2`
	_, err = dbconn.Exec(sqlStatement, followerId, followedId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAthletesFollowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var follows []int
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	var tempFollow = models.GetFollow()
	sqlStmt := `SELECT * FROM following where follower_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err2 := rows.StructScan(&tempFollow)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusBadRequest)
			return
		}
		follows = append(follows, tempFollow.FollowedId)
	}

	json.NewEncoder(w).Encode(follows)
}

func SetDB(db *sqlx.DB) {
	dbconn = db
}
