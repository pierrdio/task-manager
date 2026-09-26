package handlers

import (
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "method not allowed")
		return
	}
}
