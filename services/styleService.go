package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ronin/models"

	"github.com/gorilla/mux"
)

type StyleService struct {
	athleteScoreService *AthleteScoreService
}

func NewStyleService(athleteScoreService *AthleteScoreService) *StyleService {
	return &StyleService{athleteScoreService: athleteScoreService}
}

type RegisterStylesRequest struct {
	AthleteID int   `json:"athleteId"`
	Styles    []int `json:"styles"`
}

func GetAllStyles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var styles = models.GetStyles()

	sqlStmt := `SELECT * FROM style`
	rows, err := dbconn.Queryx(sqlStmt)

	if err == nil {
		var tempStyle = models.GetStyle()

		for rows.Next() {
			err = rows.StructScan(&tempStyle)
			styles = append(styles, tempStyle)
		}

		switch err {
		case sql.ErrNoRows:
			{
				log.Println("no rows returns.")
				http.Error(w, err.Error(), 204)
			}
		case nil:
			json.NewEncoder(w).Encode(&styles)
		default:
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func CreateStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var style = models.GetStyle()
	_ = json.NewDecoder(r.Body).Decode(&style)
	sqlStmt := `INSERT INTO style (style_name) VALUES ($1) RETURNING style_id`
	err := dbconn.QueryRowx(sqlStmt, style.StyleName).StructScan(&style)
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
	_ = json.NewDecoder(r.Body).Decode(&style)
	sqlStmt := `INSERT INTO athlete_style (style_id, athlete_id) VALUES ($1, $2) RETURNING style_id`
	err := dbconn.QueryRowx(sqlStmt, style.StyleId, id).StructScan(&style)
	//call athleteScoreService.go to create the athlete's score to be equal to 400
	intValue := 0
	_, errInt := fmt.Sscan(id, &intValue)
	if errInt != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	createErr := s.athleteScoreService.CreateAthleteScoreUponRegistration(intValue, style.StyleId)
	if createErr == nil {
		json.NewEncoder(w).Encode(&style)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (s *StyleService) RegisterMultipleStylesToAthlete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request RegisterStylesRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	athleteID := request.AthleteID
	styles := request.Styles

	for _, style := range styles {
		var returnedStyleID int
		sqlStmt := `INSERT INTO athlete_style (style_id, athlete_id) VALUES ($1, $2) RETURNING style_id`
		err := dbconn.QueryRowx(sqlStmt, style, athleteID).Scan(&returnedStyleID)
		if err == nil {
			createErr := s.athleteScoreService.CreateAthleteScoreUponRegistration(athleteID, style)
			if createErr != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			json.NewEncoder(w).Encode(&returnedStyleID)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

