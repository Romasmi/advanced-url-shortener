package application

import (
	"fmt"
	"log"
	"net/http"
	"shorturl/internal/config"
	"shorturl/internal/database"
	"shorturl/internal/routes"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	dbConn *database.DbConnection
	config *config.Config
	// TODO add logger
	router *mux.Router
}

func (app *App) InitApp(configPath string) error {
	envConfig, err := config.LoadConfig(configPath)
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
	routes.RegisterRoutes(app.router, app.dbConn.DB)
	return nil
}

func (app *App) OnStop() {
	app.dbConn.Close()
}

func (app *App) Run() {
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	})
	headers := handlers.AllowedHeaders([]string{
		"Content-Type",
		"Authorization",
	})
	origins := handlers.AllowedOrigins([]string{"*"})

	err := http.ListenAndServe(
		":"+strconv.Itoa(int(app.config.Server.Port)),
		handlers.CORS(credentials, methods, origins, headers)(app.router))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
