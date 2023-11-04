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
	styleRepo := repositories.NewStyleRepository(dbconn)
	boutRepo := repositories.NewBoutRepository(dbconn)
	outcomeRepo := repositories.NewOutcomeRepository(dbconn)
	athleteScoreRepo := repositories.NewAthleteScoreRepository(dbconn)

	services.SetAthleteRepo(athleteRepo)
	services.SetFeedRepo(feedRepo)
	services.SetStyleRepo(styleRepo)
	services.SetBoutRepo(boutRepo)
	services.SetOutcomeRepo(outcomeRepo)
	services.SetAthleteScoreRepo(athleteScoreRepo)

	outcomeService := services.NewOutcomeService(athleteScoreRepo, boutRepo)

	router.SetOutcomeService(outcomeService)

	var appRouter = router.CreateRouter()

	log.Println("listening on Port 8000")
	log.Fatal(http.ListenAndServe(":8000", appRouter))
}
