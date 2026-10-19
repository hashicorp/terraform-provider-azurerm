// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_sms

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelSMSUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.SmsChannel{
			Properties: &botservice.SmsChannelProperties{
				AccountSID:  pointer.To(d.Get("sms_channel_account_security_id").(string)),
				AuthToken:   pointer.To(d.Get("sms_channel_auth_token").(string)),
				IsValidated: pointer.To(true),
				IsEnabled:   pointer.To(true),
				Phone:       pointer.To(d.Get("phone_number").(string)),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameSmsChannel,
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, id.ResourceGroupName, id.BotServiceName, botservice.ChannelNameSmsChannel, channel); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceBotChannelSMSRead(d, meta)
}
