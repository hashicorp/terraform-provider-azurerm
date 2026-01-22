// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataplane

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

var _ client.BaseClient = &Client{}

type Client struct {
	*client.Client

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
