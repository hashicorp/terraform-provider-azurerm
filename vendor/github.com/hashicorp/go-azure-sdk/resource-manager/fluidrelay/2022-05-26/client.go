package v2022_05_26

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelaycontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	FluidRelayContainers *fluidrelaycontainers.FluidRelayContainersClient
	FluidRelayServers    *fluidrelayservers.FluidRelayServersClient
}

func NewClientWithBaseURI(api environments.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	fluidRelayContainersClient, err := fluidrelaycontainers.NewFluidRelayContainersClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building FluidRelayContainers client: %+v", err)
	}
	configureFunc(fluidRelayContainersClient.Client)

	fluidRelayServersClient, err := fluidrelayservers.NewFluidRelayServersClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building FluidRelayServers client: %+v", err)
	}
	configureFunc(fluidRelayServersClient.Client)

	return &Client{
		FluidRelayContainers: fluidRelayContainersClient,
		FluidRelayServers:    fluidRelayServersClient,
	}, nil
}
