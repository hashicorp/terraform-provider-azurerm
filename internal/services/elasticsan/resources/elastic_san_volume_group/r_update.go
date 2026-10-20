// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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
				payload.Properties.Encryption = pointer.ToEnum[volumegroups.EncryptionType](config.EncryptionType)
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
				payload.Properties.ProtocolType = pointer.ToEnum[volumegroups.StorageTargetType](config.ProtocolType)
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
