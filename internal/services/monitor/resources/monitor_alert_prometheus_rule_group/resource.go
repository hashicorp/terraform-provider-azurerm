// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_prometheus_rule_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroupresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_monitor_alert_prometheus_rule_group"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return prometheusrulegroupresources.ValidatePrometheusRuleGroupID
}
