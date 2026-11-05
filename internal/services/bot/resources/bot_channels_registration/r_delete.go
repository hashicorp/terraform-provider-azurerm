// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channels_registration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBotChannelsRegistrationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseBotServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroupName, id.BotServiceName)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Bot Channels Registration %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroupName, err)
		}
	}

	return nil
}
