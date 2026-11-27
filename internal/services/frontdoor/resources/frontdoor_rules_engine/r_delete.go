// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_rules_engine

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	frontdoorClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorRulesEngineDelete(d *pluginsdk.ResourceData, meta any) error {
	client := frontdoorClient.From(meta).FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := frontdoors.ParseRulesEngineIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if err := client.RulesEnginesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
