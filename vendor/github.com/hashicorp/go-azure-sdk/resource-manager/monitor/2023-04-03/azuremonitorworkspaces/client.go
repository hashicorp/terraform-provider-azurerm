package azuremonitorworkspaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMonitorWorkspacesClient struct {
	Client *resourcemanager.Client
}

func NewAzureMonitorWorkspacesClientWithBaseURI(sdkApi sdkEnv.Api) (*AzureMonitorWorkspacesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "azuremonitorworkspaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AzureMonitorWorkspacesClient: %+v", err)
	}

	return &AzureMonitorWorkspacesClient{
		Client: client,
	}, nil
}
