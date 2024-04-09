package gatewayapi

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayApiClient struct {
	Client *resourcemanager.Client
}

func NewGatewayApiClientWithBaseURI(sdkApi sdkEnv.Api) (*GatewayApiClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "gatewayapi", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GatewayApiClient: %+v", err)
	}

	return &GatewayApiClient{
		Client: client,
	}, nil
}
