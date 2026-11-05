// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_directline

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

func resourceBotChannelDirectlineCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := commonids.NewBotServiceChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameBasicChannelChannelNameDirectLineChannel))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, resourceId.ResourceGroupName, resourceId.BotServiceName, resourceId.ChannelType)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Directline Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			// a "Default Site" site gets created and returned.. so let's check it's not just that
			if props := existing.Properties; props != nil {
				directLineChannel, ok := props.AsDirectLineChannel()
				if ok && directLineChannel.Properties != nil {
					sites := filterSites(directLineChannel.Properties.Sites)
					if len(sites) != 0 {
						return tf.ImportAsExistsError("azurerm_bot_channel_directline", resourceId.ID())
					}
				}
			}
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.DirectLineChannel{
			Properties: &botservice.DirectLineChannelProperties{
				Sites: expandDirectlineSites(d.Get("site").(*pluginsdk.Set).List()),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameDirectLineChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroupName, resourceId.BotServiceName, botservice.ChannelNameDirectLineChannel, channel); err != nil {
		return fmt.Errorf("creating Directline Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
	}
	d.SetId(resourceId.ID())

	// Unable to add a new site with enhanced_authentication_enabled, user_upload_enabled, endpoint_parameters_enabled, storage_enabled in the same operation, so we need to make two calls
	if _, err := client.Update(ctx, resourceId.ResourceGroupName, resourceId.BotServiceName, botservice.ChannelNameDirectLineChannel, channel); err != nil {
		return fmt.Errorf("updating Directline Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
	}

	return resourceBotChannelDirectlineRead(d, meta)
}

func expandDirectlineSites(input []any) *[]botservice.DirectLineSite {
	sites := make([]botservice.DirectLineSite, 0)

	for _, element := range input {
		if element == nil {
			continue
		}

		site := element.(map[string]any)
		expanded := botservice.DirectLineSite{
			IsBlockUserUploadEnabled:    pointer.To(!site["user_upload_enabled"].(bool)),
			IsEndpointParametersEnabled: pointer.To(site["endpoint_parameters_enabled"].(bool)),
			IsNoStorageEnabled:          pointer.To(!site["storage_enabled"].(bool)),
		}

		if v, ok := site["name"].(string); ok {
			expanded.SiteName = &v
		}
		if v, ok := site["enabled"].(bool); ok {
			expanded.IsEnabled = &v
		}
		if v, ok := site["v1_allowed"].(bool); ok {
			expanded.IsV1Enabled = &v
		}
		if v, ok := site["v3_allowed"].(bool); ok {
			expanded.IsV3Enabled = &v
		}
		if v, ok := site["enhanced_authentication_enabled"].(bool); ok {
			expanded.IsSecureSiteEnabled = &v
		}
		if v, ok := site["trusted_origins"].(*pluginsdk.Set); ok {
			origins := v.List()
			items := make([]string, len(origins))
			for i, raw := range origins {
				items[i] = raw.(string)
			}
			expanded.TrustedOrigins = &items
		}

		sites = append(sites, expanded)
	}

	return &sites
}
