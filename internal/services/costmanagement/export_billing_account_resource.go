// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BillingAccountCostManagementExportResource struct {
	base costManagementExportBaseResource
}

var _ sdk.Resource = BillingAccountCostManagementExportResource{}

func (r BillingAccountCostManagementExportResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"billing_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
	}
	return r.base.arguments(schema)
}

func (r BillingAccountCostManagementExportResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r BillingAccountCostManagementExportResource) ModelObject() interface{} {
	return nil
}

func (r BillingAccountCostManagementExportResource) ResourceType() string {
	return "azurerm_billing_account_cost_management_export"
}

func (r BillingAccountCostManagementExportResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BillingAccountCostManagementExportID
}

func (r BillingAccountCostManagementExportResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "billing_account_id")
}

func (r BillingAccountCostManagementExportResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("billing_account_id")
}

func (r BillingAccountCostManagementExportResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r BillingAccountCostManagementExportResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}
