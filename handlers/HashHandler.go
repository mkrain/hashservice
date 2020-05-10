package handlers

import (
	"encoding/base64"
	"fmt"
	"hashservice/worker"
	"net/http"
	"strconv"
	"strings"
)

//HashHandler test
type hashHandler struct {
	worker *worker.HashWorker
}

//GetHashHandler API Handler for hash requests
func GetHashHandler(w *worker.HashWorker) http.HandlerFunc {
	hash := hashHandler{
		worker: w,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			//The worker was terminated, no new requests can be processed
			if !hash.worker.IsRunning() {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			fmt.Println("Generating Hash")

			password := r.URL.Query().Get("password")

			requestID := hash.worker.NewRequest(password)

			w.WriteHeader(http.StatusCreated)
			w.Header().Add("location", fmt.Sprintf("http://localhost:8080/hash/%d", requestID))
			w.Write([]byte(strconv.Itoa(requestID)))

			return
		}

		if r.Method == "GET" {
			fmt.Println("Getting Hash")

			path := r.URL.Path
			lastSlash := strings.LastIndex(path, "/")

			if lastSlash == -1 {
				fmt.Println("Missing requestID")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			segment := path[lastSlash+1]
			var encoded = ""

			if requestID, error := strconv.Atoi(string(segment)); error == nil {
				hashResponse := hash.worker.GetRequest(requestID)

				//The hash generated is the raw bytes, need to convert it to a base64 encoded string
				encoded = base64.StdEncoding.EncodeToString(hashResponse.Hash)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(encoded))

			return
		}

		fmt.Println("No Route Found")

		w.WriteHeader(http.StatusNotFound)
	}
}
