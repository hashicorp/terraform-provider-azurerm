// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

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
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
	"golang.org/x/crypto/ssh"
)

func resourceKeyVaultKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultKeyCreate,
		Read:   resourceKeyVaultKeyRead,
		Update: resourceKeyVaultKeyUpdate,
		Delete: resourceKeyVaultKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.ParseNestedItemID(id)
			return err
		}, nestedItemResourceImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read:   pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

			"key_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// turns out Azure's *really* sensitive about the casing of these
				// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.JSONWebKeyTypeEC),
					string(keyvault.JSONWebKeyTypeECHSM),
					string(keyvault.JSONWebKeyTypeRSA),
					string(keyvault.JSONWebKeyTypeRSAHSM),
				}, false),
			},

			"key_size": {
				Type:          pluginsdk.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"curve"},
			},

			"key_opts": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					// turns out Azure's *really* sensitive about the casing of these
					// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
					ValidateFunc: validation.StringInSlice([]string{
						string(keyvault.JSONWebKeyOperationDecrypt),
						string(keyvault.JSONWebKeyOperationEncrypt),
						string(keyvault.JSONWebKeyOperationSign),
						string(keyvault.JSONWebKeyOperationUnwrapKey),
						string(keyvault.JSONWebKeyOperationVerify),
						string(keyvault.JSONWebKeyOperationWrapKey),
					}, false),
				},
			},

			"curve": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
					return old == "SECP256K1" && new == string(keyvault.JSONWebKeyCurveNameP256K)
				},
				ValidateFunc: func() pluginsdk.SchemaValidateFunc {
					out := []string{
						string(keyvault.JSONWebKeyCurveNameP256),
						string(keyvault.JSONWebKeyCurveNameP256K),
						string(keyvault.JSONWebKeyCurveNameP384),
						string(keyvault.JSONWebKeyCurveNameP521),
					}
					return validation.StringInSlice(out, false)
				}(),
				// TODO: the curve name should probably be mandatory for EC in the future,
				// but handle the diff so that we don't break existing configurations and
				// imported EC keys
				ConflictsWith: []string{"key_size"},
			},

			"not_before_date": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"expiration_date": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"rotation_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"expire_after": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.ISO8601DurationBetween("P28D", "P100Y"),
							AtLeastOneOf: []string{
								"rotation_policy.0.expire_after",
								"rotation_policy.0.automatic",
							},
							RequiredWith: []string{
								"rotation_policy.0.expire_after",
								"rotation_policy.0.notify_before_expiry",
							},
						},

						// <= expiry_time - 7, >=7
						"notify_before_expiry": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.ISO8601DurationBetween("P7D", "P36493D"),
							RequiredWith: []string{
								"rotation_policy.0.expire_after",
								"rotation_policy.0.notify_before_expiry",
							},
						},

						"automatic": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"time_after_creation": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validate.ISO8601Duration,
										AtLeastOneOf: []string{
											"rotation_policy.0.automatic.0.time_after_creation",
											"rotation_policy.0.automatic.0.time_before_expiry",
										},
									},
									"time_before_expiry": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validate.ISO8601Duration,
										AtLeastOneOf: []string{
											"rotation_policy.0.automatic.0.time_after_creation",
											"rotation_policy.0.automatic.0.time_before_expiry",
										},
									},
								},
							},
						},
					},
				},
			},

			// Computed
			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"versionless_id": {
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

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("expiration_date", func(ctx context.Context, old, new, meta interface{}) bool {
				oldDateStr, ok1 := old.(string)
				newDateStr, ok2 := new.(string)
				if !ok1 || !ok2 {
					return false // If old or new values are not strings, don't force new
				}

				// Parse old and new expiration dates
				oldDate, err1 := time.Parse(time.RFC3339, oldDateStr)
				newDate, err2 := time.Parse(time.RFC3339, newDateStr)
				if err1 != nil || err2 != nil {
					return false // If there are parsing errors, don't force new
				}

				// Compare old and new expiration dates
				if newDate.After(oldDate) {
					// If the new expiration date is further in the future, allow update
					return false
				}

				// If the new expiration date is not further, force recreation
				return true
			}),
		),
	}
}

func resourceKeyVaultKeyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Key creation.")
	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Key %q vault url from id %q: %+v", name, *keyVaultId, err)
	}

	existing, err := client.GetKey(ctx, *keyVaultBaseUri, name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Key %q (Key Vault %q): %s", name, *keyVaultBaseUri, err)
		}
	}

	if existing.Key != nil && existing.Key.Kid != nil && *existing.Key.Kid != "" {
		return tf.ImportAsExistsError("azurerm_key_vault_key", *existing.Key.Kid)
	}

	keyType := d.Get("key_type").(string)
	keyOptions := expandKeyVaultKeyOptions(d)
	t := d.Get("tags").(map[string]interface{})

	// TODO: support Importing Keys once this is fixed:
	// https://github.com/Azure/azure-rest-api-specs/issues/1747
	parameters := keyvault.KeyCreateParameters{
		Kty:    keyvault.JSONWebKeyType(keyType),
		KeyOps: keyOptions,
		KeyAttributes: &keyvault.KeyAttributes{
			Enabled: utils.Bool(true),
		},

		Tags: tags.Expand(t),
	}

	if parameters.Kty == keyvault.JSONWebKeyTypeEC || parameters.Kty == keyvault.JSONWebKeyTypeECHSM {
		curveName := d.Get("curve").(string)
		parameters.Curve = keyvault.JSONWebKeyCurveName(curveName)
	} else if parameters.Kty == keyvault.JSONWebKeyTypeRSA || parameters.Kty == keyvault.JSONWebKeyTypeRSAHSM {
		keySize, ok := d.GetOk("key_size")
		if !ok {
			return fmt.Errorf("Key size is required when creating an RSA key")
		}
		parameters.KeySize = utils.Int32(int32(keySize.(int)))
	}
	// TODO: support `oct` once this is fixed
	// https://github.com/Azure/azure-rest-api-specs/issues/1739#issuecomment-332236257

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		notBeforeUnixTime := date.UnixTime(notBeforeDate)
		parameters.KeyAttributes.NotBefore = &notBeforeUnixTime
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		expirationUnixTime := date.UnixTime(expirationDate)
		parameters.KeyAttributes.Expires = &expirationUnixTime
	}

	if resp, err := client.CreateKey(ctx, *keyVaultBaseUri, name, parameters); err != nil {
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeys && utils.ResponseWasConflict(resp.Response) {
			recoveredKey, err := client.RecoverDeletedKey(ctx, *keyVaultBaseUri, name)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Recovering Key %q with ID: %q", name, *recoveredKey.Key.Kid)
			if kid := recoveredKey.Key.Kid; kid != nil {
				stateConf := &pluginsdk.StateChangeConf{
					Pending:                   []string{"pending"},
					Target:                    []string{"available"},
					Refresh:                   keyVaultChildItemRefreshFunc(*kid),
					Delay:                     30 * time.Second,
					PollInterval:              10 * time.Second,
					ContinuousTargetOccurence: 10,
					Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
				}

				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for Key Vault Secret %q to become available: %s", name, err)
				}
				log.Printf("[DEBUG] Key %q recovered with ID: %q", name, *kid)
			}
		} else {
			return fmt.Errorf("Creating Key: %+v", err)
		}
	}

	if v, ok := d.GetOk("rotation_policy"); ok {
		if respPolicy, err := client.UpdateKeyRotationPolicy(ctx, *keyVaultBaseUri, name, expandKeyVaultKeyRotationPolicy(v.([]interface{}))); err != nil {
			if utils.ResponseWasForbidden(respPolicy.Response) {
				return fmt.Errorf("current client lacks permissions to create Key Rotation Policy for Key %q (%q, Vault url: %q), please update this as described here: %s : %v", name, *keyVaultId, *keyVaultBaseUri, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
			}
			return fmt.Errorf("creating Key Rotation Policy: %+v", err)
		}
	}

	// "" indicates the latest version
	read, err := client.GetKey(ctx, *keyVaultBaseUri, name, "")
	if err != nil {
		return err
	}

	if read.Key == nil || read.Key.Kid == nil {
		return fmt.Errorf("cannot read KeyVault Key '%s' (in key vault '%s')", name, *keyVaultBaseUri)
	}
	keyId, err := parse.ParseNestedItemID(*read.Key.Kid)
	if err != nil {
		return err
	}
	d.SetId(keyId.ID())

	return resourceKeyVaultKeyRead(d, meta)
}

func resourceKeyVaultKeyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	meta.(*clients.Client).KeyVault.AddToCache(*keyVaultId, id.KeyVaultBaseUrl)

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	keyOptions := expandKeyVaultKeyOptions(d)
	t := d.Get("tags").(map[string]interface{})

	parameters := keyvault.KeyUpdateParameters{
		KeyOps: keyOptions,
		KeyAttributes: &keyvault.KeyAttributes{
			Enabled: utils.Bool(true),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		notBeforeUnixTime := date.UnixTime(notBeforeDate)
		parameters.KeyAttributes.NotBefore = &notBeforeUnixTime
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		expirationUnixTime := date.UnixTime(expirationDate)
		parameters.KeyAttributes.Expires = &expirationUnixTime
	}

	if _, err = client.UpdateKey(ctx, id.KeyVaultBaseUrl, id.Name, "", parameters); err != nil {
		return err
	}

	if d.HasChange("rotation_policy"); ok {
		if respPolicy, err := client.UpdateKeyRotationPolicy(ctx, id.KeyVaultBaseUrl, id.Name, expandKeyVaultKeyRotationPolicy(d.Get("rotation_policy").([]interface{}))); err != nil {
			if utils.ResponseWasForbidden(respPolicy.Response) {
				return fmt.Errorf("current client lacks permissions to update Key Rotation Policy for Key %q (%q, Vault url: %q), please update this as described here: %s : %v", id.Name, *keyVaultId, id.KeyVaultBaseUrl, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
			}
			return fmt.Errorf("creating Key Rotation Policy: %+v", err)
		}
	}

	return resourceKeyVaultKeyRead(d, meta)
}

func resourceKeyVaultKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	resp, err := client.GetKey(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Key %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", id.Name)

	if key := resp.Key; key != nil {
		d.Set("key_type", string(key.Kty))

		options := flattenKeyVaultKeyOptions(key.KeyOps)
		if err := d.Set("key_opts", options); err != nil {
			return err
		}

		d.Set("n", key.N)
		d.Set("e", key.E)
		d.Set("x", key.X)
		d.Set("y", key.Y)
		if key.N != nil {
			nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
			if err != nil {
				return fmt.Errorf("Could not decode N: %+v", err)
			}
			d.Set("key_size", len(nBytes)*8)
		}

		d.Set("curve", key.Crv)
	}

	if attributes := resp.Attributes; attributes != nil {
		if v := attributes.NotBefore; v != nil {
			d.Set("not_before_date", time.Time(*v).Format(time.RFC3339))
		}

		if v := attributes.Expires; v != nil {
			d.Set("expiration_date", time.Time(*v).Format(time.RFC3339))
		}
	}

	// Computed
	d.Set("version", id.Version)
	d.Set("versionless_id", id.VersionlessID())
	if key := resp.Key; key != nil {
		if key.Kty == keyvault.JSONWebKeyTypeRSA || key.Kty == keyvault.JSONWebKeyTypeRSAHSM {
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
		} else if key.Kty == keyvault.JSONWebKeyTypeEC || key.Kty == keyvault.JSONWebKeyTypeECHSM {
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
			case keyvault.JSONWebKeyCurveNameP256:
				publicKey.Curve = elliptic.P256()
			case keyvault.JSONWebKeyCurveNameP384:
				publicKey.Curve = elliptic.P384()
			case keyvault.JSONWebKeyCurveNameP521:
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

	d.Set("resource_id", parse.NewKeyID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name, id.Version).ID())
	d.Set("resource_versionless_id", parse.NewKeyVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name).ID())

	respPolicy, err := client.GetKeyRotationPolicy(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		switch {
		case utils.ResponseWasForbidden(respPolicy.Response):
			// If client is not authorized to access the policy:
			return fmt.Errorf("current client lacks permissions to read Key Rotation Policy for Key %q (%q, Vault url: %q), please update this as described here: %s : %v", id.Name, *keyVaultId, id.KeyVaultBaseUrl, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
		case utils.ResponseWasNotFound(respPolicy.Response):
			return tags.FlattenAndSet(d, resp.Tags)
		default:
			return err
		}
	}

	rotationPolicy := flattenKeyVaultKeyRotationPolicy(respPolicy)
	if err := d.Set("rotation_policy", rotationPolicy); err != nil {
		return fmt.Errorf("setting Key Vault Key Rotation Policy: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceKeyVaultKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	kv, err := keyVaultsClient.VaultsClient.Get(ctx, *keyVaultId)
	if err != nil {
		if response.WasNotFound(kv.HttpResponse) {
			log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving key vault %q properties: %+v", *keyVaultId, err)
	}

	shouldPurge := meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedKeysOnDestroy
	if shouldPurge && kv.Model != nil && utils.NormaliseNilableBool(kv.Model.Properties.EnablePurgeProtection) {
		log.Printf("[DEBUG] cannot purge key %q because vault %q has purge protection enabled", id.Name, keyVaultId.String())
		shouldPurge = false
	}

	description := fmt.Sprintf("Key %q (Key Vault %q)", id.Name, id.KeyVaultBaseUrl)
	deleter := deleteAndPurgeKey{
		client:      client,
		keyVaultUri: id.KeyVaultBaseUrl,
		name:        id.Name,
	}
	if err := deleteAndOptionallyPurge(ctx, description, shouldPurge, deleter); err != nil {
		return err
	}

	return nil
}

var _ deleteAndPurgeNestedItem = deleteAndPurgeKey{}

type deleteAndPurgeKey struct {
	client      *keyvault.BaseClient
	keyVaultUri string
	name        string
}

func (d deleteAndPurgeKey) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.DeleteKey(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func (d deleteAndPurgeKey) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetKey(ctx, d.keyVaultUri, d.name, "")
	return resp.Response, err
}

func (d deleteAndPurgeKey) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	return d.client.PurgeDeletedKey(ctx, d.keyVaultUri, d.name)
}

func (d deleteAndPurgeKey) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetDeletedKey(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func expandKeyVaultKeyOptions(d *pluginsdk.ResourceData) *[]keyvault.JSONWebKeyOperation {
	options := d.Get("key_opts").([]interface{})
	results := make([]keyvault.JSONWebKeyOperation, 0, len(options))

	for _, option := range options {
		results = append(results, keyvault.JSONWebKeyOperation(option.(string)))
	}

	return &results
}

func expandKeyVaultKeyRotationPolicy(v []interface{}) keyvault.KeyRotationPolicy {
	if len(v) == 0 {
		return keyvault.KeyRotationPolicy{LifetimeActions: &[]keyvault.LifetimeActions{}}
	}

	policy := v[0].(map[string]interface{})

	var expiryTime *string = nil // needs to be set to nil if not set
	if rawExpiryTime := policy["expire_after"]; rawExpiryTime != nil && rawExpiryTime.(string) != "" {
		expiryTime = utils.String(rawExpiryTime.(string))
	}

	lifetimeActions := make([]keyvault.LifetimeActions, 0)
	if rawNotificationTime := policy["notify_before_expiry"]; rawNotificationTime != nil && rawNotificationTime.(string) != "" {
		lifetimeActionNotify := keyvault.LifetimeActions{
			Trigger: &keyvault.LifetimeActionsTrigger{
				TimeBeforeExpiry: utils.String(rawNotificationTime.(string)), // for Type: keyvault.Notify always TimeBeforeExpiry
			},
			Action: &keyvault.LifetimeActionsType{
				Type: keyvault.ActionTypeNotify,
			},
		}
		lifetimeActions = append(lifetimeActions, lifetimeActionNotify)
	}

	if autoRotationList := policy["automatic"].([]interface{}); len(autoRotationList) == 1 && autoRotationList[0] != nil {
		lifetimeActionRotate := keyvault.LifetimeActions{
			Action: &keyvault.LifetimeActionsType{
				Type: keyvault.ActionTypeRotate,
			},
			Trigger: &keyvault.LifetimeActionsTrigger{},
		}
		autoRotationRaw := autoRotationList[0].(map[string]interface{})

		if v := autoRotationRaw["time_after_creation"]; v != nil && v.(string) != "" {
			timeAfterCreate := v.(string)
			lifetimeActionRotate.Trigger.TimeAfterCreate = &timeAfterCreate
		}

		if v := autoRotationRaw["time_before_expiry"]; v != nil && v.(string) != "" {
			timeBeforeExpiry := v.(string)
			lifetimeActionRotate.Trigger.TimeBeforeExpiry = &timeBeforeExpiry
		}

		lifetimeActions = append(lifetimeActions, lifetimeActionRotate)
	}

	return keyvault.KeyRotationPolicy{
		LifetimeActions: &lifetimeActions,
		Attributes: &keyvault.KeyRotationPolicyAttributes{
			ExpiryTime: expiryTime,
		},
	}
}

func flattenKeyVaultKeyOptions(input *[]string) []interface{} {
	results := make([]interface{}, 0, len(*input))

	for _, option := range *input {
		results = append(results, option)
	}

	return results
}

func flattenKeyVaultKeyRotationPolicy(input keyvault.KeyRotationPolicy) []interface{} {
	if input.LifetimeActions == nil && input.Attributes == nil {
		return []interface{}{}
	}

	policy := make(map[string]interface{})
	if input.Attributes != nil && input.Attributes.ExpiryTime != nil && *input.Attributes.ExpiryTime != "" {
		policy["expire_after"] = *input.Attributes.ExpiryTime
	}

	if input.LifetimeActions != nil {
		for _, ltAction := range *input.LifetimeActions {
			action := ltAction.Action
			trigger := ltAction.Trigger

			if action != nil && trigger != nil && action.Type != "" && strings.EqualFold(string(action.Type), string(keyvault.ActionTypeNotify)) && trigger.TimeBeforeExpiry != nil && *trigger.TimeBeforeExpiry != "" {
				// Somehow a default is set after creation for notify_before_expiry
				// Submitting this set value in the next run will not work though..
				if policy["expire_after"] != nil {
					policy["notify_before_expiry"] = *trigger.TimeBeforeExpiry
				}
			}

			if action != nil && trigger != nil && action.Type != "" && strings.EqualFold(string(action.Type), string(keyvault.ActionTypeRotate)) {
				autoRotation := make(map[string]interface{}, 0)
				autoRotation["time_after_creation"] = pointer.From(trigger.TimeAfterCreate)
				autoRotation["time_before_expiry"] = pointer.From(trigger.TimeBeforeExpiry)
				policy["automatic"] = []map[string]interface{}{autoRotation}
			}
		}
	}

	if len(policy) == 0 {
		return []interface{}{}
	}

	return []interface{}{policy}
}

// Credit to Hashicorp modified from https://github.com/hashicorp/terraform-provider-tls/blob/v3.1.0/internal/provider/util.go#L79-L105
func readPublicKey(d *pluginsdk.ResourceData, pubKey interface{}) error {
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
