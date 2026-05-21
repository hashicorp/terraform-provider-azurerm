package managedprivateendpoints

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedPrivateEndpointsClient struct {
	Client *dataplane.Client
}

func NewManagedPrivateEndpointsClientUnconfigured() (*ManagedPrivateEndpointsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "managedprivateendpoints", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedPrivateEndpointsClient: %+v", err)
	}

	return &ManagedPrivateEndpointsClient{
		Client: client,
	}, nil
}

func (c *ManagedPrivateEndpointsClient) ManagedPrivateEndpointsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewManagedPrivateEndpointsClientWithBaseURI(endpoint string) (*ManagedPrivateEndpointsClient, error) {
	client, err := dataplane.NewClient(endpoint, "managedprivateendpoints", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedPrivateEndpointsClient: %+v", err)
	}

	return &ManagedPrivateEndpointsClient{
		Client: client,
	}, nil
}

func (c *ManagedPrivateEndpointsClient) Clone(endpoint string) *ManagedPrivateEndpointsClient {
	return &ManagedPrivateEndpointsClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
