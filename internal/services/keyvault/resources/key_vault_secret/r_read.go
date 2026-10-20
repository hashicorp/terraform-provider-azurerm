// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	kv "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultSecretRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
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

	// we always want to get the latest version
	resp, err := client.GetSecret(ctx, id.KeyVaultBaseURL, id.Name, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[DEBUG] Secret %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseURL)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure KeyVault Secret %s: %+v", id.Name, err)
	}

	// the version may have changed, so parse the updated id
	respID, err := keyvault.ParseNestedItemID(*resp.ID, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
	if err != nil {
		return err
	}

	d.Set("name", respID.Name)
	d.Set("value", resp.Value)
	// Unset value if is a write-only value
	if _, ok := d.GetOk("value_wo_version"); ok {
		d.Set("value", nil)
	}
	d.Set("version", respID.Version)
	d.Set("content_type", resp.ContentType)
	d.Set("versionless_id", id.VersionlessID())
	d.Set("value_wo_version", d.Get("value_wo_version").(int))

	if attributes := resp.Attributes; attributes != nil {
		notBeforeDate := ""
		if v := attributes.NotBefore; v != nil {
			notBeforeDate = time.Time(*v).Format(time.RFC3339)
		}
		d.Set("not_before_date", notBeforeDate)

		expirationDate := ""
		if v := attributes.Expires; v != nil {
			expirationDate = time.Time(*v).Format(time.RFC3339)
		}
		d.Set("expiration_date", expirationDate)
	}

	d.Set("resource_id", parse.NewSecretID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name, id.Version).ID())
	d.Set("resource_versionless_id", parse.NewSecretVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name).ID())

	return tags.FlattenAndSet(d, resp.Tags)
}

var _ helpers.DeleteAndPurgeNestedItem = deleteAndPurgeSecret{}

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
