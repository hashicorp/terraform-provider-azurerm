// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_web_chat

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelWebChatRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName, string(botservice.ChannelNameWebChatChannel))
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
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

func flattenSites(input *[]botservice.WebChatSite) []any {
	results := make([]any, 0)

	for _, item := range *input {
		result := make(map[string]any)

		result["name"] = pointer.From(item.SiteName)

		userUploadEnabled := true
		if v := item.IsBlockUserUploadEnabled; v != nil {
			userUploadEnabled = !*v
		}
		result["user_upload_enabled"] = userUploadEnabled

		result["endpoint_parameters_enabled"] = pointer.From(item.IsEndpointParametersEnabled)

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
