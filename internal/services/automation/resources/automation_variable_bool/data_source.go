// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_variable_bool

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Automation Bool Variable.",
		Read:        dataSourceAutomationVariableBoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: helpers.DatasourceAutomationVariableCommonSchema(pluginsdk.TypeBool),
	}
}

func dataSourceAutomationVariableBoolRead(d *pluginsdk.ResourceData, meta any) error {
	return helpers.DataSourceAutomationVariableRead(d, meta, "Bool")
}
