package keys

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeysClient struct {
	Client *dataplane.Client
}

func NewKeysClientUnconfigured() (*KeysClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "keys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating KeysClient: %+v", err)
	}

	return &KeysClient{
		Client: client,
	}, nil
}

func (c *KeysClient) KeysClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewKeysClientWithBaseURI(endpoint string) (*KeysClient, error) {
	client, err := dataplane.NewClient(endpoint, "keys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating KeysClient: %+v", err)
	}

	return &KeysClient{
		Client: client,
	}, nil
}

func (c *KeysClient) Clone(endpoint string) *KeysClient {
	return &KeysClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
