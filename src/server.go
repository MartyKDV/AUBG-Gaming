package src

import (
	"main/src/controllers"
	"os"

	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	server.Initialise(os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME"))

	server.Run(":8080")
}
