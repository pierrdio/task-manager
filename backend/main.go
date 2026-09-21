package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
	{ID: 1, Name: "Max", Age: 20},
	{ID: 2, Name: "John", Age: 25},
	{ID: 3, Name: "Ben", Age: 17},
	{ID: 4, Name: "Martin", Age: 23},
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
		return
	}
}

const PORT = ":8080"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "task manager server")
	})

	http.HandleFunc("/users", usersHandler)

	fmt.Println("Server started on port", PORT)

	http.ListenAndServe(PORT, nil)
}
