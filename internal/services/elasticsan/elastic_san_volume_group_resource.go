// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ElasticSANVolumeGroupResource{}
var _ sdk.ResourceWithUpdate = ElasticSANVolumeGroupResource{}
var _ sdk.ResourceWithCustomizeDiff = ElasticSANVolumeGroupResource{}

type ElasticSANVolumeGroupResource struct{}

func (r ElasticSANVolumeGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumegroups.ValidateVolumeGroupID
}

func (r ElasticSANVolumeGroupResource) ResourceType() string {
	return "azurerm_elastic_san_volume_group"
}

func (r ElasticSANVolumeGroupResource) ModelObject() interface{} {
	return &ElasticSANVolumeGroupResourceModel{}
}

type ElasticSANVolumeGroupResourceModel struct {
	SanId          string                                          `tfschema:"elastic_san_id"`
	EncryptionType string                                          `tfschema:"encryption_type"`
	Encryption     []ElasticSANVolumeGroupResourceEncryptionModel  `tfschema:"encryption"`
	Identity       []identity.ModelSystemAssignedUserAssigned      `tfschema:"identity"`
	Name           string                                          `tfschema:"name"`
	NetworkRule    []ElasticSANVolumeGroupResourceNetworkRuleModel `tfschema:"network_rule"`
	ProtocolType   string                                          `tfschema:"protocol_type"`
}

type ElasticSANVolumeGroupResourceEncryptionModel struct {
	CurrentVersionedKeyExpirationTimestamp string `tfschema:"current_versioned_key_expiration_timestamp"`
	CurrentVersionedKeyId                  string `tfschema:"current_versioned_key_id"`
	UserAssignedIdentityId                 string `tfschema:"user_assigned_identity_id"`
	KeyVaultKeyId                          string `tfschema:"key_vault_key_id"`
	LastKeyRotationTimestamp               string `tfschema:"last_key_rotation_timestamp"`
}

type ElasticSANVolumeGroupResourceNetworkRuleModel struct {
	Action   string `tfschema:"action"`
	SubnetId string `tfschema:"subnet_id"`
}

func (r ElasticSANVolumeGroupResource) Arguments() map[string]*pluginsdk.Schema {
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
						ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
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

func (k ElasticSANVolumeGroupResource) CustomizeDiff() sdk.ResourceFunc {
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

func (r ElasticSANVolumeGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ElasticSANVolumeGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.VolumeGroups

			var config ElasticSANVolumeGroupResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			elasticSanId, err := volumegroups.ParseElasticSanID(config.SanId)
			if err != nil {
				return err
			}

			id := volumegroups.NewVolumeGroupID(subscriptionId, elasticSanId.ResourceGroupName, elasticSanId.ElasticSanName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			expandedIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %+v", err)
			}

			encryption, err := ExpandVolumeGroupEncryption(config.Encryption)
			if err != nil {
				return fmt.Errorf("expanding encryption: %+v", err)
			}

			payload := volumegroups.VolumeGroup{
				Identity: expandedIdentity,
				Properties: &volumegroups.VolumeGroupProperties{
					Encryption:           pointer.To(volumegroups.EncryptionType(config.EncryptionType)),
					EncryptionProperties: encryption,
					NetworkAcls:          ExpandVolumeGroupNetworkRules(config.NetworkRule),
					ProtocolType:         pointer.To(volumegroups.StorageTargetType(config.ProtocolType)),
				},
			}

			if err := client.CreateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ElasticSANVolumeGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.VolumeGroups
			schema := ElasticSANVolumeGroupResourceModel{}

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			elasticSanId := volumegroups.NewElasticSanID(id.SubscriptionId, id.ResourceGroupName, id.ElasticSanName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.SanId = elasticSanId.ID()
				schema.Name = id.VolumeGroupName

				flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening identity: %+v", err)
				}
				schema.Identity = *flattenedIdentity

				if model.Properties != nil {
					schema.EncryptionType = string(pointer.From(model.Properties.Encryption))
					schema.NetworkRule = FlattenVolumeGroupNetworkRules(model.Properties.NetworkAcls)

					if model.Properties.ProtocolType != nil {
						schema.ProtocolType = string(pointer.From(model.Properties.ProtocolType))
					}

					schema.Encryption, err = FlattenVolumeGroupEncryption(model.Properties.EncryptionProperties)
					if err != nil {
						return fmt.Errorf("flattening encryption: %+v", err)
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ElasticSANVolumeGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.VolumeGroups

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ElasticSANVolumeGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.VolumeGroups

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ElasticSANVolumeGroupResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := volumegroups.VolumeGroupUpdate{
				Properties: &volumegroups.VolumeGroupUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("encryption_type") {
				payload.Properties.Encryption = pointer.To(volumegroups.EncryptionType(config.EncryptionType))
			}

			if metadata.ResourceData.HasChange("encryption") {
				encryption, err := ExpandVolumeGroupEncryption(config.Encryption)
				if err != nil {
					return fmt.Errorf("expanding encryption: %+v", err)
				}

				payload.Properties.EncryptionProperties = encryption
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding identity: %+v", err)
				}

				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("protocol_type") {
				payload.Properties.ProtocolType = pointer.To(volumegroups.StorageTargetType(config.ProtocolType))
			}

			if metadata.ResourceData.HasChange("network_rule") {
				payload.Properties.NetworkAcls = ExpandVolumeGroupNetworkRules(config.NetworkRule)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func ExpandVolumeGroupEncryption(input []ElasticSANVolumeGroupResourceEncryptionModel) (*volumegroups.EncryptionProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}

	nestedItemId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(input[0].KeyVaultKeyId)
	if err != nil {
		return nil, err
	}

	result := volumegroups.EncryptionProperties{
		KeyVaultProperties: &volumegroups.KeyVaultProperties{
			KeyName:     pointer.To(nestedItemId.Name),
			KeyVersion:  pointer.To(nestedItemId.Version),
			KeyVaultUri: pointer.To(nestedItemId.KeyVaultBaseUrl),
		},
	}

	if input[0].UserAssignedIdentityId != "" {
		result.Identity = &volumegroups.EncryptionIdentity{
			UserAssignedIdentity: pointer.To(input[0].UserAssignedIdentityId),
		}
	}

	return &result, nil
}

func FlattenVolumeGroupEncryption(input *volumegroups.EncryptionProperties) ([]ElasticSANVolumeGroupResourceEncryptionModel, error) {
	if input == nil {
		return []ElasticSANVolumeGroupResourceEncryptionModel{}, nil
	}

	var keyVaultKeyId, currentVersionedKeyExpirationTimestamp, currentVersionedKeyId, lastKeyRotationTimestamp string
	if kv := input.KeyVaultProperties; kv != nil {
		id, err := keyVaultParse.NewNestedItemID(pointer.From(kv.KeyVaultUri), keyVaultParse.NestedItemTypeKey, pointer.From(kv.KeyName), pointer.From(kv.KeyVersion))
		if err != nil {
			return nil, fmt.Errorf("parsing Encryption Key Vault Key ID: %+v", err)
		}

		keyVaultKeyId = id.ID()

		currentVersionedKeyExpirationTimestamp = pointer.From(input.KeyVaultProperties.CurrentVersionedKeyExpirationTimestamp)
		currentVersionedKeyId = pointer.From(input.KeyVaultProperties.CurrentVersionedKeyIdentifier)
		lastKeyRotationTimestamp = pointer.From(input.KeyVaultProperties.LastKeyRotationTimestamp)
	}

	var userAssignedIdentityId string
	if input.Identity != nil && input.Identity.UserAssignedIdentity != nil {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity.UserAssignedIdentity)
		if err != nil {
			return nil, fmt.Errorf("parsing Encryption User Assigned Identity ID: %+v", err)
		}

		userAssignedIdentityId = id.ID()
	}

	return []ElasticSANVolumeGroupResourceEncryptionModel{
		{
			KeyVaultKeyId:                          keyVaultKeyId,
			UserAssignedIdentityId:                 userAssignedIdentityId,
			CurrentVersionedKeyExpirationTimestamp: currentVersionedKeyExpirationTimestamp,
			CurrentVersionedKeyId:                  currentVersionedKeyId,
			LastKeyRotationTimestamp:               lastKeyRotationTimestamp,
		},
	}, nil
}

func ExpandVolumeGroupNetworkRules(input []ElasticSANVolumeGroupResourceNetworkRuleModel) *volumegroups.NetworkRuleSet {
	// if return nil, the Network Rules will not be removed during update
	if len(input) == 0 {
		return &volumegroups.NetworkRuleSet{
			VirtualNetworkRules: &[]volumegroups.VirtualNetworkRule{},
		}
	}

	var networkRules []volumegroups.VirtualNetworkRule
	for _, rule := range input {
		networkRules = append(networkRules, volumegroups.VirtualNetworkRule{
			Id:     rule.SubnetId,
			Action: pointer.To(volumegroups.Action(rule.Action)),
		})
	}

	return &volumegroups.NetworkRuleSet{
		VirtualNetworkRules: &networkRules,
	}
}

func FlattenVolumeGroupNetworkRules(input *volumegroups.NetworkRuleSet) []ElasticSANVolumeGroupResourceNetworkRuleModel {
	if input == nil || input.VirtualNetworkRules == nil {
		return []ElasticSANVolumeGroupResourceNetworkRuleModel{}
	}

	networkRules := make([]ElasticSANVolumeGroupResourceNetworkRuleModel, 0)
	for _, rule := range *input.VirtualNetworkRules {
		networkRules = append(networkRules, ElasticSANVolumeGroupResourceNetworkRuleModel{
			SubnetId: rule.Id,
			Action:   string(pointer.From(rule.Action)),
		})
	}

	return networkRules
}
