package application

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Romasmi/advanced-url-shortener/internal/config"
	"github.com/Romasmi/advanced-url-shortener/internal/database"
	"github.com/Romasmi/advanced-url-shortener/internal/redis"
	"github.com/Romasmi/advanced-url-shortener/internal/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	DbConn    *database.DbConnection
	RedisConn *redis.RedisConnection
	Config    *config.Config
	// TODO add logger
	router *mux.Router
}

func (app *App) InitApp(configPath string) error {
	envConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading Config: %v\n", err)
	}
	app.Config = envConfig

	dbConn := &database.DbConnection{Config: &envConfig.Database}
	if err = dbConn.Connect(); err != nil {
		return fmt.Errorf("error connecting to DB: %v\n", err)
	}

	redisConn := &redis.RedisConnection{Config: &envConfig.Redis}
	redisConn.Connect()

	app.DbConn = dbConn
	app.RedisConn = redisConn
	app.router = mux.NewRouter()
	routes.RegisterRoutes(app.router, app.DbConn.DB, app.RedisConn.Rdb, app.Config)
	return nil
}

func (app *App) OnStop() {
	app.DbConn.Close()
	app.RedisConn.Close()
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
		":"+strconv.Itoa(int(app.Config.Server.Port)),
		handlers.CORS(credentials, methods, origins, headers)(app.router))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
