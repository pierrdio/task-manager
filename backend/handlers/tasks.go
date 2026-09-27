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

func TasksHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			var tasks = []models.Task{}

			rows, err := pool.Query(
				context.Background(),
				"SELECT id, title, description, completed, user_id FROM tasks")
			if err != nil {
				fmt.Println("failed query:", err)
				return
			}

			defer rows.Close()

			for rows.Next() {
				var id int
				var title string
				var description string
				var completed bool
				var user_id int

				err := rows.Scan(&id, &title, &description, &completed, &user_id)

				if err != nil {
					fmt.Println("failed scan:", err)
					return
				}

				newTask := models.Task{
					ID:          id,
					Title:       title,
					Description: description,
					Completed:   completed,
					UserID:      user_id,
				}

				tasks = append(tasks, newTask)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tasks)
			return
		}

		if r.Method == http.MethodPost {

			var newTask models.Task

			var id int
			var title string
			var description string
			var completed bool
			var user_id int

			err := json.NewDecoder(r.Body).Decode(&newTask)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "invalid JSON")
				return
			}

			if strings.TrimSpace(newTask.Title) == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "title is required")
				return
			}

			newTask.Completed = false

			err = pool.QueryRow(
				context.Background(),
				"INSERT into tasks (title, description, completed, user_id) VALUES ($1, $2, $3, $4) RETURNING id, title, description, completed, user_id", newTask.Title, newTask.Description, newTask.Completed, newTask.UserID).Scan(&id, &title, &description, &completed, &user_id)
			if errors.Is(err, pgx.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "not found")
				return
			} else if err != nil {
				fmt.Println("failed query:", err)
				return
			}

			newTask = models.Task{
				ID:          id,
				Title:       title,
				Description: description,
				Completed:   completed,
				UserID:      user_id,
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(&newTask)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "method not allowed")
		return
	}
}

func TaskHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {

			var id int
			var title string
			var description string
			var completed bool
			var user_id int

			idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
			id, err := strconv.Atoi(idStr)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "invalid ID")
				return
			}

			err = pool.QueryRow(
				context.Background(),
				"SELECT id, title, description, completed, user_id FROM tasks WHERE id = $1", id).Scan(&id, &title, &description, &completed, &user_id)
			if errors.Is(err, pgx.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "not found")
				return
			} else if err != nil {
				fmt.Println("failed query:", err)
				return
			}

			task := models.Task{
				ID:          id,
				Title:       title,
				Description: description,
				Completed:   completed,
				UserID:      user_id,
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&task)
			return
		}

		if r.Method == http.MethodPut {
			var updatedTask models.Task

			var id int
			var title string
			var description string
			var completed bool
			var user_id int

			idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
			id, err := strconv.Atoi(idStr)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "invalid ID")
				return
			}

			err = json.NewDecoder(r.Body).Decode(&updatedTask)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "invalid JSON")
				return
			}

			if strings.TrimSpace(updatedTask.Title) == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "title is required")
				return
			}

			updatedTask.ID = id

			err = pool.QueryRow(
				context.Background(),
				"UPDATE tasks SET title = $1, description = $2, completed = $3 WHERE id = $4 RETURNING id, title, description, completed, user_id", updatedTask.Title, updatedTask.Description, updatedTask.Completed, updatedTask.ID).Scan(&id, &title, &description, &completed, &user_id)
			if errors.Is(err, pgx.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "not found")
				return
			} else if err != nil {
				fmt.Println("failed query:", err)
				return
			}

			updatedTask = models.Task{
				ID:          id,
				Title:       title,
				Description: description,
				Completed:   completed,
				UserID:      user_id,
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedTask)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "method not allowed")
		return
	}
}
