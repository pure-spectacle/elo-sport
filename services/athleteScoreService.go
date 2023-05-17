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

// type AthleteScoreRepository struct {
// 	dbconn *sqlx.DB
// }

// func NewAthleteRepository(db *sqlx.DB) *AthleteScoreRepository {
// 	return &AthleteScoreRepository{
// 		dbconn: db,
// 	}
// }

// var athleteScoreRepo *repositories.AthleteScoreRepository

// func SetAthleteScoreRepo(r *repositories.AthleteScoreRepository) {
// 	athleteScoreRepo = r
// }

type AthleteScoreService struct{}

type AthleteStyleScore struct {
	Score     string `json:"score" db:"score"`
	StyleName string `json:"styleName" db:"style_name"`
}

func GetAllAthleteScoresByAthleteId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var athleteScores = models.GetAthleteScores()
	vars := mux.Vars(r)
	id := vars["athlete_id"]

	sqlStmt := `SELECT * FROM athlete_score where athlete_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)
	if err == nil {
		defer rows.Close()
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
	var athleteScores []AthleteStyleScore
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	var tempAthleteScore AthleteStyleScore
	sqlStmt := `SELECT a.score, s.style_name
	FROM athlete_score AS a
	JOIN style AS s ON s.style_id = a.style_id
	JOIN (
		SELECT athlete_id, style_id, MAX(updated_dt) AS max_updated_dt
		FROM athlete_score
		GROUP BY athlete_id, style_id
	) AS max_dt ON max_dt.athlete_id = a.athlete_id AND max_dt.style_id = a.style_id AND max_dt.max_updated_dt = a.updated_dt
	WHERE a.athlete_id = $1`
	rows, err := dbconn.Queryx(sqlStmt, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			err = rows.StructScan(&tempAthleteScore)
			athleteScores = append(athleteScores, tempAthleteScore)
		}
		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), http.StatusNoContent)
			}
		case nil:
			json.NewEncoder(w).Encode(&athleteScores)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func GetAthleteScoreByStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var athleteScores = models.GetAthleteScores()
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	style := vars["style_id"]
	var tempAthleteScore = models.GetAthleteScore()
	sqlStmt := `WITH ranked_scores AS (
		SELECT
			athlete_id,
			style_id,
			score,
			updated_dt,
			outcome_id,
			ROW_NUMBER() OVER (PARTITION BY athlete_id, style_id ORDER BY updated_dt DESC) AS rank
		FROM
			athlete_score
		WHERE
			athlete_id = $1 AND style_id = $2
	)
	SELECT
		athlete_id,
		style_id,
		score,
		updated_dt
	FROM
		ranked_scores
	WHERE
		rank = 1`
	rows, err := dbconn.Queryx(sqlStmt, id, style)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			err = rows.StructScan(&tempAthleteScore)
			athleteScores = append(athleteScores, tempAthleteScore)
		}
		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), http.StatusNoContent)
			}
		case nil:
			json.NewEncoder(w).Encode(&athleteScores)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (a *AthleteScoreService) GetAthleteScoreById(athleteId, styleId int) (models.AthleteScore, error) {
	var athleteScore = models.GetAthleteScore()
	sqlStmt := `WITH ranked_scores AS (
		SELECT
			athlete_id,
			style_id,
			score,
			updated_dt,
			outcome_id,
			ROW_NUMBER() OVER (PARTITION BY athlete_id, style_id ORDER BY updated_dt DESC) AS rank
		FROM
			athlete_score
		WHERE
			athlete_id = $1 AND style_id = $2
	)
	SELECT
		athlete_id,
		style_id,
		score,
		updated_dt
	FROM
		ranked_scores
	WHERE
		rank = 1`
	rows, err := dbconn.Queryx(sqlStmt, athleteId, styleId)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			err = rows.StructScan(&athleteScore)
			log.Println("athlete: ", athleteScore)
		}
		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				return athleteScore, err
			}
		case nil:
			// json.NewEncoder(w).Encode(&athleteScores)
			log.Println("Retrieve athlete score successfully.")
			return athleteScore, nil
		default:
			// http.Error(w, err.Error(), 400)
			log.Println("Retrieve athlete score failed.")
		}
	}
	return athleteScore, err
}

// TODO: need to add outcome_id here
func CreateAthleteScore(winnerScore, loserScore models.AthleteScore, isDraw bool, outcomeId int) {
	winnerUpdatedScore, loserUpdatedScore := CalculateScore(winnerScore, loserScore, isDraw)

	sqlStmt := `INSERT INTO athlete_score (score, athlete_id, style_id, outcome_id) VALUES ($1, $2, $3, $4)`
	_, err := dbconn.Exec(sqlStmt, winnerUpdatedScore, winnerScore.AthleteId, winnerScore.StyleId, outcomeId)
	if err != nil {
		log.Println("Update winner athlete score failed.")
		return
	}

	sqlStmt = `INSERT INTO athlete_score (score, athlete_id, style_id, outcome_id) VALUES ($1, $2, $3, $4)`
	_, err = dbconn.Exec(sqlStmt, loserUpdatedScore, loserScore.AthleteId, loserScore.StyleId, outcomeId)
	if err != nil {
		log.Println("Update winner athlete score failed.")
		return
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

func CalculateScore(winnerScore, loserScore models.AthleteScore, isDraw bool) (float64, float64) {
	expectedOutcome1 := 1 / (1 + math.Pow(10, (loserScore.Score-winnerScore.Score)/400))
	expectedOutcome2 := 1 / (1 + math.Pow(10, (winnerScore.Score-loserScore.Score)/400))

	var outcome1, outcome2 float64
	if isDraw {
		outcome1 = 0.5
		outcome2 = 0.5
	} else {
		outcome1 = 1
		outcome2 = 0
	}

	updatedScore1 := winnerScore.Score + K*(outcome1-expectedOutcome1)
	updatedScore2 := loserScore.Score + K*(outcome2-expectedOutcome2)

	return math.Round(updatedScore1), math.Round(updatedScore2)
}
