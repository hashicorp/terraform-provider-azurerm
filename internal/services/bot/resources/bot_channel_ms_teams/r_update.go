// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_ms_teams

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

func resourceBotChannelMsTeamsUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.MsTeamsChannel{
			Properties: &botservice.MsTeamsChannelProperties{
				AcceptedTerms:         pointer.To(true),
				DeploymentEnvironment: pointer.To(d.Get("deployment_environment").(string)),
				EnableCalling:         pointer.To(d.Get("calling_enabled").(bool)),
				CallingWebhook:        pointer.To(d.Get("calling_web_hook").(string)),
				IsEnabled:             pointer.To(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameMsTeamsChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameMsTeamsChannel, channel); err != nil {
		return fmt.Errorf("updating MS Teams Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroupName, err)
	}

	return resourceBotChannelMsTeamsRead(d, meta)
}
