// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_processing_rule_action_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_monitor_alert_processing_rule_action_group"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return alertprocessingrules.ValidateActionRuleID
}
