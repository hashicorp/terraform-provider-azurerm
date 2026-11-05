// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_web_app

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotWebAppUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceID(d.Id())
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = id.BotServiceName
	}

	bot := botservice.Bot{
		Properties: &botservice.BotProperties{
			DisplayName:                       pointer.To(displayName),
			Endpoint:                          pointer.To(d.Get("endpoint").(string)),
			MsaAppID:                          pointer.To(d.Get("microsoft_app_id").(string)),
			MsaAppType:                        botservice.MsaAppType(d.Get("microsoft_app_type").(string)),
			DeveloperAppInsightKey:            pointer.To(d.Get("developer_app_insights_key").(string)),
			DeveloperAppInsightsAPIKey:        pointer.To(d.Get("developer_app_insights_api_key").(string)),
			DeveloperAppInsightsApplicationID: pointer.To(d.Get("developer_app_insights_application_id").(string)),
			LuisAppIds:                        pluginsdk.ExpandStringSlice(d.Get("luis_app_ids").([]any)),
			LuisKey:                           pointer.To(d.Get("luis_key").(string)),
		},
		Location: pointer.To(d.Get("location").(string)),
		Sku: &botservice.Sku{
			Name: botservice.SkuName(d.Get("sku").(string)),
		},
		Kind: botservice.KindSdk,
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("microsoft_app_tenant_id"); ok {
		bot.Properties.MsaAppTenantID = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("microsoft_app_user_assigned_identity_id"); ok {
		bot.Properties.MsaAppMSIResourceID = pointer.To(v.(string))
	}

	if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, bot); err != nil {
		return fmt.Errorf("updating Web App Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroupName, err)
	}

	return resourceBotWebAppRead(d, meta)
}
