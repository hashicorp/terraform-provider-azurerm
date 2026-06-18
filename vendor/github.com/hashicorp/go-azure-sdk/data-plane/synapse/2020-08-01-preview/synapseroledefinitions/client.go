package synapseroledefinitions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynapseRoleDefinitionsClient struct {
	Client *dataplane.Client
}

func NewSynapseRoleDefinitionsClientUnconfigured() (*SynapseRoleDefinitionsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "synapseroledefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SynapseRoleDefinitionsClient: %+v", err)
	}

	return &SynapseRoleDefinitionsClient{
		Client: client,
	}, nil
}

func (c *SynapseRoleDefinitionsClient) SynapseRoleDefinitionsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewSynapseRoleDefinitionsClientWithBaseURI(endpoint string) (*SynapseRoleDefinitionsClient, error) {
	client, err := dataplane.NewClient(endpoint, "synapseroledefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SynapseRoleDefinitionsClient: %+v", err)
	}

	return &SynapseRoleDefinitionsClient{
		Client: client,
	}, nil
}

func (c *SynapseRoleDefinitionsClient) Clone(endpoint string) *SynapseRoleDefinitionsClient {
	return &SynapseRoleDefinitionsClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
