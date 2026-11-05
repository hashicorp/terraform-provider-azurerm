// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_batch_rule_set

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *rulesets.RuleSetId, model *rulesets.RuleSet) error {
	state := CdnFrontDoorBatchRuleSetModel{
		ID:                    id.ID(),
		Name:                  id.RuleSetName,
		CdnFrontdoorProfileID: profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID(),
	}

	if model != nil && model.Properties != nil {
		rulesetsState, err := flattenCdnFrontDoorBatchRuleSetRules(model.Properties.Rules)
		if err != nil {
			return fmt.Errorf("flattening `rules`: %+v", err)
		}
		state.Rule = rulesetsState
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
