package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2021-05-01-preview/botservice"
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
)

func resourceBotChannelWebChat() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelWebChatCreate,
		Read:   resourceBotChannelWebChatRead,
		Delete: resourceBotChannelWebChatDelete,
		Update: resourceBotChannelWebChatUpdate,

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

			"site_names": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceBotChannelWebChatCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameWebChatChannel))

	existing, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ChannelName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		// The Bot WebChat Channel would be created by default while creating Bot Registrations Channel.
		// So if the channel includes `Default Site`, it means it's default channel and delete it.
		// So if the channel includes other site, it means it's user custom channel and throws conflict error.
		if props := existing.Properties; props != nil {
			defaultChannel, ok := props.AsWebChatChannel()
			if ok && defaultChannel.Properties != nil {
				if includeDefaultWebChatSite(defaultChannel.Properties.Sites) {
					if _, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameBasicChannelChannelNameWebChatChannel)); err != nil {
						return fmt.Errorf("deleting the default Web Chat Channel %s: %+v", id, err)
					}
				} else {
					return tf.ImportAsExistsError("azurerm_bot_channel_web_chat", id.ID())
				}
			}
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.WebChatChannel{
			Properties: &botservice.WebChatChannelProperties{
				Sites: expandSiteNames(d.Get("site_names").(*pluginsdk.Set).List()),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameWebChatChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBotChannelWebChatRead(d, meta)
}

func resourceBotChannelWebChatRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameWebChatChannel))
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
		if channel, ok := props.AsWebChatChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				if err := d.Set("site_names", flattenSiteNames(channelProps.Sites)); err != nil {
					return fmt.Errorf("setting `site_names`: %+v", err)
				}
			}
		}
	}

	return nil
}

func resourceBotChannelWebChatUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.WebChatChannel{
			Properties: &botservice.WebChatChannelProperties{
				Sites: expandSiteNames(d.Get("site_names").(*pluginsdk.Set).List()),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameWebChatChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceBotChannelWebChatRead(d, meta)
}

func resourceBotChannelWebChatDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.WebChatChannel{
			Properties: &botservice.WebChatChannelProperties{
				Sites: &[]botservice.WebChatSite{
					{
						SiteName:  utils.String("Default Site"),
						IsEnabled: utils.Bool(true),
					},
				},
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameWebChatChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	// The Bot WebChat Channel would be created by default while creating Bot Registrations Channel.
	// So it has to restore the default Web Chat Channel while deleting
	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("restoring the default Web Chat Channel %s: %+v", id, err)
	}

	return nil
}

func expandSiteNames(input []interface{}) *[]botservice.WebChatSite {
	results := make([]botservice.WebChatSite, 0)

	for _, item := range input {
		results = append(results, botservice.WebChatSite{
			SiteName:  utils.String(item.(string)),
			IsEnabled: utils.Bool(true),
		})
	}

	return &results
}

func flattenSiteNames(input *[]botservice.WebChatSite) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var siteName string
		if item.SiteName != nil {
			siteName = *item.SiteName
		}

		results = append(results, siteName)
	}

	return results
}

func includeDefaultWebChatSite(sites *[]botservice.WebChatSite) bool {
	includeDefaultSite := false
	for _, site := range *sites {
		if *site.SiteName == "Default Site" {
			includeDefaultSite = true
		}
	}
	return includeDefaultSite
}
