package services

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"ronin/models"
	"ronin/repositories"

	"github.com/gorilla/mux"
)

const K float64 = 32

var outcomeRepo *repositories.OutcomeRepository

func SetOutcomeRepo(r *repositories.OutcomeRepository) {
	outcomeRepo = r
}

type OutcomeService struct {
	boutRepository         *repositories.BoutRepository
	athleteScoreRepository *repositories.AthleteScoreRepository
}

func NewOutcomeService(athleteScoreRepository *repositories.AthleteScoreRepository, boutRepo *repositories.BoutRepository) *OutcomeService {
	return &OutcomeService{
		athleteScoreRepository: athleteScoreRepository,
		boutRepository:         boutRepo,
	}
}

func GetAllOutcomes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var outcomes = models.GetOutcomes()

	outcomes, err := outcomeRepo.GetAllOutcomes()
	if err == nil {
		json.NewEncoder(w).Encode(&outcomes)
	} else {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

}

func GetOutcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcomes = models.GetOutcome()
	vars := mux.Vars(r)
	id := vars["outcome_id"]
	outcomes, err := outcomeRepo.GetOutcomeById(id)
	if err == nil {
		json.NewEncoder(w).Encode(&outcomes)
	} else {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}
}

func (o *OutcomeService) CreateOutcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcome = models.GetOutcome()
	_ = json.NewDecoder(r.Body).Decode(&outcome)

	outcome, err := outcomeRepo.CreateOutcome(outcome)
	if err == nil {
		loserScore, loserErr := o.athleteScoreRepository.GetAthleteScoreByStyle(outcome.LoserId, outcome.StyleId)
		if loserErr != nil {
			http.Error(w, loserErr.Error(), 400)
			return
		}
		winnerScore, winnerErr := o.athleteScoreRepository.GetAthleteScoreByStyle(outcome.WinnerId, outcome.StyleId)
		if winnerErr != nil {
			http.Error(w, winnerErr.Error(), 400)
			return
		}
		// o.athleteScoreService.CreateAthleteScore(winnerScore, loserScore, outcome.IsDraw, outcome.OutcomeId)
		winnerUpdatedScore, loserUpdatedScore := CalculateScoreAfterOutcome(winnerScore, loserScore, false)
		err = o.athleteScoreRepository.UpdateAthleteScore(int(winnerUpdatedScore), winnerScore.AthleteId, winnerScore.StyleId, outcome.OutcomeId)
		if err == nil {
			json.NewEncoder(w).Encode(&outcome)
		} else {
			http.Error(w, err.Error(), 400)
			return
		}
		err = o.athleteScoreRepository.UpdateAthleteScore(int(loserUpdatedScore), loserScore.AthleteId, loserScore.StyleId, outcome.OutcomeId)
		if err == nil {
			json.NewEncoder(w).Encode(&outcome)
		} else {
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func GetOutcomeByBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcomes = models.GetOutcome()
	vars := mux.Vars(r)
	id := vars["outcome_id"]
	outcomes, err := outcomeRepo.GetOutcomeById(id)
	if err == nil {
		json.NewEncoder(w).Encode(&outcomes)
	} else {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}
}

func (o *OutcomeService) CreateOutcomeByBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var outcome = models.GetOutcome()
	vars := mux.Vars(r)
	boutId := vars["bout_id"]
	_ = json.NewDecoder(r.Body).Decode(&outcome)

	err := o.insertOutcomeAndUpdateAthleteScores(&outcome, boutId)
	if err == nil {
		json.NewEncoder(w).Encode(&outcome)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (o *OutcomeService) insertOutcomeAndUpdateAthleteScores(outcome *models.Outcome, boutId string) error {
	// log.Printf("BoutRepository in OutcomeService: %v\n", o.boutRepository.DB)

	exists, err := outcomeRepo.DoesOutcomeExistByBoutId(boutId)
	if err != nil {
		return err
	}
	if exists {
		err := o.boutRepository.DeleteBout(boutId)
		if err != nil {
			log.Printf("Failed to delete bout: %v\n", err)
		}
		return errors.New("Outcome already exists for bout.")
	}
	if !outcome.IsDraw {
		err := outcomeRepo.CreateOutcomeByBoutIdNotDraw(outcome, boutId)
		if err != nil {
			return err
		}
	} else {
		err := outcomeRepo.CreateOutcomeByBoutIdDraw(outcome, boutId)
		if err != nil {
			return err
		}
	}

	loserScore, loserErr := o.athleteScoreRepository.GetAthleteScoreByStyle(outcome.LoserId, outcome.StyleId)
	if loserErr != nil {
		return loserErr
	}
	winnerScore, winnerErr := o.athleteScoreRepository.GetAthleteScoreByStyle(outcome.WinnerId, outcome.StyleId)
	if winnerErr != nil {
		return winnerErr
	}
	winnerUpdatedScore, loserUpdatedScore := CalculateScoreAfterOutcome(winnerScore, loserScore, false)
	err = o.athleteScoreRepository.UpdateAthleteScore(int(winnerUpdatedScore), winnerScore.AthleteId, winnerScore.StyleId, outcome.OutcomeId)
	if err != nil {
		return err
	}
	err = o.athleteScoreRepository.UpdateAthleteScore(int(loserUpdatedScore), loserScore.AthleteId, loserScore.StyleId, outcome.OutcomeId)
	if err != nil {
		return err
	}
	return nil
}

func CalculateScoreAfterOutcome(winnerScore, loserScore models.AthleteScore, isDraw bool) (float64, float64) {
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
