package main

import (
	"fmt"
	"time"
)

func ExpensiveFunction(n int) int {
	fmt.Printf("Calculating expensive function for %d\n", n)
	time.Sleep(5 * time.Second)
	return n
}
