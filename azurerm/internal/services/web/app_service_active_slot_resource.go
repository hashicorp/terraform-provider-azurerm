package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAppServiceActiveSlot() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceActiveSlotCreateUpdate,
		Read:   resourceAppServiceActiveSlotRead,
		Update: resourceAppServiceActiveSlotCreateUpdate,
		Delete: resourceAppServiceActiveSlotDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AppServiceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{

			"resource_group_name": azure.SchemaResourceGroupName(),

			"app_service_name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Required: true,
			},

			"app_service_slot_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAppServiceActiveSlotCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appServiceName := d.Get("app_service_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	targetSlot := d.Get("app_service_slot_name").(string)
	preserveVnet := true

	resp, err := client.Get(ctx, resGroup, appServiceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("[DEBUG] App Service %q (resource group %q) was not found.", appServiceName, resGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %q: %+v", appServiceName, err)
	}

	if _, err = client.Get(ctx, resGroup, targetSlot); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("[DEBUG] App Service Target Active Slot %q/%q (resource group %q) was not found.", appServiceName, targetSlot, resGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service Slot %q/%q: %+v", appServiceName, targetSlot, err)
	}

	cmsSlotEntity := web.CsmSlotEntity{
		TargetSlot:   &targetSlot,
		PreserveVnet: &preserveVnet,
	}

	future, err := client.SwapSlotWithProduction(ctx, resGroup, appServiceName, cmsSlotEntity)
	if err != nil {
		return fmt.Errorf("Error swapping App Service Slot %q/%q: %+v", appServiceName, targetSlot, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error swapping App Service Slot %q/%q: %+v", appServiceName, targetSlot, err)
	}
	d.SetId(*resp.ID)
	return resourceAppServiceActiveSlotRead(d, meta)
}

func resourceAppServiceActiveSlotRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service %q (resource group %q) was not found - removing from state", id.SiteName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %q: %+v", id.SiteName, err)
	}

	if resp.SiteProperties == nil || resp.SiteProperties.SlotSwapStatus == nil {
		return fmt.Errorf("App Service Slot %q: SiteProperties or SlotSwapStatus is nil", id.SiteName)
	}
	d.Set("app_service_name", resp.Name)
	d.Set("resource_group_name", resp.ResourceGroup)
	d.Set("app_service_slot_name", resp.SiteProperties.SlotSwapStatus.SourceSlotName)
	return nil
}

func resourceAppServiceActiveSlotDelete(_ *pluginsdk.ResourceData, _ interface{}) error {
	// There is nothing to delete so return nil
	return nil
}
