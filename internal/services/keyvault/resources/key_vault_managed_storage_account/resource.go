// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_storage_account

import (
	"context"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedsecrets"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type deleteAndPurgeSecret struct {
	secretsClient        *secrets.SecretsClient
	deletedSecretsClient *deletedsecrets.DeletedSecretsClient
	keyVaultUri          string
	name                 string
}

func (d deleteAndPurgeSecret) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	secretId := secrets.NewSecretID(d.keyVaultUri, d.name)
	resp, err := d.secretsClient.DeleteSecret(ctx, secretId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeSecret) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	secretVersionId := secrets.NewSecretversionID(d.keyVaultUri, d.name, "")
	resp, err := d.secretsClient.GetSecret(ctx, secretVersionId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeSecret) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	deletedSecretId := deletedsecrets.NewDeletedsecretID(d.keyVaultUri, d.name)
	resp, err := d.deletedSecretsClient.PurgeDeletedSecret(ctx, deletedSecretId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeSecret) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	deletedSecretId := deletedsecrets.NewDeletedsecretID(d.keyVaultUri, d.name)
	resp, err := d.deletedSecretsClient.GetDeletedSecret(ctx, deletedSecretId)
	return autorest.Response{Response: resp.HttpResponse}, err
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
