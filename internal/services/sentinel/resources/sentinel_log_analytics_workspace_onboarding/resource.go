// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_log_analytics_workspace_onboarding

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/sentinelonboardingstates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_sentinel_log_analytics_workspace_onboarding"
}

func (r Resource) Description() string {
	return "Manages a Security Insights Sentinel Onboarding."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sentinelonboardingstates.ValidateOnboardingStateID
}
