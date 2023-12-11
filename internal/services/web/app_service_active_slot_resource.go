// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

		DeprecationMessage: "The `azurerm_app_service_active_slot` resource has been superseded by the `azurerm_web_app_active_slot` and `azurerm_function_app_active_slot` resources. Whilst this resource will continue to be available in the 2.x and 3.x releases it is feature-frozen for compatibility purposes, will no longer receive any updates and will be removed in a future major release of the Azure Provider.",

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appServiceId := parse.NewAppServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("app_service_name").(string))
	preserveVnet := true
	id := parse.NewAppServiceSlotID(appServiceId.SubscriptionId, appServiceId.ResourceGroup, appServiceId.SiteName, d.Get("app_service_slot_name").(string))

	resp, err := client.Get(ctx, appServiceId.ResourceGroup, appServiceId.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", appServiceId)
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	if _, err = client.GetSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found.", id)
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	cmsSlotEntity := web.CsmSlotEntity{
		TargetSlot:   &id.SlotName,
		PreserveVnet: &preserveVnet,
	}

	future, err := client.SwapSlotWithProduction(ctx, id.ResourceGroup, id.SiteName, cmsSlotEntity)
	if err != nil {
		return fmt.Errorf("swapping %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("swapping %s: %+v", id, err)
	}
	d.SetId(appServiceId.ID())
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
		return fmt.Errorf("making Read request on AzureRM App Service %q: %+v", id.SiteName, err)
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
