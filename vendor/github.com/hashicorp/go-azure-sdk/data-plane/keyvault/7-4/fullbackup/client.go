package fullbackup

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FullBackupClient struct {
	Client *dataplane.Client
}

func NewFullBackupClientUnconfigured() (*FullBackupClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "fullbackup", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FullBackupClient: %+v", err)
	}

	return &FullBackupClient{
		Client: client,
	}, nil
}

func (c *FullBackupClient) FullBackupClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewFullBackupClientWithBaseURI(endpoint string) (*FullBackupClient, error) {
	client, err := dataplane.NewClient(endpoint, "fullbackup", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FullBackupClient: %+v", err)
	}

	return &FullBackupClient{
		Client: client,
	}, nil
}

func (c *FullBackupClient) Clone(endpoint string) *FullBackupClient {
	return &FullBackupClient{
		Client: c.Client.CloneClient(endpoint),
	}
}
