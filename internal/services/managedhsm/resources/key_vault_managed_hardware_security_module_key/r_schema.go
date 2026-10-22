// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"managed_hsm_id": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: managedhsms.ValidateManagedHSMID,
		},

		"key_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// turns out Azure's *really* sensitive about the casing of these
			// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.JSONWebKeyTypeECHSM),
				string(keyvault.JSONWebKeyTypeOctHSM),
				string(keyvault.JSONWebKeyTypeRSAHSM),
			}, false),
		},

		"key_size": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"curve"},
		},

		"curve": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
			ExactlyOneOf: []string{"key_size"},
		},

		"key_opts": {
			Type:     pluginsdk.TypeSet,
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
					string(keyvault.JSONWebKeyOperationImport),
				}, false),
			},
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

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"versioned_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff

			// There isn't a way to remove these dates from a key so we'll recreate the resource when they are removed from config
			for _, field := range []string{"not_before_date", "expiration_date"} {
				oldValue, newValue := diff.GetChange(field)
				if oldValue.(string) != "" && newValue.(string) == "" {
					if err := diff.ForceNew(field); err != nil {
						return err
					}
				}
			}

			// if any value has changed, we need to SetNewComputed on versioned_id as any change to the key is a new version
			if diff.HasChanges("key_opts", "not_before_date", "tags", "expiration_date") {
				return diff.SetNewComputed("versioned_id")
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}
