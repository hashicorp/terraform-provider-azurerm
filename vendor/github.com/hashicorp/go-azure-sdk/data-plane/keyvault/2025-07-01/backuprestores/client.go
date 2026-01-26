package backuprestores

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackuprestoresClient struct {
	Client *dataplane.Client
}

func NewBackuprestoresClientUnconfigured() (*BackuprestoresClient, error) {
	client, err := dataplane.NewClient("please_configure_client_endpoint", "backuprestores", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackuprestoresClient: %+v", err)
	}

	return &BackuprestoresClient{
		Client: client,
	}, nil
}

func (c *BackuprestoresClient) BackuprestoresClientSetEndpoint(endpoint string) {
	c.Client.Client.BaseUri = endpoint
}

func NewBackuprestoresClientWithBaseURI(endpoint string) (*BackuprestoresClient, error) {
	client, err := dataplane.NewClient(endpoint, "backuprestores", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackuprestoresClient: %+v", err)
	}

	return &BackuprestoresClient{
		Client: client,
	}, nil
}
