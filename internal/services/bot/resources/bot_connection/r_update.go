// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceArmBotConnectionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotConnectionID(d.Id())
	if err != nil {
		return err
	}

	connection := botservice.ConnectionSetting{
		Properties: &botservice.ConnectionSettingProperties{
			ServiceProviderDisplayName: pointer.To(d.Get("service_provider_name").(string)),
			ClientID:                   pointer.To(d.Get("client_id").(string)),
			ClientSecret:               pointer.To(d.Get("client_secret").(string)),
			Scopes:                     pointer.To(d.Get("scopes").(string)),
		},
		Kind:     botservice.KindBot,
		Location: pointer.To(d.Get("location").(string)),
	}

	if v, ok := d.GetOk("parameters"); ok {
		connection.Properties.Parameters = expandBotConnectionParameters(v.(map[string]any))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, id.ConnectionName, connection); err != nil {
		return fmt.Errorf("updating Bot Connection %q (Bot %q / Resource Group %q): %+v", id.ConnectionName, id.BotServiceName, id.ResourceGroup, err)
	}

	return resourceArmBotConnectionRead(d, meta)
}
