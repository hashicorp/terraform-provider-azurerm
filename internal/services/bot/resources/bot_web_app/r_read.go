// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_web_app

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBotWebAppRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Web App Bot %q was not found in Resource Group %q - removing from state", id.BotServiceName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Web App Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroupName, err)
	}

	d.Set("name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.Properties; props != nil {
		d.Set("microsoft_app_id", props.MsaAppID)
		d.Set("microsoft_app_type", string(props.MsaAppType))
		d.Set("microsoft_app_tenant_id", pointer.From(props.MsaAppTenantID))
		d.Set("microsoft_app_user_assigned_identity_id", pointer.From(props.MsaAppMSIResourceID))
		d.Set("endpoint", props.Endpoint)
		d.Set("display_name", props.DisplayName)
		d.Set("developer_app_insights_key", props.DeveloperAppInsightKey)
		d.Set("developer_app_insights_application_id", props.DeveloperAppInsightsApplicationID)
		d.Set("luis_app_ids", props.LuisAppIds)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
