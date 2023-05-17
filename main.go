package main

import (
	"log"
	"net/http"

	"ronin/repositories"
	"ronin/router"
	"ronin/services"
	"ronin/utils"
)

func main() {
	log.Println("In Main App")

	dbconn := utils.GetConnection()
	// services.SetDB(dbconn)

	athleteRepo := repositories.NewAthleteRepository(dbconn)
	feedRepo := repositories.NewFeedRepository(dbconn)

	services.SetAthleteRepo(athleteRepo)
	services.SetFeedRepo(feedRepo)

	var appRouter = router.CreateRouter()

	log.Println("listening on Port 8000")
	log.Fatal(http.ListenAndServe(":8000", appRouter))
}
