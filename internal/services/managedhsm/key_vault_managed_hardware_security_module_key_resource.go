// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	clientPackage "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	hsmValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
	"golang.org/x/crypto/ssh"
)

type AutomaticModel struct {
	DurationAfterCreation string `tfschema:"duration_after_creation"`
	DurationBeforeExpiry  string `tfschema:"duration_before_expiry"`
}
type RotationPolicyModel struct {
	Automatic           []AutomaticModel `tfschema:"automatic"`
	ExpireAfterDuration string           `tfschema:"expire_after_duration"`
	// NotifyBeforeExpiry string           `tfschema:"notify_before_expiry"`
}
type KeyVaultManagedHardwareSecurityModuleKeyModel struct {
	Curve                 string                `tfschema:"curve"`
	E                     string                `tfschema:"e"`
	ExpirationDate        string                `tfschema:"expiration_date"`
	KeyOptions            []string              `tfschema:"key_options"`
	KeySize               int                   `tfschema:"key_size"`
	KeyType               string                `tfschema:"key_type"`
	ManagedHsmId          string                `tfschema:"managed_hsm_id"`
	N                     string                `tfschema:"n"`
	Name                  string                `tfschema:"name"`
	NotUsableBeforeDate   string                `tfschema:"not_usable_before_date"`
	PublicKeyOpenssh      string                `tfschema:"public_key_openssh"`
	PublicKeyPem          string                `tfschema:"public_key_pem"`
	ResourceId            string                `tfschema:"resource_id"`
	ResourceVersionlessId string                `tfschema:"resource_versionless_id"`
	RotationPolicy        []RotationPolicyModel `tfschema:"rotation_policy"`
	Tags                  map[string]string     `tfschema:"tags"`
	Version               string                `tfschema:"version"`
	VersionlessId         string                `tfschema:"versionless_id"`
	X                     string                `tfschema:"x"`
	Y                     string                `tfschema:"y"`
}

type KeyVaultManagedHardwareSecurityModuleKeyResouece struct {
}

// CustomizeDiff implements sdk.ResourceWithCustomizeDiff.
func (KeyVaultManagedHardwareSecurityModuleKeyResouece) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 30,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			if meta.ResourceData == nil {
				return nil
			}

			oldRaw, newRaw := meta.ResourceData.GetChange("expiration_date")

			// Parse old and new expiration dates
			oldDate, err1 := time.Parse(time.RFC3339, oldRaw.(string))
			newDate, err2 := time.Parse(time.RFC3339, newRaw.(string))
			if err1 == nil || err2 == nil {
				// Compare old and new expiration dates
				if newDate.Before(oldDate) {
					// If the new expiration date is before te old, force new resource
					return meta.ResourceDiff.ForceNew("expiration_date")
				}
			}

			return nil
		},
	}
}

// CustomImporter implements sdk.ResourceWithCustomImporter.
func (KeyVaultManagedHardwareSecurityModuleKeyResouece) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, meta sdk.ResourceMetaData) error {
		id, err := parse.ParseNestedItemID(meta.ResourceData.Id())
		if err != nil {
			return err
		}

		subscriptionResourceId := commonids.NewSubscriptionID(meta.Client.Account.SubscriptionId)
		mHSMID, err := meta.Client.ManagedHSMs.ManagedHSMIDFromBaseUri(ctx, subscriptionResourceId, id.HSMBaseUrl)
		if err != nil {
			return fmt.Errorf("retrieving the Resource ID the managed HSM at URL %q: %s", id.HSMBaseUrl, err)
		}
		meta.ResourceData.Set("managed_hsm_id", mHSMID.ID())
		return nil
	}
}

func (KeyVaultManagedHardwareSecurityModuleKeyResouece) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hsmValidate.NestedItemName,
		},

		"managed_hsm_id": commonschema.ResourceIDReferenceRequiredForceNew(&managedhsms.ManagedHSMId{}),

		"key_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// turns out Azure's *really* sensitive about the casing of these
			// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.JSONWebKeyTypeECHSM),
				string(keyvault.JSONWebKeyTypeRSAHSM),
				string(keyvault.JSONWebKeyTypeOctHSM),
			}, false),
		},

		"key_size": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.IntBetween(0, 4096),
			ConflictsWith: []string{"curve"},
		},

		"key_options": {
			// API Response order not stable
			Type:     pluginsdk.TypeSet,
			Set:      pluginsdk.HashString,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				// turns out Azure's *really* sensitive about the casing of these
				// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.JSONWebKeyOperationDecrypt),
					string(keyvault.JSONWebKeyOperationEncrypt),
					string(keyvault.JSONWebKeyOperationImport),
					string(keyvault.JSONWebKeyOperationExport),
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
			DiffSuppressFunc: func(k, old, newVal string, d *pluginsdk.ResourceData) bool {
				return old == "SECP256K1" && newVal == string(keyvault.JSONWebKeyCurveNameP256K)
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

		"not_usable_before_date": {
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
					"expire_after_duration": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ISO8601DurationBetween("P28D", "P100Y"),
						AtLeastOneOf: []string{
							"rotation_policy.0.expire_after_duration",
							"rotation_policy.0.automatic",
						},
						RequiredWith: []string{
							"rotation_policy.0.expire_after_duration",
						},
					},

					"automatic": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"duration_after_creation": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.ISO8601Duration,
									AtLeastOneOf: []string{
										"rotation_policy.0.automatic.0.duration_after_creation",
										"rotation_policy.0.automatic.0.duration_before_expiry",
									},
								},
								"duration_before_expiry": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.ISO8601Duration,
									AtLeastOneOf: []string{
										"rotation_policy.0.automatic.0.duration_after_creation",
										"rotation_policy.0.automatic.0.duration_before_expiry",
									},
								},
							},
						},
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}

func (KeyVaultManagedHardwareSecurityModuleKeyResouece) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
	}
}

func (KeyVaultManagedHardwareSecurityModuleKeyResouece) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_key"
}

func (KeyVaultManagedHardwareSecurityModuleKeyResouece) IDValidationFunc() func(interface{}, string) ([]string, []error) {
	return hsmValidate.NestedItemId
}

func (KeyVaultManagedHardwareSecurityModuleKeyResouece) ModelObject() interface{} {
	return &KeyVaultManagedHardwareSecurityModuleKeyModel{}
}

func (k KeyVaultManagedHardwareSecurityModuleKeyResouece) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 30,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			managedHSMClient := meta.Client.ManagedHSMs
			client := managedHSMClient.DataPlaneManagedHSMClient

			var model KeyVaultManagedHardwareSecurityModuleKeyModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			name := model.Name
			managedHSMID, err := managedhsms.ParseManagedHSMID(model.ManagedHsmId)
			if err != nil {
				return err
			}

			managedHSMBaseUri, err := managedHSMClient.BaseUriForManagedHSM(ctx, *managedHSMID)
			if err != nil {
				return fmt.Errorf("looking up Key %q vault url from id %q: %+v", name, *managedHSMID, err)
			}

			existing, err := client.GetKey(ctx, *managedHSMBaseUri, name, "")
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing Key %q (Managed HSM %q): %s", name, *managedHSMBaseUri, err)
				}
			}

			if existing.Key != nil && existing.Key.Kid != nil && *existing.Key.Kid != "" {
				return tf.ImportAsExistsError("azurerm_key_vault_key", *existing.Key.Kid)
			}

			keyType := model.KeyType
			keyOptions := k.expandManagedHSMKeyOptions(&model)

			// TODO: support Importing Keys once this is fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/1747
			parameters := keyvault.KeyCreateParameters{
				Kty:    keyvault.JSONWebKeyType(keyType),
				KeyOps: keyOptions,
				KeyAttributes: &keyvault.KeyAttributes{
					Enabled: utils.Bool(true),
				},

				Tags: expandTags(model.Tags),
			}

			if parameters.Kty == keyvault.JSONWebKeyTypeEC || parameters.Kty == keyvault.JSONWebKeyTypeECHSM {
				curveName := model.Curve
				parameters.Curve = keyvault.JSONWebKeyCurveName(curveName)
			} else if parameters.Kty == keyvault.JSONWebKeyTypeRSA || parameters.Kty == keyvault.JSONWebKeyTypeRSAHSM {
				parameters.KeySize = utils.Int32(int32(model.KeySize))
			}

			if v := model.NotUsableBeforeDate; v != "" {
				notBeforeDate, _ := time.Parse(time.RFC3339, v) // validated by schema
				notBeforeUnixTime := date.UnixTime(notBeforeDate)
				parameters.KeyAttributes.NotBefore = &notBeforeUnixTime
			}

			if v := model.ExpirationDate; v != "" {
				expirationDate, _ := time.Parse(time.RFC3339, v) // validated by schema
				expirationUnixTime := date.UnixTime(expirationDate)
				parameters.KeyAttributes.Expires = &expirationUnixTime
			}

			if resp, err := client.CreateKey(ctx, *managedHSMBaseUri, name, parameters); err != nil {
				if meta.Client.Features.KeyVault.RecoverSoftDeletedKeys && utils.ResponseWasConflict(resp.Response) {
					recoveredKey, err := client.RecoverDeletedKey(ctx, *managedHSMBaseUri, name)
					if err != nil {
						return err
					}
					log.Printf("[DEBUG] Recovering Key %q with ID: %q", name, *recoveredKey.Key.Kid)
					if kid := recoveredKey.Key.Kid; kid != nil {
						recoveryPoller := custompollers.NewRecoverKeyPoller(*kid)
						poller := pollers.NewPoller(recoveryPoller, time.Second*30, pollers.DefaultNumberOfDroppedConnectionsToAllow)
						if err := poller.PollUntilDone(ctx); err != nil {
							return fmt.Errorf("waiting for Managed HSM Secret %q to become available: %s", name, err)
						}
					}
				} else {
					return fmt.Errorf("creating Key: %+v", err)
				}
			}

			if len(model.RotationPolicy) > 0 {
				policy := k.expandKeyVaultKeyRotationPolicy(model.RotationPolicy)
				if respPolicy, err := client.UpdateKeyRotationPolicy(ctx, *managedHSMBaseUri, name, policy); err != nil {
					if utils.ResponseWasForbidden(respPolicy.Response) {
						return fmt.Errorf("current client lacks permissions to create Key Rotation Policy for Key %q (%q, Vault url: %q), please update this as described here: %s : %v", name, *managedHSMID, *managedHSMBaseUri, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
					}
					return fmt.Errorf("creating Key Rotation Policy: %+v", err)
				}
			}

			// "" indicates the latest version
			read, err := client.GetKey(ctx, *managedHSMBaseUri, name, "")
			if err != nil {
				return err
			}

			if read.Key == nil || read.Key.Kid == nil {
				return fmt.Errorf("cannot read KeyVault Key '%s' (in Managed HSM '%s')", name, *managedHSMBaseUri)
			}
			keyId, err := parse.ParseNestedItemID(*read.Key.Kid)
			if err != nil {
				return err
			}
			meta.SetID(keyId)

			return nil
		},
	}
}

func (k KeyVaultManagedHardwareSecurityModuleKeyResouece) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 30,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			hsmsClient := meta.Client.ManagedHSMs
			client := hsmsClient.DataPlaneManagedHSMClient
			subscriptionId := meta.Client.Account.SubscriptionId

			id, err := parse.ParseNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
			managedHSMId, err := hsmsClient.ManagedHSMIDFromBaseUri(ctx, subscriptionResourceId, id.HSMBaseUrl)
			if err != nil {
				return fmt.Errorf("retrieving the Resource ID the Managed HSM at URL %q: %s", id.HSMBaseUrl, err)
			}
			if managedHSMId == nil {
				log.Printf("[DEBUG] Unable to determine the Resource ID for the Managed HSM at URL %q - removing from state!", id.HSMBaseUrl)
				return meta.MarkAsGone(id)
			}

			ok, err := hsmsClient.ManagedHSMExists(ctx, *managedHSMId)
			if err != nil {
				return fmt.Errorf("checking if Managed HSM %q for Key %q in Vault at url %q exists: %v", *managedHSMId, id.Name, id.HSMBaseUrl, err)
			}
			if !ok {
				log.Printf("[DEBUG] Key %q Managed HSM %q was not found in Managed HSM at URI %q - removing from state", id.Name, *managedHSMId, id.HSMBaseUrl)
				return meta.MarkAsGone(id)
			}

			resp, err := client.GetKey(ctx, id.HSMBaseUrl, id.Name, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					log.Printf("[DEBUG] Key %q was not found in Managed HSM at URI %q - removing from state", id.Name, id.HSMBaseUrl)
					return meta.MarkAsGone(id)
				}

				return err
			}

			var model KeyVaultManagedHardwareSecurityModuleKeyModel
			model.Name = id.Name
			model.ManagedHsmId = managedHSMId.ID()

			if key := resp.Key; key != nil {
				model.KeyType = string(key.Kty)
				model.KeyOptions = pointer.From(key.KeyOps)

				model.N = pointer.From(key.N)
				model.E = pointer.From(key.N)
				model.X = pointer.From(key.X)
				model.Y = pointer.From(key.Y)
				model.Curve = string(key.Crv)

				if key.N != nil {
					nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
					if err != nil {
						return fmt.Errorf("could not decode N: %+v", err)
					}
					model.KeySize = len(nBytes) * 8
				}
			}

			if attributes := resp.Attributes; attributes != nil {
				if v := attributes.NotBefore; v != nil {
					model.NotUsableBeforeDate = time.Time(*v).Format(time.RFC3339)
				}

				if v := attributes.Expires; v != nil {
					model.ExpirationDate = time.Time(*v).Format(time.RFC3339)
				}
			}

			// Computed
			model.Version = id.Version
			model.VersionlessId = id.VersionlessID()
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
					err = k.readPublicKey(&model, publicKey)
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
						err = k.readPublicKey(&model, publicKey)
						if err != nil {
							return fmt.Errorf("failed to read public key: %+v", err)
						}
					}
				}
			}

			model.ResourceId = parse.NewKeyID(managedHSMId.SubscriptionId, managedHSMId.ResourceGroupName, managedHSMId.ManagedHSMName, id.Name, id.Version).ID()

			model.ResourceVersionlessId = parse.NewKeyVersionlessID(managedHSMId.SubscriptionId, managedHSMId.ResourceGroupName, managedHSMId.ManagedHSMName, id.Name).ID()

			respPolicy, err := client.GetKeyRotationPolicy(ctx, id.HSMBaseUrl, id.Name)
			if err != nil {
				if utils.ResponseWasForbidden(respPolicy.Response) {

					// If client is not authorized to access the policy:
					return fmt.Errorf("current client lacks permissions to read Key Rotation Policy for Key %q (%q, Vault url: %q), please update this as described here: %s : %v", id.Name, *managedHSMId, id.HSMBaseUrl, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
				} else if !utils.ResponseWasNotFound(respPolicy.Response) {
					return err
				}
			}

			k.flattenKeyVaultKeyRotationPolicy(&model, respPolicy)
			model.Tags = flattenTags(resp.Tags)
			return meta.Encode(&model)
		},
	}
}

func (k KeyVaultManagedHardwareSecurityModuleKeyResouece) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 30,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			hsmsClient := meta.Client.ManagedHSMs

			id, err := parse.ParseNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KeyVaultManagedHardwareSecurityModuleKeyModel
			if err = meta.Decode(&model); err != nil {
				return err
			}

			managedHSMId, err := managedhsms.ParseManagedHSMID(model.ManagedHsmId)
			if err != nil {
				return err
			}

			clientPackage.AddToCache(*managedHSMId, id.HSMBaseUrl)

			ok, err := hsmsClient.ManagedHSMExists(ctx, *managedHSMId)
			if err != nil {
				return fmt.Errorf("checking if Managed HSM %q for Key %q in Vault at url %q exists: %v", *managedHSMId, id.Name, id.HSMBaseUrl, err)
			}
			if !ok {
				log.Printf("[DEBUG] Key %q Managed HSM %q was not found in Managed HSM at URI %q - removing from state", id.Name, *managedHSMId, id.HSMBaseUrl)
				return meta.MarkAsGone(id)
			}

			parameters := keyvault.KeyUpdateParameters{}

			if meta.ResourceData.HasChange("key_options") {
				keyOptions := k.expandManagedHSMKeyOptions(&model)
				parameters.KeyOps = keyOptions
			}

			if meta.ResourceData.HasChange("tags") {
				parameters.Tags = expandTags(model.Tags)
			}

			if meta.ResourceData.HasChanges("not_usable_before_date", "expiration_date") {
				parameters.KeyAttributes = &keyvault.KeyAttributes{
					Enabled: pointer.To(true),
				}
				if v := model.NotUsableBeforeDate; v != "" {
					notBeforeDate, _ := time.Parse(time.RFC3339, v) // validated by schema
					notBeforeUnixTime := date.UnixTime(notBeforeDate)
					parameters.KeyAttributes.NotBefore = &notBeforeUnixTime
				}

				if v := model.ExpirationDate; v != "" {
					expirationDate, _ := time.Parse(time.RFC3339, v) // validated by schema
					expirationUnixTime := date.UnixTime(expirationDate)
					parameters.KeyAttributes.Expires = &expirationUnixTime
				}
			}

			if _, err = hsmsClient.DataPlaneManagedHSMClient.UpdateKey(ctx, id.HSMBaseUrl, id.Name, "", parameters); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			if meta.ResourceData.HasChange("rotation_policy"); ok {
				policy := k.expandKeyVaultKeyRotationPolicy(model.RotationPolicy)
				if respPolicy, err := hsmsClient.DataPlaneManagedHSMClient.UpdateKeyRotationPolicy(ctx, id.HSMBaseUrl, id.Name, policy); err != nil {
					if utils.ResponseWasForbidden(respPolicy.Response) {
						return fmt.Errorf("current client lacks permissions to update Key Rotation Policy for Key %q (%q, Vault url: %q), please update this as described here: %s : %v", id.Name, *managedHSMId, id.HSMBaseUrl, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
					}
					return fmt.Errorf("creating Key Rotation Policy: %+v", err)
				}
			}
			return nil
		},
	}
}

// Delete implements sdk.ResourceWithUpdate.
func (KeyVaultManagedHardwareSecurityModuleKeyResouece) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 30,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			hsmsClient := meta.Client.ManagedHSMs
			client := hsmsClient.DataPlaneManagedHSMClient
			subscriptionId := meta.Client.Account.SubscriptionId

			id, err := parse.ParseNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
			managedHSMId, err := hsmsClient.ManagedHSMIDFromBaseUri(ctx, subscriptionResourceId, id.HSMBaseUrl)
			if err != nil {
				return fmt.Errorf("retrieving the Resource ID the Managed HSM at URL %q: %s", id.HSMBaseUrl, err)
			}
			if managedHSMId == nil {
				return fmt.Errorf("unable to determine the Resource ID for the Managed HSM at URL %q", id.HSMBaseUrl)
			}

			kv, err := hsmsClient.ManagedHsmClient.Get(ctx, *managedHSMId)
			if err != nil {
				if response.WasNotFound(kv.HttpResponse) {
					log.Printf("[DEBUG] Key %q Managed HSM %q was not found in Managed HSM at URI %q - removing from state", id.Name, *managedHSMId, id.HSMBaseUrl)
					return meta.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving Managed HSM %q properties: %+v", *managedHSMId, err)
			}

			shouldPurge := meta.Client.Features.KeyVault.PurgeSoftDeletedKeysOnDestroy
			if shouldPurge && kv.Model != nil && utils.NormaliseNilableBool(kv.Model.Properties.EnablePurgeProtection) {
				log.Printf("[DEBUG] cannot purge key %q because vault %q has purge protection enabled", id.Name, managedHSMId.String())
				shouldPurge = false
			}

			description := fmt.Sprintf("Key %q (Managed HSM %q)", id.Name, id.HSMBaseUrl)
			deleter := deleteAndPurgeKey{
				client:  client,
				mHSMUri: id.HSMBaseUrl,
				name:    id.Name,
			}
			if err := deleteAndOptionallyPurge(ctx, description, shouldPurge, deleter); err != nil {
				return err
			}

			return nil
		},
	}
}

var _ sdk.ResourceWithCustomImporter = KeyVaultManagedHardwareSecurityModuleKeyResouece{}
var _ sdk.ResourceWithUpdate = KeyVaultManagedHardwareSecurityModuleKeyResouece{}
var _ sdk.ResourceWithCustomizeDiff = KeyVaultManagedHardwareSecurityModuleKeyResouece{}

func (k KeyVaultManagedHardwareSecurityModuleKeyResouece) expandManagedHSMKeyOptions(d *KeyVaultManagedHardwareSecurityModuleKeyModel) *[]keyvault.JSONWebKeyOperation {
	results := make([]keyvault.JSONWebKeyOperation, 0)

	for _, option := range d.KeyOptions {
		results = append(results, keyvault.JSONWebKeyOperation(option))
	}

	return &results
}

func (k KeyVaultManagedHardwareSecurityModuleKeyResouece) expandKeyVaultKeyRotationPolicy(v []RotationPolicyModel) keyvault.KeyRotationPolicy {
	if len(v) == 0 {
		return keyvault.KeyRotationPolicy{LifetimeActions: &[]keyvault.LifetimeActions{}}
	}

	policy := v[0]

	var expiryTime *string = nil // needs to be set to nil if not set
	if policy.ExpireAfterDuration != "" {
		expiryTime = pointer.To(policy.ExpireAfterDuration)
	}

	lifetimeActions := make([]keyvault.LifetimeActions, 0)

	if len(policy.Automatic) == 1 {
		lifetimeActionRotate := keyvault.LifetimeActions{
			Action: &keyvault.LifetimeActionsType{
				Type: keyvault.ActionTypeRotate,
			},
			Trigger: &keyvault.LifetimeActionsTrigger{},
		}
		autoRotationRaw := policy.Automatic[0]

		if autoRotationRaw.DurationAfterCreation != "" {
			lifetimeActionRotate.Trigger.TimeAfterCreate = &autoRotationRaw.DurationAfterCreation
		}

		if autoRotationRaw.DurationBeforeExpiry != "" {
			lifetimeActionRotate.Trigger.TimeBeforeExpiry = &autoRotationRaw.DurationBeforeExpiry
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

func (k KeyVaultManagedHardwareSecurityModuleKeyResouece) flattenKeyVaultKeyRotationPolicy(model *KeyVaultManagedHardwareSecurityModuleKeyModel, input keyvault.KeyRotationPolicy) {
	if input.LifetimeActions == nil && (input.Attributes == nil || pointer.From(input.Attributes.ExpiryTime) == "") {
		return
	}

	var policy RotationPolicyModel
	if input.Attributes != nil {
		policy.ExpireAfterDuration = pointer.From(input.Attributes.ExpiryTime)
	}

	if input.LifetimeActions != nil {
		for _, ltAction := range *input.LifetimeActions {
			action := ltAction.Action
			trigger := ltAction.Trigger

			if action != nil && trigger != nil {
				if strings.EqualFold(string(action.Type), string(keyvault.ActionTypeRotate)) {
					var autoRotation AutomaticModel
					autoRotation.DurationAfterCreation = pointer.From(trigger.TimeAfterCreate)
					autoRotation.DurationBeforeExpiry = pointer.From(trigger.TimeBeforeExpiry)
					policy.Automatic = append(policy.Automatic, autoRotation)
				}
			}

		}
	}

	model.RotationPolicy = []RotationPolicyModel{policy}
}

var _ deleteAndPurgeNestedItem = deleteAndPurgeKey{}

type deleteAndPurgeKey struct {
	client  *keyvault.BaseClient
	mHSMUri string
	name    string
}

func (d deleteAndPurgeKey) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.DeleteKey(ctx, d.mHSMUri, d.name)
	return resp.Response, err
}

func (d deleteAndPurgeKey) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetKey(ctx, d.mHSMUri, d.name, "")
	return resp.Response, err
}

func (d deleteAndPurgeKey) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	return d.client.PurgeDeletedKey(ctx, d.mHSMUri, d.name)
}

func (d deleteAndPurgeKey) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetDeletedKey(ctx, d.mHSMUri, d.name)
	return resp.Response, err
}

func (k KeyVaultManagedHardwareSecurityModuleKeyResouece) readPublicKey(model *KeyVaultManagedHardwareSecurityModuleKeyModel, pubKey interface{}) error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key error: %s", err)
	}
	pubKeyPemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	model.PublicKeyPem = string(pem.EncodeToMemory(pubKeyPemBlock))

	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err == nil {
		// Not all EC types can be SSH keys, so we'll produce this only
		// if an appropriate type was selected.
		sshPubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)
		model.PublicKeyOpenssh = string(sshPubKeyBytes)
	} else {
		model.PublicKeyOpenssh = ""
	}
	return nil
}
