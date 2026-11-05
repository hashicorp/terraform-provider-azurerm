// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_variable_datetime

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceAutomationVariableDateTimeRead(d *pluginsdk.ResourceData, meta any) error {
	return helpers.ResourceAutomationVariableRead(d, meta, "Datetime")
}
