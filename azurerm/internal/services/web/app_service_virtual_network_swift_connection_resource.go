package web

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	subnetParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAppServiceVirtualNetworkSwiftConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppServiceVirtualNetworkSwiftConnectionCreateUpdate,
		Read:   resourceAppServiceVirtualNetworkSwiftConnectionRead,
		Update: resourceAppServiceVirtualNetworkSwiftConnectionCreateUpdate,
		Delete: resourceAppServiceVirtualNetworkSwiftConnectionDelete,
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
			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceAppServiceVirtualNetworkSwiftConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appID, err := parse.AppServiceID(d.Get("app_service_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing App Service Resource ID %q", appID)
	}

	subnetID, err := subnetParse.SubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Subnet Resource ID %q", subnetID)
	}

	resourceGroup := appID.ResourceGroup
	name := appID.SiteName
	subnetName := subnetID.Name
	virtualNetworkName := subnetID.VirtualNetworkName

	if d.IsNewResource() {
		existing, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed checking for presence of existing App Service Swift Network Connection %q (Resource Group %q)", name, resourceGroup)
			}
		}

		if existing.SwiftVirtualNetworkProperties.SubnetResourceID != nil && *existing.SwiftVirtualNetworkProperties.SubnetResourceID != "" {
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
			return fmt.Errorf("Error retrieving existing App Service %q (Resource Group %q): App Service not found in resource group", name, resourceGroup)
		}
		return fmt.Errorf("Error retrieving existing App Service %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	connectionEnvelope := web.SwiftVirtualNetwork{
		SwiftVirtualNetworkProperties: &web.SwiftVirtualNetworkProperties{
			SubnetResourceID: utils.String(d.Get("subnet_id").(string)),
		},
	}
	if _, err = client.CreateOrUpdateSwiftVirtualNetworkConnection(ctx, resourceGroup, name, connectionEnvelope); err != nil {
		return fmt.Errorf("Error creating/updating App Service VNet association between %q (Resource Group %q) and Virtual Network %q: %s", name, resourceGroup, virtualNetworkName, err)
	}

	read, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service VNet association between %q (Resource Group %q) and Virtual Network %q: %s", name, resourceGroup, virtualNetworkName, err)
	}
	d.SetId(*read.ID)

	return resourceAppServiceVirtualNetworkSwiftConnectionRead(d, meta)
}

func resourceAppServiceVirtualNetworkSwiftConnectionRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error retrieving existing App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
	}
	swiftVnet, err := client.GetSwiftVirtualNetworkConnection(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(swiftVnet.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving App Service VNet association for %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
	}

	if swiftVnet.SwiftVirtualNetworkProperties == nil {
		return fmt.Errorf("Error retrieving virtual network properties (App Service %q / Resource Group %q): `properties` was nil", id.SiteName, id.ResourceGroup)
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

func resourceAppServiceVirtualNetworkSwiftConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}

	subnetID, err := subnetParse.SubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Subnet Resource ID %q", subnetID)
	}
	subnetName := subnetID.Name
	virtualNetworkName := subnetID.VirtualNetworkName

	locks.ByName(virtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, network.VirtualNetworkResourceName)

	locks.ByName(subnetName, network.SubnetResourceName)
	defer locks.UnlockByName(subnetName, network.SubnetResourceName)

	read, err := client.GetSwiftVirtualNetworkConnection(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("Error making read request on virtual network properties (App Service %q / Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}
	if read.SwiftVirtualNetworkProperties == nil {
		return fmt.Errorf("Error retrieving virtual network properties (App Service %q / Resource Group %q): `properties` was nil", id.SiteName, id.ResourceGroup)
	}
	props := *read.SwiftVirtualNetworkProperties
	subnet := props.SubnetResourceID
	if subnet == nil || *subnet == "" {
		// assume deleted
		return nil
	}

	resp, err := client.DeleteSwiftVirtualNetwork(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting virtual network properties (App Service %q / Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
		}
	}

	return nil
}
