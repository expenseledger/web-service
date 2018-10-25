package main

import (
	"log"

	"github.com/expenseledger/web-service/config"
	"github.com/expenseledger/web-service/controller"
	"github.com/expenseledger/web-service/database"
)

func main() {
	var err error

	err = database.CreateTables()
	if err != nil {
		log.Fatal("Error creating tables")
	} else {
		log.Println("Successfully created tables")
	}

	router := controller.InitRoutes()

	err = router.Run(":" + config.Port)
	if err != nil {
		log.Fatal("Error running the server", err)
	}
}
