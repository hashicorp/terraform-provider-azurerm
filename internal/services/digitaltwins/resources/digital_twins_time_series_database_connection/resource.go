// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_time_series_database_connection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

func (m Resource) ResourceType() string {
	return "azurerm_digital_twins_time_series_database_connection"
}

func (m Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return timeseriesdatabaseconnections.ValidateTimeSeriesDatabaseConnectionID
}
