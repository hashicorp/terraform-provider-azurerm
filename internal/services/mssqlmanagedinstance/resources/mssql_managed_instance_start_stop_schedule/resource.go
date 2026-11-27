// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_start_stop_schedule

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/resources/mssql_managed_instance_start_stop_schedule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_mssql_managed_instance_start_stop_schedule"
}

func (r Resource) Description() string {
	return "Manages Start Stop Schedules for an MS SQL Managed Instance."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedInstanceStartStopScheduleID
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	if !features.SixPointOh() {
		return sdk.StateUpgradeData{}
	}

	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.MsSqlManagedInstanceStartStopScheduleV0ToV1{},
		},
	}
}
