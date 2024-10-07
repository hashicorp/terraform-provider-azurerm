// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"fmt"
	"log"
	"time"

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

func resourceBotChannelWebChat() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
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
			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"bot_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.BotName,
			},

			"site": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"user_upload_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"endpoint_parameters_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"storage_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		},
	}

	return resource
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
			Properties:  &botservice.WebChatChannelProperties{},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameWebChatChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("site"); ok {
		channel, _ := channel.Properties.AsWebChatChannel()
		channel.Properties.Sites = expandSites(v.(*pluginsdk.Set).List())
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Unable to add a new site with user_upload_enabled, endpoint_parameters_enabled, storage_enabled in the same operation, so we need to make two calls
	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
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
				if err := d.Set("site", flattenSites(channelProps.Sites)); err != nil {
					return fmt.Errorf("setting `site`: %+v", err)
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
			Properties:  &botservice.WebChatChannelProperties{},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameWebChatChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if d.HasChange("site") {
		channel, _ := channel.Properties.AsWebChatChannel()
		channel.Properties.Sites = expandSites(d.Get("site").(*pluginsdk.Set).List())
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// Unable to add a new site with user_upload_enabled, endpoint_parameters_enabled, storage_enabled in the same operation, so we need to make two calls
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

	existing, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameWebChatChannel))
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
		Location: utils.String(azure.NormalizeLocation(*existing.Location)),
		Kind:     botservice.KindBot,
	}

	// The Bot WebChat Channel would be created by default while creating Bot Registrations Channel.
	// So it has to restore the default Web Chat Channel while deleting
	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("restoring the default Web Chat Channel %s: %+v", id, err)
	}

	return nil
}

func expandSites(input []interface{}) *[]botservice.WebChatSite {
	results := make([]botservice.WebChatSite, 0)

	for _, item := range input {
		site := item.(map[string]interface{})
		result := botservice.WebChatSite{
			IsEnabled:                   utils.Bool(true),
			IsBlockUserUploadEnabled:    utils.Bool(!site["user_upload_enabled"].(bool)),
			IsEndpointParametersEnabled: utils.Bool(site["endpoint_parameters_enabled"].(bool)),
			IsNoStorageEnabled:          utils.Bool(!site["storage_enabled"].(bool)),
		}

		if siteName := site["name"].(string); siteName != "" {
			result.SiteName = utils.String(siteName)
		}

		results = append(results, result)
	}

	return &results
}

func flattenSites(input *[]botservice.WebChatSite) []interface{} {
	results := make([]interface{}, 0)

	for _, item := range *input {
		result := make(map[string]interface{})

		var name string
		if v := item.SiteName; v != nil {
			name = *v
		}
		result["name"] = name

		userUploadEnabled := true
		if v := item.IsBlockUserUploadEnabled; v != nil {
			userUploadEnabled = !*v
		}
		result["user_upload_enabled"] = userUploadEnabled

		var endpointParametersEnabled bool
		if v := item.IsEndpointParametersEnabled; v != nil {
			endpointParametersEnabled = *v
		}
		result["endpoint_parameters_enabled"] = endpointParametersEnabled

		storageEnabled := true
		if v := item.IsNoStorageEnabled; v != nil {
			storageEnabled = !*v
		}
		result["storage_enabled"] = storageEnabled

		results = append(results, result)
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
