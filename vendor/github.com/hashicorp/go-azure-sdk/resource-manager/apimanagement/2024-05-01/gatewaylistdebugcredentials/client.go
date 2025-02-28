package gatewaylistdebugcredentials

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayListDebugCredentialsClient struct {
	Client *resourcemanager.Client
}

func NewGatewayListDebugCredentialsClientWithBaseURI(sdkApi sdkEnv.Api) (*GatewayListDebugCredentialsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "gatewaylistdebugcredentials", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GatewayListDebugCredentialsClient: %+v", err)
	}

	return &GatewayListDebugCredentialsClient{
		Client: client,
	}, nil
}
