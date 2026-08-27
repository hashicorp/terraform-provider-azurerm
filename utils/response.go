// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package utils is deprecated and has been mostly refactored into the helpers package.
// Deprecated: Use the helpers package instead.
package utils

import (
	"net"

	"github.com/Azure/go-autorest/autorest"
)

// ResponseErrorIsRetryable is deprecated and should not be used in new code
// Deprecated
func ResponseErrorIsRetryable(err error) bool {
	if arerr, ok := err.(autorest.DetailedError); ok {
		err = arerr.Original
	}

	// nolint gocritic
	switch e := err.(type) {
	case net.Error:
		if e.Temporary() || e.Timeout() {
			return true
		}
	}

	return false
}
