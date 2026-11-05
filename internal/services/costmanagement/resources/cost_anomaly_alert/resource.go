// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cost_anomaly_alert

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_cost_anomaly_alert"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scheduledactionoperationgroup.ValidateScopedScheduledActionID
}
