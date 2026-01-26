package rbacs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RbacsClient struct {
	Client *dataplane.Client
}

func NewRbacsClientUnconfigured() (*RbacsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "rbacs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RbacsClient: %+v", err)
	}

	return &RbacsClient{
		Client: client,
	}, nil
}

func (c *RbacsClient) RbacsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewRbacsClientWithBaseURI(endpoint string) (*RbacsClient, error) {
	client, err := dataplane.NewClient(endpoint, "rbacs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RbacsClient: %+v", err)
	}

	return &RbacsClient{
		Client: client,
	}, nil
}
