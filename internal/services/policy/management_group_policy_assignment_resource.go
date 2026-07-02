// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	managementGroupValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagementGroupAssignmentResource struct {
	base assignmentBaseResource
}

var _ sdk.ResourceWithUpdate = ManagementGroupAssignmentResource{}

func (r ManagementGroupAssignmentResource) ResourceType() string {
	return "azurerm_management_group_policy_assignment"
}

func (r ManagementGroupAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagementGroupAssignmentID
}

type ManagementGroupAssignmentModel struct {
	Name                 string                                `tfschema:"name"`
	ManagementGroupId    string                                `tfschema:"management_group_id"`
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

func (r ManagementGroupAssignmentResource) ModelObject() interface{} {
	return &ManagementGroupAssignmentModel{}
}

func (r ManagementGroupAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"management_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managementGroupValidate.ManagementGroupID,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotWhiteSpace,
				validation.StringDoesNotContainAny("#<>%&:\\?/"),
				validation.StringLenBetween(1, 24),
				validation.StringMatch(regexp.MustCompile("[^ .]$"), "The name cannot end with a period or space."),
			),
		},
	}
	return r.base.arguments(schema)
}

func (r ManagementGroupAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ManagementGroupAssignmentResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "management_group_id")
}

func (r ManagementGroupAssignmentResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("management_group_id")
}

func (r ManagementGroupAssignmentResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}

func (r ManagementGroupAssignmentResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}
