// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServiceVirtualNetworkSwiftConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceVirtualNetworkSwiftConnectionCreateUpdate,
		Read:   resourceAppServiceVirtualNetworkSwiftConnectionRead,
		Update: resourceAppServiceVirtualNetworkSwiftConnectionCreateUpdate,
		Delete: resourceAppServiceVirtualNetworkSwiftConnectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualNetworkSwiftConnectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"app_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceID,
			},
			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},
		},
	}
}

func resourceAppServiceVirtualNetworkSwiftConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	subnetClient := meta.(*clients.Client).Network.Client.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appID, err := parse.AppServiceID(d.Get("app_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing App Service Resource ID %q", appID)
	}

	subnetID, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Subnet Resource ID %q", subnetID)
	}

	resourceGroup := appID.ResourceGroup
	name := appID.SiteName
	subnetName := subnetID.SubnetName
	virtualNetworkName := subnetID.VirtualNetworkName

	if d.IsNewResource() {
		existing, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed checking for presence of existing App Service Swift Network Connection %q (Resource Group %q)", name, resourceGroup)
			}
		}

		if existing.SwiftVirtualNetworkProperties != nil && existing.SwiftVirtualNetworkProperties.SubnetResourceID != nil && *existing.SwiftVirtualNetworkProperties.SubnetResourceID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_virtual_network_swift_connection", *existing.ID)
		}
	}

	locks.ByName(virtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, network.VirtualNetworkResourceName)

	locks.ByName(subnetName, network.SubnetResourceName)
	defer locks.UnlockByName(subnetName, network.SubnetResourceName)

	exists, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(exists.Response) {
			return fmt.Errorf("retrieving existing App Service %q (Resource Group %q): App Service not found in resource group", name, resourceGroup)
		}
		return fmt.Errorf("retrieving existing App Service %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	connectionEnvelope := web.SwiftVirtualNetwork{
		SwiftVirtualNetworkProperties: &web.SwiftVirtualNetworkProperties{
			SubnetResourceID: utils.String(d.Get("subnet_id").(string)),
		},
	}
	if _, err = client.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheck(ctx, resourceGroup, name, connectionEnvelope); err != nil {
		return fmt.Errorf("creating/updating App Service VNet association between %q (Resource Group %q) and Virtual Network %q: %s", name, resourceGroup, virtualNetworkName, err)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    network.SubnetProvisioningStateRefreshFunc(ctx, subnetClient, *subnetID),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of subnet for App Service VNet association between %q (Resource Group %q) and Virtual Network %q: %s", name, resourceGroup, virtualNetworkName, err)
	}

	vnetId := commonids.NewVirtualNetworkID(subnetID.SubscriptionId, subnetID.ResourceGroupName, subnetID.VirtualNetworkName)
	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    network.VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for App Service VNet association between %q (Resource Group %q) and Virtual Network %q: %s", name, resourceGroup, virtualNetworkName, err)
	}

	read, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving App Service VNet association between %q (Resource Group %q) and Virtual Network %q: %s", name, resourceGroup, virtualNetworkName, err)
	}
	swiftVirtualNetworkId, err := parse.VirtualNetworkSwiftConnectionID(*read.ID)
	if err != nil {
		return err
	}
	d.SetId(swiftVirtualNetworkId.ID())

	return resourceAppServiceVirtualNetworkSwiftConnectionRead(d, meta)
}

func resourceAppServiceVirtualNetworkSwiftConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}

	appService, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(appService.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving existing App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
	}
	swiftVnet, err := client.GetSwiftVirtualNetworkConnection(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(swiftVnet.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving App Service VNet association for %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
	}

	if swiftVnet.SwiftVirtualNetworkProperties == nil {
		return fmt.Errorf("retrieving virtual network properties (App Service %q / Resource Group %q): `properties` was nil", id.SiteName, id.ResourceGroup)
	}
	props := *swiftVnet.SwiftVirtualNetworkProperties
	subnetID := props.SubnetResourceID
	if subnetID == nil || *subnetID == "" {
		d.SetId("")
		return nil
	}
	d.Set("subnet_id", subnetID)
	d.Set("app_service_id", appService.ID)
	return nil
}

func resourceAppServiceVirtualNetworkSwiftConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.GetSwiftVirtualNetworkConnection(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("making read request on virtual network properties (App Service %q / Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}
	if read.SwiftVirtualNetworkProperties == nil {
		return fmt.Errorf("retrieving virtual network properties (App Service %q / Resource Group %q): `properties` was nil", id.SiteName, id.ResourceGroup)
	}
	props := *read.SwiftVirtualNetworkProperties
	subnet := props.SubnetResourceID
	if subnet == nil || *subnet == "" {
		// assume deleted
		return nil
	}

	subnetID, err := commonids.ParseSubnetID(pointer.From(subnet))
	if err != nil {
		return fmt.Errorf("parsing Subnet Resource ID %q", subnetID)
	}
	subnetName := subnetID.SubnetName
	virtualNetworkName := subnetID.VirtualNetworkName

	locks.ByName(virtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, network.VirtualNetworkResourceName)

	locks.ByName(subnetName, network.SubnetResourceName)
	defer locks.UnlockByName(subnetName, network.SubnetResourceName)

	resp, err := client.DeleteSwiftVirtualNetwork(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting virtual network properties (App Service %q / Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
		}
	}

	return nil
}
