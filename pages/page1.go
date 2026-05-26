package pages

import (
	"encoding/json"
	"log"
	"net/http"
	"project/sessions"
	"text/template"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Register(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("./web/page1_register.html")
			if err != nil {
				w.WriteHeader(500)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				log.Printf("failed to execute template:%v", err)
				return
			}
			return
		}
		if r.Method == http.MethodPost {
			var usr User
			err := json.NewDecoder(r.Body).Decode(&usr)
			if err != nil {
				w.WriteHeader(400)
				return
			}
			var exists string
			err = db.QueryRow(r.Context(), "SELECT login FROM users WHERE login=$1", usr.Login).Scan(&exists)
			if err == nil {
				w.WriteHeader(409)
				return
			}

			_, err = db.Exec(r.Context(), "INSERT INTO users(login, password) VALUES($1,$2)", usr.Login, usr.Password)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			Idsession := sessions.Session()

			_, err = db.Exec(r.Context(),
				"INSERT INTO sessions(session_id, login) VALUES($1,$2) ON CONFLICT (login) DO UPDATE SET session_id = EXCLUDED.session_id",
				Idsession, usr.Login)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: Idsession,
				Path:  "/",
			})
		}
	}
}

func Login(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("./web/page1_login.html")
			if err != nil {
				w.WriteHeader(500)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				log.Printf("failed to execute template:%v", err)
				return
			}
			return
		}
		if r.Method == http.MethodPost {
			var usr User
			err := json.NewDecoder(r.Body).Decode(&usr)
			if err != nil {
				w.WriteHeader(400)
				return
			}
			var pass string
			err = db.QueryRow(r.Context(), "SELECT password FROM users WHERE login=$1", usr.Login).Scan(&pass)
			if err != nil {
				w.WriteHeader(404)
				return
			}
			if pass != usr.Password {
				http.Error(w, "wrong password", 401)
				return
			}
			Idsession := sessions.Session()

			_, err = db.Exec(r.Context(),
				"INSERT INTO sessions(session_id, login) VALUES($1,$2) ON CONFLICT (login) DO UPDATE SET session_id = EXCLUDED.session_id",
				Idsession, usr.Login)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: Idsession,
				Path:  "/",
			})
		}
	}
}
