package gatewayinvalidatedebugcredentials

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayInvalidateDebugCredentialsClient struct {
	Client *resourcemanager.Client
}

func NewGatewayInvalidateDebugCredentialsClientWithBaseURI(sdkApi sdkEnv.Api) (*GatewayInvalidateDebugCredentialsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "gatewayinvalidatedebugcredentials", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GatewayInvalidateDebugCredentialsClient: %+v", err)
	}

	return &GatewayInvalidateDebugCredentialsClient{
		Client: client,
	}, nil
}
