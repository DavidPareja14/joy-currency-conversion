package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Connect opens DB connection using env vars
func Connect() error {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || host == "" || name == "" {
		return fmt.Errorf("database configuration missing")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, pass, host, port, name)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// reasonable connection settings
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Wait for db to be ready (optional)
	for i := 0; i < 10; i++ {
		if err = DB.Ping(); err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("could not connect to db: %w", err)
}
