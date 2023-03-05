package main

import (
	"log"
	"net/http"

	"ronin/router"
	"ronin/services"
	"ronin/utils"

)

func main() {
	log.Println("In Main App")

	var dbconn = utils.GetConnection()
	services.SetDB(dbconn)
	var appRouter = router.CreateRouter()

	log.Println("listening on Port 8000")
	log.Fatal(http.ListenAndServe(":8000", appRouter))
}
