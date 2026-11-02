// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

const operationInProgressMessage = `unexpected status 409 (409 Conflict) with response: {"Code":"Conflict","Message":"Cannot modify this site because another operation is in progress. Details: Id: 1, OperationName: ContainerAppUpdateSiteConfig"}`

func TestIsOperationInProgressConflict(t *testing.T) {
	testCases := []struct {
		name     string
		status   int
		err      error
		expected bool
	}{
		{name: "operation in progress", status: http.StatusConflict, err: errors.New(operationInProgressMessage), expected: true},
		{name: "different casing", status: http.StatusConflict, err: errors.New("ANOTHER OPERATION IS IN PROGRESS"), expected: true},
		{name: "a different conflict", status: http.StatusConflict, err: errors.New("Server farm cannot be deleted because it has web apps assigned to it"), expected: false},
		{name: "same text without a conflict status", status: http.StatusBadRequest, err: errors.New(operationInProgressMessage), expected: false},
		{name: "no error", status: http.StatusConflict, err: nil, expected: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: testCase.status}
			if actual := helpers.IsOperationInProgressConflict(resp, testCase.err); actual != testCase.expected {
				t.Fatalf("expected %t but got %t", testCase.expected, actual)
			}
		})
	}

	if helpers.IsOperationInProgressConflict(nil, errors.New(operationInProgressMessage)) {
		t.Fatalf("a nil response must not be treated as a conflict")
	}
}

func TestRetryWhileOperationInProgress(t *testing.T) {
	conflict := func() (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusConflict}, errors.New(operationInProgressMessage)
	}

	t.Run("retries until the operation has finished", func(t *testing.T) {
		attempts := 0
		if err := helpers.RetryWhileOperationInProgress(context.Background(), time.Millisecond, func() (*http.Response, error) {
			attempts++
			if attempts < 4 {
				return conflict()
			}
			return &http.Response{StatusCode: http.StatusOK}, nil
		}); err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if attempts != 4 {
			t.Fatalf("expected 4 attempts but got %d", attempts)
		}
	})

	t.Run("succeeds first time without waiting", func(t *testing.T) {
		attempts := 0
		err := helpers.RetryWhileOperationInProgress(context.Background(), time.Hour, func() (*http.Response, error) {
			attempts++
			return &http.Response{StatusCode: http.StatusOK}, nil
		})
		if err != nil || attempts != 1 {
			t.Fatalf("expected one successful attempt but got %d attempts, err: %+v", attempts, err)
		}
	})

	t.Run("does not retry other errors", func(t *testing.T) {
		attempts := 0
		err := helpers.RetryWhileOperationInProgress(context.Background(), time.Millisecond, func() (*http.Response, error) {
			attempts++
			return &http.Response{StatusCode: http.StatusForbidden}, errors.New("AuthorizationFailed")
		})
		if err == nil || !strings.Contains(err.Error(), "AuthorizationFailed") {
			t.Fatalf("expected the original error but got %+v", err)
		}
		if attempts != 1 {
			t.Fatalf("expected 1 attempt but got %d", attempts)
		}
	})

	t.Run("does not retry a different conflict", func(t *testing.T) {
		attempts := 0
		err := helpers.RetryWhileOperationInProgress(context.Background(), time.Millisecond, func() (*http.Response, error) {
			attempts++
			return &http.Response{StatusCode: http.StatusConflict}, errors.New("Server farm cannot be deleted because it has web apps assigned to it")
		})
		if err == nil || attempts != 1 {
			t.Fatalf("expected 1 attempt and an error but got %d attempts, err: %+v", attempts, err)
		}
	})

	t.Run("stops when the context ends and keeps the original error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := helpers.RetryWhileOperationInProgress(ctx, time.Millisecond, conflict)
		if err == nil {
			t.Fatalf("expected an error")
		}
		if !strings.Contains(err.Error(), "another operation is in progress") {
			t.Fatalf("expected the original conflict to be reported but got %+v", err)
		}
		if !strings.Contains(err.Error(), "gave up waiting") {
			t.Fatalf("expected the error to say the retry was abandoned but got %+v", err)
		}
	})
}
