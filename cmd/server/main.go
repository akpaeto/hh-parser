package main

import (
	"encoding/json"
	"fmt"
	"hh-parser/internal/hh"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)
	http.HandleFunc("/api/vacancies", vaсanciesHandler)
	fmt.Println("Server started on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ошибка сервера:", err)
	}
}

func vaсanciesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	vacancies, err := hh.SearchVacancies(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vacancies)
}
