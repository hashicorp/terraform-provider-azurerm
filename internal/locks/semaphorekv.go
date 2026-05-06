package locks

import (
	"context"
	"log"
	"sync"

	"golang.org/x/sync/semaphore"
)

// semaphoreKV is a simple key/value store for weighted semaphores. It can be
// used to limit concurrent operations across arbitrary collaborators that share
// knowledge of the keys they must coordinate on.
type semaphoreKV struct {
	lock  sync.Mutex
	store map[string]*semaphore.Weighted
}

// NewWeighted initializes a weighted semaphore for the given key when one does
// not already exist.
func (s *semaphoreKV) NewWeighted(key string, max int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.store[key]
	if !ok {
		s.store[key] = semaphore.NewWeighted(max)
	}
}

// Acquire acquires the provided weight on the semaphore for the given key.
func (s *semaphoreKV) Acquire(ctx context.Context, key string, weight int64) error {
	log.Printf("[DEBUG] acquiring %q with weight %d", key, weight)
	if err := s.get(key).Acquire(ctx, weight); err != nil {
		return err
	}
	log.Printf("[DEBUG] acquired %q with weight %d", key, weight)
	return nil
}

// Release releases the provided weight on the semaphore for the given key.
func (m *semaphoreKV) Release(key string, weight int64) {
	log.Printf("[DEBUG] releasing %q with weight %d", key, weight)
	m.get(key).Release(weight)
	log.Printf("[DEBUG] released %q with weight %d", key, weight)
}

// get returns the semaphore for the given key, initializing it with a default
// maximum weight of 1 when missing.
func (s *semaphoreKV) get(key string) *semaphore.Weighted {
	s.lock.Lock()
	defer s.lock.Unlock()
	sp, ok := s.store[key]
	if !ok {
		sp = semaphore.NewWeighted(1)
		s.store[key] = sp
	}
	return sp
}

// newSemaphoreKV returns a properly initialized semaphoreKV.
func newSemaphoreKV() *semaphoreKV {
	return &semaphoreKV{
		store: make(map[string]*semaphore.Weighted),
	}
}
