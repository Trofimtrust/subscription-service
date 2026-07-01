package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDB() (*pgx.Conn, error) {
	_ = godotenv.Load()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var conn *pgx.Conn
	var err error

	for i := 0; i < 10; i++ {
		conn, err = pgx.Connect(context.Background(), connString)
		if err == nil {
			fmt.Println("Connected to PostgreSQL")
			return conn, nil
		}

		fmt.Printf("Waiting for PostgreSQL... (%d/10)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
