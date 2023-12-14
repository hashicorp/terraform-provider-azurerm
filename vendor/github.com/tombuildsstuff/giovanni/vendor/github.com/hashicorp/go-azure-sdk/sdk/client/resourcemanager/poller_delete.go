// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

var _ pollers.PollerType = &deletePoller{}

type deletePoller struct {
	apiVersion           string
	client               *Client
	initialRetryDuration time.Duration
	originalUri          string
	resourcePath         string
}

func deletePollerFromResponse(response *client.Response, client *Client, pollingInterval time.Duration) (*deletePoller, error) {
	// if we've gotten to this point then we're polling against a Resource Manager resource/operation of some kind
	// we next need to determine if the current URI is a Resource Manager resource, or if we should be polling on the
	// resource (e.g. `/my/resource`) rather than an operation on the resource (e.g. `/my/resource/start`)
	if response.Request == nil {
		return nil, fmt.Errorf("request was nil")
	}
	originalUri := response.Request.URL.RequestURI()
	if response.Request.URL == nil {
		return nil, fmt.Errorf("request url was nil")
	}

	// all Resource Manager operations require the `api-version` querystring
	apiVersion := response.Request.URL.Query().Get("api-version")
	if apiVersion == "" {
		return nil, fmt.Errorf("unable to determine `api-version` from %q", originalUri)
	}

	resourcePath, err := resourceManagerResourcePathFromUri(originalUri)
	if err != nil {
		return nil, fmt.Errorf("determining Resource Manager Resource Path from %q: %+v", originalUri, err)
	}

	return &deletePoller{
		apiVersion:           apiVersion,
		client:               client,
		initialRetryDuration: pollingInterval,
		originalUri:          originalUri,
		resourcePath:         *resourcePath,
	}, nil
}

func (p deletePoller) Poll(ctx context.Context) (result *pollers.PollResult, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusNotFound,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: deleteOptions{
			apiVersion: p.apiVersion,
		},
		Path: p.resourcePath,
	}
	req, err := p.client.NewRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}
	resp, err := p.client.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %+v", err)
	}
	if resp == nil {
		return nil, pollers.PollingDroppedConnectionError{}
	}

	result = &pollers.PollResult{
		PollInterval: p.initialRetryDuration,
	}

	if resp.Response != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			{
				result.Status = pollers.PollingStatusSucceeded
				return
			}

		case http.StatusOK:
			{
				result.Status = pollers.PollingStatusInProgress
				return
			}
		}

		err = fmt.Errorf("unexpected status code when polling for resource after deletion, expected a 200/204 but got %d", resp.StatusCode)
	}

	return
}

var _ client.Options = deleteOptions{}

type deleteOptions struct {
	apiVersion string
}

func (p deleteOptions) ToHeaders() *client.Headers {
	return &client.Headers{}
}

func (p deleteOptions) ToOData() *odata.Query {
	return &odata.Query{}
}

func (p deleteOptions) ToQuery() *client.QueryParams {
	q := client.QueryParams{}
	q.Append("api-version", p.apiVersion)
	return &q
}
