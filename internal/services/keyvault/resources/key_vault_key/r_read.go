// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_key

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedkeys"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/keys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"golang.org/x/crypto/ssh"
)

func resourceKeyVaultKeyRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyvault.ParseNestedItemID(d.Id(), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
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
		return fmt.Errorf("checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseURL, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseURL)
		d.SetId("")
		return nil
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Keys.Clone(id.KeyVaultBaseURL)
	keyVersionId := keys.NewKeyversionID(id.KeyVaultBaseURL, id.Name, "")
	resp, err := client.GetKey(ctx, keyVersionId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Key %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseURL)
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", id.Name)

	if resp.Model != nil {
		if key := resp.Model.Key; key != nil {
			d.Set("key_type", pointer.FromEnum(key.Kty))

			if err := d.Set("key_opts", flattenKeyVaultKeyOptions(key.KeyOps)); err != nil {
				return err
			}

			d.Set("n", key.N)
			d.Set("e", key.E)
			d.Set("x", key.X)
			d.Set("y", key.Y)
			if key.N != nil {
				nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
				if err != nil {
					return fmt.Errorf("could not decode N: %+v", err)
				}
				d.Set("key_size", len(nBytes)*8)
			}

			d.Set("curve", pointer.FromEnum(key.Crv))
		}

		data, err := flattenKeyVaultKeyReleasePolicy(resp.Model.ReleasePolicy)
		if err != nil {
			return fmt.Errorf("flattening `release_policy`: %+v", err)
		}
		d.Set("release_policy", data)

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

		// Computed
		d.Set("version", id.Version)
		d.Set("versionless_id", id.VersionlessID())
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

	d.Set("resource_id", commonids.NewKeyVaultKeyVersionID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name, id.Version).ID())
	d.Set("resource_versionless_id", commonids.NewKeyVaultKeyID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name).ID())

	keyIdForPolicy := keys.NewKeyID(id.KeyVaultBaseURL, id.Name)
	respPolicy, err := client.GetKeyRotationPolicy(ctx, keyIdForPolicy)
	if err != nil {
		switch {
		case response.WasForbidden(respPolicy.HttpResponse):
			// If client is not authorized to access the policy:
			return fmt.Errorf("current client lacks permissions to read Key Rotation Policy for Key %q (%q, Vault url: %q), please update this as described here: %s : %v", id.Name, *keyVaultId, id.KeyVaultBaseURL, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
		case response.WasNotFound(respPolicy.HttpResponse):
			if resp.Model != nil {
				return tags.FlattenAndSet(d, tags.FromTypedObject(pointer.From(resp.Model.Tags)))
			}
			return nil
		default:
			return err
		}
	}

	if respPolicy.Model != nil {
		if err := d.Set("rotation_policy", flattenKeyVaultKeyRotationPolicy(*respPolicy.Model)); err != nil {
			return fmt.Errorf("setting Key Vault Key Rotation Policy: %+v", err)
		}
	}

	if resp.Model != nil {
		return tags.FlattenAndSet(d, tags.FromTypedObject(pointer.From(resp.Model.Tags)))
	}
	return nil
}

var _ helpers.DeleteAndPurgeNestedItem = deleteAndPurgeKey{}

type deleteAndPurgeKey struct {
	keysClient        *keys.KeysClient
	deletedKeysClient *deletedkeys.DeletedKeysClient
	keyVaultUri       string
	name              string
}

func (d deleteAndPurgeKey) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	keyId := keys.NewKeyID(d.keyVaultUri, d.name)
	resp, err := d.keysClient.DeleteKey(ctx, keyId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeKey) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	keyVersionId := keys.NewKeyversionID(d.keyVaultUri, d.name, "")
	resp, err := d.keysClient.GetKey(ctx, keyVersionId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeKey) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	deletedKeyId := deletedkeys.NewDeletedkeyID(d.keyVaultUri, d.name)
	resp, err := d.deletedKeysClient.PurgeDeletedKey(ctx, deletedKeyId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

func (d deleteAndPurgeKey) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	deletedKeyId := deletedkeys.NewDeletedkeyID(d.keyVaultUri, d.name)
	resp, err := d.deletedKeysClient.GetDeletedKey(ctx, deletedKeyId)
	return autorest.Response{Response: resp.HttpResponse}, err
}

// Credit to Hashicorp modified from https://github.com/hashicorp/terraform-provider-tls/blob/v3.1.0/internal/provider/util.go#L79-L105
func readPublicKey(d *pluginsdk.ResourceData, pubKey any) error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key error: %s", err)
	}
	pubKeyPemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	d.Set("public_key_pem", string(pem.EncodeToMemory(pubKeyPemBlock)))

	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err == nil {
		// Not all EC types can be SSH keys, so we'll produce this only
		// if an appropriate type was selected.
		sshPubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)
		d.Set("public_key_openssh", string(sshPubKeyBytes))
	} else {
		d.Set("public_key_openssh", "")
	}
	return nil
}
