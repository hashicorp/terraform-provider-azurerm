package deletedsecrets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedSecretsClient struct {
	Client *dataplane.Client
}

func NewDeletedSecretsClientUnconfigured() (*DeletedSecretsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "deletedsecrets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedSecretsClient: %+v", err)
	}

	return &DeletedSecretsClient{
		Client: client,
	}, nil
}

func (c *DeletedSecretsClient) DeletedSecretsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewDeletedSecretsClientWithBaseURI(endpoint string) (*DeletedSecretsClient, error) {
	client, err := dataplane.NewClient(endpoint, "deletedsecrets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedSecretsClient: %+v", err)
	}

	return &DeletedSecretsClient{
		Client: client,
	}, nil
}

func (c *DeletedSecretsClient) Clone(endpoint string) *DeletedSecretsClient {
	return &DeletedSecretsClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
