package handlers

import (
	"encoding/json"
	"net/http"
	"project/middleware"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Posts struct {
	Userpost string `json:"userpost"`
	Login    string `json:"login"`
}

var tmpl = template.Must(template.ParseFiles("./web/createpost.html"))

func CreatePost(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		IdUser := r.Context().Value(middleware.IdUserKey)

		if r.Method == http.MethodPost {
			var posts Posts
			err := json.NewDecoder(r.Body).Decode(&posts)
			if err != nil {
				w.WriteHeader(500)
				return
			}

			_, err = db.Exec(r.Context(), "INSERT INTO posts(post_text,id_user) VALUES($1,$2)", posts.Userpost, IdUser)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}

		if r.Method == http.MethodGet {

			row, err := db.Query(r.Context(), "SELECT post_text,login FROM posts JOIN users ON posts.id_user=users.id")
			if err != nil {
				w.WriteHeader(500)
				return
			}
			defer row.Close()

			var PostList []Posts
			for row.Next() {
				var p Posts
				err = row.Scan(&p.Userpost, &p.Login)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				PostList = append(PostList, p)
			}
			if err = row.Err(); err != nil {
				w.WriteHeader(500)
				return
			}

			err = tmpl.Execute(w, PostList)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			return
		}
	}
}
