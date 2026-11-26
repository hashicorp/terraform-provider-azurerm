// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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
					Encryption:           pointer.ToEnum[volumegroups.EncryptionType](config.EncryptionType),
					EncryptionProperties: encryption,
					NetworkAcls:          ExpandVolumeGroupNetworkRules(config.NetworkRule),
					ProtocolType:         pointer.ToEnum[volumegroups.StorageTargetType](config.ProtocolType),
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
