package store

import (
	"github.com/go-pg/pg/v10"
	"log"
)

// конфигурация для подключения к БД
type DBConfig struct {
	Address  string
	User     string
	Password string
	Database string
}

// создаём новое подключение к БД
func NewDB(cfg DBConfig) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.Address,
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
	})

	if db == nil {
		log.Fatalf("Failed to connect to database.")
	}

	// Проверка конекта
	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database.")
	return db
}
