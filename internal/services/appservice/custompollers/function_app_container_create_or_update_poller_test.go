// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
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

func TestNewFunctionAppContainerCreateOrUpdatePoller(t *testing.T) {
	const location = "https://management.azure.com/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Web/locations/westeurope/operationResults/1?api-version=2023-12-01"

	testCases := []struct {
		name       string
		statusCode int
		location   string
		body       string
		wantPoller bool
		wantFailed bool
		wantCancel bool
	}{
		{name: "ok site body", statusCode: http.StatusOK, body: `{"kind":"functionapp,linux,container,azurecontainerapps"}`},
		{name: "ok empty body", statusCode: http.StatusOK},
		{name: "ok succeeded", statusCode: http.StatusOK, body: `{"properties":{"provisioningState":"Succeeded"}}`},
		{name: "ok failed", statusCode: http.StatusOK, body: `{"properties":{"provisioningState":"Failed"}}`, wantFailed: true},
		{name: "ok canceled", statusCode: http.StatusOK, body: `{"properties":{"provisioningState":"Canceled"}}`, wantCancel: true},
		{name: "ok updating without location polls the resource", statusCode: http.StatusOK, body: `{"properties":{"provisioningState":"Updating"}}`, wantPoller: true},
		{name: "ok in progress without location polls the resource", statusCode: http.StatusOK, body: `{"status":"InProgress"}`, wantPoller: true},
		{name: "ok with location", statusCode: http.StatusOK, location: location, body: `{"kind":"x"}`, wantPoller: true},
		{name: "ok failed with location", statusCode: http.StatusOK, location: location, body: `{"status":"Failed"}`, wantFailed: true},
		{name: "accepted with location", statusCode: http.StatusAccepted, location: location, wantPoller: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			requestURL, err := url.Parse("https://management.azure.com/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Web/sites/site?api-version=2023-12-01")
			if err != nil {
				t.Fatalf("parsing request URL: %+v", err)
			}

			response := &http.Response{
				StatusCode: testCase.statusCode,
				Header:     http.Header{},
				Body:       io.NopCloser(strings.NewReader(testCase.body)),
				Request:    &http.Request{URL: requestURL},
			}
			if testCase.location != "" {
				response.Header.Set("Location", testCase.location)
			}

			poller, err := NewFunctionAppContainerCreateOrUpdatePoller(nil, response)

			switch {
			case testCase.wantFailed:
				if _, ok := errors.AsType[pollers.PollingFailedError](err); !ok {
					t.Fatalf("expected a PollingFailedError but got %+v", err)
				}
			case testCase.wantCancel:
				if _, ok := errors.AsType[pollers.PollingCancelledError](err); !ok {
					t.Fatalf("expected a PollingCancelledError but got %+v", err)
				}
			default:
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}
				if (poller != nil) != testCase.wantPoller {
					t.Fatalf("expected a poller: %t but got a poller: %t", testCase.wantPoller, poller != nil)
				}
			}
		})
	}
}

func TestFunctionAppContainerPollTerminalResponse(t *testing.T) {
	testCases := []struct {
		name      string
		body      string
		expected  pollers.PollingStatus
		failed    bool
		cancelled bool
		message   string
	}{
		{name: "site body", body: `{"kind":"functionapp,linux,container,azurecontainerapps"}`, expected: pollers.PollingStatusSucceeded},
		{name: "empty body", expected: pollers.PollingStatusSucceeded},
		{name: "succeeded status", body: `{"status":"Succeeded"}`, expected: pollers.PollingStatusSucceeded},
		{name: "succeeded provisioning state", body: `{"properties":{"provisioningState":"Succeeded"}}`, expected: pollers.PollingStatusSucceeded},
		{name: "in progress status", body: `{"status":"InProgress"}`, expected: pollers.PollingStatusInProgress},
		{name: "updating provisioning state", body: `{"properties":{"provisioningState":"Updating"}}`, expected: pollers.PollingStatusInProgress},
		{name: "runtime state is not an operation state", body: `{"properties":{"state":"Running"}}`, expected: pollers.PollingStatusSucceeded},
		{name: "failed status", body: `{"status":"Failed","error":{"code":"BadRequest","message":"something went wrong"}}`, failed: true, message: "BadRequest: something went wrong"},
		{name: "failed provisioning state", body: `{"properties":{"provisioningState":"Failed"}}`, failed: true},
		{name: "canceled status", body: `{"status":"Canceled"}`, cancelled: true},
		{name: "cancelled status", body: `{"status":"Cancelled"}`, cancelled: true},
		{name: "canceled provisioning state", body: `{"properties":{"provisioningState":"Canceled"}}`, cancelled: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(testCase.body))
			}))
			defer server.Close()

			pollingURL, err := url.Parse(server.URL + "/operations/1?api-version=2023-12-01")
			if err != nil {
				t.Fatalf("parsing test polling URL: %+v", err)
			}

			poller := FunctionAppContainerCreateOrUpdatePoller{
				client: &resourcemanager.Client{
					Client: client.NewClient(server.URL, "webapps", "2023-12-01"),
				},
				pollingURL: pollingURL,
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			result, err := poller.Poll(ctx)

			switch {
			case testCase.failed:
				failedErr, ok := errors.AsType[pollers.PollingFailedError](err)
				if !ok {
					t.Fatalf("expected a PollingFailedError but got %+v", err)
				}
				if testCase.message != "" && failedErr.Message != testCase.message {
					t.Fatalf("expected message %q but got %q", testCase.message, failedErr.Message)
				}
			case testCase.cancelled:
				if _, ok := errors.AsType[pollers.PollingCancelledError](err); !ok {
					t.Fatalf("expected a PollingCancelledError but got %+v", err)
				}
			default:
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}
				if result.Status != testCase.expected {
					t.Fatalf("expected status %q but got %q", testCase.expected, result.Status)
				}
			}
		})
	}
}

func TestFunctionAppContainerPollOKWithLocation(t *testing.T) {
	var requests []string
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests = append(requests, request.URL.RequestURI())
		writer.Header().Set("Content-Type", "application/json")
		switch len(requests) {
		case 1:
			// a `200` pointing at a different operation URL is a continuation
			writer.Header().Set("Location", server.URL+"/operations/2?api-version=2023-12-01&token=second")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"kind":"functionapp,linux,container,azurecontainerapps"}`))
		case 2:
			// a `Location` referring back to the URL just polled is not
			writer.Header().Set("Location", server.URL+"/operations/2?api-version=2023-12-01&token=second")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"kind":"functionapp,linux,container,azurecontainerapps"}`))
		}
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
		pollingURL: pollingURL,
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

func TestFunctionAppContainerPollContinuationWithTerminalBody(t *testing.T) {
	testCases := []struct {
		name       string
		body       string
		wantFailed bool
		wantCancel bool
		expected   pollers.PollingStatus
	}{
		{name: "failed", body: `{"status":"Failed","error":{"code":"DeploymentFailed","message":"deployment failed"}}`, wantFailed: true},
		{name: "canceled", body: `{"properties":{"provisioningState":"Canceled"}}`, wantCancel: true},
		{name: "site body continues", body: `{"kind":"functionapp,linux,container,azurecontainerapps"}`, expected: pollers.PollingStatusInProgress},
		{name: "in progress continues", body: `{"status":"InProgress"}`, expected: pollers.PollingStatusInProgress},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.Header().Set("Location", server.URL+"/operations/2?api-version=2023-12-01&token=second")
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(testCase.body))
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
				pollingURL: pollingURL,
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			result, err := poller.Poll(ctx)

			switch {
			case testCase.wantFailed:
				failedErr, ok := errors.AsType[pollers.PollingFailedError](err)
				if !ok {
					t.Fatalf("expected a PollingFailedError but got %+v", err)
				}
				if failedErr.Message != "DeploymentFailed: deployment failed" {
					t.Fatalf("expected the operation error to be preserved but got %q", failedErr.Message)
				}
			case testCase.wantCancel:
				if _, ok := errors.AsType[pollers.PollingCancelledError](err); !ok {
					t.Fatalf("expected a PollingCancelledError but got %+v", err)
				}
			default:
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}
				if result.Status != testCase.expected {
					t.Fatalf("expected status %q but got %q", testCase.expected, result.Status)
				}
				if poller.pollingURL.Query().Get("token") != "second" {
					t.Fatalf("expected the continuation location to be followed but polling %q", poller.pollingURL)
				}
			}
		})
	}
}

func TestCreateOrUpdateFunctionAppContainer(t *testing.T) {
	siteBody := `{"kind":"functionapp,linux,container,azurecontainerapps"}`

	testCases := []struct {
		name string
		// handler receives the number of requests seen, with the PUT as request 1
		handler       func(server *httptest.Server, writer http.ResponseWriter, count int32)
		wantFailed    bool
		wantCancelled bool
		wantMessage   string
		wantRequests  int32
	}{
		{
			name: "ok site body returns without polling",
			handler: func(_ *httptest.Server, writer http.ResponseWriter, _ int32) {
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(siteBody))
			},
			wantRequests: 1,
		},
		{
			name: "ok failed provisioning state is an error",
			handler: func(_ *httptest.Server, writer http.ResponseWriter, _ int32) {
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(`{"properties":{"provisioningState":"Failed"}}`))
			},
			wantFailed:   true,
			wantRequests: 1,
		},
		{
			name: "ok canceled provisioning state is an error",
			handler: func(_ *httptest.Server, writer http.ResponseWriter, _ int32) {
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(`{"properties":{"provisioningState":"Canceled"}}`))
			},
			wantCancelled: true,
			wantRequests:  1,
		},
		{
			name: "ok updating provisioning state polls the resource until it completes",
			handler: func(_ *httptest.Server, writer http.ResponseWriter, count int32) {
				writer.Header().Set("Retry-After", "0")
				writer.WriteHeader(http.StatusOK)
				if count < 3 {
					_, _ = writer.Write([]byte(`{"properties":{"provisioningState":"Updating"}}`))
					return
				}
				_, _ = writer.Write([]byte(`{"properties":{"provisioningState":"Succeeded"}}`))
			},
			wantRequests: 3,
		},
		{
			name: "accepted follows the location to a completed site",
			handler: func(server *httptest.Server, writer http.ResponseWriter, count int32) {
				if count == 1 {
					writer.Header().Set("Location", server.URL+"/operations/1?api-version=2023-12-01")
					writer.Header().Set("Retry-After", "0")
					writer.WriteHeader(http.StatusAccepted)
					return
				}
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(siteBody))
			},
			wantRequests: 2,
		},
		{
			name: "accepted then a failed operation preserves the error",
			handler: func(server *httptest.Server, writer http.ResponseWriter, count int32) {
				if count == 1 {
					writer.Header().Set("Location", server.URL+"/operations/1?api-version=2023-12-01")
					writer.Header().Set("Retry-After", "0")
					writer.WriteHeader(http.StatusAccepted)
					return
				}
				writer.Header().Set("Location", server.URL+"/operations/2?api-version=2023-12-01")
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(`{"status":"Failed","error":{"code":"DeploymentFailed","message":"deployment failed"}}`))
			},
			wantFailed:   true,
			wantMessage:  "DeploymentFailed: deployment failed",
			wantRequests: 2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var requests atomic.Int32
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				testCase.handler(server, writer, requests.Add(1))
			}))
			defer server.Close()

			webAppsClient := newTestWebAppsClient(t, server.URL)
			id := commonids.NewAppServiceID("00000000-0000-0000-0000-000000000000", "rg", "site")

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			err := CreateOrUpdateFunctionAppContainer(ctx, webAppsClient, id, webapps.Site{}, nil)

			// the errors are wrapped with `%+v`, so the type is asserted via the message
			switch {
			case testCase.wantFailed:
				if err == nil || !strings.Contains(err.Error(), "polling failed") {
					t.Fatalf("expected a polling failed error but got %+v", err)
				}
				if testCase.wantMessage != "" && !strings.Contains(err.Error(), testCase.wantMessage) {
					t.Fatalf("expected the error to contain %q but got %+v", testCase.wantMessage, err)
				}
			case testCase.wantCancelled:
				if err == nil || !strings.Contains(err.Error(), "polling was cancelled") {
					t.Fatalf("expected a polling cancelled error but got %+v", err)
				}
			default:
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}
			}

			if actual := requests.Load(); actual != testCase.wantRequests {
				t.Fatalf("expected %d requests but got %d", testCase.wantRequests, actual)
			}
		})
	}
}

// newTestWebAppsClient builds a client the same way the provider does, so the PUT carries an `api-version`, but
// pointed at a local server and without authorization.
func newTestWebAppsClient(t *testing.T, baseURI string) *webapps.WebAppsClient {
	t.Helper()

	rmClient, err := resourcemanager.NewClient(environments.NewApiEndpoint("test", baseURI, nil), "webapps", "2023-12-01")
	if err != nil {
		t.Fatalf("building client: %+v", err)
	}
	rmClient.AuthorizeRequest = nil

	return &webapps.WebAppsClient{Client: rmClient}
}
