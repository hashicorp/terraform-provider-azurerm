// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	managementGroupValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ManagementGroupAssignmentResource{}

type ManagementGroupAssignmentResource struct {
	base assignmentBaseResource
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
				validation.StringDoesNotContainAny("/"),
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

func (r ManagementGroupAssignmentResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ManagementGroupAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagementGroupAssignmentID
}

func (r ManagementGroupAssignmentResource) ModelObject() interface{} {
	return nil
}

func (r ManagementGroupAssignmentResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("management_group_id")
}

func (r ManagementGroupAssignmentResource) ResourceType() string {
	return "azurerm_management_group_policy_assignment"
}

func (r ManagementGroupAssignmentResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}
