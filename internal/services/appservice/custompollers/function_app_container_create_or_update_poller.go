// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

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

func NewFunctionAppContainerCreateOrUpdatePoller(client *resourcemanager.Client, response *http.Response) (*pollers.Poller, error) {
	pollingURL, err := functionAppContainerPollingURL(response)
	if err != nil {
		return nil, err
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
	if response.Response.StatusCode == http.StatusAccepted {
		if location := response.Response.Header.Get("Location"); location != "" {
			pollingURL, err := functionAppContainerPollingURL(response.Response)
			if err != nil {
				return nil, err
			}
			p.pollingURL = pollingURL
		}

		return &pollers.PollResult{
			HttpResponse: response,
			PollInterval: p.pollInterval,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	return &pollers.PollResult{
		HttpResponse: response,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}

func functionAppContainerPollingURL(response *http.Response) (*url.URL, error) {
	if response == nil {
		return nil, fmt.Errorf("polling response was nil")
	}

	location := response.Header.Get("Location")
	if location == "" {
		return nil, fmt.Errorf("polling response did not contain a Location header")
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
