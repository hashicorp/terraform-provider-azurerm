// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthbot

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2025-05-25/healthbots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthbotServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.HealthBotClient.HealthBots
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := healthbots.ParseHealthBotID(d.Id())
	if err != nil {
		return err
	}

	if err := client.BotsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
