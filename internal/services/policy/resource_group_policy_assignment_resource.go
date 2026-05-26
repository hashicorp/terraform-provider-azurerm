// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	resourceValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ResourceGroupAssignmentResource struct {
	base assignmentBaseResource
}

var _ sdk.ResourceWithUpdate = ResourceGroupAssignmentResource{}

func (r ResourceGroupAssignmentResource) ResourceType() string {
	return "azurerm_resource_group_policy_assignment"
}

func (r ResourceGroupAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceGroupAssignmentID
}

type ResourceGroupAssignmentModel struct {
	Name                 string                                `tfschema:"name"`
	ResourceGroupId      string                                `tfschema:"resource_group_id"`
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

func (r ResourceGroupAssignmentResource) ModelObject() interface{} {
	return &ResourceGroupAssignmentModel{}
}

func (r ResourceGroupAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
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
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceValidate.ResourceGroupID,
		},
	}
	return r.base.arguments(schema)
}

func (r ResourceGroupAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ResourceGroupAssignmentResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "resource_group_id")
}

func (r ResourceGroupAssignmentResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("resource_group_id")
}

func (r ResourceGroupAssignmentResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}

func (r ResourceGroupAssignmentResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}
