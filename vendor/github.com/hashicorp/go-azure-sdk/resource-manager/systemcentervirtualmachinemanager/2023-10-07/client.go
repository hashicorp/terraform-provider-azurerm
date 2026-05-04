package v2023_10_07

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/clouds"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/guestagents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachinetemplates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vminstancehybrididentitymetadatas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AvailabilitySets                  *availabilitysets.AvailabilitySetsClient
	Clouds                            *clouds.CloudsClient
	GuestAgents                       *guestagents.GuestAgentsClient
	InventoryItems                    *inventoryitems.InventoryItemsClient
	VMInstanceHybridIdentityMetadatas *vminstancehybrididentitymetadatas.VMInstanceHybridIdentityMetadatasClient
	VMmServers                        *vmmservers.VMmServersClient
	VirtualMachineInstances           *virtualmachineinstances.VirtualMachineInstancesClient
	VirtualMachineTemplates           *virtualmachinetemplates.VirtualMachineTemplatesClient
	VirtualNetworks                   *virtualnetworks.VirtualNetworksClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	availabilitySetsClient, err := availabilitysets.NewAvailabilitySetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AvailabilitySets client: %+v", err)
	}
	configureFunc(availabilitySetsClient.Client)

	cloudsClient, err := clouds.NewCloudsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Clouds client: %+v", err)
	}
	configureFunc(cloudsClient.Client)

	guestAgentsClient, err := guestagents.NewGuestAgentsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GuestAgents client: %+v", err)
	}
	configureFunc(guestAgentsClient.Client)

	inventoryItemsClient, err := inventoryitems.NewInventoryItemsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building InventoryItems client: %+v", err)
	}
	configureFunc(inventoryItemsClient.Client)

	vMInstanceHybridIdentityMetadatasClient, err := vminstancehybrididentitymetadatas.NewVMInstanceHybridIdentityMetadatasClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VMInstanceHybridIdentityMetadatas client: %+v", err)
	}
	configureFunc(vMInstanceHybridIdentityMetadatasClient.Client)

	vMmServersClient, err := vmmservers.NewVMmServersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VMmServers client: %+v", err)
	}
	configureFunc(vMmServersClient.Client)

	virtualMachineInstancesClient, err := virtualmachineinstances.NewVirtualMachineInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachineInstances client: %+v", err)
	}
	configureFunc(virtualMachineInstancesClient.Client)

	virtualMachineTemplatesClient, err := virtualmachinetemplates.NewVirtualMachineTemplatesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachineTemplates client: %+v", err)
	}
	configureFunc(virtualMachineTemplatesClient.Client)

	virtualNetworksClient, err := virtualnetworks.NewVirtualNetworksClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VirtualNetworks client: %+v", err)
	}
	configureFunc(virtualNetworksClient.Client)

	return &Client{
		AvailabilitySets:                  availabilitySetsClient,
		Clouds:                            cloudsClient,
		GuestAgents:                       guestAgentsClient,
		InventoryItems:                    inventoryItemsClient,
		VMInstanceHybridIdentityMetadatas: vMInstanceHybridIdentityMetadatasClient,
		VMmServers:                        vMmServersClient,
		VirtualMachineInstances:           virtualMachineInstancesClient,
		VirtualMachineTemplates:           virtualMachineTemplatesClient,
		VirtualNetworks:                   virtualNetworksClient,
	}, nil
}
