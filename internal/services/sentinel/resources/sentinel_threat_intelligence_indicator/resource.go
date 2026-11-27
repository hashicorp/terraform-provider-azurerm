// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_threat_intelligence_indicator

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/threatintelligence"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_sentinel_threat_intelligence_indicator"
}

func (r Resource) Description() string {
	return "Manages a Sentinel Threat Intelligence Indicator."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return threatintelligence.ValidateIndicatorID
}
