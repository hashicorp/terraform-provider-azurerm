// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_direct_line_speech

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

func resourceBotChannelDirectLineSpeechUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.DirectLineSpeechChannel{
			Properties: &botservice.DirectLineSpeechChannelProperties{
				IsDefaultBotForCogSvcAccount: pointer.To(false),
				IsEnabled:                    pointer.To(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameDirectLineSpeechChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	directLineSpeechChannel, _ := channel.Properties.AsDirectLineSpeechChannel()
	directLineSpeechChannel.Properties.CognitiveServiceRegion = pointer.To(d.Get("cognitive_service_location").(string))
	directLineSpeechChannel.Properties.CognitiveServiceSubscriptionKey = pointer.To(d.Get("cognitive_service_access_key").(string))

	if v, ok := d.GetOk("cognitive_account_id"); ok {
		directLineSpeechChannel.Properties.CognitiveServiceResourceID = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("custom_speech_model_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CustomSpeechModelID = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("custom_voice_deployment_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CustomVoiceDeploymentID = pointer.To(v.(string))
	}

	if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameDirectLineSpeechChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceBotChannelDirectLineSpeechRead(d, meta)
}
