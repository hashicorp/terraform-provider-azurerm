// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataplane

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type Client struct {
	Client *client.Client

	// ApiVersion specifies the version of the API being used, which (by design) will be consistent across a client
	// as we intentionally split out multiple API Versions into different clients, rather than using composite API
	// Versions/packages which can cause confusion about which version is being used.
	ApiVersion string
}

func NewDataPlaneClient(baseUri string, serviceName, apiVersion string) *Client {
	client := &Client{
		Client:     client.NewClient(baseUri, serviceName, apiVersion),
		ApiVersion: apiVersion,
	}
	return client
}

func (c *Client) Execute(ctx context.Context, req *client.Request) (*client.Response, error) {
	return c.Client.Execute(ctx, req)
}

func (c *Client) ExecutePaged(ctx context.Context, req *client.Request) (*client.Response, error) {
	return c.Client.ExecutePaged(ctx, req)
}
