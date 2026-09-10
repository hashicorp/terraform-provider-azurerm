// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
)

func TestFunctionAppContainerPollingURL(t *testing.T) {
	testCases := []struct {
		name     string
		location string
		expected string
		wantErr  bool
	}{
		{
			name:     "absolute",
			location: "https://management.azure.com/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Web/locations/westeurope/operationResults/123?api-version=2023-12-01",
			expected: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Web/locations/westeurope/operationResults/123?api-version=2023-12-01",
		},
		{
			name:     "relative",
			location: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Web/locations/westeurope/operationResults/123?api-version=2023-12-01",
			wantErr:  true,
		},
		{
			name:    "missing",
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			response := &http.Response{Header: http.Header{}}
			if testCase.location != "" {
				response.Header.Set("Location", testCase.location)
			}

			actual, err := functionAppContainerPollingURL(response)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if actual.RequestURI() != testCase.expected {
				t.Fatalf("expected %q but got %q", testCase.expected, actual.RequestURI())
			}
		})
	}
}

func TestFunctionAppContainerPoll(t *testing.T) {
	var requests []string
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests = append(requests, request.URL.RequestURI())
		writer.Header().Set("Content-Type", "application/json")
		if len(requests) == 1 {
			writer.Header().Set("Location", server.URL+"/operations/2?api-version=2023-12-01&token=second")
			writer.Header().Set("Retry-After", "0")
			writer.WriteHeader(http.StatusAccepted)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"kind":"functionapp,linux,container,azurecontainerapps"}`))
	}))
	defer server.Close()

	pollingURL, err := url.Parse(server.URL + "/operations/1?api-version=2023-12-01&token=first")
	if err != nil {
		t.Fatalf("parsing test polling URL: %+v", err)
	}

	poller := FunctionAppContainerCreateOrUpdatePoller{
		client: &resourcemanager.Client{
			Client: client.NewClient(server.URL, "webapps", "2023-12-01"),
		},
		pollingURL:   pollingURL,
		pollInterval: 0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	result, err := poller.Poll(ctx)
	if err != nil {
		t.Fatalf("first poll failed: %+v", err)
	}
	if result.Status != pollers.PollingStatusInProgress {
		t.Fatalf("expected first poll status %q but got %q", pollers.PollingStatusInProgress, result.Status)
	}

	result, err = poller.Poll(ctx)
	if err != nil {
		t.Fatalf("second poll failed: %+v", err)
	}
	if result.Status != pollers.PollingStatusSucceeded {
		t.Fatalf("expected second poll status %q but got %q", pollers.PollingStatusSucceeded, result.Status)
	}

	expectedRequests := []string{
		"/operations/1?api-version=2023-12-01&token=first",
		"/operations/2?api-version=2023-12-01&token=second",
	}
	for index, expected := range expectedRequests {
		if requests[index] != expected {
			t.Fatalf("expected request %d to be %q but got %q", index, expected, requests[index])
		}
	}
}

func TestFunctionAppContainerPollingInterval(t *testing.T) {
	response := &http.Response{Header: http.Header{}}
	if actual := functionAppContainerPollingInterval(response); actual != 10*time.Second {
		t.Fatalf("expected default interval of 10s but got %s", actual)
	}

	response.Header.Set("Retry-After", "15")
	if actual := functionAppContainerPollingInterval(response); actual != 15*time.Second {
		t.Fatalf("expected interval of 15s but got %s", actual)
	}
}
