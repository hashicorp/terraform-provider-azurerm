package synonymmaps

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynonymMapsClient struct {
	Client *dataplane.Client
}

func NewSynonymMapsClientUnconfigured() (*SynonymMapsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "synonymmaps", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SynonymMapsClient: %+v", err)
	}

	return &SynonymMapsClient{
		Client: client,
	}, nil
}

func (c *SynonymMapsClient) SynonymMapsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewSynonymMapsClientWithBaseURI(endpoint string) (*SynonymMapsClient, error) {
	client, err := dataplane.NewClient(endpoint, "synonymmaps", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SynonymMapsClient: %+v", err)
	}

	return &SynonymMapsClient{
		Client: client,
	}, nil
}

func (c *SynonymMapsClient) Clone(endpoint string) *SynonymMapsClient {
	return &SynonymMapsClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
