package routes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutesClient struct {
	Client *resourcemanager.Client
}

func NewRoutesClientWithBaseURI(api environments.Api) (*RoutesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "routes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RoutesClient: %+v", err)
	}

	return &RoutesClient{
		Client: client,
	}, nil
}
