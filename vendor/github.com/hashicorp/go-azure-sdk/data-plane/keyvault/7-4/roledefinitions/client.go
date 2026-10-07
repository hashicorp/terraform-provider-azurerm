package roledefinitions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleDefinitionsClient struct {
	Client *dataplane.Client
}

func NewRoleDefinitionsClientUnconfigured() (*RoleDefinitionsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "roledefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoleDefinitionsClient: %+v", err)
	}

	return &RoleDefinitionsClient{
		Client: client,
	}, nil
}

func (c *RoleDefinitionsClient) RoleDefinitionsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewRoleDefinitionsClientWithBaseURI(endpoint string) (*RoleDefinitionsClient, error) {
	client, err := dataplane.NewClient(endpoint, "roledefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoleDefinitionsClient: %+v", err)
	}

	return &RoleDefinitionsClient{
		Client: client,
	}, nil
}

func (c *RoleDefinitionsClient) Clone(endpoint string) *RoleDefinitionsClient {
	return &RoleDefinitionsClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
