package middleware

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Idcontext int

const IdUserKey Idcontext = 123

func MiddleWare(db *pgxpool.Pool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var IdUser string
		err = db.QueryRow(r.Context(), "SELECT id FROM users WHERE session=$1", cookie.Value).Scan(&IdUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), IdUserKey, IdUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
