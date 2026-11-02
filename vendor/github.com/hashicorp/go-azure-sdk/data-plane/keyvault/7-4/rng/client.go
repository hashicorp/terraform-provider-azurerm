package rng

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RNGClient struct {
	Client *dataplane.Client
}

func NewRNGClientUnconfigured() (*RNGClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "rng", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RNGClient: %+v", err)
	}

	return &RNGClient{
		Client: client,
	}, nil
}

func (c *RNGClient) RNGClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewRNGClientWithBaseURI(endpoint string) (*RNGClient, error) {
	client, err := dataplane.NewClient(endpoint, "rng", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RNGClient: %+v", err)
	}

	return &RNGClient{
		Client: client,
	}, nil
}

func (c *RNGClient) Clone(endpoint string) *RNGClient {
	return &RNGClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
