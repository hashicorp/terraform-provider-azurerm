// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type testManagedEnvironmentReader struct {
	status int
	err    error
}

func (c testManagedEnvironmentReader) Get(context.Context, managedenvironments.ManagedEnvironmentId) (managedenvironments.GetOperationResponse, error) {
	result := managedenvironments.GetOperationResponse{}
	if c.status != 0 {
		result.HttpResponse = &http.Response{StatusCode: c.status}
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
