// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_encryption_scope

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/encryptionscopes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name storage_encryption_scope -service-package-name storage -properties "name" -compare-values "subscription_id:storage_account_id,resource_group_name:storage_account_id,storage_account_name:storage_account_id" -test-name "keyVaultKey"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageEncryptionScopeCreate,
		Read:   resourceStorageEncryptionScopeRead,
		Update: resourceStorageEncryptionScopeUpdate,
		Delete: resourceStorageEncryptionScopeDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&encryptionscopes.EncryptionScopeId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&encryptionscopes.EncryptionScopeId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: storageEncryptionScopeSchema(),
	}
}
