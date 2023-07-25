// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/globalschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/labs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GlobalLabSchedulesClient *globalschedules.GlobalSchedulesClient
	LabsClient               *labs.LabsClient
	LabSchedulesClient       *schedules.SchedulesClient
	PoliciesClient           *policies.PoliciesClient
	VirtualMachinesClient    *virtualmachines.VirtualMachinesClient
	VirtualNetworksClient    *virtualnetworks.VirtualNetworksClient
}

func NewClient(o *common.ClientOptions) *Client {
	LabsClient := labs.NewLabsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LabsClient.Client, o.ResourceManagerAuthorizer)

	PoliciesClient := policies.NewPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&PoliciesClient.Client, o.ResourceManagerAuthorizer)

	VirtualMachinesClient := virtualmachines.NewVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&VirtualMachinesClient.Client, o.ResourceManagerAuthorizer)

	VirtualNetworksClient := virtualnetworks.NewVirtualNetworksClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&VirtualNetworksClient.Client, o.ResourceManagerAuthorizer)

	LabSchedulesClient := schedules.NewSchedulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LabSchedulesClient.Client, o.ResourceManagerAuthorizer)

	GlobalLabSchedulesClient := globalschedules.NewGlobalSchedulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&GlobalLabSchedulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		GlobalLabSchedulesClient: &GlobalLabSchedulesClient,
		LabsClient:               &LabsClient,
		LabSchedulesClient:       &LabSchedulesClient,
		PoliciesClient:           &PoliciesClient,
		VirtualMachinesClient:    &VirtualMachinesClient,
		VirtualNetworksClient:    &VirtualNetworksClient,
	}
}
