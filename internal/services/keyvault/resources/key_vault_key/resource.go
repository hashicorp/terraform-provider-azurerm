// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_key

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/keys"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
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
		Create: resourceKeyVaultKeyCreate,
		Read:   resourceKeyVaultKeyRead,
		Update: resourceKeyVaultKeyUpdate,
		Delete: resourceKeyVaultKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := keyvault.ParseNestedItemID(id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
			return err
		}, helpers.NestedItemResourceImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read:   pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: keyVaultKeySchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("expiration_date", func(ctx context.Context, old, new, meta any) bool {
				oldDateStr, ok1 := old.(string)
				newDateStr, ok2 := new.(string)
				if !ok1 || !ok2 {
					return false // If old or new values are not strings, don't force new
				}

				// There isn't a way to remove the expiration date from a key so we'll recreate the resource when it's removed from config
				if oldDateStr != "" && newDateStr == "" {
					return true
				}

				return false
			}),
			func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
				if p := d.Get("release_policy").([]any); len(p) > 0 {
					if t := d.Get("key_type").(string); !slices.Contains([]string{string(keys.JsonWebKeyTypeRSANegativeHSM), string(keys.JsonWebKeyTypeECNegativeHSM)}, t) {
						return fmt.Errorf("when `release_policy` is set, `key_type` must be `RSA-HSM` or `EC-HSM`, got `%s`", t)
					}

					if d.Id() != "" {
						if d.HasChange("release_policy.0.json") && (!d.HasChange("release_policy.0.immutable") && d.Get("release_policy.0.immutable").(bool)) {
							// before triggering a resource recreation, confirm it's not just a whitespace difference
							o, n := d.GetChange("release_policy.0.json")
							if o != "" {
								oldMap, err := structure.ExpandJsonFromString(o.(string))
								if err != nil {
									return fmt.Errorf("expanding `release_policy.json`: %+v", err)
								}

								newMap, err := structure.ExpandJsonFromString(n.(string))
								if err != nil {
									return fmt.Errorf("expanding updated `release_policy.json`: %+v", err)
								}

								if !reflect.DeepEqual(oldMap, newMap) {
									if err := d.ForceNew("release_policy.0.json"); err != nil {
										return err
									}
								}
							}
						}

						if d.HasChange("release_policy.0.immutable") {
							o, n := d.GetChange("release_policy.0.immutable")
							if o.(bool) && !n.(bool) {
								if err := d.ForceNew("release_policy.0.immutable"); err != nil {
									return err
								}
							}
						}
					}
				}

				return nil
			},
		),
	}
}
