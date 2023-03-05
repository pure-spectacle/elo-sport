package router

import (
	"github.com/gorilla/mux"

	"ronin/services"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/athletes", services.GetAllAthletes).Methods("GET")
	router.HandleFunc("/athlete/{athlete_id}", services.GetAthlete).Methods("GET")
	router.HandleFunc("/athlete", services.CreateAthlete).Methods("POST")

	return router
}
