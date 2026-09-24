package main

import (
	"backend/database"
	"backend/handlers"
	"fmt"
	"net/http"
	"os"
)

const PORT = ":8080"

func main() {

	pool, err := database.DB()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "task manager server")
	})

	http.HandleFunc("/users", handlers.UsersHandler)
	http.HandleFunc("/users/", handlers.UserHandler)

	fmt.Println("Server started on port", PORT)

	http.ListenAndServe(PORT, nil)
}
