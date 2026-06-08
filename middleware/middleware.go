package middleware

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func MiddleWare(db *pgx.Conn, n http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
