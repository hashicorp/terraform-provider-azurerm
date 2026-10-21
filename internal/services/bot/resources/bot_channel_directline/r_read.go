// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_directline

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

func resourceBotChannelDirectlineRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName, string(botservice.ChannelNameBasicChannelChannelNameDirectLineChannel))
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Directline Channel for Bot %q (Resource Group %q) was not found - removing from state!", id.ResourceGroupName, id.BotServiceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Channel Directline for Bot %q (Resource Group %q): %+v", id.ResourceGroupName, id.BotServiceName, err)
	}

	channelsResp, err := client.ListWithKeys(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameDirectLineChannel)
	if err != nil {
		return fmt.Errorf("listing Keys for Directline Channel for Bot %q (Resource Group %q): %+v", id.ResourceGroupName, id.BotServiceName, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := channelsResp.Properties; props != nil {
		if channel, ok := props.AsDirectLineChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("site", flattenDirectlineSites(filterSites(channelProps.Sites)))

				if channelProps.ExtensionKey1 != nil {
					d.Set("extension_key_1", channelProps.ExtensionKey1)
				}

				if channelProps.ExtensionKey2 != nil {
					d.Set("extension_key_2", channelProps.ExtensionKey2)
				}
			}
		}
	}

	return nil
}

func flattenDirectlineSites(input []botservice.DirectLineSite) []any {
	sites := make([]any, len(input))

	for i, element := range input {
		site := make(map[string]any)

		if v := element.SiteName; v != nil {
			site["name"] = *v
		}

		userUploadEnabled := true
		if v := element.IsBlockUserUploadEnabled; v != nil {
			userUploadEnabled = !*v
		}
		site["user_upload_enabled"] = userUploadEnabled

		site["endpoint_parameters_enabled"] = pointer.From(element.IsEndpointParametersEnabled)

		storageEnabled := true
		if v := element.IsNoStorageEnabled; v != nil {
			storageEnabled = !*v
		}
		site["storage_enabled"] = storageEnabled

		if element.Key != nil {
			site["key"] = *element.Key
		}

		if element.Key2 != nil {
			site["key2"] = *element.Key2
		}

		if element.IsEnabled != nil {
			site["enabled"] = *element.IsEnabled
		}

		if element.IsV1Enabled != nil {
			site["v1_allowed"] = *element.IsV1Enabled
		}

		if element.IsV3Enabled != nil {
			site["v3_allowed"] = *element.IsV3Enabled
		}

		if element.IsSecureSiteEnabled != nil {
			site["enhanced_authentication_enabled"] = *element.IsSecureSiteEnabled
		}

		if element.TrustedOrigins != nil {
			site["trusted_origins"] = *element.TrustedOrigins
		}

		sites[i] = site
	}

	return sites
}

// When creating a new directline channel, a Default Site is created
// There is a race condition where this site is not removed before the create request is completed
func filterSites(sites *[]botservice.DirectLineSite) []botservice.DirectLineSite {
	filtered := make([]botservice.DirectLineSite, 0)
	for _, site := range *sites {
		if *site.SiteName == "Default Site" {
			continue
		}
		filtered = append(filtered, site)
	}
	return filtered
}
