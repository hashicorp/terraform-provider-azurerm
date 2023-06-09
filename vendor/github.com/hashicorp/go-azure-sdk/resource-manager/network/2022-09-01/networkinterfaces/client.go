package networkinterfaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfacesClient struct {
	Client *resourcemanager.Client
}

func NewNetworkInterfacesClientWithBaseURI(api environments.Api) (*NetworkInterfacesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "networkinterfaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkInterfacesClient: %+v", err)
	}

	return &NetworkInterfacesClient{
		Client: client,
	}, nil
}
