package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project/auth"
	"project/sessions"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Login(db *pgxpool.Pool) http.HandlerFunc {
	var tmpl = template.Must(template.ParseFiles("./web/page1_login.html"))
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
			var usr UserLogin
			err := json.NewDecoder(r.Body).Decode(&usr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var pass string
			err = db.QueryRow(r.Context(), "SELECT password FROM users WHERE login=$1", usr.Login).Scan(&pass)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect login or password"})
				return
			}

			if !auth.HashCheck(pass, usr.Password) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect login or password"})
				return
			}

			Idsession := sessions.Session()

			_, err = db.Exec(r.Context(), "UPDATE users SET session = $2 WHERE login=$1", usr.Login, Idsession)
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
