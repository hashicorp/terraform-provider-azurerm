// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_customer_managed_key

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -parent-id "storage_account_id"

var (
	storageAccountResourceName  = "azurerm_storage_account"
	storageKindsSupportsSkuTier = map[storageaccounts.Kind]struct{}{
		storageaccounts.KindBlobStorage: {},
		storageaccounts.KindFileStorage: {},
		storageaccounts.KindStorageVTwo: {},
	}
	storageKindsSupportHns = map[storageaccounts.Kind]struct{}{
		storageaccounts.KindBlobStorage:      {},
		storageaccounts.KindBlockBlobStorage: {},
		storageaccounts.KindStorageVTwo:      {},
	}
	storageKindsSupportLargeFileShares = map[storageaccounts.Kind]struct{}{
		storageaccounts.KindFileStorage: {},
		storageaccounts.KindStorageVTwo: {},
	}
	storageReplicationZonalMigrationPairs = map[string]string{
		"LRS":    "ZRS",
		"GRS":    "GZRS",
		"RAGRS":  "RAGZRS",
		"ZRS":    "LRS",
		"GZRS":   "GRS",
		"RAGZRS": "RAGRS",
	}
)

var storageAccountCustomerManagedKeyResourceName = "azurerm_storage_account_customer_managed_key"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountCustomerManagedKeyCreateUpdate,
		Read:   resourceStorageAccountCustomerManagedKeyRead,
		Update: resourceStorageAccountCustomerManagedKeyCreateUpdate,
		Delete: resourceStorageAccountCustomerManagedKeyDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.StorageAccountId{}, pluginsdk.ResourceTypeForIdentityVirtual),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.StorageAccountId{}, pluginsdk.ResourceTypeForIdentityVirtual),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: storageAccountCustomerManagedKeySchema(),
	}
}
