// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channels_registration

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
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelsRegistrationRead(d *pluginsdk.ResourceData, meta any) error {
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
			log.Printf("[INFO] Bot Channels Registration %q (Resource Group %q) was not found - removing from state", id.BotServiceName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Bot Channels Registration %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroupName, err)
	}

	d.Set("name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.Properties; props != nil {
		d.Set("cmk_key_vault_url", props.CmekKeyVaultURL)
		d.Set("microsoft_app_id", props.MsaAppID)
		d.Set("microsoft_app_type", string(props.MsaAppType))
		d.Set("microsoft_app_tenant_id", pointer.From(props.MsaAppTenantID))
		d.Set("microsoft_app_user_assigned_identity_id", pointer.From(props.MsaAppMSIResourceID))
		d.Set("endpoint", props.Endpoint)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("developer_app_insights_key", props.DeveloperAppInsightKey)
		d.Set("developer_app_insights_application_id", props.DeveloperAppInsightsApplicationID)
		d.Set("icon_url", props.IconURL)
		d.Set("streaming_endpoint_enabled", props.IsStreamingSupported)

		// `PublicNetworkAccess` is empty string when `public_network_access_enabled` isn't specified. So `public_network_access_enabled` shouldn't be set at this time to avoid diff
		if props.PublicNetworkAccess != "" {
			d.Set("public_network_access_enabled", props.PublicNetworkAccess == botservice.PublicNetworkAccessEnabled)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
