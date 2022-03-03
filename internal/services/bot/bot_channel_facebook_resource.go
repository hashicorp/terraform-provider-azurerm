package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2021-05-01-preview/botservice"
	"github.com/hashicorp/go-azure-helpers/lang/response"
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

			"facebook_application_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"facebook_application_secret": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"page": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"access_token": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
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
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_bot_channel_facebook", id.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.FacebookChannel{
			Properties: &botservice.FacebookChannelProperties{
				AppID:     utils.String(d.Get("facebook_application_id").(string)),
				AppSecret: utils.String(d.Get("facebook_application_secret").(string)),
				Pages:     expandFacebookPage(d.Get("page").(*pluginsdk.Set).List()),
				IsEnabled: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameFacebookChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
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

	channelsResp, err := client.ListWithKeys(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameFacebookChannel)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := channelsResp.Properties; props != nil {
		if channel, ok := props.AsFacebookChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("facebook_application_id", channelProps.AppID)
				d.Set("facebook_application_secret", channelProps.AppSecret)

				if err := d.Set("page", flattenFacebookPage(channelProps.Pages)); err != nil {
					return fmt.Errorf("setting `page`: %+v", err)
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
				AppID:     utils.String(d.Get("facebook_application_id").(string)),
				AppSecret: utils.String(d.Get("facebook_application_secret").(string)),
				Pages:     expandFacebookPage(d.Get("page").(*pluginsdk.Set).List()),
				IsEnabled: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameFacebookChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
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

func expandFacebookPage(input []interface{}) *[]botservice.FacebookPage {
	results := make([]botservice.FacebookPage, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := botservice.FacebookPage{
			AccessToken: utils.String(v["access_token"].(string)),
			ID:          utils.String(v["id"].(string)),
		}

		results = append(results, result)
	}

	return &results
}

func flattenFacebookPage(input *[]botservice.FacebookPage) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var id string
		if item.ID != nil {
			id = *item.ID
		}

		var accessToken string
		if item.AccessToken != nil {
			accessToken = *item.AccessToken
		}

		results = append(results, map[string]interface{}{
			"id":           id,
			"access_token": accessToken,
		})
	}

	return results
}
