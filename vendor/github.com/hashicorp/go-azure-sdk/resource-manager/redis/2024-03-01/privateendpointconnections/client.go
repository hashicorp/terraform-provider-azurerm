package privateendpointconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewPrivateEndpointConnectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*PrivateEndpointConnectionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "privateendpointconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrivateEndpointConnectionsClient: %+v", err)
	}

	return &PrivateEndpointConnectionsClient{
		Client: client,
	}, nil
}
