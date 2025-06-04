package db

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	instance *sql.DB
	once     sync.Once
)

// This function creates a DB connection
/* This fucntion follows singleton pattern
to ensure only one connection is created after restarting server */
func GetDB() *sql.DB {
	once.Do(func() {
		db_url := os.Getenv("DB_URL")
		if db_url == "" {
			log.Printf("Empty DB_URL in .env file")
			return
		}
		var err error
		instance, err = sql.Open("mysql", db_url)
		if err != nil {
			log.Fatalf("DB connection error: %v\n", err)
			return
		}
		log.Println("Connection opened")

		err = instance.Ping()
		if err != nil {
			log.Printf("DB ping failed %v\n", err)
			return
		}
	})
	return instance
}
