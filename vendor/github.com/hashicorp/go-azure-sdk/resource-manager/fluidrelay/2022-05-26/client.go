package v2022_05_26

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelaycontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Client struct {
	FluidRelayContainers *fluidrelaycontainers.FluidRelayContainersClient
	FluidRelayServers    *fluidrelayservers.FluidRelayServersClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	fluidRelayContainersClient := fluidrelaycontainers.NewFluidRelayContainersClientWithBaseURI(endpoint)
	configureAuthFunc(&fluidRelayContainersClient.Client)

	fluidRelayServersClient := fluidrelayservers.NewFluidRelayServersClientWithBaseURI(endpoint)
	configureAuthFunc(&fluidRelayServersClient.Client)

	return Client{
		FluidRelayContainers: &fluidRelayContainersClient,
		FluidRelayServers:    &fluidRelayServersClient,
	}
}
