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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceBotChannelAlexa() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelAlexaCreate,
		Read:   resourceBotChannelAlexaRead,
		Delete: resourceBotChannelAlexaDelete,
		Update: resourceBotChannelAlexaUpdate,

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

			"skill_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceBotChannelAlexaCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameAlexaChannel))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ChannelName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_channel_alexa", id.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.AlexaChannel{
			Properties: &botservice.AlexaChannelProperties{
				AlexaSkillID: utils.String(d.Get("skill_id").(string)),
				IsEnabled:    utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameAlexaChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameAlexaChannel, channel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBotChannelAlexaRead(d, meta)
}

func resourceBotChannelAlexaRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameAlexaChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsAlexaChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("skill_id", channelProps.AlexaSkillID)
			}
		}
	}

	return nil
}

func resourceBotChannelAlexaUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.AlexaChannel{
			Properties: &botservice.AlexaChannelProperties{
				AlexaSkillID: utils.String(d.Get("skill_id").(string)),
				IsEnabled:    utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameAlexaChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	_, err = client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameAlexaChannel, channel)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceBotChannelAlexaRead(d, meta)
}

func resourceBotChannelAlexaDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameAlexaChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
