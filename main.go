package main

import (
	"fmt"
	"net/http"
	"project/database"
	"project/handlers"
	"project/middleware"
)

func main() {

	db := database.Connect()
	defer db.Close()
	http.HandleFunc("/register", handlers.Register(db))
	http.HandleFunc("/login", handlers.Login(db))
	http.HandleFunc("/feed", middleware.MiddleWare(db, handlers.CreatePost(db)))
	fmt.Println("start!!!!!")
	http.ListenAndServe(":8080", nil)
}
