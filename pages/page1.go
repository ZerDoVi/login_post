package pages

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/jackc/pgx/v5"
)

type User struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func Page1(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("./web/h.html")
			if err != nil {
				w.WriteHeader(500)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				log.Printf("failed to execute template:%v", err)
				return
			}
		}
		if r.Method == http.MethodPost {
			var usr User
			err := json.NewDecoder(r.Body).Decode(&usr)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			_, err = db.Exec(r.Context(), "INSERT INTO users (user, password) VALUES ($1,$2)", usr.User, usr.Password)
			if err != nil {
				w.WriteHeader(500)
				return
			}
		}
	}
}
