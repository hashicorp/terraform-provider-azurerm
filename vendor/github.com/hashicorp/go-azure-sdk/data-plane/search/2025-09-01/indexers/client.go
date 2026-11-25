package indexers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndexersClient struct {
	Client *dataplane.Client
}

func NewIndexersClientUnconfigured() (*IndexersClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "indexers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IndexersClient: %+v", err)
	}

	return &IndexersClient{
		Client: client,
	}, nil
}

func (c *IndexersClient) IndexersClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewIndexersClientWithBaseURI(endpoint string) (*IndexersClient, error) {
	client, err := dataplane.NewClient(endpoint, "indexers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IndexersClient: %+v", err)
	}

	return &IndexersClient{
		Client: client,
	}, nil
}

func (c *IndexersClient) Clone(endpoint string) *IndexersClient {
	return &IndexersClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
