package fullrestore

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FullRestoreClient struct {
	Client *dataplane.Client
}

func NewFullRestoreClientUnconfigured() (*FullRestoreClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "fullrestore", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FullRestoreClient: %+v", err)
	}

	return &FullRestoreClient{
		Client: client,
	}, nil
}

func (c *FullRestoreClient) FullRestoreClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewFullRestoreClientWithBaseURI(endpoint string) (*FullRestoreClient, error) {
	client, err := dataplane.NewClient(endpoint, "fullrestore", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FullRestoreClient: %+v", err)
	}

	return &FullRestoreClient{
		Client: client,
	}, nil
}

func (c *FullRestoreClient) Clone(endpoint string) *FullRestoreClient {
	return &FullRestoreClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
