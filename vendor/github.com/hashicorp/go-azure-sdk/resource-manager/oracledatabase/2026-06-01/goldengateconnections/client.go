package goldengateconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GoldenGateConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewGoldenGateConnectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*GoldenGateConnectionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "goldengateconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GoldenGateConnectionsClient: %+v", err)
	}

	return &GoldenGateConnectionsClient{
		Client: client,
	}, nil
}
