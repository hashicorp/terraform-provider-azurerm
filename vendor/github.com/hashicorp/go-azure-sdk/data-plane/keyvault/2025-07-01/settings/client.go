package settings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SettingsClient struct {
	Client *dataplane.Client
}

func NewSettingsClientUnconfigured() (*SettingsClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "settings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SettingsClient: %+v", err)
	}

	return &SettingsClient{
		Client: client,
	}, nil
}

func (c *SettingsClient) SettingsClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewSettingsClientWithBaseURI(endpoint string) (*SettingsClient, error) {
	client, err := dataplane.NewClient(endpoint, "settings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SettingsClient: %+v", err)
	}

	return &SettingsClient{
		Client: client,
	}, nil
}
