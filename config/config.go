package config

import (
	"flag"
	"fmt"
	"gotoko/data"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv string
	AppPort string
}
type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}
func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("welcome to "+ appConfig.AppName)
	
	server.Router = mux.NewRouter()
	server.initializeRoutes()
	server.InitializeDB(dbConfig)
}
func (server *Server) InitializeDB(dbConfig DBConfig) {
	var err error
	dsn :=  fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed on connecting the database server")
	}
	
}
func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func Run() {
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

	err := godotenv .Load()
	if err !=nil {
		log.Fatalf("Error on loading .env file")
	}

	appConfig.AppName = os.Getenv("APP_NAME")
	appConfig.AppEnv = os.Getenv("APP_ENV")
	appConfig.AppPort = os.Getenv("APP_PORT")

	dbConfig.DBHost = os.Getenv("DB_HOST")
	dbConfig.DBUser = os.Getenv("DB_USER")
	dbConfig.DBPassword = os.Getenv("DB_PASSWORD")
	dbConfig.DBName = os.Getenv("DB_NAME")
	dbConfig.DBPort = os.Getenv("DB_PORT")
	
	server.Initialize(appConfig, dbConfig)

	// check if migrate flag is set
	migrate := flag.Bool("migrate", false, "migrate database tables")
	flag.Parse()
	if *migrate {
		for _, model := range data.RegisterModels() {
			err = server.DB.Debug().AutoMigrate(model.Model)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Database migrated success")
	}
	server.Run(":" + appConfig.AppPort)
}