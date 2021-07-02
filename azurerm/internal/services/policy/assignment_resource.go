package policy

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ResourceAssignmentResource{}

type ResourceAssignmentResource struct {
	base assignmentBaseResource
}

func (r ResourceAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"resource_id": {
			//TODO: tests for this
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

func (r ResourceAssignmentResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ResourceAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceAssignmentId()
}

func (r ResourceAssignmentResource) ModelObject() interface{} {
	return nil
}

func (r ResourceAssignmentResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("resource_id")
}

func (r ResourceAssignmentResource) ResourceType() string {
	return "azurerm_resource_policy_assignment"
}

func (r ResourceAssignmentResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}
