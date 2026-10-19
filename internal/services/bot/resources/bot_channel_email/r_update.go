// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_email

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/botservice/2022-09-15/channel"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBotChannelEmailUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.EmailChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
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

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating Email Channel for Bot %q (Resource Group %q): %+v", id.ResourceGroupName, id.BotServiceName, err)
	}

	return resourceBotChannelEmailRead(d, meta)
}
