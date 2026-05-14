// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package locks

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSemaphoreKV(t *testing.T) {
	testCases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "new weighted then acquire and release",
			run: func(t *testing.T) {
				store := newSemaphoreKV()
				const key = "basic"

				store.NewWeighted(key, 1)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				if err := store.Acquire(ctx, key, 1); err != nil {
					t.Fatalf("acquiring semaphore: %+v", err)
				}

				store.Release(key, 1)
			},
		},
		{
			name: "acquire without NewWeighted panics",
			run: func(t *testing.T) {
				defer func() {
					if r := recover(); r == nil {
						t.Fatal("expected panic for uninitialized key")
					}
				}()
				store := newSemaphoreKV()
				ctx := context.Background()
				_ = store.Acquire(ctx, "missing", 1)
			},
		},
		{
			name: "concurrent acquire respects capacity",
			run: func(t *testing.T) {
				const (
					key      = "concurrent"
					capacity = 4
				)

				store := newSemaphoreKV()
				store.NewWeighted(key, capacity)

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				var inFlight int
				var maxInFlight int
				var mu sync.Mutex
				wg := sync.WaitGroup{}

				for i := 0; i < 20; i++ {
					wg.Add(1)
					go func(worker int) {
						defer wg.Done()

						if err := store.Acquire(ctx, key, 1); err != nil {
							t.Errorf("worker %d acquiring semaphore: %+v", worker, err)
							return
						}

						mu.Lock()
						inFlight++
						if inFlight > maxInFlight {
							maxInFlight = inFlight
						}
						mu.Unlock()

						time.Sleep(10 * time.Millisecond)

						mu.Lock()
						inFlight--
						mu.Unlock()

						store.Release(key, 1)
					}(i)
				}

				wg.Wait()

				if maxInFlight > capacity {
					t.Fatalf("expected max in-flight <= %d, got %d", capacity, maxInFlight)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.run)
	}
}
