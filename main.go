package main

import (
	"fmt"
	"sync"
	"time"
)

func ExpensiveFunction(n int) int {
	fmt.Printf("Calculating expensive function for %d\n", n)
	time.Sleep(5 * time.Second)
	return n
}

type Service struct {
	InProgress map[int]bool
	IsPending  map[int][]chan int
	Mutex      sync.RWMutex
}

func (s *Service) Work(job int) {
	s.Mutex.RLock()
	exist := s.InProgress[job]
	if exist {
		s.Mutex.Unlock()
		response := make(chan int)
		defer close(response)
		s.Mutex.Lock()
		s.IsPending[job] = append(s.IsPending[job], response)
		s.Mutex.Unlock()
		fmt.Printf("Waiting for Response job: %d\n", job)
		resp := <-response
		fmt.Printf("Response Done, received %d\n", resp)
		return
	}
	s.Mutex.RUnlock()
	s.Mutex.Lock()
	s.InProgress[job] = true
	s.Mutex.Unlock()
	fmt.Printf("Calculating expensive function for %d\n", job)
	result := ExpensiveFunction(job)
	s.Mutex.RLock()
	pendingWorkers, exist := s.IsPending[job]
	s.Mutex.RUnlock()
	if exist {
		for _, pendingWorker := range pendingWorkers {
			pendingWorker <- result
		}
		fmt.Printf("Result send, all pending workers ready for job: %d\n", job)
	}
	s.Mutex.Lock()
	s.InProgress[job] = false
	s.IsPending[job] = make([]chan int, 0)
	s.Mutex.Unlock()
}

// Constructor of Service
func NewService() *Service {
	return &Service{
		InProgress: make(map[int]bool),
		IsPending:  make(map[int][]chan int),
	}
}

func main() {
	service := NewService()
	jobs := []int{3, 4, 5, 5, 4, 8, 8, 8}
	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, element := range jobs {
		go func(job int) {
			defer wg.Done()
			service.Work(job)
		}(element)
	}
	wg.Wait()
}
