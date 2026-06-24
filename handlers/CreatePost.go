package handlers

import (
	"encoding/json"
	"net/http"
	"project/middleware"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Posts struct {
	Userpost string
	Login    string
}
type NewPost struct {
	NewUserpost string `json:"posttext"`
}

func CreatePost(db *pgxpool.Pool) http.HandlerFunc {
	var tmpl = template.Must(template.ParseFiles("./web/createpost.html"))
	return func(w http.ResponseWriter, r *http.Request) {

		IdUser := r.Context().Value(middleware.IdUserKey)

		if r.Method == http.MethodGet {

			row, err := db.Query(r.Context(), "SELECT post_text,login FROM posts JOIN users ON posts.id_user=users.id")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer row.Close()

			var PostList []Posts
			for row.Next() {
				var p Posts
				err = row.Scan(&p.Userpost, &p.Login)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				PostList = append(PostList, p)
			}
			if err = row.Err(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var Username string
			err = db.QueryRow(r.Context(), "SELECT login FROM users WHERE id=$1", IdUser).Scan(&Username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			Pagedata := map[string]any{
				"Posts":    PostList,
				"Username": Username,
			}

			err = tmpl.Execute(w, Pagedata)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if r.Method == http.MethodPost {
			var posts NewPost
			err := json.NewDecoder(r.Body).Decode(&posts)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err = db.Exec(r.Context(), "INSERT INTO posts(post_text,id_user) VALUES($1,$2)", posts.NewUserpost, IdUser)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)

		}
	}
}
