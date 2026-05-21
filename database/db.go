package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect() *pgx.Conn {
	gtnv := os.Getenv("DB_PASSWORD")
	cnn := fmt.Sprintf("host=db user=user password=%s dbname=mydb sslmode=disable", gtnv)
	conn, err := pgx.Connect(context.Background(), cnn)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
