package onlineendpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OnlineEndpointClient struct {
	Client *resourcemanager.Client
}

func NewOnlineEndpointClientWithBaseURI(api environments.Api) (*OnlineEndpointClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "onlineendpoint", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OnlineEndpointClient: %+v", err)
	}

	return &OnlineEndpointClient{
		Client: client,
	}, nil
}
