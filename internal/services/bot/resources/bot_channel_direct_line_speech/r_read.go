// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_direct_line_speech

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelDirectLineSpeechRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName, string(botservice.ChannelNameDirectLineSpeechChannel))
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

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsDirectLineSpeechChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("cognitive_account_id", channelProps.CognitiveServiceResourceID)

				// cognitive_service_access_key is currently not returned by the API.
				d.Set("cognitive_service_access_key", d.Get("cognitive_service_access_key").(string))
				if channelProps.CognitiveServiceSubscriptionKey != nil {
					d.Set("cognitive_service_access_key", channelProps.CognitiveServiceSubscriptionKey)
				}

				// cognitive_service_location is currently not returned by the API.
				d.Set("cognitive_service_location", d.Get("cognitive_service_location").(string))
				if channelProps.CognitiveServiceRegion != nil {
					d.Set("cognitive_service_location", channelProps.CognitiveServiceRegion)
				}

				d.Set("custom_speech_model_id", channelProps.CustomSpeechModelID)
				d.Set("custom_voice_deployment_id", channelProps.CustomVoiceDeploymentID)
			}
		}
	}

	return nil
}
