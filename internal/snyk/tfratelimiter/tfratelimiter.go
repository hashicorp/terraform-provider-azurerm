package tfratelimiter

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// A tenant, a subscription, account, whatever the specific cloud API uses for
// rate limits.
type Scope = string

// A prefix, specified as the prefix of the terraform names of the corresponding
// resources
type Service = string

// A kind of operation performed against the cloud API.
type Operation int

const (
	// A raw single list query against an API.
	List = iota

	// A raw single read query against an API.
	Read
)

func (op Operation) String() string {
	switch op {
	case List:
		return "list"
	case Read:
		return "read"
	}
	return "unknown"
}

type key struct {
	scope     Scope
	service   Service
	operation Operation
}

type limiter struct {
	// Actual rate limiter.
	limiter *rate.Limiter

	// When can we throw this away?
	expires time.Time
}

// Global limiter
var mutex sync.Mutex
var limiters map[key]limiter = map[key]limiter{}

const expiry = 2 * time.Hour
const epsilon = 10 * time.Millisecond

// Lookup a limiter, creating it if necessary.  If this returns nil, there is
// no limit currently in place.
func getLimiter(
	scope Scope,
	service Service,
	operation Operation,
) *rate.Limiter {
	key := key{scope, service, operation}
	if limiter, ok := limiters[key]; ok {
		limiter.expires = time.Now().Add(expiry)
		return limiter.limiter
	}

	ratePerSecond, burst := getLimit(service, operation)
	if ratePerSecond < 0 {
		return nil
	}

	mutex.Lock()
	defer mutex.Unlock()
	limiter := limiter{
		expires: time.Now().Add(expiry),
		limiter: rate.NewLimiter(rate.Limit(ratePerSecond), burst),
	}
	limiters[key] = limiter
	go reaper(key)
	return limiter.limiter
}

func reaper(key key) {
	time.Sleep(expiry + epsilon) // Start with initial sleep
	done := false
	for !done {
		if limiter, ok := limiters[key]; ok {
			if time.Now().After(limiter.expires) {
				mutex.Lock()
				delete(limiters, key)
				mutex.Unlock()
				done = true
			} else {
				time.Sleep(time.Until(limiter.expires) + epsilon)
			}
		} else {
			done = true // Has already been deleted
		}
	}
}

// Main entry point for raw calls.  Blocks until the requested call can be made.
func WaitForService(
	scope Scope,
	service Service,
	op Operation,
	numCallsToReserve int,
) error {
	limiter := getLimiter(scope, service, op)
	if limiter == nil {
		return nil
	}

	reservation := limiter.ReserveN(time.Now(), numCallsToReserve)
	if !reservation.OK() {
		return fmt.Errorf(
			"tfratelimiter: reservation for %d %s operations on %s failed: burst too low",
			numCallsToReserve,
			op.String(),
			service,
		)
	}

	time.Sleep(reservation.Delay())
	return nil
}

func WaitForResourceList(
	scope Scope,
	resourceType string,
) error {
	service := serviceForResourceType(resourceType)
	if service == Service_unknown {
		return nil
	}

	return WaitForService(scope, service, List, 1)
}

func WaitForResourceRefresh(
	scope Scope,
	resourceType string,
) error {
	service := serviceForResourceType(resourceType)
	n := readsForResourceType(resourceType)
	if service == Service_unknown {
		return nil
	}

	return WaitForService(scope, service, Read, n)
}
