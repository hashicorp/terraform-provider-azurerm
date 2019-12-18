package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceVirtualNetworkSwiftConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceVirtualNetworkSwiftConnectionCreateUpdate,
		Read:   resourceArmAppServiceVirtualNetworkSwiftConnectionRead,
		Update: resourceArmAppServiceVirtualNetworkSwiftConnectionCreateUpdate,
		Delete: resourceArmAppServiceVirtualNetworkSwiftConnectionDelete,
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
				ForceNew:     false,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmAppServiceVirtualNetworkSwiftConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Get("app_service_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
	subnetID, err := azure.ParseAzureResourceID(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]
	subnetName := subnetID.Path["subnets"]
	virtualNetworkName := subnetID.Path["virtualNetworks"]

	locks.ByName(virtualNetworkName, virtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, virtualNetworkResourceName)

	locks.ByName(subnetName, subnetResourceName)
	defer locks.UnlockByName(subnetName, subnetResourceName)

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

	return resourceArmAppServiceVirtualNetworkSwiftConnectionRead(d, meta)
}

func resourceArmAppServiceVirtualNetworkSwiftConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
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
	resp, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving App Service VNet association for %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	if resp.SwiftVirtualNetworkProperties == nil {
		return fmt.Errorf("Error retrieving virtual network properties (App Service %q / Resource Group %q): `properties` was nil", name, resourceGroup)
	}
	props := *resp.SwiftVirtualNetworkProperties
	subnetID := props.SubnetResourceID
	if subnetID == nil || *subnetID == "" {
		d.SetId("")
		return nil
	}
	d.Set("subnet_id", subnetID)
	d.Set("app_service_id", appService.ID)
	return nil
}

func resourceArmAppServiceVirtualNetworkSwiftConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Get("app_service_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
	subnetID, err := azure.ParseAzureResourceID(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]
	subnetName := subnetID.Path["subnets"]
	virtualNetworkName := subnetID.Path["virtualNetworks"]

	locks.ByName(virtualNetworkName, virtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, virtualNetworkResourceName)

	locks.ByName(subnetName, subnetResourceName)
	defer locks.UnlockByName(subnetName, subnetResourceName)

	read, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making read request on virtual network properties (App Service %q / Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.SwiftVirtualNetworkProperties == nil {
		return fmt.Errorf("Error retrieving virtual network properties (App Service %q / Resource Group %q): `properties` was nil", name, resourceGroup)
	}
	props := *read.SwiftVirtualNetworkProperties
	subnet := props.SubnetResourceID
	if subnet == nil || *subnet == "" {
		// assume deleted
		return nil
	}

	resp, err := client.DeleteSwiftVirtualNetwork(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting virtual network properties (App Service %q / Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
