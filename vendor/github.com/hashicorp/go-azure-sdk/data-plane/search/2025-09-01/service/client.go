package service

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceClient struct {
	Client *dataplane.Client
}

func NewServiceClientUnconfigured() (*ServiceClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "service", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServiceClient: %+v", err)
	}

	return &ServiceClient{
		Client: client,
	}, nil
}

func (c *ServiceClient) ServiceClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewServiceClientWithBaseURI(endpoint string) (*ServiceClient, error) {
	client, err := dataplane.NewClient(endpoint, "service", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServiceClient: %+v", err)
	}

	return &ServiceClient{
		Client: client,
	}, nil
}

func (c *ServiceClient) Clone(endpoint string) *ServiceClient {
	return &ServiceClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
