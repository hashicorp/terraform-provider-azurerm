// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_set_definition

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
)

type PolicySetDefinitionResourceModel struct {
	Name                      string                                   `tfschema:"name"`
	PolicyType                string                                   `tfschema:"policy_type"`
	DisplayName               string                                   `tfschema:"display_name"`
	Description               string                                   `tfschema:"description"`
	Metadata                  string                                   `tfschema:"metadata"`
	Parameters                string                                   `tfschema:"parameters"`
	PolicyDefinitionReference []helpers.PolicyDefinitionReferenceModel `tfschema:"policy_definition_reference"`
	PolicyDefinitionGroup     []helpers.PolicyDefinitionGroupModel     `tfschema:"policy_definition_group"`
}

func (r Resource) ModelObject() any {
	return &PolicySetDefinitionResourceModel{}
}
