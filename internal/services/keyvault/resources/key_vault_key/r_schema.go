// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemName,
		},

		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

		"key_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
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
					string(keys.JsonWebKeyOperationDecrypt),
					string(keys.JsonWebKeyOperationEncrypt),
					string(keys.JsonWebKeyOperationSign),
					string(keys.JsonWebKeyOperationUnwrapKey),
					string(keys.JsonWebKeyOperationVerify),
					string(keys.JsonWebKeyOperationWrapKey),
				}, false),
			},
		},

		"curve": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				return old == "SECP256K1" && new == string(keys.JsonWebKeyCurveNamePNegativeTwoFiveSixK)
			},
			ValidateFunc: func() pluginsdk.SchemaValidateFunc {
				out := []string{
					string(keys.JsonWebKeyCurveNamePNegativeTwoFiveSix),
					string(keys.JsonWebKeyCurveNamePNegativeTwoFiveSixK),
					string(keys.JsonWebKeyCurveNamePNegativeThreeEightFour),
					string(keys.JsonWebKeyCurveNamePNegativeFiveTwoOne),
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

		"release_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"json": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
					},

					"immutable": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.ISO8601DurationBetween("P28D", "P100Y"),
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
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.ISO8601Duration,
									AtLeastOneOf: []string{
										"rotation_policy.0.automatic.0.time_after_creation",
										"rotation_policy.0.automatic.0.time_before_expiry",
									},
								},
								"time_before_expiry": {
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
	}
}
