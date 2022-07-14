package keyvault

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceKeyVaultKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultKeyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.VaultID,
			},

			"key_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"key_size": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"key_opts": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"curve": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"n": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"e": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"x": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"y": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_key_pem": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_key_openssh": {
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

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceKeyVaultKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := parse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Key %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	resp, err := client.GetKey(ctx, *keyVaultBaseUri, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Key %q was not found in Key Vault at URI %q", name, *keyVaultBaseUri)
		}

		return err
	}

	id := *resp.Key.Kid
	parsedId, err := parse.ParseNestedItemID(id)
	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("key_vault_id", keyVaultId.ID())
	d.Set("versionless_id", parsedId.VersionlessID())

	if key := resp.Key; key != nil {
		d.Set("key_type", string(key.Kty))

		options := flattenKeyVaultKeyDataSourceOptions(key.KeyOps)
		if err := d.Set("key_opts", options); err != nil {
			return err
		}

		d.Set("n", key.N)
		d.Set("e", key.E)
		d.Set("x", key.X)
		d.Set("y", key.Y)
		d.Set("curve", key.Crv)

		if key := resp.Key; key != nil {
			if key.Kty == keyvault.RSA || key.Kty == keyvault.RSAHSM {
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
			} else if key.Kty == keyvault.EC || key.Kty == keyvault.ECHSM {
				// do ec keys
				xBytes, err := base64.RawURLEncoding.DecodeString(*key.X)
				if err != nil {
					return fmt.Errorf("failed to decode X: %+v", err)
				}
				yBytes, err := base64.RawURLEncoding.DecodeString(*key.Y)
				if err != nil {
					return fmt.Errorf("failed to decode Y: %+v", err)
				}
				publicKey := &ecdsa.PublicKey{
					X: big.NewInt(0).SetBytes(xBytes),
					Y: big.NewInt(0).SetBytes(yBytes),
				}
				switch key.Crv {
				case keyvault.P256:
					publicKey.Curve = elliptic.P256()
				case keyvault.P384:
					publicKey.Curve = elliptic.P384()
				case keyvault.P521:
					publicKey.Curve = elliptic.P521()
				}
				if publicKey.Curve != nil {
					err = readPublicKey(d, publicKey)
					if err != nil {
						return fmt.Errorf("failed to read public key: %+v", err)
					}
				}
			}
		}
	}

	d.Set("version", parsedId.Version)

	d.Set("resource_id", parse.NewKeyID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, parsedId.Name, parsedId.Version).ID())
	d.Set("resource_versionless_id", parse.NewKeyVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, parsedId.Name).ID())

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenKeyVaultKeyDataSourceOptions(input *[]string) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		for _, option := range *input {
			results = append(results, option)
		}
	}

	return results
}
