// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_built_in

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_sentinel_alert_rule_anomaly_built_in"
}

func (r Resource) Description() string {
	return "Manages a Built-in Anomaly Alert Rule."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return securitymlanalyticssettings.ValidateSecurityMLAnalyticsSettingID
}
