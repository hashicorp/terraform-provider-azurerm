// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceKeyVaultSecret() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultSecretRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

			"value": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"content_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"not_before_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"expiration_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceKeyVaultSecretRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Secret %q vault url from id %q: %+v", name, *keyVaultId, err)
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Secrets.Clone(*keyVaultBaseUri)
	secretVersionId := secrets.NewSecretversionID(*keyVaultBaseUri, name, version)

	resp, err := client.GetSecret(ctx, secretVersionId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("KeyVault Secret %q (KeyVault URI %q) does not exist", name, *keyVaultBaseUri)
		}
		return fmt.Errorf("making Read request on Azure KeyVault Secret %s: %+v", name, err)
	}

	if resp.Model == nil || resp.Model.Id == nil {
		return fmt.Errorf("reading KeyVault Secret %q: response model was nil", name)
	}

	// the version may have changed, so parse the updated id
	secretId, err := parse.ParseNestedItemID(*resp.Model.Id)
	if err != nil {
		return err
	}

	d.SetId(secretId.ID())

	d.Set("name", secretId.Name)
	d.Set("key_vault_id", keyVaultId.ID())
	d.Set("value", resp.Model.Value)
	d.Set("version", secretId.Version)
	d.Set("content_type", resp.Model.ContentType)
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
	d.Set("versionless_id", secretId.VersionlessID())

	d.Set("resource_id", parse.NewSecretID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, secretId.Name, secretId.Version).ID())
	d.Set("resource_versionless_id", parse.NewSecretVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, secretId.Name).ID())

	return tags.FlattenAndSet(d, tags.FromTypedObject(pointer.From(resp.Model.Tags)))
}
