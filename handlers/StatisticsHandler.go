package handlers

import (
	"encoding/json"
	"hashservice/repository"
	"net/http"
)

//NewStatisticsHandler API Handler for API statistics
func NewStatisticsHandler(repo *repository.StatisticsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		stats := repo.GetStatistics()
		json.NewEncoder(w).Encode(stats)
	}
}
