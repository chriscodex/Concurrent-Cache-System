package main

import (
	"fmt"
	"sync"
	"time"
)

func ExpensiveFunction(n int) int {
	//fmt.Printf("Calculating expensive function for %d\n", n)
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
