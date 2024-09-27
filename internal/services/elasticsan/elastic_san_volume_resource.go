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
	diskSnapshots "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/restorepoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource                  = ElasticSANVolumeResource{}
	_ sdk.ResourceWithUpdate        = ElasticSANVolumeResource{}
	_ sdk.ResourceWithCustomizeDiff = ElasticSANVolumeResource{}
)

type ElasticSANVolumeResource struct{}

func (r ElasticSANVolumeResource) ModelObject() interface{} {
	return &ElasticSANVolumeResourceModel{}
}

type ElasticSANVolumeResourceModel struct {
	CreateSource         []ElasticSANVolumeCreateSource `tfschema:"create_source"`
	Name                 string                         `tfschema:"name"`
	SizeInGiB            int64                          `tfschema:"size_in_gib"`
	TargetIqn            string                         `tfschema:"target_iqn"`
	TargetPortalHostname string                         `tfschema:"target_portal_hostname"`
	TargetPortalPort     int64                          `tfschema:"target_portal_port"`
	VolumeGroupId        string                         `tfschema:"volume_group_id"`
	VolumeId             string                         `tfschema:"volume_id"`
}

type ElasticSANVolumeCreateSource struct {
	SourceType string `tfschema:"source_type"`
	SourceId   string `tfschema:"source_id"`
}

func (r ElasticSANVolumeResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumes.ValidateVolumeID
}

func (r ElasticSANVolumeResource) ResourceType() string {
	return "azurerm_elastic_san_volume"
}

func (r ElasticSANVolumeResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticSanVolumeName,
		},

		"volume_group_id": commonschema.ResourceIDReferenceRequiredForceNew(&volumes.VolumeGroupId{}),

		"size_in_gib": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65536),
		},

		"create_source": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.Any(
							commonids.ValidateManagedDiskID,
							restorepoints.ValidateRestorePointID,
							diskSnapshots.ValidateSnapshotID,
							snapshots.ValidateSnapshotID,
						),
					},
					"source_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							// None is not exposed
							string(volumes.VolumeCreateOptionDisk),
							string(volumes.VolumeCreateOptionDiskRestorePoint),
							string(volumes.VolumeCreateOptionDiskSnapshot),
							string(volumes.VolumeCreateOptionVolumeSnapshot),
						},
							false),
					},
				},
			},
		},
	}
}

func (r ElasticSANVolumeResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"target_iqn": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},

		"target_portal_hostname": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},

		"target_portal_port": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"volume_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r ElasticSANVolumeResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("size_in_gib"); newVal.(int) < oldVal.(int) {
				return fmt.Errorf("new size_in_gib should be greater than the existing one")
			}

			return nil
		},
	}
}

func (r ElasticSANVolumeResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Volumes

			var config ElasticSANVolumeResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			volumeGroupId, err := volumes.ParseVolumeGroupID(config.VolumeGroupId)
			if err != nil {
				return err
			}

			id := volumes.NewVolumeID(subscriptionId, volumeGroupId.ResourceGroupName, volumeGroupId.ElasticSanName, volumeGroupId.VolumeGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := volumes.Volume{
				Properties: volumes.VolumeProperties{
					CreationData: ExpandElasticSANVolumeCreateSource(config.CreateSource),
					SizeGiB:      config.SizeInGiB,
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

func (r ElasticSANVolumeResource) Read() sdk.ResourceFunc {
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

func (r ElasticSANVolumeResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Volumes

			id, err := volumes.ParseVolumeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, volumes.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ElasticSANVolumeResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Volumes

			id, err := volumes.ParseVolumeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ElasticSANVolumeResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := volumes.VolumeUpdate{}
			if metadata.ResourceData.HasChange("size_in_gib") {
				if payload.Properties == nil {
					payload.Properties = &volumes.VolumeUpdateProperties{}
				}
				payload.Properties.SizeGiB = pointer.To(config.SizeInGiB)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func ExpandElasticSANVolumeCreateSource(input []ElasticSANVolumeCreateSource) *volumes.SourceCreationData {
	if len(input) == 0 {
		return nil
	}

	return &volumes.SourceCreationData{
		SourceId:     pointer.To(input[0].SourceId),
		CreateSource: pointer.To(volumes.VolumeCreateOption(input[0].SourceType)),
	}
}

func FlattenElasticSANVolumeCreateSource(input *volumes.SourceCreationData) []ElasticSANVolumeCreateSource {
	// the response might return the block but with only sourceType=None in the block`
	if input == nil || input.SourceId == nil {
		return []ElasticSANVolumeCreateSource{}
	}

	return []ElasticSANVolumeCreateSource{
		{
			SourceType: string(pointer.From(input.CreateSource)),
			SourceId:   pointer.From(input.SourceId),
		},
	}
}
