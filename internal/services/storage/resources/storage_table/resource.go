// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_table

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_table/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageTableCreate,
		Read:   resourceStorageTableRead,
		Delete: resourceStorageTableDelete,
		Update: resourceStorageTableUpdate,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := tables.ParseTableID(id)
			return err
		}),

		SchemaVersion: 3,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.TableV0ToV1{},
			1: migration.TableV1ToV2{},
			2: migration.TableV2ToV3{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: storageTableSchema(),
	}
}
