package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project/auth"
	"project/sessions"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRegister struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Register(db *pgxpool.Pool) http.HandlerFunc {
	var tmpl = template.Must(template.ParseFiles("./web/page1_register.html"))
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {

			err := tmpl.Execute(w, nil)
			if err != nil {
				log.Printf("failed to execute template:%v", err)
				return
			}
			return
		}

		if r.Method == http.MethodPost {

			var usr UserRegister

			err := json.NewDecoder(r.Body).Decode(&usr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			usr.Login = strings.TrimSpace(usr.Login)

			if usr.Login == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "login cannot be empty"})
				return
			}

			var exists string

			err = db.QueryRow(r.Context(), "SELECT login FROM users WHERE login=$1",
				usr.Login).Scan(&exists)

			if err == nil {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(map[string]string{"error": "existing login"})
				return
			}

			Idsession := sessions.Session()

			bytes := auth.Hash(usr.Password)

			_, err = db.Exec(r.Context(), "INSERT INTO users(login, password, session) VALUES($1,$2,$3)",
				usr.Login, bytes, Idsession)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    Idsession,
				Path:     "/",
				HttpOnly: true,
			})

			json.NewEncoder(w).Encode(map[string]string{"href": "/posts"})
		}
	}
}
