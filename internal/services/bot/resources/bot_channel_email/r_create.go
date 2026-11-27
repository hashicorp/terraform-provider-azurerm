// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bot_channel_email

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/botservice/2022-09-15/channel"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	botClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBotChannelEmailCreate(d *pluginsdk.ResourceData, meta any) error {
	client := botClient.From(meta).EmailChannelClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	resourceId := commonids.NewBotServiceChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(channel.BotServiceChannelTypeEmailChannel))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, resourceId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Email Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_bot_channel_email", resourceId.ID())
		}
	}

	parameters := channel.BotChannel{
		Properties: channel.EmailChannel{
			Properties: &channel.EmailChannelProperties{
				EmailAddress: d.Get("email_address").(string),
				IsEnabled:    true,
			},
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     pointer.To(channel.KindBot),
	}

	if v, ok := d.GetOk("email_password"); ok {
		channelProps := parameters.Properties.(channel.EmailChannel)
		channelProps.Properties.AuthMethod = pointer.To(channel.EmailChannelAuthMethodZero)
		channelProps.Properties.Password = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("magic_code"); ok {
		channelProps := parameters.Properties.(channel.EmailChannel)
		channelProps.Properties.AuthMethod = pointer.To(channel.EmailChannelAuthMethodOne)
		channelProps.Properties.MagicCode = pointer.To(v.(string))
	}

	if _, err := client.Create(ctx, resourceId, parameters); err != nil {
		return fmt.Errorf("creating Email Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroupName, err)
	}

	d.SetId(resourceId.ID())
	return resourceBotChannelEmailRead(d, meta)
}
