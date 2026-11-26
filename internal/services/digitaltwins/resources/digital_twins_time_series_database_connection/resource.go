// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_time_series_database_connection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (m Resource) ResourceType() string {
	return "azurerm_digital_twins_time_series_database_connection"
}

func (m Resource) Description() string {
	return "Manages a Digital Twins Time Series Database Connection."
}

func (m Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return timeseriesdatabaseconnections.ValidateTimeSeriesDatabaseConnectionID
}
