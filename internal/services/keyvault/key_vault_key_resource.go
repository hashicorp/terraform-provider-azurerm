// Copyright IBM Corp. 2014, 2025
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
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/keys"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
	"golang.org/x/crypto/ssh"
)

const vaultBaseURLFmt string = "https://%s.vault.azure.net"

func resourceKeyVaultKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultKeyCreate,
		Read:   resourceKeyVaultKeyRead,
		Update: resourceKeyVaultKeyUpdate,
		Delete: resourceKeyVaultKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := keys.ParseKeyversionID(id)
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

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("expiration_date", func(ctx context.Context, old, new, meta interface{}) bool {
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
		),
	}
}

func resourceKeyVaultKeyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Key creation.")
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	keyVaultBaseUri := fmt.Sprintf(vaultBaseURLFmt, keyVaultId.VaultName)

	id := keys.NewKeyID(keyVaultBaseUri, name)

	client := meta.(*clients.Client).KeyVault.DataPlaneClient.Keys.Clone(keyVaultBaseUri)

	keyVersionID := keys.NewKeyversionID(id.BaseURI, id.KeyName, "")
	existing, err := client.GetKey(ctx, keyVersionID)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_key_vault_key", keyVersionID.ID())
	}

	keyType := d.Get("key_type").(string)
	keyOptions := expandKeyVaultKeyOptions(d)
	t := d.Get("tags").(map[string]interface{})

	// TODO: support Importing Keys once this is fixed:
	// https://github.com/Azure/azure-rest-api-specs/issues/1747
	parameters := keys.KeyCreateParameters{
		Kty:    keys.JsonWebKeyType(keyType),
		KeyOps: keyOptions,
		Attributes: &keys.KeyAttributes{
			Enabled: pointer.To(true),
		},

		Tags: tags.Expand(t),
	}

	switch parameters.Kty {
	case keys.JsonWebKeyTypeEC, keys.JsonWebKeyTypeECNegativeHSM:
		curveName := d.Get("curve").(string)
		parameters.Crv = pointer.ToEnum[keys.JsonWebKeyCurveName](curveName)
	case keys.JsonWebKeyTypeRSA, keys.JsonWebKeyTypeRSANegativeHSM:
		keySize, ok := d.GetOk("key_size")
		if !ok {
			return errors.New("key_size is required when creating an RSA key")
		}
		parameters.KeySize = pointer.To(int64(keySize.(int)))
	case keys.JsonWebKeyTypeOct, keys.JsonWebKeyTypeOctNegativeHSM:
		// TODO: add `oct` support
	}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Nbf = pointer.To(notBeforeDate.UnixMilli())
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Exp = pointer.To(expirationDate.UnixMilli())
	}

	if resp, err := client.CreateKey(ctx, id, parameters); err != nil {
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeys && response.WasConflict(resp.HttpResponse) {
			deletedId := keys.NewDeletedkeyID(id.BaseURI, id.KeyName)
			recoveredKey, err := client.RecoverDeletedKey(ctx, deletedId)
			if err != nil {
				return err
			}
			if recoveredKey.Model == nil || recoveredKey.Model.Key == nil || recoveredKey.Model.Key.Kid == nil {
				return fmt.Errorf("cannot recover deleted key %q", deletedId)
			}
			log.Printf("[DEBUG] Recovering Key %q with ID: %q", name, *recoveredKey.Model.Key.Kid)
			if kid := recoveredKey.Model.Key.Kid; kid != nil {
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
			return fmt.Errorf("creating Key: %+v", err)
		}
	}

	if v, ok := d.GetOk("rotation_policy"); ok {
		if respPolicy, err := client.UpdateKeyRotationPolicy(ctx, id, expandKeyVaultKeyRotationPolicy(v.([]interface{}))); err != nil {
			if response.WasForbidden(respPolicy.HttpResponse) {
				return fmt.Errorf("current client lacks permissions to create Key Rotation Policy for %s, please update this as described here: %s : %v", id, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
			}
			return fmt.Errorf("creating Key Rotation Policy: %+v", err)
		}
	}

	read, err := client.GetKey(ctx, keyVersionID)
	if err != nil {
		return err
	}

	if read.Model == nil || read.Model.Key == nil {
		return fmt.Errorf("cannot read KeyVault Key '%s' (in key vault '%s')", name, keyVaultBaseUri)
	}

	kvId, err := keys.ParseKeyversionID(*read.Model.Key.Kid)
	if err != nil {
		return err
	}
	d.SetId(kvId.ID())

	return resourceKeyVaultKeyRead(d, meta)
}

func resourceKeyVaultKeyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keys.ParseKeyversionID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneClient.Keys.Clone(id.BaseURI)

	keyOptions := expandKeyVaultKeyOptions(d)
	t := d.Get("tags").(map[string]interface{})

	parameters := keys.KeyUpdateParameters{
		KeyOps: keyOptions,
		Attributes: &keys.KeyAttributes{
			Enabled: pointer.To(true),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Nbf = pointer.To(notBeforeDate.UnixMilli())
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Exp = pointer.To(expirationDate.UnixMilli())
	}

	if _, err = client.UpdateKey(ctx, *id, parameters); err != nil {
		return err
	}

	if d.HasChange("rotation_policy") {
		kId := keys.NewKeyID(id.BaseURI, id.KeyName)
		if respPolicy, err := client.UpdateKeyRotationPolicy(ctx, kId, expandKeyVaultKeyRotationPolicy(d.Get("rotation_policy").([]interface{}))); err != nil {
			if response.WasForbidden(respPolicy.HttpResponse) {
				return fmt.Errorf("current client lacks permissions to update Key Rotation Policy for Key %s, please update this as described here: %s : %v", *id, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
			}
			return fmt.Errorf("creating Key Rotation Policy: %+v", err)
		}
	}

	return resourceKeyVaultKeyRead(d, meta)
}

func resourceKeyVaultKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keys.ParseKeyversionID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneClient.Keys.Clone(id.BaseURI)

	resp, err := client.GetKey(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", id.KeyName)
	d.Set("version", id.Keyversion)
	d.Set("versionless_id", keys.NewKeyID(id.BaseURI, id.KeyName).ID())

	if model := resp.Model; model != nil {
		if key := model.Key; key != nil {
			d.Set("key_type", pointer.From(key.Kty))

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
					return fmt.Errorf("could not decode N: %+v", err)
				}
				d.Set("key_size", len(nBytes)*8)
			}

			d.Set("curve", key.Crv)
		}

		if attributes := model.Attributes; attributes != nil {
			notBeforeDate := ""
			if v := attributes.Nbf; v != nil {
				notBeforeDate = time.UnixMilli(*v).UTC().Format(time.RFC3339)
			}
			d.Set("not_before_date", notBeforeDate)

			expirationDate := ""
			if v := attributes.Exp; v != nil {
				expirationDate = time.UnixMilli(*v).UTC().Format(time.RFC3339)
			}
			d.Set("expiration_date", expirationDate)
		}
		// Computed

		if key := model.Key; key != nil {
			switch *key.Kty {
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
				switch *key.Crv {
				case keys.JsonWebKeyCurveNamePNegativeTwoFiveSix:
					publicKey.Curve = elliptic.P256()
				case keys.JsonWebKeyCurveNamePNegativeThreeEightFour:
					publicKey.Curve = elliptic.P384()
				case keys.JsonWebKeyCurveNamePNegativeFiveTwoOne:
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
		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	// Not recoverable for import.
	// d.Set("resource_id", parse.NewKeyID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name, id.Version).ID())
	// d.Set("resource_versionless_id", parse.NewKeyVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, id.Name).ID())

	respPolicy, err := client.GetKeyRotationPolicy(ctx, keys.NewKeyID(id.BaseURI, id.KeyName))
	if err != nil {
		switch {
		case response.WasForbidden(respPolicy.HttpResponse):
			// If client is not authorized to access the policy:
			return fmt.Errorf("current client lacks permissions to update Key Rotation Policy for Key %s, please update this as described here: %s : %v", *id, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
		case !response.WasNotFound(respPolicy.HttpResponse):
			if model := respPolicy.Model; model != nil {
				rotationPolicy := flattenKeyVaultKeyRotationPolicy(*respPolicy.Model)
				if err := d.Set("rotation_policy", rotationPolicy); err != nil {
					return fmt.Errorf("setting Key Vault Key Rotation Policy: %+v", err)
				}
			}
		default:
			return err
		}
	}

	return nil
}

func resourceKeyVaultKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keys.ParseKeyversionID(d.Id())
	keyId := keys.NewKeyID(id.BaseURI, id.KeyName)
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneClient.Keys.Clone(id.BaseURI)

	resp, err := client.DeleteKey(ctx, keyId)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if !response.WasNotFound(resp.HttpResponse) {
		// TODO - This needs a custom poller, maybe a generic one?
		log.Printf("[STEBUG] waiting for delete of %s: %q", *id, time.Now().String())
		time.Sleep(30 * time.Second)
	}

	if meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedKeysOnDestroy {
		deletedId := keys.NewDeletedkeyID(id.BaseURI, id.KeyName)
		if _, err := client.PurgeDeletedKey(ctx, deletedId); err != nil {
			// TODO - response codes / handling of purge protected items?
			return fmt.Errorf("purging %s: %+v", *id, err)
		}
	}

	return nil
}

func expandKeyVaultKeyOptions(d *pluginsdk.ResourceData) *[]keys.JsonWebKeyOperation {
	options := d.Get("key_opts").([]interface{})
	results := make([]keys.JsonWebKeyOperation, 0, len(options))

	for _, option := range options {
		results = append(results, keys.JsonWebKeyOperation(option.(string)))
	}

	return &results
}

func expandKeyVaultKeyRotationPolicy(v []interface{}) keys.KeyRotationPolicy {
	if len(v) == 0 {
		return keys.KeyRotationPolicy{LifetimeActions: &[]keys.LifetimeActions{}}
	}

	policy := v[0].(map[string]interface{})

	var expiryTime *string = nil // needs to be set to nil if not set
	if rawExpiryTime := policy["expire_after"]; rawExpiryTime != nil && rawExpiryTime.(string) != "" {
		expiryTime = pointer.To(rawExpiryTime.(string))
	}

	lifetimeActions := make([]keys.LifetimeActions, 0)
	if rawNotificationTime := policy["notify_before_expiry"]; rawNotificationTime != nil && rawNotificationTime.(string) != "" {
		lifetimeActionNotify := keys.LifetimeActions{
			Trigger: &keys.LifetimeActionsTrigger{
				TimeBeforeExpiry: pointer.To(rawNotificationTime.(string)), // for Type: keyvault.Notify always TimeBeforeExpiry
			},
			Action: &keys.LifetimeActionsType{
				Type: pointer.To(keys.KeyRotationPolicyActionNotify),
			},
		}
		lifetimeActions = append(lifetimeActions, lifetimeActionNotify)
	}

	if autoRotationList := policy["automatic"].([]interface{}); len(autoRotationList) == 1 && autoRotationList[0] != nil {
		lifetimeActionRotate := keys.LifetimeActions{
			Action: &keys.LifetimeActionsType{
				Type: pointer.To(keys.KeyRotationPolicyActionRotate),
			},
			Trigger: &keys.LifetimeActionsTrigger{},
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

	return keys.KeyRotationPolicy{
		LifetimeActions: &lifetimeActions,
		Attributes: &keys.KeyRotationPolicyAttributes{
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

func flattenKeyVaultKeyRotationPolicy(input keys.KeyRotationPolicy) []interface{} {
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

			if action != nil && trigger != nil && action.Type != nil && *action.Type == keys.KeyRotationPolicyActionNotify && trigger.TimeBeforeExpiry != nil && *trigger.TimeBeforeExpiry != "" {
				// Somehow a default is set after creation for notify_before_expiry
				// Submitting this set value in the next run will not work though..
				if policy["expire_after"] != nil {
					policy["notify_before_expiry"] = *trigger.TimeBeforeExpiry
				}
			}

			if action != nil && trigger != nil && action.Type != nil && *action.Type == keys.KeyRotationPolicyActionRotate {
				autoRotation := make(map[string]interface{})
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
