// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume_snapshot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

type ElasticSANVolumeSnapshotDataSourceModel struct {
	Name                  string `tfschema:"name"`
	SourceId              string `tfschema:"source_id"`
	SourceVolumeSizeInGiB int64  `tfschema:"source_volume_size_in_gib"`
	VolumeGroupId         string `tfschema:"volume_group_id"`
	VolumeName            string `tfschema:"volume_name"`
}

func (r DataSource) ResourceType() string {
	return "azurerm_elastic_san_volume_snapshot"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about an existing Elastic SAN Volume Snapshot."
}

func (r DataSource) ModelObject() any {
	return &ElasticSANVolumeSnapshotDataSourceModel{}
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Elastic SAN Volume Snapshot.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ElasticSanSnapshotName,
		},

		"volume_group_id": {
			Description:  "The Elastic SAN Volume Group ID within which the Elastic SAN Volume Snapshot exists.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: snapshots.ValidateVolumeGroupID,
		},
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"source_id": {
			Description: "The resource ID from which the Snapshot is created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"source_volume_size_in_gib": {
			Description: "The size of source volume.",
			Computed:    true,
			Type:        pluginsdk.TypeInt,
		},

		"volume_name": {
			Description: "The source volume name of the Snapshot.",
			Computed:    true,
			Type:        pluginsdk.TypeString,
		},
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Snapshots

			var state ElasticSANVolumeSnapshotDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			volumeGroupId, err := snapshots.ParseVolumeGroupID(state.VolumeGroupId)
			if err != nil {
				return err
			}

			id := snapshots.NewSnapshotID(volumeGroupId.SubscriptionId, volumeGroupId.ResourceGroupName, volumeGroupId.ElasticSanName, volumeGroupId.VolumeGroupName, state.Name)

			resp, err := client.VolumeSnapshotsGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.VolumeGroupId = volumeGroupId.ID()
			state.Name = id.SnapshotName
			if model := resp.Model; model != nil {
				// these properties are not pointer so we don't need to check for nil
				state.SourceVolumeSizeInGiB = pointer.From(model.Properties.SourceVolumeSizeGiB)
				state.VolumeName = pointer.From(model.Properties.VolumeName)

				// only ElasticSAN Volumes are supported for now
				volumeId, err := volumes.ParseVolumeIDInsensitively(model.Properties.CreationData.SourceId)
				if err != nil {
					return err
				}

				state.SourceId = volumeId.ID()
			}
			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
