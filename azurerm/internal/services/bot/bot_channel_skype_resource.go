package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/botservice/mgmt/2021-03-01/botservice"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceBotChannelSkype() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelSkypeCreate,
		Read:   resourceBotChannelSkypeRead,
		Delete: resourceBotChannelSkypeDelete,
		Update: resourceBotChannelSkypeUpdate,

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

			// issue: https://github.com/Azure/azure-rest-api-specs/issues/15170
			// this field could not update to empty, so add `Computed: true` to avoid diff
			"calling_web_hook": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.BotCallingWebHook(),
			},

			"enable_calling": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"enable_groups": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"enable_media_cards": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"enable_messaging": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceBotChannelSkypeCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameSkypeChannel))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ChannelName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			// As Bot Skype Channel would be created by default while creating Bot Registrations Channel
			// So it has to delete default one
			resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSkypeChannel))
			if err != nil {
				if !response.WasNotFound(resp.Response) {
					return fmt.Errorf("deleting default Bot Skype Channel %s: %+v", id, err)
				}
			}
		}
	}

	parameters := botservice.BotChannel{
		Properties: botservice.SkypeChannel{
			Properties: &botservice.SkypeChannelProperties{
				EnableCalling:    utils.Bool(d.Get("enable_calling").(bool)),
				EnableGroups:     utils.Bool(d.Get("enable_groups").(bool)),
				EnableMediaCards: utils.Bool(d.Get("enable_media_cards").(bool)),
				EnableMessaging:  utils.Bool(d.Get("enable_messaging").(bool)),
				IsEnabled:        utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameSkypeChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("calling_web_hook"); ok {
		channel, _ := parameters.Properties.AsSkypeChannel()
		channel.Properties.CallingWebHook = utils.String(v.(string))
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameSkypeChannel, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBotChannelSkypeRead(d, meta)
}

func resourceBotChannelSkypeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSkypeChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsSkypeChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("calling_web_hook", channelProps.CallingWebHook)
				d.Set("enable_calling", channelProps.EnableCalling)
				d.Set("enable_groups", channelProps.EnableGroups)
				d.Set("enable_media_cards", channelProps.EnableMediaCards)
				d.Set("enable_messaging", channelProps.EnableMessaging)
			}
		}
	}

	return nil
}

func resourceBotChannelSkypeUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	parameters := botservice.BotChannel{
		Properties: botservice.SkypeChannel{
			Properties: &botservice.SkypeChannelProperties{
				EnableCalling:    utils.Bool(d.Get("enable_calling").(bool)),
				EnableGroups:     utils.Bool(d.Get("enable_groups").(bool)),
				EnableMediaCards: utils.Bool(d.Get("enable_media_cards").(bool)),
				EnableMessaging:  utils.Bool(d.Get("enable_messaging").(bool)),
				IsEnabled:        utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameSkypeChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("calling_web_hook"); ok {
		channel, _ := parameters.Properties.AsSkypeChannel()
		channel.Properties.CallingWebHook = utils.String(v.(string))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameSkypeChannel, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceBotChannelSkypeRead(d, meta)
}

func resourceBotChannelSkypeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSkypeChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}
