// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_line

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelLineRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName, string(botservice.ChannelNameLineChannel))
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	channelsResp, err := client.ListWithKeys(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameLineChannel)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	if props := channelsResp.Properties; props != nil {
		if channel, ok := props.AsLineChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				if err := d.Set("line_channel", flattenLineChannel(channelProps.LineRegistrations)); err != nil {
					return fmt.Errorf("setting `line_channel`: %+v", err)
				}
			}
		}
	}

	return nil
}

func flattenLineChannel(input *[]botservice.LineRegistration) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"access_token": pointer.From(item.ChannelAccessToken),
			"secret":       pointer.From(item.ChannelSecret),
		})
	}

	return results
}
