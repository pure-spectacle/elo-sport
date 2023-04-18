package router

import (
	"github.com/gorilla/mux"

	"ronin/services"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/athletes", services.GetAllAthletes).Methods("GET");
	router.HandleFunc("/athlete/{athlete_id}", services.GetAthlete).Methods("GET");
	router.HandleFunc("/athlete", services.CreateAthlete).Methods("POST");
	router.HandleFunc("/athlete/{athlete_id}", services.UpdateAthlete).Methods("PUT");
	router.HandleFunc("/athlete/{athlete_id}", services.DeleteAthlete).Methods("DELETE");

	router.HandleFunc("/bouts", services.GetAllBouts).Methods("GET");
	router.HandleFunc("/bout/{bout_id}", services.GetBout).Methods("GET");
	router.HandleFunc("/bout", services.CreateBout).Methods("POST");
	router.HandleFunc("/bout/{bout_id}", services.UpdateBout).Methods("PUT");
	router.HandleFunc("/bout/{bout_id}", services.DeleteBout).Methods("DELETE");
	return router
}
