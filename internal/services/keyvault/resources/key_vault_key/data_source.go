// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/keys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Key Vault Key.",
		Read:        dataSourceKeyVaultKeyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Key Vault Key.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyvault.ValidateNestedItemName,
			},

			"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}, "Key Vault"),

			"key_type": {
				Description: "The Key Type of this Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"key_size": {
				Description: "The Size of this Key Vault Key.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"key_opts": {
				Description: "A list of JSON web key operations assigned to this Key Vault Key.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"version": {
				Description: "The current version of the Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"versionless_id": {
				Description: "The Base ID of the Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"curve": {
				Description: "The EC Curve name of this Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"n": {
				Description: "The RSA modulus of this Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"e": {
				Description: "The RSA public exponent of this Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"x": {
				Description: "The EC X component of this Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"y": {
				Description: "The EC Y component of this Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"public_key_pem": {
				Description: "The PEM encoded public key of this Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"public_key_openssh": {
				Description: "The OpenSSH encoded public key of this Key Vault Key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"resource_id": {
				Description: "The (Versioned) ID for this Key Vault Key. This property points to a specific version of a Key Vault Key, as such using this won't auto-rotate values if used in other Azure Services.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"resource_versionless_id": {
				Description: "The Versionless ID of the Key Vault Key. This property allows other Azure Services (that support it) to auto-rotate their value when the Key Vault Key is updated.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceKeyVaultKeyRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Key %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	client := keyVaultsClient.DataPlaneKeyVaultClient.Keys.Clone(*keyVaultBaseUri)
	resp, err := client.GetKey(ctx, keys.NewKeyversionID(*keyVaultBaseUri, name, ""))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("key %q was not found in Key Vault at URI %q", name, *keyVaultBaseUri)
		}

		return err
	}

	if resp.Model == nil || resp.Model.Key == nil || resp.Model.Key.Kid == nil {
		return fmt.Errorf("retrieving key %q from Key Vault at URI %q: unable to determine key ID", name, *keyVaultBaseUri)
	}

	id := *resp.Model.Key.Kid
	parsedId, err := keyvault.ParseNestedItemID(id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}

	d.SetId(id)

	d.Set("key_vault_id", keyVaultId.ID())
	d.Set("versionless_id", parsedId.VersionlessID())

	if key := resp.Model.Key; key != nil {
		d.Set("key_type", pointer.FromEnum(key.Kty))

		if err := d.Set("key_opts", flattenKeyVaultKeyDataSourceOptions(key.KeyOps)); err != nil {
			return err
		}

		d.Set("n", key.N)
		d.Set("e", key.E)
		d.Set("x", key.X)
		d.Set("y", key.Y)
		d.Set("curve", pointer.FromEnum(key.Crv))

		if key := resp.Model.Key; key != nil {
			switch pointer.From(key.Kty) {
			case keys.JsonWebKeyTypeRSA, keys.JsonWebKeyTypeRSANegativeHSM:
				nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
				if err != nil {
					return fmt.Errorf("failed to decode N: %+v", err)
				}
				eBytes, err := base64.RawURLEncoding.DecodeString(*key.E)
				if err != nil {
					return fmt.Errorf("failed to decode E: %+v", err)
				}
				publicKey := &rsa.PublicKey{
					N: big.NewInt(0).SetBytes(nBytes),
					E: int(big.NewInt(0).SetBytes(eBytes).Uint64()),
				}
				err = readPublicKey(d, publicKey)
				if err != nil {
					return fmt.Errorf("failed to read public key: %+v", err)
				}
			case keys.JsonWebKeyTypeEC, keys.JsonWebKeyTypeECNegativeHSM:
				// do ec keys
				xBytes, err := base64.RawURLEncoding.DecodeString(*key.X)
				if err != nil {
					return fmt.Errorf("failed to decode X: %+v", err)
				}
				yBytes, err := base64.RawURLEncoding.DecodeString(*key.Y)
				if err != nil {
					return fmt.Errorf("failed to decode Y: %+v", err)
				}
				// X/Y are deprecated since go 1.26 in favour of ecdsa.ParseUncompressedPublicKey, which needs the curve up front
				publicKey := &ecdsa.PublicKey{
					X: big.NewInt(0).SetBytes(xBytes), //nolint:staticcheck
					Y: big.NewInt(0).SetBytes(yBytes), //nolint:staticcheck
				}
				switch pointer.From(key.Crv) {
				case keys.JsonWebKeyCurveNamePNegativeTwoFiveSix:
					publicKey.Curve = elliptic.P256()
				case keys.JsonWebKeyCurveNamePNegativeThreeEightFour:
					publicKey.Curve = elliptic.P384()
				case keys.JsonWebKeyCurveNamePNegativeFiveTwoOne:
					publicKey.Curve = elliptic.P521()
				}
				if publicKey.Curve != nil {
					if err = readPublicKey(d, publicKey); err != nil {
						return fmt.Errorf("failed to read public key: %+v", err)
					}
				}
			}
		}
	}

	d.Set("version", parsedId.Version)

	d.Set("resource_id", commonids.NewKeyVaultKeyVersionID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, parsedId.Name, parsedId.Version).ID())
	d.Set("resource_versionless_id", commonids.NewKeyVaultKeyID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, parsedId.Name).ID())

	return tags.FlattenAndSet(d, tags.FromTypedObject(pointer.From(resp.Model.Tags)))
}

func flattenKeyVaultKeyDataSourceOptions(input *[]string) []any {
	results := make([]any, 0)

	if input != nil {
		for _, option := range *input {
			results = append(results, option)
		}
	}

	return results
}
