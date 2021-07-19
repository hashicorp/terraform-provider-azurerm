package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/botservice/mgmt/2021-03-01/botservice"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceBotChannelSMS() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelSMSCreate,
		Read:   resourceBotChannelSMSRead,
		Delete: resourceBotChannelSMSDelete,
		Update: resourceBotChannelSMSUpdate,

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
			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"bot_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.BotName,
			},

			"account_sid": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"auth_token": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"phone": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"is_validated": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceBotChannelSMSCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameSmsChannel))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ChannelName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_channel_sms", id.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.SmsChannel{
			Properties: &botservice.SmsChannelProperties{
				AccountSID:  utils.String(d.Get("account_sid").(string)),
				AuthToken:   utils.String(d.Get("auth_token").(string)),
				IsValidated: utils.Bool(d.Get("is_validated").(bool)),
				IsEnabled:   utils.Bool(true),
				Phone:       utils.String(d.Get("phone").(string)),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameSmsChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameSmsChannel, channel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBotChannelSMSRead(d, meta)
}

func resourceBotChannelSMSRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSmsChannel))
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
		if channel, ok := props.AsSmsChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("account_sid", channelProps.AccountSID)
				d.Set("auth_token", channelProps.AuthToken)
				d.Set("phone", channelProps.Phone)
				d.Set("is_validated", channelProps.IsValidated)
			}
		}
	}

	return nil
}

func resourceBotChannelSMSUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.SmsChannel{
			Properties: &botservice.SmsChannelProperties{
				AccountSID:  utils.String(d.Get("account_sid").(string)),
				AuthToken:   utils.String(d.Get("auth_token").(string)),
				IsValidated: utils.Bool(d.Get("is_validated").(bool)),
				IsEnabled:   utils.Bool(true),
				Phone:       utils.String(d.Get("phone").(string)),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameSmsChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameSmsChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceBotChannelSMSRead(d, meta)
}

func resourceBotChannelSMSDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSmsChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}
