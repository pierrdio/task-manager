package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintln(w, "method not allowed")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		var user User

		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid ID")
			return
		}

		for _, u := range users {
			if u.ID == id {
				user = u
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(user)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "not found")
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintln(w, "method not allowed")
	return
}

const PORT = ":8080"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "task manager server")
	})

	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)

	fmt.Println("Server started on port", PORT)

	http.ListenAndServe(PORT, nil)
}
