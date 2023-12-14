// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

var storageDefaultRetryFunctions = []client.RequestRetryFunc{
	// TODO: stuff n tings
}

type BaseClient struct {
	Client *dataplane.Client
}

func NewBaseClient(baseUri string, componentName, apiVersion string) (*BaseClient, error) {
	// NOTE: both the domain name _and_ the domain format can change entirely depending on the type of storage account being used
	// when provisioned in an edge zone, and when AzureDNSZone is used, as such we require the baseUri is provided here
	return &BaseClient{
		Client: dataplane.NewDataPlaneClient(baseUri, fmt.Sprintf("storage/%s", componentName), apiVersion),
	}, nil
}

func (c *BaseClient) NewRequest(ctx context.Context, input client.RequestOptions) (*client.Request, error) {
	// TODO move these validations to base client method
	if _, ok := ctx.Deadline(); !ok {
		return nil, fmt.Errorf("the context used must have a deadline attached for polling purposes, but got no deadline")
	}
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("pre-validating request payload: %+v", err)
	}

	req, err := c.Client.Client.NewRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building %s request: %+v", input.HttpMethod, err)
	}

	req.Client = c
	req.Header.Add("x-ms-version", c.Client.ApiVersion)

	query := url.Values{}
	if input.OptionsObject != nil {
		if h := input.OptionsObject.ToHeaders(); h != nil {
			for k, v := range h.Headers() {
				req.Header[k] = v
			}
		}

		if q := input.OptionsObject.ToQuery(); q != nil {
			query = q.Values()
		}

		if o := input.OptionsObject.ToOData(); o != nil {
			req.Header = o.AppendHeaders(req.Header)
			query = o.AppendValues(query)
		}
	}

	req.URL.RawQuery = query.Encode()
	req.RetryFunc = client.RequestRetryAny(storageDefaultRetryFunctions...)
	req.ValidStatusCodes = input.ExpectedStatusCodes

	return req, nil
}

func (c *BaseClient) Execute(ctx context.Context, req *client.Request) (*client.Response, error) {
	return c.Client.Execute(ctx, req)
}

func (c *BaseClient) ExecutePaged(ctx context.Context, req *client.Request) (*client.Response, error) {
	return c.Client.ExecutePaged(ctx, req)
}

func (c *BaseClient) WithAuthorizer(auth auth.Authorizer) {
	c.Client.Client.Authorizer = auth
}
