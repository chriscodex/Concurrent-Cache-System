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
	Lock       sync.RWMutex
}
