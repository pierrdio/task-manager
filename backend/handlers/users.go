package handlers

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var users = []models.User{
	{ID: 1, Name: "Max", Age: 20},
	{ID: 2, Name: "John", Age: 25},
	{ID: 3, Name: "Ben", Age: 17},
	{ID: 4, Name: "Martin", Age: 23},
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
		return
	}

	if r.Method == http.MethodPost {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid JSON")
			return
		}

		if strings.TrimSpace(user.Name) == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "name is required")
			return
		}

		if user.Age <= 0 || user.Age > 120 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "incorrect age")
			return
		}

		newUser := models.User{
			ID:   len(users) + 1,
			Name: strings.TrimSpace(user.Name),
			Age:  user.Age,
		}

		users = append(users, newUser)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintln(w, "method not allowed")
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		var user models.User

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

	if r.Method == http.MethodDelete {
		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid ID")
			return
		}

		for i, user := range users {
			if user.ID == id {
				users = append(users[:i], users[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "not found")
		return
	}

	if r.Method == http.MethodPut {
		var updatedUser models.User

		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid ID")
			return
		}

		err = json.NewDecoder(r.Body).Decode(&updatedUser)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "invalid JSON")
			return
		}

		for index, user := range users {
			if user.ID == id {

				if strings.TrimSpace(updatedUser.Name) == "" {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, "name is required")
					return
				}

				if updatedUser.Age <= 0 || updatedUser.Age > 120 {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, "incorrect age")
					return
				}

				users[index].Name = strings.TrimSpace(updatedUser.Name)
				users[index].Age = updatedUser.Age

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(users[index])
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
