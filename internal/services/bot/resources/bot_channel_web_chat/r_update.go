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

func resourceBotChannelWebChatUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.WebChatChannel{
			Properties:  &botservice.WebChatChannelProperties{},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameWebChatChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if d.HasChange("site") {
		channel, _ := channel.Properties.AsWebChatChannel()
		channel.Properties.Sites = expandSites(d.Get("site").(*pluginsdk.Set).List())
	}

	if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// Unable to add a new site with user_upload_enabled, endpoint_parameters_enabled, storage_enabled in the same operation, so we need to make two calls
	if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceBotChannelWebChatRead(d, meta)
}
