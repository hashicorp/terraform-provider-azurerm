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

func resourceBotChannelFacebook() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelFacebookCreate,
		Read:   resourceBotChannelFacebookRead,
		Delete: resourceBotChannelFacebookDelete,
		Update: resourceBotChannelFacebookUpdate,

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

			"app_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"app_secret": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"pages": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"access_token": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceBotChannelFacebookCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameFacebookChannel))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ChannelName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_channel_facebook", id.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.FacebookChannel{
			Properties: &botservice.FacebookChannelProperties{
				IsEnabled: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameFacebookChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("app_id"); ok {
		channel, _ := channel.Properties.AsFacebookChannel()
		channel.Properties.AppID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("app_secret"); ok {
		channel, _ := channel.Properties.AsFacebookChannel()
		channel.Properties.AppSecret = utils.String(v.(string))
	}

	if v, ok := d.GetOk("pages"); ok {
		channel, _ := channel.Properties.AsFacebookChannel()
		channel.Properties.Pages = expandFacebookPages(v.(*pluginsdk.Set).List())
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameFacebookChannel, channel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBotChannelFacebookRead(d, meta)
}

func resourceBotChannelFacebookRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameFacebookChannel))
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
		if channel, ok := props.AsFacebookChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				if err := d.Set("pages", flattenFacebookPages(channelProps.Pages)); err != nil {
					return fmt.Errorf("setting `pages`: %+v", err)
				}
			}
		}
	}

	return nil
}

func resourceBotChannelFacebookUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.FacebookChannel{
			Properties: &botservice.FacebookChannelProperties{
				IsEnabled: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameFacebookChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("app_id"); ok {
		channel, _ := channel.Properties.AsFacebookChannel()
		channel.Properties.AppID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("app_secret"); ok {
		channel, _ := channel.Properties.AsFacebookChannel()
		channel.Properties.AppSecret = utils.String(v.(string))
	}

	if v, ok := d.GetOk("pages"); ok {
		channel, _ := channel.Properties.AsFacebookChannel()
		channel.Properties.Pages = expandFacebookPages(v.(*pluginsdk.Set).List())
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameFacebookChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceBotChannelFacebookRead(d, meta)
}

func resourceBotChannelFacebookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameFacebookChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func expandFacebookPages(input []interface{}) *[]botservice.FacebookPage {
	results := make([]botservice.FacebookPage, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		result := botservice.FacebookPage{
			ID: utils.String(v["id"].(string)),
		}

		if accessToken := v["access_token"].(string); accessToken != "" {
			result.AccessToken = utils.String(accessToken)
		}

		results = append(results, result)
	}
	return &results
}

func flattenFacebookPages(input *[]botservice.FacebookPage) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var accessToken string
		if item.AccessToken != nil {
			accessToken = *item.AccessToken
		}

		var id string
		if item.ID != nil {
			id = *item.ID
		}

		results = append(results, map[string]interface{}{
			"access_token": accessToken,
			"id":           id,
		})
	}
	return results
}
