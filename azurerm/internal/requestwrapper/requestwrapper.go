package requestwrapper

import (
	"context"
	"log"
	"time"
)

const (
	// Default values for the number of retries
	DefaultRetries int = 10

	// Default values for request timeout in seconds
	DefaultRequestTimeout int = 60
)

// Wrapper for Get request that empose a timeout on the reply and will retry the request if the reploy takes too long
func GetWithTimeoutsAndRetries(
	numRetries int,
	requestTimeout int,
	ctx context.Context,
	getMethod func(ctx context.Context) (interface{}, error)) (read interface{}, err error) {

	// Retry the request based on <numRetries> should the request return an error or be not reploy afer <requestTimeout> seconds
	for attempt := 0; attempt < numRetries; attempt++ {

		var timeTillContextDeadline = time.Now().Add(time.Duration(requestTimeout) * time.Second)

		// Create a context based on the duration we will wait for a reply
		getCtx, ctxCancelFunc := context.WithDeadline(ctx, timeTillContextDeadline)

		read, err = getMethod(getCtx)
		defer ctxCancelFunc()

		// Check if the request was sucessful
		if err == nil {
			// Request completed
			break
		}

		log.Printf("[DEBUG] Retrying request due to timeout")
	}

	return read, err
}
