package services

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"ronin/models"
	"ronin/repositories"

	"github.com/gorilla/mux"
)

var outcomeRepo *repositories.OutcomeRepository

func SetOutcomeRepo(r *repositories.OutcomeRepository) {
	outcomeRepo = r
}

type OutcomeService struct {
	athleteScoreService *AthleteScoreService
	boutRepository      *repositories.BoutRepository
}

func NewOutcomeService(athleteScoreService *AthleteScoreService) *OutcomeService {
	return &OutcomeService{athleteScoreService: athleteScoreService}
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
		loserScore, loserErr := o.athleteScoreService.GetAthleteScoreById(outcome.LoserId, outcome.StyleId)
		winnerScore, winnerErr := o.athleteScoreService.GetAthleteScoreById(outcome.WinnerId, outcome.StyleId)

		o.athleteScoreService.CreateAthleteScore(winnerScore, loserScore, outcome.IsDraw, outcome.OutcomeId)

		if winnerErr == nil && loserErr == nil {
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

	err := o.boutRepository.CompleteBoutByBoutId(boutId)
	if err != nil {
		return err
	}

	loserScore, loserErr := o.athleteScoreService.GetAthleteScoreById(outcome.LoserId, outcome.StyleId)
	winnerScore, winnerErr := o.athleteScoreService.GetAthleteScoreById(outcome.WinnerId, outcome.StyleId)
	if loserErr != nil || winnerErr != nil {
		return errors.New("Error fetching athlete scores")
	}

	o.athleteScoreService.CreateAthleteScore(winnerScore, loserScore, outcome.IsDraw, outcome.OutcomeId)
	return nil
}
