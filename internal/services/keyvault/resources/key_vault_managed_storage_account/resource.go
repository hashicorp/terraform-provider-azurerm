// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_storage_account

import (
	"context"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	kv "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

type deleteAndPurgeSecret struct {
	client      *kv.BaseClient
	keyVaultUri string
	name        string
}

func (d deleteAndPurgeSecret) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.DeleteSecret(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func (d deleteAndPurgeSecret) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetSecret(ctx, d.keyVaultUri, d.name, "")
	return resp.Response, err
}

func (d deleteAndPurgeSecret) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	return d.client.PurgeDeletedSecret(ctx, d.keyVaultUri, d.name)
}

func (d deleteAndPurgeSecret) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetDeletedSecret(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultManagedStorageAccountCreateUpdate,
		Read:   resourceKeyVaultManagedStorageAccountRead,
		Update: resourceKeyVaultManagedStorageAccountCreateUpdate,
		Delete: resourceKeyVaultManagedStorageAccountDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := keyvault.ParseNestedItemID(id, keyvault.VersionTypeAny, keyvault.NestedItemTypeStorage)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: keyVaultManagedStorageAccountSchema(),
	}
}
