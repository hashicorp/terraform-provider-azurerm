// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/botservice/2022-09-15/channel"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceBotChannelEmail() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelEmailCreate,
		Read:   resourceBotChannelEmailRead,
		Delete: resourceBotChannelEmailDelete,
		Update: resourceBotChannelEmailUpdate,

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

			"email_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"email_password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"email_password", "magic_code"},
			},

			"magic_code": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"email_password", "magic_code"},
			},
		},
	}
}

func resourceBotChannelEmailCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.EmailChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := commonids.NewBotServiceChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(channel.BotServiceChannelTypeEmailChannel))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Email Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_bot_channel_email", resourceId.ID())
		}
	}

	parameters := channel.BotChannel{
		Properties: channel.EmailChannel{
			Properties: &channel.EmailChannelProperties{
				EmailAddress: d.Get("email_address").(string),
				IsEnabled:    true,
			},
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     pointer.To(channel.KindBot),
	}

	if v, ok := d.GetOk("email_password"); ok {
		channelProps := parameters.Properties.(channel.EmailChannel)
		channelProps.Properties.AuthMethod = pointer.To(channel.EmailChannelAuthMethodZero)
		channelProps.Properties.Password = utils.String(v.(string))
	}

	if v, ok := d.GetOk("magic_code"); ok {
		channelProps := parameters.Properties.(channel.EmailChannel)
		channelProps.Properties.AuthMethod = pointer.To(channel.EmailChannelAuthMethodOne)
		channelProps.Properties.MagicCode = utils.String(v.(string))
	}

	if _, err := client.Create(ctx, resourceId, parameters); err != nil {
		return fmt.Errorf("creating Email Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
	}

	d.SetId(resourceId.ID())
	return resourceBotChannelEmailRead(d, meta)
}

func resourceBotChannelEmailRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.EmailChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Email Channel for Bot %q (Resource Group %q) was not found - removing from state!", id.BotServiceName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Email Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroupName, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if channel, ok := props.(channel.EmailChannel); ok {
				if channelProps := channel.Properties; channelProps != nil {
					d.Set("email_address", channelProps.EmailAddress)
				}
			}
		}
	}

	return nil
}

func resourceBotChannelEmailUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.EmailChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	parameters := channel.BotChannel{
		Properties: channel.EmailChannel{
			Properties: &channel.EmailChannelProperties{
				EmailAddress: d.Get("email_address").(string),
				IsEnabled:    true,
			},
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     pointer.To(channel.KindBot),
	}

	if v, ok := d.GetOk("email_password"); ok {
		channelProps := parameters.Properties.(channel.EmailChannel)
		channelProps.Properties.AuthMethod = pointer.To(channel.EmailChannelAuthMethodZero)
		channelProps.Properties.Password = utils.String(v.(string))
	}

	if v, ok := d.GetOk("magic_code"); ok {
		channelProps := parameters.Properties.(channel.EmailChannel)
		channelProps.Properties.AuthMethod = pointer.To(channel.EmailChannelAuthMethodOne)
		channelProps.Properties.MagicCode = utils.String(v.(string))
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating Email Channel for Bot %q (Resource Group %q): %+v", id.ResourceGroupName, id.BotServiceName, err)
	}

	return resourceBotChannelEmailRead(d, meta)
}

func resourceBotChannelEmailDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.EmailChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Email Channel for Bot %q (Resource Group %q): %+v", id.ResourceGroupName, id.BotServiceName, err)
		}
	}

	return nil
}
