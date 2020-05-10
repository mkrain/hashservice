package middleware

import (
	"fmt"
	"hashservice/repository"
	"net/http"
	"strings"
	"time"
)

//NewStatisticsMiddleware API Handler for service statistics
func NewStatisticsMiddleware(repo *repository.StatisticsRepository, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		fmt.Println("Start time", startTime)

		handler(w, r)

		endTime := time.Now().Sub(startTime)

		fmt.Println("End time", endTime)

		//This is a hack as the standard router doesn't handle routes per HTTP method
		//For now only keep track of the POST /hash routes
		if strings.Index(r.URL.Path, "/hash") == -1 || r.Method == "GET" {
			return
		}

		repo.IncrementRequestCount(1)
		repo.AddRequestDuration(endTime)
	}
}
