package main

import (
	"fmt"
	"net/http"
)

const PORT = ":8080"

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Task manager api server")
}

func main() {
	http.HandleFunc("/", testHandler)
	fmt.Println("Server started on port", PORT)
	http.ListenAndServe(PORT, nil)
}
