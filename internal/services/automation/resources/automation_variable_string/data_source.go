// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_variable_string

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAutomationVariableStringRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: helpers.DatasourceAutomationVariableCommonSchema(pluginsdk.TypeString),
	}
}

func dataSourceAutomationVariableStringRead(d *pluginsdk.ResourceData, meta any) error {
	return helpers.DataSourceAutomationVariableRead(d, meta, "String")
}
