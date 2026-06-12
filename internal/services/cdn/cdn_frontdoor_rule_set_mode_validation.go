// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func ensureNonBatchRuleSetMode(ctx context.Context, client *clients.Client, id rulesets.RuleSetId) error {
	resp, err := client.Cdn.FrontDoorRuleSetsClient_v2025_12_01.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model != nil && resp.Model.Properties != nil && pointer.From(resp.Model.Properties.BatchMode) {
		return fmt.Errorf("creating or updating an individually managed Front Door Rule in %s is not supported when `batch_mode_enabled` is `true` on the parent Rule Set; use `azurerm_cdn_frontdoor_batch_rule_set` instead or create a non-batch Rule Set with `azurerm_cdn_frontdoor_rule_set`", id)
	}

	return nil
}
