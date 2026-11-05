// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_slack

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

func resourceBotChannelSlackUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.SlackChannel{
			Properties: &botservice.SlackChannelProperties{
				ClientID:                pointer.To(d.Get("client_id").(string)),
				ClientSecret:            pointer.To(d.Get("client_secret").(string)),
				VerificationToken:       pointer.To(d.Get("verification_token").(string)),
				LandingPageURL:          pointer.To(d.Get("landing_page_url").(string)),
				IsEnabled:               pointer.To(true),
				RegisterBeforeOAuthFlow: pointer.To(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameSlackChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("signing_secret"); ok {
		channel, _ := channel.Properties.AsSlackChannel()
		channel.Properties.SigningSecret = pointer.To(v.(string))
	}

	if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameSlackChannel, channel); err != nil {
		return fmt.Errorf("updating Slack Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroupName, err)
	}

	return resourceBotChannelSlackRead(d, meta)
}
