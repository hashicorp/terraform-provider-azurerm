// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Volumes
			schema := ElasticSANVolumeResourceModel{}

			id, err := volumes.ParseVolumeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			volumeGroupId := volumes.NewVolumeGroupID(id.SubscriptionId, id.ResourceGroupName, id.ElasticSanName, id.VolumeGroupName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema.Name = id.VolumeName
			schema.VolumeGroupId = volumeGroupId.ID()

			if model := resp.Model; model != nil {
				// model.Properties is not a pointer
				props := model.Properties

				schema.SizeInGiB = props.SizeGiB
				schema.VolumeId = pointer.From(props.VolumeId)
				schema.CreateSource = FlattenElasticSANVolumeCreateSource(props.CreationData)

				if storageTarget := props.StorageTarget; storageTarget != nil {
					schema.TargetIqn = pointer.From(storageTarget.TargetIqn)
					schema.TargetPortalPort = pointer.From(storageTarget.TargetPortalPort)
					schema.TargetPortalHostname = pointer.From(storageTarget.TargetPortalHostname)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

type ElasticSANVolumeCreateSource struct {
	SourceType string `tfschema:"source_type"`
	SourceId   string `tfschema:"source_id"`
}

func ExpandElasticSANVolumeCreateSource(input []ElasticSANVolumeCreateSource) *volumes.SourceCreationData {
	if len(input) == 0 {
		return nil
	}

	return &volumes.SourceCreationData{
		SourceId:     pointer.To(input[0].SourceId),
		CreateSource: pointer.ToEnum[volumes.VolumeCreateOption](input[0].SourceType),
	}
}

func FlattenElasticSANVolumeCreateSource(input *volumes.SourceCreationData) []ElasticSANVolumeCreateSource {
	// the response might return the block but with only sourceType=None in the block`
	if input == nil || input.SourceId == nil {
		return []ElasticSANVolumeCreateSource{}
	}

	return []ElasticSANVolumeCreateSource{
		{
			SourceType: pointer.FromEnum(input.CreateSource),
			SourceId:   pointer.From(input.SourceId),
		},
	}
}
