// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
					schema.EncryptionType = pointer.FromEnum(model.Properties.Encryption)
					schema.NetworkRule = FlattenVolumeGroupNetworkRules(model.Properties.NetworkAcls)

					if model.Properties.ProtocolType != nil {
						schema.ProtocolType = pointer.FromEnum(model.Properties.ProtocolType)
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

func ExpandVolumeGroupEncryption(input []ElasticSANVolumeGroupResourceEncryptionModel) (*volumegroups.EncryptionProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}

	nestedItemId, err := keyvault.ParseNestedItemID(input[0].KeyVaultKeyId, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return nil, err
	}

	result := volumegroups.EncryptionProperties{
		KeyVaultProperties: &volumegroups.KeyVaultProperties{
			KeyName:     pointer.To(nestedItemId.Name),
			KeyVersion:  pointer.To(nestedItemId.Version),
			KeyVaultUri: pointer.To(nestedItemId.KeyVaultBaseURL),
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
		id, err := keyvault.NewNestedItemID(pointer.From(kv.KeyVaultUri), keyvault.NestedItemTypeKey, pointer.From(kv.KeyName), pointer.From(kv.KeyVersion))
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

	networkRules := make([]volumegroups.VirtualNetworkRule, 0, len(input))
	for _, rule := range input {
		networkRules = append(networkRules, volumegroups.VirtualNetworkRule{
			Id:     rule.SubnetId,
			Action: pointer.ToEnum[volumegroups.Action](rule.Action),
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
			Action:   pointer.FromEnum(rule.Action),
		})
	}

	return networkRules
}
