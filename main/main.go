package main

import (
	"os"
	"upsizeAPI/restapi"
)

func main() {
	restapi.SetupEnv()
	a := restapi.Api{}
	a.Initialize(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	a.Run(":8000")
}
