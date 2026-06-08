package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	gtnv := os.Getenv("DB_PASSWORD")
	cnn := fmt.Sprintf("host=db user=user password=%s dbname=mydb sslmode=disable", gtnv)
	conn, err := pgxpool.New(context.Background(), cnn)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
