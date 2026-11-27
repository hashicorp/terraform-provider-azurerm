// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_schedule

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_job_schedule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithStateMigration = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_stream_analytics_job_schedule"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.StreamingJobScheduleID
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsJobScheduleV0ToV1{},
		},
	}
}
