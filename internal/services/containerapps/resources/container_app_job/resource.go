// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_job

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app_job/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

var _ sdk.ResourceWithStateMigration = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_container_app_job"
}

func (r Resource) Description() string {
	return "Manages a Container App Job."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return jobs.ValidateJobID
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ContainerAppJobV0ToV1{},
		},
	}
}
