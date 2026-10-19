// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_email

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/botservice/2022-09-15/channel"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBotChannelEmailRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.EmailChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Email Channel for Bot %q (Resource Group %q) was not found - removing from state!", id.BotServiceName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Email Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroupName, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if channel, ok := props.(channel.EmailChannel); ok {
				if channelProps := channel.Properties; channelProps != nil {
					d.Set("email_address", channelProps.EmailAddress)
				}
			}
		}
	}

	return nil
}
