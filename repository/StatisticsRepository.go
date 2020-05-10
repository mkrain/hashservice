package repository

import (
	"fmt"
	"time"
)

type statisticsService struct {
	totalRequests     int64
	totalCallDuration time.Duration
}

//StatisticsRepository represents a store for statistics information
type StatisticsRepository struct {
	service statisticsService
}

//NewStatisticsRepository represents a store for statistics information
func NewStatisticsRepository() *StatisticsRepository {
	repository := StatisticsRepository{
		service: statisticsService{
			totalRequests:     0,
			totalCallDuration: 0,
		},
	}

	return &repository
}

//GetStatistics gets the current statistics for the API
func (s *StatisticsRepository) GetStatistics() Statistics {
	return Statistics{
		Total:   s.service.totalRequests,
		Average: s.getDurationAverageMilliseconds(),
	}
}

func (s *StatisticsRepository) getDurationAverageMilliseconds() float64 {
	return (float64(s.service.totalCallDuration) / float64(time.Millisecond)) /
		float64(s.service.totalRequests)
}

//IncrementRequestCount increments the total request count by the specified amount
func (s *StatisticsRepository) IncrementRequestCount(count int64) {
	fmt.Println(fmt.Sprintf("Incrementing request count by %d", count))
	s.service.totalRequests++
}

//AddRequestDuration adds the specified duration to the cumulative duration
func (s *StatisticsRepository) AddRequestDuration(requestDuration time.Duration) {
	fmt.Println(fmt.Sprintf("Incrementing request duration by %f ms", float64(requestDuration)/float64(time.Millisecond)))
	s.service.totalCallDuration =
		s.service.totalCallDuration + requestDuration
}
