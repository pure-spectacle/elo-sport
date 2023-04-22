package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"ronin/models"

	"github.com/gorilla/mux"
)

const K float64 = 32

type AthleteScoreService struct {}

func GetAllAthleteScores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var athleteScores = models.GetAthleteScores()

	sqlStmt := `SELECT * FROM athlete_score`
	rows, err := dbconn.Queryx(sqlStmt)

	if err == nil {
		var tempAthleteScore = models.GetAthleteScore()

		for rows.Next() {
			err = rows.StructScan(&tempAthleteScore)
			athleteScores = append(athleteScores, tempAthleteScore)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&athleteScores)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func GetAthleteScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var athleteScores = models.GetAthleteScores()
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	var tempAthleteScore = models.GetAthleteScore()
	sqlStmt := `SELECT * FROM athlete_score where athlete_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)
	if err == nil {
		for rows.Next() {
			err = rows.StructScan(&tempAthleteScore)
			athleteScores = append(athleteScores, tempAthleteScore)
		}
		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&athleteScores)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (a *AthleteScoreService) GetAthleteScoreById(id int) ([]models.AthleteScore, error) {
	var athleteScores = models.GetAthleteScores()
	var tempAthleteScore = models.GetAthleteScore()
	sqlStmt := `SELECT * FROM athlete_score where athlete_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)
	log.Println("err: ", err)
	if err == nil {
		for rows.Next() {
			err = rows.StructScan(&tempAthleteScore)
			log.Println("athlete: ", tempAthleteScore)
			athleteScores = append(athleteScores, tempAthleteScore)
		}
		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				return athleteScores, err
			}
		case nil:
			// json.NewEncoder(w).Encode(&athleteScores)
			log.Println("Retrieve athlete score successfully.")
			return athleteScores, nil
		default:
			// http.Error(w, err.Error(), 400)
			log.Println("Retrieve athlete score failed.")
		}
	}
	return athleteScores, err
}

func UpdateOrCreateAthleteScore(winnerScore, loserScore models.AthleteScore) {
	// w.Header().Set("Content-Type", "application/json")
	// vars := mux.Vars(r)
	// athleteId := vars["athlete_id"]
	var athleteScore = models.GetAthleteScore()
	sqlStmt := `SELECT * FROM athlete_score where athlete_id = $1 and style_id = $2`
	winnerRows, winErr := dbconn.Queryx(sqlStmt, winnerScore.AthleteId, winnerScore.StyleId)
	loserRows, losErr := dbconn.Queryx(sqlStmt, loserScore.AthleteId, loserScore.StyleId)

	winnerUpdatedScore, loserUpdatedScore := CalculateScore(winnerScore, loserScore)
	if winErr == nil {
		for winnerRows.Next() {
			winErr = winnerRows.StructScan(&athleteScore)
		}

		switch winErr {
		case sql.ErrNoRows:
			{
				sqlStmt := `INSERT INTO athlete_score (athlete_id, style_id, score) VALUES ($1, $2, $3)`
				_, err := dbconn.Exec(sqlStmt, winnerScore.AthleteId, winnerScore.StyleId, winnerUpdatedScore)
				if err != nil {
					log.Println("Insert winner athlete score failed.")
					return
				}
			}
		case nil:
			sqlStmt := `UPDATE athlete_score SET score = $1 WHERE athlete_id = $2 and style_id = $3`
			_, err := dbconn.Exec(sqlStmt, winnerUpdatedScore, winnerScore.AthleteId, winnerScore.StyleId)
			if err != nil {
				log.Println("Update winner athlete score failed.")
				return
			}
		default:
			return
		}
	}
	if losErr == nil {
		for loserRows.Next() {
			losErr = loserRows.StructScan(&athleteScore)
		}

		switch losErr {
		case sql.ErrNoRows:
			{
				sqlStmt := `INSERT INTO athlete_score (athlete_id, style_id, score) VALUES ($1, $2, $3)`
				_, err := dbconn.Exec(sqlStmt, loserScore.AthleteId, loserScore.StyleId, loserUpdatedScore)
				if err != nil {
					log.Println("Insert loser athlete score failed.")
					return
				}
			}
		case nil:
			sqlStmt := `UPDATE athlete_score SET score = $1 WHERE athlete_id = $2 and style_id = $3`
			_, err := dbconn.Exec(sqlStmt, loserUpdatedScore, loserScore.AthleteId, loserScore.StyleId)
			if err != nil {
				log.Println("Update loser athlete score failed.")
				return
			}
		default:
			return
		}
	}
}

func (a *AthleteScoreService) CreateAthleteScoreUponRegistration(athleteId, styleId int) error {
	sqlStmt := `INSERT INTO athlete_score (athlete_id, style_id, score) VALUES ($1, $2, $3)`
	_, err := dbconn.Exec(sqlStmt, athleteId, styleId, 400)
	if err != nil {
		log.Println("Insert athlete score failed.")
		return err
	}
	return nil
}

func CalculateScore(winnerScore, loserScore models.AthleteScore) (float64, float64) {
	expectedOutcome1 := 1 / (1 + math.Pow(10, (loserScore.Score-winnerScore.Score)/400))
	expectedOutcome2 := 1 / (1 + math.Pow(10, (winnerScore.Score-loserScore.Score)/400))

	var outcome1, outcome2 float64
	outcome1 = 1
	outcome2 = 0

	updatedScore1 := winnerScore.Score + K*(outcome1-expectedOutcome1)
	updatedScore2 := loserScore.Score + K*(outcome2-expectedOutcome2)

	return math.Round(updatedScore1), math.Round(updatedScore2)
}
