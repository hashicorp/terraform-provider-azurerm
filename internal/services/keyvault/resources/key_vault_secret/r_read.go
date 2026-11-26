// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_secret

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedsecrets"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/secrets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resources"
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

	secretVersionId := secrets.NewSecretversionID(id.KeyVaultBaseURL, id.Name, "")

	resourceVersionlessId := parse.NewSecretVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name).ID()
	resp, err := getLatestSecret(ctx, d, meta.(*clients.Client), secretVersionId, resourceVersionlessId)
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
	// Note that the value will be unset if it is a write-only value because getSecretWithoutValue does not set the value
	d.Set("value", resp.Model.Value)
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
	d.Set("resource_versionless_id", resourceVersionlessId)

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

func getLatestSecret(ctx context.Context, d *pluginsdk.ResourceData, client *clients.Client, secretVersionId secrets.SecretversionId, resourceVersionlessId string) (secrets.GetSecretOperationResponse, error) {
	if _, ok := d.GetOk("value_wo_version"); ok {
		return getSecretWithoutValue(ctx, client.Resource.ResourcesClient, resourceVersionlessId)
	}

	// an empty version in the ID indicates the latest version
	return client.KeyVault.DataPlaneKeyVaultClient.Secrets.Clone(secretVersionId.BaseURI).GetSecret(ctx, secretVersionId)
}

func getSecretWithoutValue(ctx context.Context, client *resources.ResourcesClient, resourceID string) (secrets.GetSecretOperationResponse, error) {
	resp, err := client.Get(ctx, commonids.NewScopeID(resourceID))
	secret := secrets.GetSecretOperationResponse{
		HttpResponse: resp.HttpResponse,
	}
	if err != nil {
		return secret, err
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return secret, fmt.Errorf("getting Key Vault Secret %q: missing properties", resourceID)
	}

	prop, ok := (*resp.Model.Properties).(map[string]any)
	if !ok {
		return secret, fmt.Errorf("getting Key Vault Secret %q: unexpected property type %T", resourceID, *resp.Model.Properties)
	}

	model := secrets.SecretBundle{}

	if v, ok := prop["secretUriWithVersion"].(string); ok {
		model.Id = &v
	}
	if v, ok := prop["contentType"].(string); ok {
		model.ContentType = &v
	}
	if attrs, ok := prop["attributes"].(map[string]any); ok {
		model.Attributes = &secrets.SecretAttributes{}
		if v, ok := attrs["enabled"].(bool); ok {
			model.Attributes.Enabled = &v
		}
		if v, ok := attrs["created"].(float64); ok {
			model.Attributes.Created = pointer.To(int64(v))
		}
		if v, ok := attrs["updated"].(float64); ok {
			model.Attributes.Updated = pointer.To(int64(v))
		}
		if v, ok := attrs["exp"].(float64); ok {
			model.Attributes.Exp = pointer.To(int64(v))
		}
		if v, ok := attrs["nbf"].(float64); ok {
			model.Attributes.Nbf = pointer.To(int64(v))
		}
	}

	model.Tags = resp.Model.Tags
	secret.Model = &model

	return secret, nil
}
