// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume_group

import (
	"context"
	"errors"
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
			Description:  "The name of this Elastic SAN Volume Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticSanVolumeGroupName,
		},

		"elastic_san_id": commonschema.ResourceIDReferenceRequiredForceNew(&volumegroups.ElasticSanId{}, "Elastic SAN"),

		"encryption_type": {
			Description:  "The type of the key used to encrypt the data of the disk.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(volumegroups.EncryptionTypeEncryptionAtRestWithPlatformKey),
			ValidateFunc: validation.StringInSlice(volumegroups.PossibleValuesForEncryptionType(), false),
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Description:  "The Key Vault Key URI for Customer Managed Key encryption, which can be either a full URI or a versionless URI.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},
					"user_assigned_identity_id": {
						Description:  "The ID of the User Assigned Identity used by this Elastic SAN Volume Group.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
					"current_versioned_key_expiration_timestamp": {
						Description: "The timestamp of the expiration time for the current version of the customer managed key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"current_versioned_key_id": {
						Description: "The ID of the current versioned Key Vault Key in use.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"last_key_rotation_timestamp": {
						Description: "The timestamp of the last rotation of the Key Vault Key.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
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
						Description:  "The ID of the Subnet which should be allowed to access this Elastic SAN Volume Group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"action": {
						Description:  "The action to take when the Subnet attempts to access this Elastic SAN Volume Group.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(volumegroups.ActionAllow),
						ValidateFunc: validation.StringInSlice(volumegroups.PossibleValuesForAction(), false),
					},
				},
			},
		},

		"protocol_type": {
			Description: "The type of the storage target.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(volumegroups.StorageTargetTypeIscsi),
			ValidateFunc: validation.StringInSlice([]string{
				// None is not a valid value and service team will consider removing it in future versions.
				string(volumegroups.StorageTargetTypeIscsi),
			}, false),
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
				return errors.New("encryption can only be set if encryption_type is EncryptionAtRestWithCustomerManagedKey")
			}

			if len(config.Encryption) == 0 && config.EncryptionType == string(volumegroups.EncryptionTypeEncryptionAtRestWithCustomerManagedKey) {
				return errors.New("encryption must be set if encryption_type is EncryptionAtRestWithCustomerManagedKey")
			}

			return nil
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
