package secrets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretsClient struct {
	Client *dataplane.Client
}

func NewSecretsClientUnconfigured() (*SecretsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "secrets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SecretsClient: %+v", err)
	}

	return &SecretsClient{
		Client: client,
	}, nil
}

func (c *SecretsClient) SecretsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewSecretsClientWithBaseURI(endpoint string) (*SecretsClient, error) {
	client, err := dataplane.NewClient(endpoint, "secrets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SecretsClient: %+v", err)
	}

	return &SecretsClient{
		Client: client,
	}, nil
}

func (c *SecretsClient) Clone(endpoint string) *SecretsClient {
	return &SecretsClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
