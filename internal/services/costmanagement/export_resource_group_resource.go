// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	resourceValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ResourceGroupCostManagementExportResource struct {
	base costManagementExportBaseResource
}

var _ sdk.Resource = ResourceGroupCostManagementExportResource{}

func (r ResourceGroupCostManagementExportResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
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

func (r ResourceGroupCostManagementExportResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ResourceGroupCostManagementExportResource) ModelObject() interface{} {
	return nil
}

func (r ResourceGroupCostManagementExportResource) ResourceType() string {
	return "azurerm_resource_group_cost_management_export"
}

func (r ResourceGroupCostManagementExportResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceGroupCostManagementExportID
}

func (r ResourceGroupCostManagementExportResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "resource_group_id")
}

func (r ResourceGroupCostManagementExportResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("resource_group_id")
}

func (r ResourceGroupCostManagementExportResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ResourceGroupCostManagementExportResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}
