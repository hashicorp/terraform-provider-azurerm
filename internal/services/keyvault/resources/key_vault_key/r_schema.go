// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_key

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/keys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func keyVaultKeySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Key Vault Key.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemName,
		},

		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}, "Key Vault"),

		"key_type": {
			Description: "The Key Type to use for this Key Vault Key.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			// turns out Azure's *really* sensitive about the casing of these
			// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
			ValidateFunc: validation.StringInSlice([]string{
				string(keys.JsonWebKeyTypeEC),
				string(keys.JsonWebKeyTypeECNegativeHSM),
				string(keys.JsonWebKeyTypeRSA),
				string(keys.JsonWebKeyTypeRSANegativeHSM),
			}, false),
		},

		"key_size": {
			Description:   "The Size of the RSA key to create in bytes. For example, 1024 or 2048. *Note*: This field is required if `key_type` is `RSA` or `RSA-HSM`.",
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"curve"},
		},

		"key_opts": {
			Description: "A list of JSON web key operations. These values are case-sensitive.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				// turns out Azure's *really* sensitive about the casing of these
				// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
				ValidateFunc: validation.StringInSlice([]string{
					string(keys.JsonWebKeyOperationDecrypt),
					string(keys.JsonWebKeyOperationEncrypt),
					string(keys.JsonWebKeyOperationImport),
					string(keys.JsonWebKeyOperationSign),
					string(keys.JsonWebKeyOperationUnwrapKey),
					string(keys.JsonWebKeyOperationVerify),
					string(keys.JsonWebKeyOperationWrapKey),
				}, false),
			},
		},

		"curve": {
			Description: "The curve to use when creating an `EC` key. This field will be required in a future release if `key_type` is `EC` or `EC-HSM`. The API will.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			ValidateFunc: func() pluginsdk.SchemaValidateFunc {
				out := []string{
					string(keys.JsonWebKeyCurveNamePNegativeTwoFiveSix),
					string(keys.JsonWebKeyCurveNamePNegativeTwoFiveSixK),
					string(keys.JsonWebKeyCurveNamePNegativeThreeEightFour),
					string(keys.JsonWebKeyCurveNamePNegativeFiveTwoOne),
				}
				return validation.StringInSlice(out, false)
			}(),
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				return old == "SECP256K1" && new == string(keys.JsonWebKeyCurveNamePNegativeTwoFiveSixK)
			},
			// TODO: the curve name should probably be mandatory for EC in the future,
			// but handle the diff so that we don't break existing configurations and
			// imported EC keys
			ConflictsWith: []string{"key_size"},
		},

		"not_before_date": {
			Description:  "Key not usable before the provided UTC datetime (Y-m-d'T'H:M:S'Z').",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"expiration_date": {
			Description:  "Expiration UTC datetime (Y-m-d'T'H:M:S'Z').",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"release_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"json": {
						Description:      "The policy contents in JSON format.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
					},

					"immutable": {
						Description: "Whether this policy is immutable.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"rotation_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expire_after": {
						Description:  "Expire a Key Vault Key after given duration as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.ISO8601DurationBetween("P28D", "P100Y"),
						RequiredWith: []string{
							"rotation_policy.0.expire_after",
							"rotation_policy.0.notify_before_expiry",
						},
						AtLeastOneOf: []string{
							"rotation_policy.0.expire_after",
							"rotation_policy.0.automatic",
						},
					},

					// <= expiry_time - 7, >=7
					"notify_before_expiry": {
						Description:  "Notify at a given duration before expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.ISO8601DurationBetween("P7D", "P36493D"),
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
									Description:  "Rotate automatically at a duration after create as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.ISO8601Duration,
									AtLeastOneOf: []string{
										"rotation_policy.0.automatic.0.time_after_creation",
										"rotation_policy.0.automatic.0.time_before_expiry",
									},
								},
								"time_before_expiry": {
									Description:  "Rotate automatically at a duration before expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.ISO8601Duration,
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
			Description: "The current version of the Key Vault Key.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"versionless_id": {
			Description: "The Base ID of the Key Vault Key.",
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

		"tags": commonschema.Tags(),
	}
}
