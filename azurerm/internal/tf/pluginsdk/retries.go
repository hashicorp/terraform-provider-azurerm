package pluginsdk

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO: work through and switch these out for WaitForState funcs

// RetryFunc is the function retried until it succeeds.
type RetryFunc = resource.RetryFunc

// RetryError is the required return type of RetryFunc. It forces client code
// to choose whether or not a given error is retryable.
// TODO: deprecate this in the future
type RetryError = resource.RetryError

// Retry is a basic wrapper around StateChangeConf that will just retry
// a function until it no longer returns an error.
func Retry(timeout time.Duration, f RetryFunc) error {
	// TODO: deprecate this
	return resource.Retry(timeout, f) //nolint:SA1019
}

// RetryableError is a helper to create a RetryError that's retryable from a
// given error.
func RetryableError(err error) *RetryError {
	// TODO: deprecate this in the future
	return resource.RetryableError(err)
}

// NonRetryableError is a helper to create a RetryError that's _not_ retryable
// from a given error.
func NonRetryableError(err error) *RetryError {
	// TODO: deprecate this in the future
	return resource.NonRetryableError(err)
}
