package router

import (
	"github.com/gorilla/mux"

	"ronin/services"
)

const base_url = "/api/v1"

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(base_url+"/athletes", services.GetAllAthletes).Methods("GET")
	router.HandleFunc(base_url+"/athlete/{athlete_id}", services.GetAthlete).Methods("GET")
	router.HandleFunc(base_url+"/athlete", services.CreateAthlete).Methods("POST")
	router.HandleFunc(base_url+"/athlete/{athlete_id}", services.UpdateAthlete).Methods("PUT")
	router.HandleFunc(base_url+"/athlete/{athlete_id}", services.DeleteAthlete).Methods("DELETE")
	router.HandleFunc(base_url+"/athlete/all/usernames", services.GetAllAthleteUsernames).Methods("GET")
	router.HandleFunc(base_url+"/athlete/{athlete_id}/record", services.GetAthleteRecord).Methods("GET")
	router.HandleFunc(base_url+"/athlete/authorize", services.IsAuthorizedUser).Methods("POST")

	router.HandleFunc(base_url+"/bouts", services.GetAllBouts).Methods("GET")
	router.HandleFunc(base_url+"/bout/{bout_id}", services.GetBout).Methods("GET")
	router.HandleFunc(base_url+"/bout", services.CreateBout).Methods("POST")
	router.HandleFunc(base_url+"/bout/{bout_id}", services.UpdateBout).Methods("PUT")
	router.HandleFunc(base_url+"/bout/{bout_id}", services.DeleteBout).Methods("DELETE")
	router.HandleFunc(base_url+"/bout/{bout_id}/accept", services.AcceptBout).Methods("PUT")
	router.HandleFunc(base_url+"/bout/{bout_id}/decline", services.DeclineBout).Methods("PUT")
	router.HandleFunc(base_url+"/bout/{bout_id}/complete/{referee_id}", services.CompleteBout).Methods("PUT")
	router.HandleFunc(base_url+"/bout/pending/{athlete_id}", services.GetPendingBouts).Methods("GET")
	router.HandleFunc(base_url+"/bout/incomplete/{athlete_id}", services.GetIncompleteBouts).Methods("GET")

	router.HandleFunc(base_url+"/gyms", services.GetAllGyms).Methods("GET")
	router.HandleFunc(base_url+"/gym", services.CreateGym).Methods("POST")
	router.HandleFunc(base_url+"/gym/{gym_id}", services.GetGym).Methods("GET")

	outcomeService := services.OutcomeService{}
	router.HandleFunc(base_url+"/outcome", outcomeService.CreateOutcome).Methods("POST")
	router.HandleFunc(base_url+"/outcome/{outcome_id}", services.GetOutcome).Methods("GET")
	router.HandleFunc(base_url+"/outcome/bout/{bout_id}", services.GetOutcomeByBout).Methods("GET")
	router.HandleFunc(base_url+"/outcome/bout/{bout_id}", outcomeService.CreateOutcomeByBout).Methods("POST")

	styleService := services.StyleService{}
	router.HandleFunc(base_url+"/styles", services.GetAllStyles).Methods("GET")
	router.HandleFunc(base_url+"/style", services.CreateStyle).Methods("POST")
	router.HandleFunc(base_url+"/style/athlete/{athlete_id}", styleService.RegisterAthleteToStyle).Methods("POST")
	router.HandleFunc(base_url+"/styles/athlete/{athlete_id}", styleService.RegisterMultipleStylesToAthlete).Methods("POST")
	router.HandleFunc(base_url+"/styles/common/{athlete_id}/{challenger_id}", services.GetCommonStyles).Methods("GET")

	router.HandleFunc(base_url+"/score/{athlete_id}", services.GetAthleteScore).Methods("GET")
	router.HandleFunc(base_url+"/score/{athlete_id}/all", services.GetAllAthleteScoresByAthleteId).Methods("GET")
	router.HandleFunc(base_url+"/score/{athlete_id}/style/{style_id}", services.GetAthleteScoreByStyle).Methods("GET")

	return router
}
