// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// TODO: work through and switch these out for WaitForState funcs

// RetryFunc is the function retried until it succeeds.
type RetryFunc = retry.RetryFunc

// RetryError is the required return type of RetryFunc. It forces client code
// to choose whether or not a given error is retryable.
// TODO: deprecate this in the future
type RetryError = retry.RetryError

// Retry is a basic wrapper around StateChangeConf that will just retry
// a function until it no longer returns an error.
func Retry(timeout time.Duration, f RetryFunc) error {
	// TODO: deprecate this
	// lint:ignore SA1019 SDKv2 migration - staticcheck's own linter directives are currently being ignored under golanci-lint
	return retry.Retry(timeout, f) //nolint:staticcheck
}

// RetryableError is a helper to create a RetryError that's retryable from a
// given error.
func RetryableError(err error) *RetryError {
	// TODO: deprecate this in the future
	return retry.RetryableError(err)
}

// NonRetryableError is a helper to create a RetryError that's _not_ retryable
// from a given error.
func NonRetryableError(err error) *RetryError {
	// TODO: deprecate this in the future
	return retry.NonRetryableError(err)
}
