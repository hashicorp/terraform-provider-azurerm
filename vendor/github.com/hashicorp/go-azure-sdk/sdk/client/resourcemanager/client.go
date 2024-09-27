// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

var _ client.BaseClient = &Client{}

type Client struct {
	*client.Client

	// apiVersion specifies the version of the API being used, which (by design) will be consistent across a client
	// as we intentionally split out multiple API Versions into different clients, rather than using composite API
	// Versions/packages which can cause confusion about which version is being used.
	apiVersion string
}

func NewClient(api environments.Api, serviceName, apiVersion string) (*Client, error) {
	endpoint, ok := api.Endpoint()
	if !ok {
		return nil, fmt.Errorf("no `endpoint` was returned for this environment")
	}

	baseClient := client.NewClient(*endpoint, serviceName, apiVersion)
	baseClient.AuthorizeRequest = AuthorizeResourceManagerRequest

	return &Client{
		Client:     baseClient,
		apiVersion: apiVersion,
	}, nil
}

// Deprecated: use NewClient instead
func NewResourceManagerClient(api environments.Api, serviceName, apiVersion string) (*Client, error) {
	return NewClient(api, serviceName, apiVersion)
}

func (c *Client) NewRequest(ctx context.Context, input client.RequestOptions) (*client.Request, error) {
	// TODO move these validations to base client method
	if _, ok := ctx.Deadline(); !ok {
		return nil, fmt.Errorf("the context used must have a deadline attached for polling purposes, but got no deadline")
	}
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("pre-validating request payload: %+v", err)
	}
	if input.ContentType == "" {
		return nil, fmt.Errorf("pre-validating request payload: missing `ContentType`")
	}

	req, err := c.Client.NewRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building %s request: %+v", input.HttpMethod, err)
	}

	req.Client = c
	query := url.Values{}

	// there's a handful of cases (e.g. Network, LRO's) where we want to override this on a per-request basis via the options object
	if c.apiVersion != "" {
		query.Set("api-version", c.apiVersion)
	}

	if input.OptionsObject != nil {
		if h := input.OptionsObject.ToHeaders(); h != nil {
			for k, v := range h.Headers() {
				req.Header[k] = v
			}
		}

		if q := input.OptionsObject.ToQuery(); q != nil {
			for k, v := range q.Values() {
				// we intentionally only add one of each type
				query.Del(k)
				query.Add(k, v[0])
			}
		}

		if o := input.OptionsObject.ToOData(); o != nil {
			req.Header = o.AppendHeaders(req.Header)
			query = o.AppendValues(query)
		}
	}

	req.URL.RawQuery = query.Encode()
	req.Pager = input.Pager
	req.RetryFunc = client.RequestRetryAny(append(defaultRetryFunctions, input.RetryFunc)...)
	req.ValidStatusCodes = input.ExpectedStatusCodes

	return req, nil
}
