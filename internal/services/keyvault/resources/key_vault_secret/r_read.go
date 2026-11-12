// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_secret

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedsecrets"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultSecretRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyvault.ParseNestedItemID(d.Id(), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseURL)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseURL, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseURL)
		d.SetId("")
		return nil
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseURL, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseURL)
		d.SetId("")
		return nil
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Secrets.Clone(id.KeyVaultBaseURL)
	secretVersionId := secrets.NewSecretversionID(id.KeyVaultBaseURL, id.Name, "")

	// we always want to get the latest version
	resp, err := client.GetSecret(ctx, secretVersionId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Secret %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseURL)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure KeyVault Secret %s: %+v", id.Name, err)
	}

	if resp.Model == nil || resp.Model.Id == nil {
		return fmt.Errorf("retrieving %s: unable to determine secret ID", secrets.NewSecretID(id.KeyVaultBaseURL, id.Name))
	}

	// the version may have changed, so parse the updated id
	respID, err := keyvault.ParseNestedItemID(*resp.Model.Id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
	if err != nil {
		return err
	}

	d.Set("name", respID.Name)
	d.Set("value", resp.Model.Value)
	// Unset value if is a write-only value
	if _, ok := d.GetOk("value_wo_version"); ok {
		d.Set("value", nil)
	}
	d.Set("version", respID.Version)
	d.Set("content_type", resp.Model.ContentType)
	d.Set("versionless_id", id.VersionlessID())
	d.Set("value_wo_version", d.Get("value_wo_version").(int))

	if attributes := resp.Model.Attributes; attributes != nil {
		notBeforeDate := ""
		if v := attributes.Nbf; v != nil {
			notBeforeDate = time.Unix(*v, 0).UTC().Format(time.RFC3339)
		}
		d.Set("not_before_date", notBeforeDate)

		expirationDate := ""
		if v := attributes.Exp; v != nil {
			expirationDate = time.Unix(*v, 0).UTC().Format(time.RFC3339)
		}
		d.Set("expiration_date", expirationDate)
	}

	d.Set("resource_id", parse.NewSecretID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name, id.Version).ID())
	d.Set("resource_versionless_id", parse.NewSecretVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name).ID())

	return tags.FlattenAndSet(d, resp.Model.Tags)
}

var _ helpers.DeleteAndPurgeNestedItem = deleteAndPurgeSecret{}

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
