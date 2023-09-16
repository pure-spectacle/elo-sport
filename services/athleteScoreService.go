package services

import (	"encoding/json"
	"log"
	"math"
	"net/http"
	"ronin/models"
	"ronin/repositories"
	"strconv"

	"github.com/gorilla/mux"
)

var athleteScoreRepo *repositories.AthleteScoreRepository

func SetAthleteScoreRepo(r *repositories.AthleteScoreRepository) {
	athleteScoreRepo = r
}

func NewAthleteScoreService() *AthleteScoreService {
	return &AthleteScoreService{}
}

const K float64 = 32

type AthleteScoreService struct{}

func GetAllAthleteScoresByAthleteId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var athleteScores = models.GetAthleteScores()
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	athleteScores, err := athleteScoreRepo.GetAllAthleteScoresByAthleteId(id)
	if err == nil {
		json.NewEncoder(w).Encode(&athleteScores)
	} else {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}
}

func GetAthleteScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var athleteScores []models.AthleteStyleScore
	vars := mux.Vars(r)
	id := vars["athlete_id"]

	athleteScores, err := athleteScoreRepo.GetAthleteStyleScoresById(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	} else {
		json.NewEncoder(w).Encode(&athleteScores)
	}
}

func GetAthleteScoreByStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var athleteScore models.AthleteScore
	vars := mux.Vars(r)
	idStr := vars["athlete_id"]
	styleStr := vars["style_id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid athlete_id:", err.Error())
		http.Error(w, "Invalid athlete_id", http.StatusBadRequest)
		return
	}

	style, err := strconv.Atoi(styleStr)
	if err != nil {
		log.Println("Invalid style_id:", err.Error())
		http.Error(w, "Invalid style_id", http.StatusBadRequest)
		return
	}
	athleteScore, err = athleteScoreRepo.GetAthleteScoreByStyle(id, style)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	} else {
		json.NewEncoder(w).Encode(&athleteScore)
	}
}

func (a *AthleteScoreService) GetAthleteScoreById(athleteId, styleId int) (models.AthleteScore, error) {
	var athleteScore = models.GetAthleteScore()
	athleteScore, err := athleteScoreRepo.GetAthleteScoreByStyle(athleteId, styleId)
	if err != nil {
		log.Println("Retrieve athlete score by id (athlete and style id) failed.")
		return athleteScore, err
	}

	return athleteScore, err

}

func (a *AthleteScoreService) CreateAthleteScore(winnerScore, loserScore models.AthleteScore, isDraw bool, outcomeId int) {
	winnerUpdatedScore, loserUpdatedScore := CalculateScore(winnerScore, loserScore, isDraw)

	err := athleteScoreRepo.UpdateAthleteScore(int(winnerUpdatedScore), winnerScore.AthleteId, winnerScore.StyleId, outcomeId)
	if err != nil {
		log.Println("Update winner athlete score failed.")
		return
	}

	err = athleteScoreRepo.UpdateAthleteScore(int(loserUpdatedScore), loserScore.AthleteId, loserScore.StyleId, outcomeId)
	if err != nil {
		log.Println("Update loser athlete score failed.")
		return
	}

}

func (a *AthleteScoreService) CreateAthleteScoreUponRegistration(athleteId, styleId int) error {
	err := athleteScoreRepo.CreateAthleteScoreUponRegistration(athleteId, styleId)
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
