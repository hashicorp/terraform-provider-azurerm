package gatewaylistkeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayListKeysClient struct {
	Client *resourcemanager.Client
}

func NewGatewayListKeysClientWithBaseURI(sdkApi sdkEnv.Api) (*GatewayListKeysClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "gatewaylistkeys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GatewayListKeysClient: %+v", err)
	}

	return &GatewayListKeysClient{
		Client: client,
	}, nil
}
