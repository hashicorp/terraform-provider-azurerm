// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SubscriptionAssignmentResource struct {
	base assignmentBaseResource
}

var _ sdk.ResourceWithUpdate = SubscriptionAssignmentResource{}

func (r SubscriptionAssignmentResource) ResourceType() string {
	return "azurerm_subscription_policy_assignment"
}

func (r SubscriptionAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SubscriptionAssignmentID
}

type SubscriptionAssignmentModel struct {
	Name                 string                                `tfschema:"name"`
	SubscriptionId       string                                `tfschema:"subscription_id"`
	PolicyDefinitionId   string                                `tfschema:"policy_definition_id"`
	Description          string                                `tfschema:"description"`
	DisplayName          string                                `tfschema:"display_name"`
	Location             string                                `tfschema:"location"`
	Enforce              bool                                  `tfschema:"enforce"`
	Metadata             string                                `tfschema:"metadata"`
	Parameters           string                                `tfschema:"parameters"`
	NotScopes            []string                              `tfschema:"not_scopes"`
	Identity             []assignmentIdentityModel             `tfschema:"identity"`
	NonComplianceMessage []assignmentNonComplianceMessageModel `tfschema:"non_compliance_message"`
	Overrides            []assignmentOverrideModel             `tfschema:"overrides"`
	ResourceSelectors    []assignmentResourceSelectorModel     `tfschema:"resource_selectors"`
}

func (r SubscriptionAssignmentResource) ModelObject() interface{} {
	return &SubscriptionAssignmentModel{}
}

func (r SubscriptionAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotWhiteSpace,
				validation.StringDoesNotContainAny("#<>%&:\\?/"),
				validation.StringLenBetween(1, 64),
				validation.StringMatch(regexp.MustCompile("[^ .]$"), "The name cannot end with a period or space."),
			),
		},
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},
	}
	return r.base.arguments(schema)
}

func (r SubscriptionAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r SubscriptionAssignmentResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "subscription_id")
}

func (r SubscriptionAssignmentResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("subscription_id")
}

func (r SubscriptionAssignmentResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}

func (r SubscriptionAssignmentResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}
