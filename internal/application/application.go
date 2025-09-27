package application

import (
	"fmt"
	"shorturl/internal/config"
	"shorturl/internal/database"
	"shorturl/internal/routes"
)

type App struct {
	dbConn *database.DbConnection
	config *config.Config
	// TODO add logger
}

func (app *App) InitApp() error {
	envConfig, err := config.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("error loading config: %v\n", err)
	}
	app.config = envConfig

	dbConn := &database.DbConnection{Config: envConfig}
	err = dbConn.Connect()
	if err != nil {
		return fmt.Errorf("error connecting to DB: %v\n", err)
	}
	app.dbConn = dbConn
	app.RegisterRoutes()
	return nil
}

func (app *App) OnStop() {
	app.dbConn.Close()
}

func (app *App) RegisterRoutes() {
	routes.RegisterUrlRoutes()
}
