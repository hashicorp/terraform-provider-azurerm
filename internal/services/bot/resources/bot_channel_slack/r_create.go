// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_slack

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

func resourceBotChannelSlackCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := commonids.NewBotServiceChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameSlackChannel))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, resourceId.ResourceGroupName, resourceId.BotServiceName, resourceId.ChannelType)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Slack Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_bot_channel_slack", resourceId.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.SlackChannel{
			Properties: &botservice.SlackChannelProperties{
				ClientID:                pointer.To(d.Get("client_id").(string)),
				ClientSecret:            pointer.To(d.Get("client_secret").(string)),
				VerificationToken:       pointer.To(d.Get("verification_token").(string)),
				LandingPageURL:          pointer.To(d.Get("landing_page_url").(string)),
				IsEnabled:               pointer.To(true),
				RegisterBeforeOAuthFlow: pointer.To(true),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameSlackChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if v, ok := d.GetOk("signing_secret"); ok {
		channel, _ := channel.Properties.AsSlackChannel()
		channel.Properties.SigningSecret = pointer.To(v.(string))
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroupName, resourceId.BotServiceName, botservice.ChannelNameSlackChannel, channel); err != nil {
		return fmt.Errorf("creating Slack Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
	}

	d.SetId(resourceId.ID())
	return resourceBotChannelSlackRead(d, meta)
}
