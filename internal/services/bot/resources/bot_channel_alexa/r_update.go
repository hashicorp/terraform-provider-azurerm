// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_alexa

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelAlexaUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.AlexaChannel{
			Properties: &botservice.AlexaChannelProperties{
				AlexaSkillID: pointer.To(d.Get("skill_id").(string)),
				IsEnabled:    pointer.To(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameAlexaChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err = client.Update(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameAlexaChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/15298
	// There is a long running issue on updating the Alexa Channel.
	// The Bot Channel API cannot update Skill ID immediately after the Alexa Channel is updated.
	// It has to wait a while after updated. Then the skill ID can be updated successfully.
	stateConf := &pluginsdk.StateChangeConf{
		Delay:      5 * time.Minute,
		Pending:    []string{"204", "404"},
		Target:     []string{"200", "202"},
		Refresh:    botChannelAlexaStateRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceBotChannelAlexaRead(d, meta)
}
