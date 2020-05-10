package handlers

import (
	"hashservice/worker"
	"net/http"
)

type shutdownService struct {
	worker *worker.HashWorker
}

//NewShutdownHandler API Handler for service shutdown
func NewShutdownHandler(worker *worker.HashWorker) http.HandlerFunc {
	service := shutdownService{
		worker: worker,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		service.worker.Stop()
	}
}
