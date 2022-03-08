package main

import (
	"github.com/joho/godotenv"
	"github.com/superbkibbles/realestate_complex-api/application"
)

func main() {
	godotenv.Load()
	application.StartApplication()
}
