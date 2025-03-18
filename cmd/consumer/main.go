package main

import (
	"github.com/Diyarjan/BankSystem/constants"
	"github.com/Diyarjan/BankSystem/pkg/repository/db"
	"github.com/Diyarjan/BankSystem/pkg/repository/listeners"
	"github.com/Diyarjan/BankSystem/third_party/cachePart"
	"github.com/Diyarjan/BankSystem/third_party/database"
	"github.com/Diyarjan/BankSystem/third_party/kafkaPart"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL драйвер
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

	transactionRepo := db.NewTransactionPostgres(dbClient, cacheClient)
	broker := viper.GetString("KAFKA.Brokers")
	groupID := "bank-system-consumer"

	depositCons := kafkaPart.NewConsumer(broker, groupID, []string{constants.Deposit})
	withdrawCons := kafkaPart.NewConsumer(broker, groupID, []string{constants.Withdraw})
	transactionCons := kafkaPart.NewConsumer(broker, groupID, []string{constants.Transfer})

	depositConsumer := listeners.NewDepositConsumer(transactionRepo, groupID, depositCons)
	withdrawConsumer := listeners.NewWithdrawConsumer(transactionRepo, groupID, withdrawCons)
	transferConsumer := listeners.NewTransferConsumer(transactionRepo, groupID, transactionCons)

	go depositConsumer.StartListening()
	go withdrawConsumer.StartListening()
	go transferConsumer.StartListening()

	select {}

}

func InitConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	return viper.ReadInConfig()
}
