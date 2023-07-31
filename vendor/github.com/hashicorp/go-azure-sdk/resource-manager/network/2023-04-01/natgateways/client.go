package natgateways

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NatGatewaysClient struct {
	Client *resourcemanager.Client
}

func NewNatGatewaysClientWithBaseURI(api environments.Api) (*NatGatewaysClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "natgateways", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NatGatewaysClient: %+v", err)
	}

	return &NatGatewaysClient{
		Client: client,
	}, nil
}
