package main

import (
	"math/rand"
	"sync"
)

// generateNumbers -генерирует числа в горутине
func generateNumbers(limit int, goNum int, numbers chan int) {
	numbersSet := Set[int]{data: make(map[int]bool)}
	var lock sync.Mutex
	for i := 0; i < goNum; i++ {
		go func() {
			for {
				randomNumber := rand.Intn(limit + 1)
				lock.Lock()
				numbersLen := numbersSet.Len()
				if numbersLen >= limit {
					close(numbers)
					break
				}
				numberIsAdded := numbersSet.Add(randomNumber)
				if numberIsAdded {
					numbers <- randomNumber
				}
				lock.Unlock()
			}
		}()
	}
}