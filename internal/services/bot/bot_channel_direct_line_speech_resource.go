// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelDirectLineSpeech() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelDirectLineSpeechCreate,
		Read:   resourceBotChannelDirectLineSpeechRead,
		Delete: resourceBotChannelDirectLineSpeechDelete,
		Update: resourceBotChannelDirectLineSpeechUpdate,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.BotChannelID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"bot_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.BotName,
			},

			"cognitive_service_access_key": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"cognitive_service_location": commonschema.LocationWithoutForceNew(),

			"cognitive_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
			},

			"custom_speech_model_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"custom_voice_deployment_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceBotChannelDirectLineSpeechCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameDirectLineSpeechChannel))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ChannelName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_bot_channel_direct_line_speech", id.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.DirectLineSpeechChannel{
			Properties: &botservice.DirectLineSpeechChannelProperties{
				CognitiveServiceRegion:          utils.String(d.Get("cognitive_service_location").(string)),
				CognitiveServiceSubscriptionKey: utils.String(d.Get("cognitive_service_access_key").(string)),
				IsDefaultBotForCogSvcAccount:    utils.Bool(false),
				IsEnabled:                       utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameDirectLineSpeechChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("cognitive_account_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CognitiveServiceResourceID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("custom_speech_model_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CustomSpeechModelID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("custom_voice_deployment_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CustomVoiceDeploymentID = utils.String(v.(string))
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameDirectLineSpeechChannel, channel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBotChannelDirectLineSpeechRead(d, meta)
}

func resourceBotChannelDirectLineSpeechRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameDirectLineSpeechChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsDirectLineSpeechChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("cognitive_account_id", channelProps.CognitiveServiceResourceID)
				d.Set("custom_speech_model_id", channelProps.CustomSpeechModelID)
				d.Set("custom_voice_deployment_id", channelProps.CustomVoiceDeploymentID)
			}
		}
	}

	return nil
}

func resourceBotChannelDirectLineSpeechUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.DirectLineSpeechChannel{
			Properties: &botservice.DirectLineSpeechChannelProperties{
				CognitiveServiceRegion:          utils.String(d.Get("cognitive_service_location").(string)),
				CognitiveServiceSubscriptionKey: utils.String(d.Get("cognitive_service_access_key").(string)),
				IsDefaultBotForCogSvcAccount:    utils.Bool(false),
				IsEnabled:                       utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameDirectLineSpeechChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if d.HasChange("cognitive_account_id") {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CognitiveServiceResourceID = utils.String(d.Get("cognitive_account_id").(string))
	}

	if v, ok := d.GetOk("custom_speech_model_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CustomSpeechModelID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("custom_voice_deployment_id"); ok {
		channel, _ := channel.Properties.AsDirectLineSpeechChannel()
		channel.Properties.CustomVoiceDeploymentID = utils.String(v.(string))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameDirectLineSpeechChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceBotChannelDirectLineSpeechRead(d, meta)
}

func resourceBotChannelDirectLineSpeechDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameDirectLineSpeechChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}
