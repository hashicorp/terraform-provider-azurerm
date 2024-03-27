package networkmanagerconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkManagerConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewNetworkManagerConnectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkManagerConnectionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "networkmanagerconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkManagerConnectionsClient: %+v", err)
	}

	return &NetworkManagerConnectionsClient{
		Client: client,
	}, nil
}
