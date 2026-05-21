package main

import (
	"context"
	"fmt"
	"net/http"
	"project/database"
	"project/pages"
)

func main() {

	db := database.Connect()
	defer db.Close(context.Background())
	http.HandleFunc("/login", pages.Page1(db))
	fmt.Println("start!!!!!")
	http.ListenAndServe(":8080", nil)
}
