package main

import (
	"log"
	"server/internal/appcontext"
	"server/internal/controller"
)

func main() {
	appContext := appcontext.BuildAppContext()
	r := controller.SetupRouter(appContext)
	if err := r.Run(":8080"); err != nil {
		log.Panic("Failed to run server", err)
	}
}
