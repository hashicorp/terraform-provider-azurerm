package natgateways

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NatGatewaysClient struct {
	Client *resourcemanager.Client
}

func NewNatGatewaysClientWithBaseURI(sdkApi sdkEnv.Api) (*NatGatewaysClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "natgateways", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NatGatewaysClient: %+v", err)
	}

	return &NatGatewaysClient{
		Client: client,
	}, nil
}
