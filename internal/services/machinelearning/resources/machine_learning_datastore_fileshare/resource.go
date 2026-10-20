// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_fileshare

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/datastore"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_datastore_fileshare/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_machine_learning_datastore_fileshare"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datastore.ValidateDataStoreID
}

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	// rewrite `storage_fileshare_id` to account for azurerm_storage_share.resource_manager_id deprecation
	// in favour of azurerm_storage_share.id in 5.0
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.MachineLearningDataStoreFileShareV0ToV1{},
		},
	}
}
