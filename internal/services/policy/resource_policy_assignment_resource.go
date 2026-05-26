// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ResourceAssignmentResource struct {
	base assignmentBaseResource
}

var _ sdk.ResourceWithUpdate = ResourceAssignmentResource{}

func (r ResourceAssignmentResource) ResourceType() string {
	return "azurerm_resource_policy_assignment"
}

func (r ResourceAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceAssignmentId()
}

type ResourceAssignmentModel struct {
	Name                 string                                `tfschema:"name"`
	ResourceId           string                                `tfschema:"resource_id"`
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

func (r ResourceAssignmentResource) ModelObject() interface{} {
	return &ResourceAssignmentModel{}
}

func (r ResourceAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
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
		"resource_id": {
			// TODO: tests for this
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ResourceAssignmentId(),
		},
	}
	return r.base.arguments(schema)
}

func (r ResourceAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ResourceAssignmentResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "resource_id")
}

func (r ResourceAssignmentResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("resource_id")
}

func (r ResourceAssignmentResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}

func (r ResourceAssignmentResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}
