package main

import (
	"EWallet/api/app/handlers"
	"EWallet/api/pkg/store"
	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	// Загрузка конфигураций
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Подключение к бд
	db := pg.Connect(&pg.Options{
		Addr:     viper.GetString("database.address"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Database: viper.GetString("database.dbname"),
	})
	defer db.Close()

	// Инициализация сторов
	// repository лучше инициализировать в service на будующее
	walletStore := store.NewWalletStore(db)
	transactionStore := store.NewTransactionStore(db)

	// Настройка маршрутов с передачей сторов
	api.SetupRoutes(walletStore, transactionStore)

	// Запуск сервера
	log.Println("Starting server on", viper.GetString("server.host")+":"+viper.GetString("server.port"))
	http.ListenAndServe(viper.GetString("server.host")+":"+viper.GetString("server.port"), nil)

}
