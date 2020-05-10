package worker

import (
	"fmt"
	"hashservice/repository"
	"time"
)

//HashWorker t
type HashWorker struct {
	repo                  *repository.HashRepository
	fullfilledRequests    map[int]HashResponse
	pendingRequests       map[string]HashRequest
	requestCounter        int
	pendingRequestCount   int
	requestQueue          chan HashRequest
	quit                  chan bool
	shouldProcessRequests bool
	waitDurationSeconds   int64
}

//HashRequest represents a request, pending or fullfilled
type HashRequest struct {
	Password  string
	RequestID int
}

//HashResponse contains the hash value and requestID from the original HashRequest
type HashResponse struct {
	Hash      []byte
	RequestID int
}

//NewHashWorker returns a HashWorker
func NewHashWorker(repo *repository.HashRepository, waitSeconds int64) *HashWorker {
	hashWorker := HashWorker{
		repo:                  repo,
		fullfilledRequests:    make(map[int]HashResponse),
		pendingRequests:       make(map[string]HashRequest),
		requestCounter:        0,
		requestQueue:          make(chan HashRequest),
		quit:                  make(chan bool),
		shouldProcessRequests: true,
		waitDurationSeconds:   waitSeconds,
	}

	hashWorker.start()

	return &hashWorker
}

func (worker *HashWorker) queueNextRequest(nextRequest HashRequest) {
	fmt.Println(fmt.Sprintf("Queueing request %d", nextRequest.RequestID))

	worker.pendingRequests[nextRequest.Password] = nextRequest

	timeout := time.Duration(int64(time.Second) * worker.waitDurationSeconds)

	time.AfterFunc(timeout, func() {
		worker.processRequest(nextRequest)
	})
}

func (worker *HashWorker) processRequest(nextRequest HashRequest) {
	fmt.Printf("Processing request %d -> %s\n",
		nextRequest.RequestID,
		time.Now().Format(time.RFC850))

	var hashResponse = HashResponse{
		RequestID: nextRequest.RequestID,
	}

	//Existing key has already been processed
	if worker.repo.ContainsKey(nextRequest.Password) {
		fmt.Println("There was an existing hash, reusing")
		hashResponse.Hash = worker.repo.GetHash(nextRequest.Password)
	} else {
		//Create new entry
		fmt.Println("Creating an new hash entry")
		hashResponse.Hash = worker.repo.CreateHash(nextRequest.Password)
	}

	worker.fullfilledRequests[nextRequest.RequestID] = hashResponse
	worker.pendingRequestCount--
	delete(worker.pendingRequests, nextRequest.Password)
}

func (worker *HashWorker) start() {
	go func() {
		for {
			select {
			case nextRequest := <-worker.requestQueue:
				worker.queueNextRequest(nextRequest)
			case <-worker.quit:
				fmt.Println("Worker shutting down")
				return
			}
		}
	}()
}

//IsRunning returns true if the worker is not shutdown and can process requests, false otherwise
func (worker *HashWorker) IsRunning() bool {
	return worker.shouldProcessRequests
}

//Stop stops the worker after pending requests have been fullfilled
func (worker *HashWorker) Stop() {
	fmt.Println("Stop initiated")

	//Periodically check for pending requests
	timeout := time.Duration(int64(time.Second) * 1)

	time.AfterFunc(timeout, func() {
		//Wait for requests to finish
		if worker.pendingRequestCount == 0 {
			//Signal worker to completely stop
			worker.quit <- true
		} else {
			fmt.Println(
				fmt.Sprintf("There are %d pending requests, waiting",
					worker.pendingRequestCount))
			//Recheck
			worker.Stop()
		}

		worker.shouldProcessRequests = false
	})
}

//GetRequest returns a fullfilled request or an empty one
func (worker *HashWorker) GetRequest(requestID int) HashResponse {
	if request, ok := worker.fullfilledRequests[requestID]; ok {
		return request
	}

	return HashResponse{RequestID: requestID}
}

//NewRequest submits a new request for fullfillment and returns the requestID
func (worker *HashWorker) NewRequest(value string) int {
	request := HashRequest{
		Password:  value,
		RequestID: worker.requestCounter + 1,
	}

	worker.requestCounter++
	worker.pendingRequestCount++

	worker.requestQueue <- request

	return request.RequestID
}
