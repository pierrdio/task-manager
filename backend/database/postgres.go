package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func DB() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())

	if err != nil {
		fmt.Println("database error:", err)
		return nil, err
	}

	fmt.Println("database connected")
	return pool, nil
}
