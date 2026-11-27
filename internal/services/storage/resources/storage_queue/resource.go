// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_queue

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storagequeues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_queue/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageQueueCreate,
		Read:   resourceStorageQueueRead,
		Update: resourceStorageQueueUpdate,
		Delete: resourceStorageQueueDelete,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := storagequeues.ParseQueueID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.QueueV0ToV1{},
			1: migration.StorageQueueV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: storageQueueSchema(),
	}
}
