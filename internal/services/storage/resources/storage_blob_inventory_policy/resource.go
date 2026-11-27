// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob_inventory_policy

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_blob_inventory_policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -parent-id "storage_account_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageBlobInventoryPolicyCreateUpdate,
		Read:   resourceStorageBlobInventoryPolicyRead,
		Update: resourceStorageBlobInventoryPolicyCreateUpdate,
		Delete: resourceStorageBlobInventoryPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.StorageAccountId{}, pluginsdk.ResourceTypeForIdentityVirtual),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.StorageAccountId{}, pluginsdk.ResourceTypeForIdentityVirtual),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.BlobInventoryPolicyV0ToV1{},
		}),

		Schema: storageBlobInventoryPolicySchema(),
	}
}
