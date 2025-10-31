package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Connect opens DB connection using env vars
func Connect() error {
	// Comment the way we are getting the environment variables to use the local database
	/*
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || host == "" || name == "" {
		return fmt.Errorf("database configuration missing")
	}
	*/

	fmt.Println("Fetching DB credentials from Parameter Store...")
	user, pass, host, port, name, err := fetchDBCredentials()
	if err != nil {
		return fmt.Errorf("fetching credentials: %w", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, pass, host, port, name)

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

func fetchDBCredentials() (user, pass, host, port, name string, err error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", "", "", "", "", err
	}

	client := ssm.NewFromConfig(cfg)

	params := []string{
		"/exchange-rate/mysql/host",
		"/exchange-rate/mysql/port",
		"/exchange-rate/mysql/user",
		"/exchange-rate/mysql/password",
	}

	out, err := client.GetParameters(ctx, &ssm.GetParametersInput{
		Names:          params,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", "", "", "", "", err
	}

	values := map[string]string{}
	for _, p := range out.Parameters {
		values[*p.Name] = *p.Value
	}

	return values["/exchange-rate/mysql/user"],
		values["/exchange-rate/mysql/password"],
		values["/exchange-rate/mysql/host"],
		values["/exchange-rate/mysql/port"],
		"exchange_db", nil
}