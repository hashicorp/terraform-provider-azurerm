// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_ms_teams

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

func resourceBotChannelMsTeamsCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := commonids.NewBotServiceChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameMsTeamsChannel))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, resourceId.ResourceGroupName, resourceId.BotServiceName, resourceId.ChannelType)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for the presence of existing MS Teams Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_bot_channel_ms_teams", resourceId.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.MsTeamsChannel{
			Properties: &botservice.MsTeamsChannelProperties{
				AcceptedTerms:         pointer.To(true),
				DeploymentEnvironment: pointer.To(d.Get("deployment_environment").(string)),
				EnableCalling:         pointer.To(d.Get("calling_enabled").(bool)),
				IsEnabled:             pointer.To(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameMsTeamsChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("calling_web_hook"); ok {
		channel, _ := channel.Properties.AsMsTeamsChannel()
		channel.Properties.CallingWebhook = pointer.To(v.(string))
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroupName, resourceId.BotServiceName, botservice.ChannelNameMsTeamsChannel, channel); err != nil {
		return fmt.Errorf("creating MS Teams Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
	}

	d.SetId(resourceId.ID())
	return resourceBotChannelMsTeamsRead(d, meta)
}
