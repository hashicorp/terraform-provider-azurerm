// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_threat_intelligence

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithUpdate         = Resource{}
)

func (a Resource) ResourceType() string {
	return "azurerm_sentinel_alert_rule_threat_intelligence"
}

func (a Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return alertrules.ValidateAlertRuleID
}

func (a Resource) CustomImporter() sdk.ResourceRunFunc {
	return helpers.ImportSentinelAlertRuleForTypedSdk(alertrules.AlertRuleKindThreatIntelligence)
}
