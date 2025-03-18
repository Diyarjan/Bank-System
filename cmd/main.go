package main

import (
	"github.com/Diyarjan/BankSystem/third_party/cachePart"
	"github.com/Diyarjan/BankSystem/third_party/kafkaPart"
	_ "github.com/lib/pq" // PostgreSQL драйвер

	"github.com/Diyarjan/BankSystem/pkg/handler"
	"github.com/Diyarjan/BankSystem/pkg/repository"
	"github.com/Diyarjan/BankSystem/pkg/service"
	"github.com/Diyarjan/BankSystem/third_party/database"
	"github.com/Diyarjan/BankSystem/third_party/server"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {

	if err := InitConfig(); err != nil {
		log.Fatalf("Error initializing configs %s", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	//	init db postgres connection
	dbClient, err := database.NewPostgresDB(database.ConfigDB{
		Host:     viper.GetString("DB.Host"),
		Port:     viper.GetString("DB.Port"),
		UserName: viper.GetString("DB.UserName"),
		Password: os.Getenv("DB_Password"),
		DBName:   viper.GetString("DB.DBName"),
		SSLMode:  viper.GetString("DB.SSLMode"),
	})
	if err != nil {
		log.Fatalf("Error initializing Postgres %s", err)
	}

	// init cachePart connection
	cacheClient, err := cachePart.NewRedis(cachePart.Params{
		Host:     viper.GetString("REDIS.Host"),
		Port:     viper.GetString("REDIS.Port"),
		Password: viper.GetString("REDIS.Password"),
		DB:       viper.GetInt("REDIS.DB"),
	})
	if err != nil {
		log.Fatalf("Error initializing Redis %s", err)
	}

	//init KafkaProducer connection
	producer := kafkaPart.NewProducer(viper.GetString("KAFKA.Brokers"))

	repos := repository.NewRepository(dbClient, cacheClient)
	services := service.NewService(repos, producer, cacheClient)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("PORT"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error starting server - %s", err)
	}

}

func InitConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	return viper.ReadInConfig()
}
