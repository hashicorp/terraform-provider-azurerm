// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package healthbot

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2025-05-25/healthbots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	botClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthbotServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := botClient.From(meta).HealthBotClient.HealthBots
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := healthbots.ParseHealthBotID(d.Id())
	if err != nil {
		return err
	}

	payload := healthbots.HealthBotUpdateParameters{}
	if d.HasChange("sku_name") {
		payload.Sku = &healthbots.Sku{
			Name: healthbots.SkuName(d.Get("sku_name").(string)),
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.BotsUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	return resourceHealthbotServiceRead(d, meta)
}
