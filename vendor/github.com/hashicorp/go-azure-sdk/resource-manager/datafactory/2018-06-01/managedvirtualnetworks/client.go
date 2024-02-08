package managedvirtualnetworks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedVirtualNetworksClient struct {
	Client *resourcemanager.Client
}

func NewManagedVirtualNetworksClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedVirtualNetworksClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "managedvirtualnetworks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedVirtualNetworksClient: %+v", err)
	}

	return &ManagedVirtualNetworksClient{
		Client: client,
	}, nil
}
