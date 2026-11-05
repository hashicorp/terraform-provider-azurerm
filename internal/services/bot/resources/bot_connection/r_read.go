// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_connection

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceArmBotConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ConnectionName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Bot Connection %q (Bot %q / Resource Group %q) was not found - removing from state!", id.ConnectionName, id.BotServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Bot Connection %q (Bot %q / Resource Group %q): %+v", id.ConnectionName, id.BotServiceName, id.ResourceGroup, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		d.Set("client_id", props.ClientID)
		d.Set("scopes", props.Scopes)
		if err := d.Set("parameters", flattenBotConnectionParameters(props.Parameters)); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}
	}

	return nil
}

func flattenBotConnectionParameters(input *[]botservice.ConnectionSettingParameter) map[string]any {
	output := make(map[string]any)
	if input == nil {
		return output
	}

	for _, parameter := range *input {
		if key := parameter.Key; key != nil {
			// We disregard the clientSecret and clientId as one is sensitive and the other is returned in the ClientId attribute.
			if *key != "clientSecret" && *key != "clientId" && *key != "scopes" {
				if value := parameter.Value; value != nil {
					output[*key] = *value
				}
			}
		}
	}

	return output
}
