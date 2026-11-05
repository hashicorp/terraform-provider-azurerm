// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_line

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

func resourceBotChannelLineCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewBotServiceChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameLineChannel))
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName, id.ChannelType)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_bot_channel_line", id.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.LineChannel{
			Properties: &botservice.LineChannelProperties{
				LineRegistrations: expandLineChannel(d.Get("line_channel").(*pluginsdk.Set).List()),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameLineChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameLineChannel, channel); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBotChannelLineRead(d, meta)
}

func expandLineChannel(input []any) *[]botservice.LineRegistration {
	results := make([]botservice.LineRegistration, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, botservice.LineRegistration{
			ChannelSecret:      pointer.To(v["secret"].(string)),
			ChannelAccessToken: pointer.To(v["access_token"].(string)),
		})
	}

	return &results
}
