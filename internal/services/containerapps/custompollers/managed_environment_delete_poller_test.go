// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type testManagedEnvironmentReader struct {
	status int
	state  managedenvironments.EnvironmentProvisioningState
	err    error
}

func (c testManagedEnvironmentReader) Get(context.Context, managedenvironments.ManagedEnvironmentId) (managedenvironments.GetOperationResponse, error) {
	result := managedenvironments.GetOperationResponse{}
	if c.status != 0 {
		result.HttpResponse = &http.Response{StatusCode: c.status}
	}
	if c.state != "" {
		result.Model = &managedenvironments.ManagedEnvironment{
			Properties: &managedenvironments.ManagedEnvironmentProperties{ProvisioningState: &c.state},
		}
	}
	return result, c.err
}

func TestManagedEnvironmentDeletePoller(t *testing.T) {
	for _, test := range []struct {
		name    string
		reader  testManagedEnvironmentReader
		status  pollers.PollingStatus
		wantErr bool
	}{
		{"still exists", testManagedEnvironmentReader{status: http.StatusOK}, pollers.PollingStatusInProgress, false},
		{"deleted", testManagedEnvironmentReader{status: http.StatusNotFound, err: errors.New("not found")}, pollers.PollingStatusSucceeded, false},
		{"still deleting", testManagedEnvironmentReader{status: http.StatusOK, state: managedenvironments.EnvironmentProvisioningStateScheduledForDelete}, pollers.PollingStatusInProgress, false},
		{"delete failed", testManagedEnvironmentReader{status: http.StatusOK, state: managedenvironments.EnvironmentProvisioningStateFailed}, "", true},
		{"delete canceled", testManagedEnvironmentReader{status: http.StatusOK, state: managedenvironments.EnvironmentProvisioningStateCanceled}, "", true},
		{"forbidden", testManagedEnvironmentReader{status: http.StatusForbidden, err: errors.New("forbidden")}, "", true},
		{"transport failure", testManagedEnvironmentReader{err: errors.New("connection reset")}, "", true},
		{"missing response", testManagedEnvironmentReader{}, "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			p := managedEnvironmentDeletePoller{client: test.reader}
			result, err := p.Poll(context.Background())
			if (err != nil) != test.wantErr {
				t.Fatalf("want error %t, got %v", test.wantErr, err)
			}
			if err == nil && result.Status != test.status {
				t.Fatalf("want status %s, got %s", test.status, result.Status)
			}
		})
	}
}

func TestIsUndecodableOperationResponse(t *testing.T) {
	for _, test := range []struct {
		name     string
		err      error
		expected bool
	}{
		{"sdk error", errors.New(`internal-error: polling support for the Content-Type "" was not implemented: <nil>`), true},
		{"authorization failure", errors.New("unexpected status 403 (403 Forbidden) with error: AuthorizationFailed"), false},
		{"transport failure", errors.New("connection reset by peer"), false},
		{"context cancelled", context.Canceled, false},
		{"deadline exceeded", context.DeadlineExceeded, false},
		{"polling failed", pollers.PollingFailedError{Message: "delete failed"}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if actual := isUndecodableOperationResponse(test.err); actual != test.expected {
				t.Fatalf("expected %t but got %t", test.expected, actual)
			}
		})
	}
}

// TestWaitForManagedEnvironmentDeletion drives the real SDK client and poller against a local server.
func TestWaitForManagedEnvironmentDeletion(t *testing.T) {
	const operationPath = "/operations/1"

	for _, test := range []struct {
		name string
		// operation handles requests to the operation URL, and environment handles GETs of the environment
		operation   func(writer http.ResponseWriter)
		environment int
		wantErr     string
		wantNoErr   bool
		// wantEnvironmentReads is the number of times the environment must have been read
		wantEnvironmentReads int32
	}{
		{
			name: "operation completes",
			operation: func(writer http.ResponseWriter) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(`{"status":"Succeeded"}`))
			},
			environment:          http.StatusOK,
			wantNoErr:            true,
			wantEnvironmentReads: 0,
		},
		{
			name: "operation failure is returned without checking the environment",
			operation: func(writer http.ResponseWriter) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(`{"id":"x","status":"Failed","error":{"code":"DeleteFailed","message":"the delete failed"}}`))
			},
			environment:          http.StatusNotFound,
			wantErr:              "the delete failed",
			wantEnvironmentReads: 0,
		},
		{
			// the review's reproduction: the operation endpoint is unauthorized, but the environment is gone
			name: "operation authorization failure is not suppressed when the environment is gone",
			operation: func(writer http.ResponseWriter) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusForbidden)
				_, _ = writer.Write([]byte(`{"error":{"code":"AuthorizationFailed","message":"not authorized"}}`))
			},
			environment:          http.StatusNotFound,
			wantErr:              "AuthorizationFailed",
			wantEnvironmentReads: 0,
		},
		{
			name: "undecodable operation response falls back to the environment being deleted",
			operation: func(writer http.ResponseWriter) {
				writer.Header()["Content-Type"] = nil
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(`{}`))
			},
			environment:          http.StatusNotFound,
			wantNoErr:            true,
			wantEnvironmentReads: 1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var environmentReads atomic.Int32
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				switch {
				case request.Method == http.MethodDelete:
					writer.Header().Set("Location", server.URL+operationPath+"?api-version=2025-07-01")
					writer.Header().Set("Retry-After", "0")
					writer.WriteHeader(http.StatusAccepted)
				case strings.HasPrefix(request.URL.Path, operationPath):
					test.operation(writer)
				default:
					environmentReads.Add(1)
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(test.environment)
					_, _ = writer.Write([]byte(`{}`))
				}
			}))
			defer server.Close()

			api := environments.NewApiEndpoint("test", server.URL, nil)
			environmentsClient, err := managedenvironments.NewManagedEnvironmentsClientWithBaseURI(api)
			if err != nil {
				t.Fatalf("building client: %+v", err)
			}
			environmentsClient.Client.Client.AuthorizeRequest = nil

			id := managedenvironments.NewManagedEnvironmentID("00000000-0000-0000-0000-000000000000", "rg", "env")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := environmentsClient.Delete(ctx, id)
			if err != nil {
				t.Fatalf("deleting: %+v", err)
			}

			err = WaitForManagedEnvironmentDeletion(ctx, &result.Poller, environmentsClient, id)
			switch {
			case test.wantNoErr:
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}
			default:
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("expected an error containing %q but got %+v", test.wantErr, err)
				}
			}

			if actual := environmentReads.Load(); actual != test.wantEnvironmentReads {
				t.Fatalf("expected the environment to be read %d times but got %d", test.wantEnvironmentReads, actual)
			}
		})
	}
}
