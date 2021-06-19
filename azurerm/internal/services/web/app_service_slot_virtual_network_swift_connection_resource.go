package web

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	subnetParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAppServiceSlotVirtualNetworkSwiftConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceSlotVirtualNetworkSwiftConnectionCreateUpdate,
		Read:   resourceAppServiceSlotVirtualNetworkSwiftConnectionRead,
		Update: resourceAppServiceSlotVirtualNetworkSwiftConnectionCreateUpdate,
		Delete: resourceAppServiceSlotVirtualNetworkSwiftConnectionDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				ValidateFunc: networkValidate.SubnetID,
			},
			"slot_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceName,
			},
		},
	}
}

func resourceAppServiceSlotVirtualNetworkSwiftConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appID, err := parse.AppServiceID(d.Get("app_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing app service ID %+v", err)
	}
	subnetID, err := subnetParse.SubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("parsing subnet ID %+v", err)
	}

	resourceGroup := appID.ResourceGroup
	name := appID.SiteName
	subnetName := subnetID.Name
	virtualNetworkName := subnetID.VirtualNetworkName
	slotName := d.Get("slot_name").(string)

	if d.IsNewResource() {
		existing, err := client.GetSwiftVirtualNetworkConnectionSlot(ctx, resourceGroup, name, slotName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed checking for presence of existing App Service Slot Swift Network Connection %q (Resource Group %q)", name, resourceGroup)
			}
		}

		if existing.SwiftVirtualNetworkProperties != nil && existing.SwiftVirtualNetworkProperties.SubnetResourceID != nil && *existing.SwiftVirtualNetworkProperties.SubnetResourceID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_slot_virtual_network_swift_connection", *existing.ID)
		}
	}

	locks.ByName(virtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, network.VirtualNetworkResourceName)

	locks.ByName(subnetName, network.SubnetResourceName)
	defer locks.UnlockByName(subnetName, network.SubnetResourceName)

	appServiceExists, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(appServiceExists.Response) {
			return fmt.Errorf("Error retrieving existing App Service %q (Resource Group %q): App Service not found in resource group", name, resourceGroup)
		}
		return fmt.Errorf("Error retrieving existing App Service %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	slotExists, err := client.GetSlot(ctx, resourceGroup, name, slotName)
	if err != nil {
		if utils.ResponseWasNotFound(slotExists.Response) {
			return fmt.Errorf("Error retrieving existing App Service Slot %q (App Service %q / Resource Group %q): App Service not found in resource group", slotName, name, resourceGroup)
		}
		return fmt.Errorf("Error retrieving existing App Service Slot %q (App Service %q / Resource Group %q): %s", slotName, name, resourceGroup, err)
	}

	connectionEnvelope := web.SwiftVirtualNetwork{
		SwiftVirtualNetworkProperties: &web.SwiftVirtualNetworkProperties{
			SubnetResourceID: utils.String(d.Get("subnet_id").(string)),
		},
	}
	if _, err = client.CreateOrUpdateSwiftVirtualNetworkConnectionSlot(ctx, resourceGroup, name, connectionEnvelope, slotName); err != nil {
		return fmt.Errorf("Error creating/updating App Service Slot VNet association between %q (App Service %q / Resource Group %q) and Virtual Network %q: %s", slotName, name, resourceGroup, virtualNetworkName, err)
	}

	read, err := client.GetSwiftVirtualNetworkConnectionSlot(ctx, resourceGroup, name, slotName)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service Slot VNet association between %q (App Service %q / Resource Group %q) and Virtual Network %q: %s", slotName, name, resourceGroup, virtualNetworkName, err)
	}
	d.SetId(*read.ID)

	return resourceAppServiceSlotVirtualNetworkSwiftConnectionRead(d, meta)
}

func resourceAppServiceSlotVirtualNetworkSwiftConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SlotVirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}

	slot, err := client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(slot.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving existing App Service Slot %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}
	appService, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(appService.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving existing App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
	}
	swiftVnet, err := client.GetSwiftVirtualNetworkConnectionSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if utils.ResponseWasNotFound(swiftVnet.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving App Service Slot VNet association for %q (App Service %q / Resource Group %q): %s", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}

	if swiftVnet.SwiftVirtualNetworkProperties == nil {
		return fmt.Errorf("Error retrieving virtual network properties (Slot Name %q / App Service %q / Resource Group %q): `properties` was nil", id.SlotName, id.SiteName, id.ResourceGroup)
	}
	props := *swiftVnet.SwiftVirtualNetworkProperties
	subnetID := props.SubnetResourceID
	if subnetID == nil || *subnetID == "" {
		d.SetId("")
		return nil
	}
	d.Set("subnet_id", subnetID)
	d.Set("app_service_id", appService.ID)
	d.Set("slot_name", id.SlotName)
	return nil
}

func resourceAppServiceSlotVirtualNetworkSwiftConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SlotVirtualNetworkSwiftConnectionID(d.Id())
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

	read, err := client.GetSwiftVirtualNetworkConnectionSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		return fmt.Errorf("Error making read request on virtual network properties (Slot Name %q / App Service %q / Resource Group %q): %+v", id.SlotName, id.SiteName, id.ResourceGroup, err)
	}
	if read.SwiftVirtualNetworkProperties == nil {
		return fmt.Errorf("Error retrieving virtual network properties (Slot Name %q / App Service %q / Resource Group %q): `properties` was nil", id.SlotName, id.SiteName, id.ResourceGroup)
	}
	props := *read.SwiftVirtualNetworkProperties
	subnet := props.SubnetResourceID
	if subnet == nil || *subnet == "" {
		// assume deleted
		return nil
	}

	resp, err := client.DeleteSwiftVirtualNetworkSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting virtual network properties (Slot Name %q / App Service %q / Resource Group %q): %+v", id.SlotName, id.SiteName, id.ResourceGroup, err)
		}
	}

	return nil
}
