// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticSanVolumeGroupName,
		},

		"elastic_san_id": commonschema.ResourceIDReferenceRequiredForceNew(&volumegroups.ElasticSanId{}),

		"encryption_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(volumegroups.PossibleValuesForEncryptionType(), false),
			Default:      string(volumegroups.EncryptionTypeEncryptionAtRestWithPlatformKey),
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Required:     true,
						Type:         pluginsdk.TypeString,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},
					"user_assigned_identity_id": {
						Optional:     true,
						Type:         pluginsdk.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
					"current_versioned_key_expiration_timestamp": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"current_versioned_key_id": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"last_key_rotation_timestamp": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
				},
			},
		},

		"network_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Required:     true,
						Type:         pluginsdk.TypeString,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"action": {
						Optional:     true,
						Type:         pluginsdk.TypeString,
						Default:      string(volumegroups.ActionAllow),
						ValidateFunc: validation.StringInSlice(volumegroups.PossibleValuesForAction(), false),
					},
				},
			},
		},

		"protocol_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				// None is not a valid value and service team will consider removing it in future versions.
				string(volumegroups.StorageTargetTypeIscsi),
			}, false),
			Default: string(volumegroups.StorageTargetTypeIscsi),
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),
	}
}

func (k Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config ElasticSANVolumeGroupResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(config.Encryption) > 0 && config.EncryptionType != string(volumegroups.EncryptionTypeEncryptionAtRestWithCustomerManagedKey) {
				return fmt.Errorf("encryption can only be set if encryption_type is EncryptionAtRestWithCustomerManagedKey")
			}

			if len(config.Encryption) == 0 && config.EncryptionType == string(volumegroups.EncryptionTypeEncryptionAtRestWithCustomerManagedKey) {
				return fmt.Errorf("encryption must be set if encryption_type is EncryptionAtRestWithCustomerManagedKey")
			}

			return nil
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
