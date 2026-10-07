package deletedstorage

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedStorageClient struct {
	Client *dataplane.Client
}

func NewDeletedStorageClientUnconfigured() (*DeletedStorageClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "deletedstorage", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedStorageClient: %+v", err)
	}

	return &DeletedStorageClient{
		Client: client,
	}, nil
}

func (c *DeletedStorageClient) DeletedStorageClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewDeletedStorageClientWithBaseURI(endpoint string) (*DeletedStorageClient, error) {
	client, err := dataplane.NewClient(endpoint, "deletedstorage", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedStorageClient: %+v", err)
	}

	return &DeletedStorageClient{
		Client: client,
	}, nil
}

func (c *DeletedStorageClient) Clone(endpoint string) *DeletedStorageClient {
	return &DeletedStorageClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
