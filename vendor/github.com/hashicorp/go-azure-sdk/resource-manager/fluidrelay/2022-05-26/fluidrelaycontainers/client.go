package fluidrelaycontainers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluidRelayContainersClient struct {
	Client *resourcemanager.Client
}

func NewFluidRelayContainersClientWithBaseURI(api environments.Api) (*FluidRelayContainersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "fluidrelaycontainers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FluidRelayContainersClient: %+v", err)
	}

	return &FluidRelayContainersClient{
		Client: client,
	}, nil
}
