package storage

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageClient struct {
	Client *dataplane.Client
}

func NewStorageClientUnconfigured() (*StorageClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "storage", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageClient: %+v", err)
	}

	return &StorageClient{
		Client: client,
	}, nil
}

func (c *StorageClient) StorageClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewStorageClientWithBaseURI(endpoint string) (*StorageClient, error) {
	client, err := dataplane.NewClient(endpoint, "storage", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageClient: %+v", err)
	}

	return &StorageClient{
		Client: client,
	}, nil
}

func (c *StorageClient) Clone(endpoint string) *StorageClient {
	return &StorageClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
