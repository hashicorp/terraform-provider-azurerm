// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_rule

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

type Resource struct{}

func (c Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, rmd sdk.ResourceMetaData) error {
		client := rmd.Client.Cdn.FrontDoorRuleSetsClient

		id, err := rules.ParseRuleID(rmd.ResourceData.Id())
		if err != nil {
			return err
		}
		ruleSetID := rulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)

		resp, err := client.Get(ctx, ruleSetID)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", ruleSetID, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil`", id)
		}

		if resp.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: `properties` was nil`", id)
		}

		if pointer.From(resp.Model.Properties.BatchMode) {
			return fmt.Errorf("the parent ruleset (%s) was provisioned using batch mode, and individual rules for this cannot be managed by this resource, use `azurerm_cdn_frontdoor_batch_rule_set` instead, or create a non-batch Rule Set with `azurerm_cdn_frontdoor_rule_set`", ruleSetID)
		}

		return nil
	}
}

func (c Resource) ResourceType() string {
	return "azurerm_cdn_frontdoor_rule"
}

func (c Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return rules.ValidateRuleID
}
