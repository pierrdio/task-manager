package handlers

import (
	"backend/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UsersHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {

			var users = []models.User{}

			rows, err := pool.Query(
				context.Background(),
				"SELECT id, name, age FROM users")
			if err != nil {
				fmt.Println("failed query:", err)
				return
			}

			defer rows.Close()

			for rows.Next() {
				var id int
				var name string
				var age int

				err := rows.Scan(&id, &name, &age)

				if err != nil {
					fmt.Println("failed scan:", err)
					return
				}

				user := models.User{
					ID:   id,
					Name: name,
					Age:  age,
				}

				users = append(users, user)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(users)
			return
		}

		if r.Method == http.MethodPost {

			var id int
			var name string
			var age int
			var newUser models.User

			err := json.NewDecoder(r.Body).Decode(&newUser)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "invalid JSON")
				return
			}

			if strings.TrimSpace(newUser.Name) == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "name is required")
				return
			}

			if newUser.Age <= 0 || newUser.Age > 120 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "incorrect age")
				return
			}

			err = pool.QueryRow(
				context.Background(),
				"INSERT into users (name, age) VALUES ($1, $2) RETURNING id, name, age", newUser.Name, newUser.Age).Scan(&id, &name, &age)
			if err != nil {
				fmt.Println("failed query:", err)
				return
			}

			newUser = models.User{
				ID:   id,
				Name: strings.TrimSpace(newUser.Name),
				Age:  newUser.Age,
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newUser)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
	}
}

func UserHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {

			var name string
			var age int

			idStr := strings.TrimPrefix(r.URL.Path, "/users/")
			id, err := strconv.Atoi(idStr)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "invalid ID")
				return
			}

			err = pool.QueryRow(
				context.Background(),
				"SELECT id, name, age FROM users WHERE id = $1", id).Scan(&id, &name, &age)
			if errors.Is(err, pgx.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "not found")
				return
			} else if err != nil {
				fmt.Println("failed query:", err)
				return
			}

			user := models.User{
				ID:   id,
				Name: name,
				Age:  age,
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
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

			var users = []models.User{}

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

			var users = []models.User{}

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
}
