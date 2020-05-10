package main

import (
	"fmt"
	"hashservice/handlers"
	hashMiddleware "hashservice/middleware"
	"hashservice/repository"
	"hashservice/worker"
	"net/http"
)

var port = 8080
var delayPerRequest int64 = 5

func main() {
	repo := repository.NewHashRepository()
	statsRepo := repository.NewStatisticsRepository()
	worker := worker.NewHashWorker(repo, delayPerRequest)
	hashHandler := hashMiddleware.NewStatisticsMiddleware(statsRepo, handlers.GetHashHandler(worker))
	shutdownHandlers := handlers.NewShutdownHandler(worker)
	statisticsHandler := handlers.NewStatisticsHandler(statsRepo)

	http.HandleFunc("/hash/", hashHandler)
	http.HandleFunc("/hash", hashHandler)
	http.HandleFunc("/shutdown/", shutdownHandlers)
	http.HandleFunc("/shutdown", shutdownHandlers)
	http.HandleFunc("/stats", statisticsHandler)
	http.HandleFunc("/stats/", statisticsHandler)

	fmt.Println("Starting Server on Port ", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
