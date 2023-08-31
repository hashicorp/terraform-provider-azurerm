// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machineextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MachineExtensionsClient          *machineextensions.MachineExtensionsClient
	MachinesClient                   *machines.MachinesClient
	PrivateEndpointConnectionsClient *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkScopesClient          *privatelinkscopes.PrivateLinkScopesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	machineExtensionsClient, err := machineextensions.NewMachineExtensionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building MachineExtensions client: %+v", err)
	}
	o.Configure(machineExtensionsClient.Client, o.Authorizers.ResourceManager)

	machinesClient, err := machines.NewMachinesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Machines client: %+v", err)
	}
	o.Configure(machinesClient.Client, o.Authorizers.ResourceManager)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	o.Configure(privateEndpointConnectionsClient.Client, o.Authorizers.ResourceManager)

	privateLinkScopesClient, err := privatelinkscopes.NewPrivateLinkScopesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PrivateLinkScopes client: %+v", err)
	}
	o.Configure(privateLinkScopesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MachineExtensionsClient:          machineExtensionsClient,
		MachinesClient:                   machinesClient,
		PrivateEndpointConnectionsClient: privateEndpointConnectionsClient,
		PrivateLinkScopesClient:          privateLinkScopesClient,
	}, nil
}
