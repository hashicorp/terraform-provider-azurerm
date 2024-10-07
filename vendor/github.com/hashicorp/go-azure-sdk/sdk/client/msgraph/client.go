// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package msgraph

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

var _ client.BaseClient = &Client{}

type ApiVersion string

const (
	VersionOnePointZero ApiVersion = "v1.0"
	VersionBeta         ApiVersion = "beta"
)

type Client struct {
	*client.Client

	// EnableRetries allows reattempting failed requests to work around eventual consistency issues
	// Note that 429 responses are always handled by the base client regardless of this setting
	EnableRetries bool

	// apiVersion specifies the version of the API being used, either "beta" or "v1.0"
	apiVersion ApiVersion

	// tenantId is the tenant ID to use in requests
	tenantId string
}

func NewClient(api environments.Api, serviceName string, apiVersion ApiVersion) (*Client, error) {
	endpoint, ok := api.Endpoint()
	if !ok {
		return nil, fmt.Errorf("no `endpoint` was returned for this environment")
	}
	baseUri := fmt.Sprintf("%s/%s", *endpoint, apiVersion)
	baseClient := client.NewClient(baseUri, fmt.Sprintf("MicrosoftGraph-%s", serviceName), string(apiVersion))
	return &Client{
		Client:        baseClient,
		EnableRetries: true,
		apiVersion:    apiVersion,
	}, nil
}

// Deprecated: use NewClient instead
func NewMsGraphClient(api environments.Api, serviceName string, apiVersion ApiVersion) (*Client, error) {
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
	//req.RetryFunc = client.RequestRetryAny(defaultRetryFunctions...)
	req.ValidStatusCodes = input.ExpectedStatusCodes

	return req, nil
}
