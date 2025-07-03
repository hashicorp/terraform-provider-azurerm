package networkconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewNetworkConnectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkConnectionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networkconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkConnectionsClient: %+v", err)
	}

	return &NetworkConnectionsClient{
		Client: client,
	}, nil
}
