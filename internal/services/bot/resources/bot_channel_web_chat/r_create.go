// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_web_chat

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelWebChatCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewBotServiceChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameWebChatChannel))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName, id.ChannelType)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			// The Bot WebChat Channel would be created by default while creating Bot Registrations Channel.
			// So if the channel includes `Default Site`, it means it's default channel and delete it.
			// So if the channel includes other site, it means it's user custom channel and throws conflict error.
			if props := existing.Properties; props != nil {
				defaultChannel, ok := props.AsWebChatChannel()
				if ok && defaultChannel.Properties != nil {
					if includeDefaultWebChatSite(defaultChannel.Properties.Sites) {
						if _, err := client.Delete(ctx, id.ResourceGroupName, id.BotServiceName, string(botservice.ChannelNameBasicChannelChannelNameWebChatChannel)); err != nil {
							return fmt.Errorf("deleting the default Web Chat Channel %s: %+v", id, err)
						}
					} else {
						return tf.ImportAsExistsError("azurerm_bot_channel_web_chat", id.ID())
					}
				}
			}
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.WebChatChannel{
			Properties:  &botservice.WebChatChannelProperties{},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameWebChatChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("site"); ok {
		channel, _ := channel.Properties.AsWebChatChannel()
		channel.Properties.Sites = expandSites(v.(*pluginsdk.Set).List())
	}

	if _, err := client.Create(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// Unable to add a new site with user_upload_enabled, endpoint_parameters_enabled, storage_enabled in the same operation, so we need to make two calls
	if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameWebChatChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceBotChannelWebChatRead(d, meta)
}

func expandSites(input []any) *[]botservice.WebChatSite {
	results := make([]botservice.WebChatSite, 0)

	for _, item := range input {
		site := item.(map[string]any)
		result := botservice.WebChatSite{
			IsEnabled:                   pointer.To(true),
			IsBlockUserUploadEnabled:    pointer.To(!site["user_upload_enabled"].(bool)),
			IsEndpointParametersEnabled: pointer.To(site["endpoint_parameters_enabled"].(bool)),
			IsNoStorageEnabled:          pointer.To(!site["storage_enabled"].(bool)),
		}

		if siteName := site["name"].(string); siteName != "" {
			result.SiteName = pointer.To(siteName)
		}

		results = append(results, result)
	}

	return &results
}
