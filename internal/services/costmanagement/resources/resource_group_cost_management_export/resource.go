// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource_group_cost_management_export

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct {
	base helpers.CostManagementExportBaseResource
}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_resource_group_cost_management_export"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceGroupCostManagementExportID
}
