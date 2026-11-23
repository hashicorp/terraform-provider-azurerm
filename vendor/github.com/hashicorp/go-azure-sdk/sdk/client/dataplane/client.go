// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataplane

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-uuid"
)

var _ client.BaseClient = &Client{}

type Client struct {
	*client.Client

	// ApiVersion specifies the version of the API being used, which (by design) will be consistent across a client
	// as we intentionally split out multiple API Versions into different clients, rather than using composite API
	// Versions/packages which can cause confusion about which version is being used.
	ApiVersion string
}

var freshTransport = func() *http.Transport {
	opts := client.DefaultTransportOptions()
	opts.DisableKeepAlives = true
	return client.NewTransport(opts)
}()

// NewClient returns a pointer to a Client configured with AuthorizeDataPlaneRequest and a fresh-connection
// transport (DisableKeepAlives: true) so per-endpoint DNS cutovers (such as Private Endpoints) do not reuse
// pre-cutover sockets. If connection pooling is desired for a specific client, call c.SetTransport(client.NewDefaultTransport()).
// If the endpoint value is available at call-time, it will be set into the BaseUri property of the underlying BaseClient.
func NewClient(endpoint string, serviceName, apiVersion string) (*Client, error) {
	c := NewDataPlaneClient(endpoint, serviceName, apiVersion)
	c.AuthorizeRequest = AuthorizeDataPlaneRequest
	c.SetTransport(freshTransport)

	return c, nil
}

// NewDataPlaneClient returns a pointer to a Client configured with the shared pooled transport
// (used by storage data-plane clients to avoid SNAT exhaustion).
// To disable HTTP keep-alives for a specific client, call c.SetTransport(client.NewTransport(client.TransportOptions{DisableKeepAlives: true})).
func NewDataPlaneClient(baseUri string, serviceName, apiVersion string) *Client {
	return &Client{
		ApiVersion: apiVersion,
		Client:     client.NewClient(baseUri, serviceName, apiVersion),
	}
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
	if c.ApiVersion != "" {
		query.Set("api-version", c.ApiVersion)
	}

	clientReqID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("generating client request ID: %+v", err)
	}

	query.Set("client-request-id", clientReqID)

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

// CloneClient will take Client pointer a string value named endpoint and returns a deep copy of the provided Client
// pointer into a new pointer and sets the baseClient.BaseUri value to the string value provided in endpoint.
// This function should be called for a new *Client wherever the client may be used for multiple endpoints
// simultaneously
func (c *Client) CloneClient(endpoint string) *Client {
	baseClient := *c.Client
	if c.RequestMiddlewares != nil {
		middlewares := make([]client.RequestMiddleware, len(*c.RequestMiddlewares))
		copy(middlewares, *c.RequestMiddlewares)
		baseClient.RequestMiddlewares = &middlewares
	}
	if c.ResponseMiddlewares != nil {
		middlewares := make([]client.ResponseMiddleware, len(*c.ResponseMiddlewares))
		copy(middlewares, *c.ResponseMiddlewares)
		baseClient.ResponseMiddlewares = &middlewares
	}
	baseClient.BaseUri = endpoint

	return &Client{
		Client:     &baseClient,
		ApiVersion: c.ApiVersion,
	}
}
