// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_direct_line_speech

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelDirectLineSpeechCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewBotServiceChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameDirectLineSpeechChannel))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName, id.ChannelType)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_bot_channel_direct_line_speech", id.ID())
		}
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

	// Set cognitive service properties if provided
	if v, ok := d.GetOk("cognitive_service_location"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CognitiveServiceRegion = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("cognitive_service_access_key"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CognitiveServiceSubscriptionKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("cognitive_account_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CognitiveServiceResourceID = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("custom_speech_model_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CustomSpeechModelID = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("custom_voice_deployment_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CustomVoiceDeploymentID = pointer.To(v.(string))
	}

	if _, err := client.Create(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameDirectLineSpeechChannel, channel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBotChannelDirectLineSpeechRead(d, meta)
}
