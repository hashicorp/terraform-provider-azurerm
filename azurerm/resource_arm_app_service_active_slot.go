package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceActiveSlot() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceActiveSlotCreateUpdate,
		Read:   resourceArmAppServiceActiveSlotRead,
		Update: resourceArmAppServiceActiveSlotCreateUpdate,
		Delete: resourceArmAppServiceActiveSlotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"resource_group_name": azure.SchemaResourceGroupName(),

			"app_service_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"app_service_slot_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmAppServiceActiveSlotCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

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

	_, err = client.Get(ctx, resGroup, targetSlot)
	if err != nil {
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
	return resourceArmAppServiceActiveSlotRead(d, meta)
}

func resourceArmAppServiceActiveSlotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["sites"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service %q (resource group %q) was not found - removing from state", name, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM App Service %q: %+v", name, err)
	}

	d.Set("app_service_name", resp.Name)
	d.Set("resource_group_name", resp.ResourceGroup)
	d.Set("app_service_slot_name", resp.SiteProperties.SlotSwapStatus.SourceSlotName)
	return nil
}

func resourceArmAppServiceActiveSlotDelete(_ *schema.ResourceData, _ interface{}) error {
	// There is nothing to delete so return nil
	return nil
}
