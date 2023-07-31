package capacities

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacitiesClient struct {
	Client *resourcemanager.Client
}

func NewCapacitiesClientWithBaseURI(api environments.Api) (*CapacitiesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "capacities", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CapacitiesClient: %+v", err)
	}

	return &CapacitiesClient{
		Client: client,
	}, nil
}
