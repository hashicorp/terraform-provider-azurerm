// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
)

var _ pollers.PollerType = &FunctionAppContainerCreateOrUpdatePoller{}

type FunctionAppContainerCreateOrUpdatePoller struct {
	client       *resourcemanager.Client
	pollingURL   *url.URL
	pollInterval time.Duration
}

// CreateOrUpdateFunctionAppContainer performs CreateOrUpdate for a Function App hosted in a Container Apps Environment,
// runs the optional callback, then waits for the operation to complete. The generated SDK poller can't be used here
// because this API reports completion with a `200` Site body that has no operation state.
func CreateOrUpdateFunctionAppContainer(ctx context.Context, client *webapps.WebAppsClient, id commonids.AppServiceId, input webapps.Site, callback func() error) error {
	result, err := client.CreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CreateOrUpdate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	poller, err := NewFunctionAppContainerCreateOrUpdatePoller(client.Client, result.HttpResponse)
	if err != nil {
		return fmt.Errorf("building Function App container poller: %+v", err)
	}
	if poller == nil {
		return nil
	}
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CreateOrUpdate: %+v", err)
	}

	return nil
}

// NewFunctionAppContainerCreateOrUpdatePoller returns a poller for the initial CreateOrUpdate response, or nil when
// that response already represents a completed operation.
//
// A `200 OK` is only treated as complete when it has no `Location` header and doesn't report an unfinished
// operation state. A failed or cancelled state is returned as an error. A `200` that reports an unfinished state
// without a `Location` header is polled on the resource itself, which is how the SDK handles that response.
func NewFunctionAppContainerCreateOrUpdatePoller(client *resourcemanager.Client, response *http.Response) (*pollers.Poller, error) {
	if response == nil {
		return nil, errors.New("polling response was nil")
	}

	var pollingURL *url.URL
	if response.StatusCode == http.StatusOK {
		status, err := functionAppContainerOperationStatus(response)
		if err != nil {
			return nil, err
		}

		switch {
		case response.Header.Get("Location") != "":
		case status == pollers.PollingStatusInProgress && response.Request != nil && response.Request.URL != nil:
			pollingURL = response.Request.URL
		default:
			return nil, nil
		}
	}

	if pollingURL == nil {
		var err error
		if pollingURL, err = functionAppContainerPollingURL(response); err != nil {
			return nil, err
		}
	}

	pollerType := &FunctionAppContainerCreateOrUpdatePoller{
		client:       client,
		pollingURL:   pollingURL,
		pollInterval: functionAppContainerPollingInterval(response),
	}
	poller := pollers.NewPoller(pollerType, pollerType.pollInterval, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	return &poller, nil
}

func (p *FunctionAppContainerCreateOrUpdatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	request, err := p.client.NewRequest(ctx, client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusAccepted,
		},
		HttpMethod: http.MethodGet,
		Path:       p.pollingURL.Path,
	})
	if err != nil {
		return nil, fmt.Errorf("building polling request: %+v", err)
	}
	request.URL.RawQuery = p.pollingURL.RawQuery

	response, err := request.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("executing polling request: %+v", err)
	}

	p.pollInterval = functionAppContainerPollingInterval(response.Response)

	// A `200` can report a failed or cancelled operation, including when it also points at another URL.
	status := pollers.PollingStatusInProgress
	if response.StatusCode == http.StatusOK {
		var err error
		status, err = functionAppContainerOperationStatus(response.Response)
		if err != nil {
			var failed pollers.PollingFailedError
			var cancelled pollers.PollingCancelledError
			if errors.As(err, &failed) || errors.As(err, &cancelled) || response.Header.Get("Location") == "" {
				return nil, err
			}
			status = pollers.PollingStatusInProgress
		}
	}

	// A `200` that points at a different operation URL is still in progress, matching how the SDK treats
	// the initial response. A `Location` that refers back to the URL just polled is not a continuation.
	if location := response.Header.Get("Location"); location != "" && (response.StatusCode == http.StatusAccepted || response.StatusCode == http.StatusOK) {
		pollingURL, err := functionAppContainerPollingURL(response.Response)
		if err != nil {
			return nil, err
		}
		if response.StatusCode == http.StatusAccepted || pollingURL.String() != p.pollingURL.String() {
			p.pollingURL = pollingURL
			status = pollers.PollingStatusInProgress
		}
	}

	return &pollers.PollResult{
		HttpResponse: response,
		PollInterval: p.pollInterval,
		Status:       status,
	}, nil
}

type functionAppContainerOperationResult struct {
	Status     string `json:"status"`
	Properties struct {
		ProvisioningState string `json:"provisioningState"`
	} `json:"properties"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// functionAppContainerOperationStatus inspects the body of a `200 OK`. A completed operation returns the Site,
// which has neither a `status` nor a `provisioningState`, so that is treated as success. Operation responses
// can also report a `status` or `provisioningState` that is failed, cancelled or still in progress.
func functionAppContainerOperationStatus(response *http.Response) (pollers.PollingStatus, error) {
	if response.Body == nil {
		return pollers.PollingStatusSucceeded, nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("reading polling response body: %+v", err)
	}
	response.Body.Close()
	response.Body = io.NopCloser(bytes.NewReader(body))

	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return pollers.PollingStatusSucceeded, nil
	}

	var result functionAppContainerOperationResult
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshalling polling response body: %+v", err)
	}

	message := result.Error.Message
	if message == "" {
		message = string(body)
	} else if result.Error.Code != "" {
		message = fmt.Sprintf("%s: %s", result.Error.Code, message)
	}

	status := pollers.PollingStatusSucceeded
	for _, state := range []string{result.Status, result.Properties.ProvisioningState} {
		switch {
		case state == "":
		case strings.EqualFold(state, "Failed"):
			return "", pollers.PollingFailedError{HttpResponse: &client.Response{Response: response}, Message: message}
		case strings.EqualFold(state, "Canceled"), strings.EqualFold(state, "Cancelled"):
			return "", pollers.PollingCancelledError{HttpResponse: &client.Response{Response: response}, Message: message}
		case !strings.EqualFold(state, "Succeeded"):
			status = pollers.PollingStatusInProgress
		}
	}

	return status, nil
}

func functionAppContainerPollingURL(response *http.Response) (*url.URL, error) {
	if response == nil {
		return nil, errors.New("polling response was nil")
	}

	location := response.Header.Get("Location")
	if location == "" {
		return nil, errors.New("polling response did not contain a Location header")
	}

	parsed, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("parsing polling URL %q: %+v", location, err)
	}
	if !parsed.IsAbs() {
		return nil, fmt.Errorf("polling URL %q was not absolute", location)
	}
	return parsed, nil
}

func functionAppContainerPollingInterval(response *http.Response) time.Duration {
	const defaultPollingInterval = 10 * time.Second

	if response == nil {
		return defaultPollingInterval
	}

	retryAfter, err := strconv.ParseInt(response.Header.Get("Retry-After"), 10, 64)
	if err != nil {
		return defaultPollingInterval
	}
	return time.Duration(retryAfter) * time.Second
}
