package deletedkeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedKeysClient struct {
	Client *dataplane.Client
}

func NewDeletedKeysClientUnconfigured() (*DeletedKeysClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "deletedkeys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedKeysClient: %+v", err)
	}

	return &DeletedKeysClient{
		Client: client,
	}, nil
}

func (c *DeletedKeysClient) DeletedKeysClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewDeletedKeysClientWithBaseURI(endpoint string) (*DeletedKeysClient, error) {
	client, err := dataplane.NewClient(endpoint, "deletedkeys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedKeysClient: %+v", err)
	}

	return &DeletedKeysClient{
		Client: client,
	}, nil
}

func (c *DeletedKeysClient) Clone(endpoint string) *DeletedKeysClient {
	return &DeletedKeysClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
