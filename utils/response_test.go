// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"fmt"
	"testing"

	"github.com/Azure/go-autorest/autorest"
)

type testNetError struct {
	timeout   bool
	temporary bool
}

// testNetError fulfils net.Error interface
func (e testNetError) Error() string   { return "testError" }
func (e testNetError) Timeout() bool   { return e.timeout }
func (e testNetError) Temporary() bool { return e.temporary }

func TestResponseErrorIsRetryable(t *testing.T) {
	testCases := []struct {
		desc           string
		err            error
		expectedResult bool
	}{
		{"Unhandled error types are not retryable", fmt.Errorf("Some other error"), false},
		{"Temporary AND timeout errors are retryable", testNetError{true, true}, true},
		{"Timeout errors are retryable", testNetError{true, false}, true},
		{"Temporary errors are retryable", testNetError{false, true}, true},
		{"net.Errors that are neither temporary nor timeouts are not retryable", testNetError{false, false}, false},
		{"Retryable error nested in autorest.DetailedError is retryable", autorest.DetailedError{
			Original: testNetError{true, true},
		}, true},
		{"Unhandled error nested in autorest.DetailedError is not retryable", autorest.DetailedError{
			Original: fmt.Errorf("Some other error"),
		}, false},
		{"nil is handled as non-retryable", nil, false},
	}

	for _, test := range testCases {
		result := ResponseErrorIsRetryable(test.err)
		if test.expectedResult != result {
			t.Errorf("Expected '%v' for case '%s' - got '%v'",
				test.expectedResult, test.desc, result)
		}
	}
}
