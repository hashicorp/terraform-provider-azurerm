package web

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceVirtualNetworkGatewayConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceVirtualNetworkGatewayConnectionCreateUpdate,
		Read:   resourceArmAppServiceVirtualNetworkGatewayConnectionRead,
		Update: resourceArmAppServiceVirtualNetworkGatewayConnectionCreateUpdate,
		Delete: resourceArmAppServiceVirtualNetworkGatewayConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"app_service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"vnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"vpn_gateway_package_uri": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				DiffSuppressFunc: func(_, old, _ string, _ *schema.ResourceData) bool {
					return old != ""
				},
			},
		},
	}
}

func resourceArmAppServiceVirtualNetworkGatewayConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Get("app_service_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
	vnetIDRaw := d.Get("vnet_id").(string)
	vnetID, err := azure.ParseAzureResourceID(vnetIDRaw)
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
	packageURI := d.Get("vpn_gateway_package_uri").(string)

	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]
	virtualNetworkName := vnetID.Path["virtualNetworks"]

	locks.ByName(virtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, network.VirtualNetworkResourceName)

	exists, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(exists.Response) {
			return fmt.Errorf("Error retrieving existing App Service %q (Resource Group %q): App Service not found in resource group", name, resourceGroup)
		}
		return fmt.Errorf("Error retrieving existing App Service %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	vnetInfo := web.VnetInfo{
		VnetInfoProperties: &web.VnetInfoProperties{
			VnetResourceID: utils.String(vnetIDRaw),
		},
	}
	if _, err = client.CreateOrUpdateVnetConnection(ctx, resourceGroup, name, virtualNetworkName, vnetInfo); err != nil {
		return fmt.Errorf("Error creating/updating App Service VNet association between %q (Resource Group %q) and Virtual Network %q: %s", name, resourceGroup, virtualNetworkName, err)
	}

	connectionEnvelope := web.VnetGateway{
		VnetGatewayProperties: &web.VnetGatewayProperties{
			VnetName:      utils.String(virtualNetworkName),
			VpnPackageURI: utils.String(packageURI),
		},
	}
	gatewayName := "primary"
	if _, err = client.CreateOrUpdateVnetConnectionGateway(ctx, resourceGroup, name, virtualNetworkName, gatewayName, connectionEnvelope); err != nil {
		return fmt.Errorf("Error creating/updating App Service VNet Gateway association between %q (Resource Group %q) and Virtual Network %q: %s", name, resourceGroup, virtualNetworkName, err)
	}
	read, err := client.GetVnetConnectionGateway(ctx, resourceGroup, name, virtualNetworkName, gatewayName)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service VNet Gateway association between %q (Resource Group %q) and Virtual Network %q: %s", name, resourceGroup, virtualNetworkName, err)
	}
	d.SetId(*read.ID)

	return resourceArmAppServiceVirtualNetworkGatewayConnectionRead(d, meta)
}

func resourceArmAppServiceVirtualNetworkGatewayConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
	virtualNetworkName := id.Path["virtualNetworkConnections"]
	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]

	appService, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(appService.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving existing App Service %q (Resource Group %q): %s", name, resourceGroup, err)
	}
	vnetConnection, err := client.GetVnetConnection(ctx, resourceGroup, name, virtualNetworkName)
	if err != nil {
		if utils.ResponseWasNotFound(vnetConnection.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving App Service VNet association for %q (Resource Group %q): %s", name, resourceGroup, err)
	}
	gatewayName := "primary"
	vnetGateway, err := client.GetVnetConnectionGateway(ctx, resourceGroup, name, virtualNetworkName, gatewayName)
	if err != nil {
		if utils.ResponseWasNotFound(vnetGateway.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving App Service VNet Gateway association for %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	d.Set("vnet_id", vnetConnection.VnetInfoProperties.VnetResourceID)
	d.Set("app_service_id", appService.ID)
	return nil
}

func resourceArmAppServiceVirtualNetworkGatewayConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Get("app_service_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
	vnetID, err := azure.ParseAzureResourceID(d.Get("vnet_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", vnetID)
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]
	virtualNetworkName := vnetID.Path["virtualNetworks"]

	locks.ByName(virtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, network.VirtualNetworkResourceName)

	gatewayName := "primary"
	read, err := client.GetVnetConnectionGateway(ctx, resourceGroup, name, virtualNetworkName, gatewayName)
	if err != nil {
		return fmt.Errorf("Error making read request on virtual network properties (App Service %q / Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.VnetGatewayProperties == nil {
		return fmt.Errorf("Error retrieving virtual network properties (App Service %q / Resource Group %q): `properties` was nil", name, resourceGroup)
	}
	props := *read.VnetGatewayProperties
	vnet := props.VnetName
	if vnet == nil || *vnet == "" {
		// assume deleted
		return nil
	}

	resp, err := client.DeleteVnetConnection(ctx, resourceGroup, name, virtualNetworkName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting virtual network properties (App Service %q / Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
