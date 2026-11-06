package indexes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndexesClient struct {
	Client *dataplane.Client
}

func NewIndexesClientUnconfigured() (*IndexesClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "indexes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IndexesClient: %+v", err)
	}

	return &IndexesClient{
		Client: client,
	}, nil
}

func (c *IndexesClient) IndexesClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewIndexesClientWithBaseURI(endpoint string) (*IndexesClient, error) {
	client, err := dataplane.NewClient(endpoint, "indexes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IndexesClient: %+v", err)
	}

	return &IndexesClient{
		Client: client,
	}, nil
}

func (c *IndexesClient) Clone(endpoint string) *IndexesClient {
	return &IndexesClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
