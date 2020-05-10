package repository

//Statistics represent statistics for HTTP POST requests
type Statistics struct {
	Total   int64   `json:"total"`
	Average float64 `json:"average"`
}
