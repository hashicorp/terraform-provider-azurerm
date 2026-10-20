// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_secret

import (
	"context"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DeleteAndPurgeNestedItem interface {
	DeleteNestedItem(ctx context.Context) (autorest.Response, error)
	NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error)
	PurgeNestedItem(ctx context.Context) (autorest.Response, error)
	NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error)
}

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultSecretCreate,
		Read:   resourceKeyVaultSecretRead,
		Update: resourceKeyVaultSecretUpdate,
		Delete: resourceKeyVaultSecretDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := keyvault.ParseNestedItemID(id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
			return err
		}, helpers.NestedItemResourceImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read:   pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: keyVaultSecretSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
			// Updating `value` or `value_wo` creates a new version of the secret, so the version and the attributes derived from it are unknown until after apply.
			// `value_wo_version` is checked because WriteOnly attributes (`value_wo`) carry no prior state.
			if d.Id() != "" && d.HasChanges("value", "value_wo_version") {
				for _, field := range []string{"version", "resource_id"} {
					if err := d.SetNewComputed(field); err != nil {
						return err
					}
				}
			}

			return nil
		}),
	}
}
