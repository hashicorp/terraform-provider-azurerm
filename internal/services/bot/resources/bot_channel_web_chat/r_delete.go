// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_web_chat

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelWebChatDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName, string(botservice.ChannelNameWebChatChannel))
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.WebChatChannel{
			Properties: &botservice.WebChatChannelProperties{
				Sites: &[]botservice.WebChatSite{
					{
						SiteName:  pointer.To("Default Site"),
						IsEnabled: pointer.To(true),
					},
				},
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameWebChatChannel,
		},
		Location: pointer.To(location.Normalize(*existing.Location)),
		Kind:     botservice.KindBot,
	}

	// The Bot WebChat Channel would be created by default while creating Bot Registrations Channel.
	// So it has to restore the default Web Chat Channel while deleting
	if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("restoring the default Web Chat Channel %s: %+v", id, err)
	}

	return nil
}
