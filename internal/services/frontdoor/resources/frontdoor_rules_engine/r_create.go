// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor_rules_engine

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorRulesEngineCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rules := d.Get("rule").([]any)

	id := frontdoors.NewRulesEngineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("frontdoor_name").(string), d.Get("name").(string))

	frontdoorRulesEngine := frontdoors.RulesEngine{
		Name: pointer.To(id.RulesEngineName),
		Properties: &frontdoors.RulesEngineProperties{
			Rules: expandFrontDoorRulesEngineRules(rules),
		},
	}

	if d.IsNewResource() {
		if err := client.RulesEnginesCreateOrUpdateCallbackThenPoll(ctx, id, frontdoorRulesEngine, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.RulesEnginesCreateOrUpdateThenPoll(ctx, id, frontdoorRulesEngine); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceFrontDoorRulesEngineRead(d, meta)
}
