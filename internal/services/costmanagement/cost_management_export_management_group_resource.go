package costmanagement

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	mgmtGrpValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagementGroupCostManagementExportResource struct {
	base costManagementExportBaseResource
}

var _ sdk.Resource = ManagementGroupCostManagementExportResource{}

func (r ManagementGroupCostManagementExportResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"management_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: mgmtGrpValidate.ManagementGroupID,
		},
	}
	return r.base.arguments(schema)
}

func (r ManagementGroupCostManagementExportResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ManagementGroupCostManagementExportResource) ModelObject() interface{} {
	return nil
}

func (r ManagementGroupCostManagementExportResource) ResourceType() string{
	return "azurerm_cost_management_export_management_group"
}

func (r ManagementGroupCostManagementExportResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagementGroupCostManagementExportID
}

func (r ManagementGroupCostManagementExportResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "management_group_id")
}

func (r ManagementGroupCostManagementExportResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("management_group_id")
}

func (r ManagementGroupCostManagementExportResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ManagementGroupCostManagementExportResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}