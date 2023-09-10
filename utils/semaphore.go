package utils

import (
	"sync"
)

type Semaphore struct {
	capacity int
	tokens   chan struct{}
	mu       sync.Mutex
}

func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		capacity: capacity,
		tokens:   make(chan struct{}, capacity),
	}
}

func (s *Semaphore) Acquire() {
	s.tokens <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.tokens
}
