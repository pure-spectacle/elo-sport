package services

import (
	"encoding/json"
	"log"
	"net/http"
	"ronin/models"
	"ronin/repositories"
	"strconv"

	"github.com/gorilla/mux"
)

var styleRepo *repositories.StyleRepository

func SetStyleRepo(r *repositories.StyleRepository) {
	styleRepo = r
}

type StyleService struct {
	athleteScoreService *AthleteScoreService
}

func NewStyleService(athleteScoreService *AthleteScoreService) *StyleService {
	return &StyleService{athleteScoreService: athleteScoreService}
}

func GetAllStyles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var styles = models.GetStyles()
	styles, err := styleRepo.GetAllStyles()
	if err == nil {
		json.NewEncoder(w).Encode(&styles)
	} else {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}
}

func CreateStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var style = models.GetStyle()
	_ = json.NewDecoder(r.Body).Decode(&style)
	err := styleRepo.CreateStyle(style)
	if err == nil {
		json.NewEncoder(w).Encode(&style)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (s *StyleService) RegisterAthleteToStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var style = models.GetStyle()
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	intAthleteId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewDecoder(r.Body).Decode(&style)
	err = styleRepo.RegisterAthleteToStyle(intAthleteId, style.StyleId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//call athleteScoreService.go to create the athlete's score to be equal to 400
	createErr := s.athleteScoreService.CreateAthleteScoreUponRegistration(intAthleteId, style.StyleId)
	if createErr == nil {
		json.NewEncoder(w).Encode(&style)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (s *StyleService) RegisterMultipleStylesToAthlete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request models.RegisterStylesRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	athleteID := request.AthleteID
	styles := request.Styles

	for _, style := range styles {
		err := styleRepo.RegisterAthleteToStyle(athleteID, style)
		if err == nil {
			createErr := s.athleteScoreService.CreateAthleteScoreUponRegistration(athleteID, style)
			if createErr != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

}

func GetCommonStyles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	acceptorId := vars["athlete_id"]
	challengerId := vars["challenger_id"]

	var styles = models.GetStyles()

	styles, err := styleRepo.GetCommonStyles(acceptorId, challengerId)
	if err == nil {
		json.NewEncoder(w).Encode(&styles)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}
