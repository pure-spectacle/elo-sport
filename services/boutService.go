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

type OutboundBout struct {
	BoutId              int    `json:"boutId"`
	ChallengerId        int    `json:"challengerId"`
	ChallengerFirstName string `json:"challengerFirstName"`
	ChallengerLastName  string `json:"challengerLastName"`
	Style               string `json:"style"`
	ChallengerScore     int    `json:"challengerScore"`
	AcceptorId          int    `json:"acceptorId"`
	AcceptorFirstName   string `json:"acceptorFirstName"`
	AcceptorLastName    string `json:"acceptorLastName"`
	AcceptorScore       int    `json:"acceptorScore"`
	RefereeId           int    `json:"refereeId"`
	RefereeFirstName    string `json:"refereeFirstName"`
	RefereeLastName     string `json:"refereeLastName"`
}

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
	if bout.ChallengerId != bout.AcceptorId {
		sqlStmt := `INSERT INTO bout (challenger_id, acceptor_id, referee_id, accepted, completed, points, style_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING bout_id`
		err := dbconn.QueryRow(sqlStmt, bout.ChallengerId, bout.AcceptorId, bout.RefereeId, bout.Accepted, bout.Completed, bout.Points, bout.StyleId).Scan(&bout.BoutId)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		sqlOutboundBout := `
			SELECT 
				b.bout_id AS "boutId",
				b.challenger_id AS "challengerId",
				c.first_name AS "challengerFirstName",
				c.last_name AS "challengerLastName",
				s.style_name AS "style",
				cs.score AS "challengerScore",
				b.acceptor_id AS "acceptorId",
				a.first_name AS "acceptorFirstName",
				a.last_name AS "acceptorLastName",
				ascore.score AS "acceptorScore",
				r.athlete_id AS "refereeId",
				r.first_name AS "refereeFirstName",
				r.last_name AS "refereeLastName"
			FROM 
				bout b
			JOIN 
				athlete c ON b.challenger_id = c.athlete_id
			JOIN 
				athlete a ON b.acceptor_id = a.athlete_id
			JOIN 
				athlete_score cs ON b.challenger_id = cs.athlete_id AND b.style_id = cs.style_id
			JOIN 
				athlete_score ascore ON b.acceptor_id = ascore.athlete_id AND b.style_id = ascore.style_id
			JOIN 
				athlete r ON b.referee_id = r.athlete_id
			JOIN 
				style s ON b.style_id = s.style_id
			WHERE 
				b.bout_id = $1;
		`
		var outboundBout OutboundBout
		err = dbconn.QueryRow(sqlOutboundBout, bout.BoutId).Scan(&outboundBout.BoutId, &outboundBout.ChallengerId, &outboundBout.ChallengerFirstName, &outboundBout.ChallengerLastName, &outboundBout.Style, &outboundBout.ChallengerScore, &outboundBout.AcceptorId, &outboundBout.AcceptorFirstName, &outboundBout.AcceptorLastName, &outboundBout.AcceptorScore, &outboundBout.RefereeId, &outboundBout.RefereeFirstName, &outboundBout.RefereeLastName)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		json.NewEncoder(w).Encode(&outboundBout)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := map[string]string{
			"error": "ChallengerId and AcceptorId must be different. You cannot create a bout against yourself.",
		}
		json.NewEncoder(w).
			Encode(errorMessage)
	}
}

func UpdateBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bout = models.GetBout()
	_ = json.NewDecoder(r.Body).Decode(&bout)
	vars := mux.Vars(r)
	id := vars["bout_id"]
	sqlStmt := `UPDATE bout SET challenger_id = $1, acceptor_id = $2, referee_id =$3, accepted = $4, points = $5, style_id = $6 WHERE bout_id = $7`
	_, err := dbconn.Exec(sqlStmt, bout.ChallengerId, bout.AcceptorId, bout.RefereeId, bout.Accepted, bout.Points, bout.StyleId, id)
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

func AcceptBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["bout_id"]
	sqlStmt := `UPDATE bout SET accepted = true WHERE bout_id = $1`
	_, err := dbconn.Exec(sqlStmt, id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&id)
}

func DeclineBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["bout_id"]
	sqlStmt := `UPDATE bout SET accepted = false, completed = true WHERE bout_id = $1`
	_, err := dbconn.Exec(sqlStmt, id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&id)
}

func CompleteBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	boutId := vars["bout_id"]
	refereeId := vars["referee_id"]
	var refereeIdReturned string
	sqlStmt := `SELECT referee_id FROM bout WHERE bout_id = $1`

	err := dbconn.QueryRowx(sqlStmt, boutId).Scan(&refereeIdReturned)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if refereeId == refereeIdReturned {
		sqlStmt := `UPDATE bout SET completed = true WHERE bout_id = $1`
		_, err := dbconn.Exec(sqlStmt, boutId)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		json.NewEncoder(w).Encode(&boutId)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := map[string]string{
			"error": "You are not the referee of this bout.",
		}
		json.NewEncoder(w).Encode(errorMessage)
	}
}

func GetPendingBouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bouts = models.GetBouts()
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	sqlStmt := `SELECT * FROM bout where accepted = false and acceptor_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)

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

func GetIncompleteBouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bouts = models.GetBouts()
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	sqlStmt := `SELECT * FROM bout where completed = false and (challenger_id = $1 or acceptor_id = $1 or referee_id = $1)`
	rows, err := dbconn.Queryx(sqlStmt, id)

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
