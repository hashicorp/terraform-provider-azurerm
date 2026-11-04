package storagediscoveryworkspaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageDiscoveryWorkspacesClient struct {
	Client *resourcemanager.Client
}

func NewStorageDiscoveryWorkspacesClientWithBaseURI(sdkApi sdkEnv.Api) (*StorageDiscoveryWorkspacesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "storagediscoveryworkspaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageDiscoveryWorkspacesClient: %+v", err)
	}

	return &StorageDiscoveryWorkspacesClient{
		Client: client,
	}, nil
}
