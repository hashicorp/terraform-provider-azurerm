package skillsets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkillsetsClient struct {
	Client *dataplane.Client
}

func NewSkillsetsClientUnconfigured() (*SkillsetsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "skillsets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SkillsetsClient: %+v", err)
	}

	return &SkillsetsClient{
		Client: client,
	}, nil
}

func (c *SkillsetsClient) SkillsetsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewSkillsetsClientWithBaseURI(endpoint string) (*SkillsetsClient, error) {
	client, err := dataplane.NewClient(endpoint, "skillsets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SkillsetsClient: %+v", err)
	}

	return &SkillsetsClient{
		Client: client,
	}, nil
}

func (c *SkillsetsClient) Clone(endpoint string) *SkillsetsClient {
	return &SkillsetsClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
