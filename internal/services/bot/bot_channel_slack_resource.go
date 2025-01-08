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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelSlack() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelSlackCreate,
		Read:   resourceBotChannelSlackRead,
		Delete: resourceBotChannelSlackDelete,
		Update: resourceBotChannelSlackUpdate,

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_secret": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"verification_token": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"landing_page_url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"signing_secret": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceBotChannelSlackCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameSlackChannel))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ChannelName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Slack Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_bot_channel_slack", resourceId.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.SlackChannel{
			Properties: &botservice.SlackChannelProperties{
				ClientID:                utils.String(d.Get("client_id").(string)),
				ClientSecret:            utils.String(d.Get("client_secret").(string)),
				VerificationToken:       utils.String(d.Get("verification_token").(string)),
				LandingPageURL:          utils.String(d.Get("landing_page_url").(string)),
				IsEnabled:               utils.Bool(true),
				RegisterBeforeOAuthFlow: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameSlackChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("signing_secret"); ok {
		channel, _ := channel.Properties.AsSlackChannel()
		channel.Properties.SigningSecret = utils.String(v.(string))
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, botservice.ChannelNameSlackChannel, channel); err != nil {
		return fmt.Errorf("creating Slack Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceBotChannelSlackRead(d, meta)
}

func resourceBotChannelSlackRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSlackChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Slack Channel for Bot %q (Resource Group %q) was not found - removing from state!", id.BotServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retreving Slack Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsSlackChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("client_id", d.Get("client_id"))
				d.Set("signing_secret", d.Get("signing_secret"))
			}
		}
	}

	return nil
}

func resourceBotChannelSlackUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.SlackChannel{
			Properties: &botservice.SlackChannelProperties{
				ClientID:                utils.String(d.Get("client_id").(string)),
				ClientSecret:            utils.String(d.Get("client_secret").(string)),
				VerificationToken:       utils.String(d.Get("verification_token").(string)),
				LandingPageURL:          utils.String(d.Get("landing_page_url").(string)),
				IsEnabled:               utils.Bool(true),
				RegisterBeforeOAuthFlow: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameSlackChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("signing_secret"); ok {
		channel, _ := channel.Properties.AsSlackChannel()
		channel.Properties.SigningSecret = utils.String(v.(string))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameSlackChannel, channel); err != nil {
		return fmt.Errorf("updating Slack Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
	}

	return resourceBotChannelSlackRead(d, meta)
}

func resourceBotChannelSlackDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSlackChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Slack Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}
