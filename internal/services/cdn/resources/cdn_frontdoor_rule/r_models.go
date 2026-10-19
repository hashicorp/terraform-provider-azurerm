// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_rule

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
)

type CdnFrontDoorRuleResourceModel struct {
	Name                    string                                    `tfschema:"name"`
	CdnFrontDoorRuleSetID   string                                    `tfschema:"cdn_frontdoor_rule_set_id"`
	BehaviourOnMatch        string                                    `tfschema:"behaviour_on_match"`
	Order                   int64                                     `tfschema:"order"`
	Actions                 []helpers.CdnFrontDoorRuleActionsModel    `tfschema:"actions"`
	Conditions              []helpers.CdnFrontDoorRuleConditionsModel `tfschema:"conditions"`
	CdnFrontDoorRuleSetName string                                    `tfschema:"cdn_frontdoor_rule_set_name"`
}

func (c Resource) ModelObject() any {
	return &CdnFrontDoorRuleResourceModel{}
}
