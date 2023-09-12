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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelMsTeams() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelMsTeamsCreate,
		Read:   resourceBotChannelMsTeamsRead,
		Delete: resourceBotChannelMsTeamsDelete,
		Update: resourceBotChannelMsTeamsUpdate,

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

			// issue: https://github.com/Azure/azure-rest-api-specs/issues/9809
			// this field could not update to empty, so add `Computed: true` to avoid diff
			"calling_web_hook": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.BotMSTeamsCallingWebHook(),
			},

			"deployment_environment": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "CommercialDeployment",
				ValidateFunc: validation.StringInSlice([]string{
					"CommercialDeployment",
					"GCCModerateDeployment",
				}, false),
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_calling": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceBotChannelMsTeamsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameMsTeamsChannel))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ChannelName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for the presence of existing MS Teams Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_bot_channel_ms_teams", resourceId.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.MsTeamsChannel{
			Properties: &botservice.MsTeamsChannelProperties{
				AcceptedTerms:         utils.Bool(true),
				DeploymentEnvironment: utils.String(d.Get("deployment_environment").(string)),
				EnableCalling:         utils.Bool(d.Get("enable_calling").(bool)),
				IsEnabled:             utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameMsTeamsChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("calling_web_hook"); ok {
		channel, _ := channel.Properties.AsMsTeamsChannel()
		channel.Properties.CallingWebhook = utils.String(v.(string))
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, botservice.ChannelNameMsTeamsChannel, channel); err != nil {
		return fmt.Errorf("creating MS Teams Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceBotChannelMsTeamsRead(d, meta)
}

func resourceBotChannelMsTeamsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameMsTeamsChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] MS Teams Channel for Bot %q (Resource Group %q) was not found - removing from state!", id.BotServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving MS Teams Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsMsTeamsChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("calling_web_hook", channelProps.CallingWebhook)
				d.Set("deployment_environment", channelProps.DeploymentEnvironment)
				d.Set("enable_calling", channelProps.EnableCalling)
			}
		}
	}

	return nil
}

func resourceBotChannelMsTeamsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.MsTeamsChannel{
			Properties: &botservice.MsTeamsChannelProperties{
				AcceptedTerms:         utils.Bool(true),
				DeploymentEnvironment: utils.String(d.Get("deployment_environment").(string)),
				EnableCalling:         utils.Bool(d.Get("enable_calling").(bool)),
				CallingWebhook:        utils.String(d.Get("calling_web_hook").(string)),
				IsEnabled:             utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameMsTeamsChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameMsTeamsChannel, channel); err != nil {
		return fmt.Errorf("updating MS Teams Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
	}

	return resourceBotChannelMsTeamsRead(d, meta)
}

func resourceBotChannelMsTeamsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameMsTeamsChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting MS Teams Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}
