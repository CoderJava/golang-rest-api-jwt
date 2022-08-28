package main

import (
	"golang-rest-api-jwt/config"
)

func main() {
	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)
}
