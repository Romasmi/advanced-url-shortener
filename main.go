package main

import (
	"fmt"
	"shorturl/internal/config"
	"shorturl/internal/database"
)

func main() {
	envConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		return
	}
	fmt.Println(envConfig)

	dbConn := &database.DbConnection{Config: envConfig}
	err = dbConn.Connect()
	if err != nil {
		fmt.Printf("error connecting to DB: %v\n", err)
	}

}
