// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription_policy_assignment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
)

type SubscriptionAssignmentModel struct {
	Name                 string                                        `tfschema:"name"`
	SubscriptionId       string                                        `tfschema:"subscription_id"`
	PolicyDefinitionId   string                                        `tfschema:"policy_definition_id"`
	Description          string                                        `tfschema:"description"`
	DisplayName          string                                        `tfschema:"display_name"`
	Location             string                                        `tfschema:"location"`
	Enforce              bool                                          `tfschema:"enforce"`
	Metadata             string                                        `tfschema:"metadata"`
	Parameters           string                                        `tfschema:"parameters"`
	NotScopes            []string                                      `tfschema:"not_scopes"`
	Identity             []identity.ModelSystemAssignedUserAssigned    `tfschema:"identity"`
	NonComplianceMessage []helpers.AssignmentNonComplianceMessageModel `tfschema:"non_compliance_message"`
	Overrides            []helpers.AssignmentOverrideModel             `tfschema:"overrides"`
	ResourceSelectors    []helpers.AssignmentResourceSelectorModel     `tfschema:"resource_selectors"`
}

func (r Resource) ModelObject() any {
	return &SubscriptionAssignmentModel{}
}
