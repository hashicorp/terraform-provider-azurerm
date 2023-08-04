// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func dataSourceAutomationVariableInt() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAutomationVariableIntRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: datasourceAutomationVariableCommonSchema(pluginsdk.TypeInt),
	}
}

func dataSourceAutomationVariableIntRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "Int")
}
