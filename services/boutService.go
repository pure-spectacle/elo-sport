package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"ronin/models"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	// "github.com/jmoiron/sqlx"
)

var dbconn *sqlx.DB

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
		sqlStmt := `INSERT INTO bout (challenger_id, acceptor_id, referee_id, accepted, completed, points, style_id, cancelled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING bout_id`
		err := dbconn.QueryRow(sqlStmt, bout.ChallengerId, bout.AcceptorId, bout.RefereeId, bout.Accepted, bout.Completed, bout.Points, bout.StyleId, false).Scan(&bout.BoutId)
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
				r.last_name AS "refereeLastName",
				s.style_id AS "styleId"
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
		var outboundBout models.OutboundBout
		err = dbconn.QueryRow(sqlOutboundBout, bout.BoutId).Scan(&outboundBout.BoutId, &outboundBout.ChallengerId, &outboundBout.ChallengerFirstName, &outboundBout.ChallengerLastName, &outboundBout.Style, &outboundBout.ChallengerScore, &outboundBout.AcceptorId, &outboundBout.AcceptorFirstName, &outboundBout.AcceptorLastName, &outboundBout.AcceptorScore, &outboundBout.RefereeId, &outboundBout.RefereeFirstName, &outboundBout.RefereeLastName, &outboundBout.StyleId)
		// rows, err = dbconn.Queryx(sqlOutboundBout, bout.BoutId).Scan(&outboundBout)
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
	var bouts []models.OutboundBout
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	sqlStmt := `WITH latest_scores AS (
		SELECT 
			athlete_id, 
			style_id, 
			score, 
			updated_dt,
			ROW_NUMBER() OVER (PARTITION BY athlete_id, style_id ORDER BY updated_dt DESC) as row_num
		FROM athlete_score
	)
	SELECT 
		b.bout_id AS "boutId",
		b.challenger_id AS "challengerId",
		c.first_name AS "challengerFirstName",
		c.last_name AS "challengerLastName",
		s.style_name AS "style",
		s.style_id AS "styleId",
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
		latest_scores cs ON b.challenger_id = cs.athlete_id AND b.style_id = cs.style_id AND cs.row_num = 1
	JOIN 
		latest_scores ascore ON b.acceptor_id = ascore.athlete_id AND b.style_id = ascore.style_id AND ascore.row_num = 1
	JOIN 
		athlete r ON b.referee_id = r.athlete_id
	JOIN 
		style s ON b.style_id = s.style_id
	WHERE 
		b.accepted = false AND b.cancelled = false AND b.completed = false AND (b.challenger_id = $1 OR b.acceptor_id = $1 OR b.referee_id = 6)`
	rows, err := dbconn.Queryx(sqlStmt, id)

	if err == nil {
		var tempBout models.OutboundBout

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
	var bouts []models.OutboundBout
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	sqlStmt := `WITH latest_scores AS (
		SELECT 
			athlete_id, 
			style_id, 
			score, 
			updated_dt,
			ROW_NUMBER() OVER (PARTITION BY athlete_id, style_id ORDER BY updated_dt DESC) as row_num
		FROM athlete_score
	)
	SELECT 
		b.bout_id AS "boutId",
		b.challenger_id AS "challengerId",
		c.first_name AS "challengerFirstName",
		c.last_name AS "challengerLastName",
		s.style_name AS "style",
		s.style_id AS "styleId",
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
		latest_scores cs ON b.challenger_id = cs.athlete_id AND b.style_id = cs.style_id AND cs.row_num = 1
	JOIN 
		latest_scores ascore ON b.acceptor_id = ascore.athlete_id AND b.style_id = ascore.style_id AND ascore.row_num = 1
	JOIN 
		athlete r ON b.referee_id = r.athlete_id
	JOIN 
		style s ON b.style_id = s.style_id
	WHERE 
		b.accepted = true AND b.cancelled = false AND b.completed = false AND (b.challenger_id = $1 OR b.acceptor_id = $1 OR b.referee_id = $1)`
	rows, err := dbconn.Queryx(sqlStmt, id)

	if err == nil {
		var tempBout models.OutboundBout

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

func CancelBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	boutId := vars["bout_id"]
	athleteId := vars["challenger_id"]
	var athleteIdReturned string
	sqlStmt := `SELECT challenger_id FROM bout WHERE bout_id = $1`

	err := dbconn.QueryRowx(sqlStmt, boutId).Scan(&athleteIdReturned)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if athleteId == athleteIdReturned {
		sqlStmt := `UPDATE bout set cancelled=true where bout_id = $1`
		_, err := dbconn.Exec(sqlStmt, boutId)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		json.NewEncoder(w).Encode(&boutId)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := map[string]string{
			"error": "You are not the challenger of this bout.",
		}
		json.NewEncoder(w).Encode(errorMessage)
	}
}
