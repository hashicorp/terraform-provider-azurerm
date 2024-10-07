// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

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

func NewClient(o *common.ClientOptions) (*Client, error) {
	globalLabSchedulesClient, err := globalschedules.NewGlobalSchedulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GlobalSchedules client: %+v", err)
	}
	o.Configure(globalLabSchedulesClient.Client, o.Authorizers.ResourceManager)

	labsClient, err := labs.NewLabsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Labs client: %+v", err)
	}
	o.Configure(labsClient.Client, o.Authorizers.ResourceManager)

	labSchedulesClient, err := schedules.NewSchedulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building LabSchedules client: %+v", err)
	}
	o.Configure(labSchedulesClient.Client, o.Authorizers.ResourceManager)

	policiesClient, err := policies.NewPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Policies client: %+v", err)
	}
	o.Configure(policiesClient.Client, o.Authorizers.ResourceManager)

	virtualMachinesClient, err := virtualmachines.NewVirtualMachinesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachines client: %+v", err)
	}
	o.Configure(virtualMachinesClient.Client, o.Authorizers.ResourceManager)

	virtualNetworksClient, err := virtualnetworks.NewVirtualNetworksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualNetworks client: %+v", err)
	}
	o.Configure(virtualNetworksClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		GlobalLabSchedulesClient: globalLabSchedulesClient,
		LabsClient:               labsClient,
		LabSchedulesClient:       labSchedulesClient,
		PoliciesClient:           policiesClient,
		VirtualMachinesClient:    virtualMachinesClient,
		VirtualNetworksClient:    virtualNetworksClient,
	}, nil
}
