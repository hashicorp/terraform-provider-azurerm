package elasticsan

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

var _ sdk.Resource = ElasticSANSnapshotResource{}

type ElasticSANSnapshotResource struct{}

func (r ElasticSANSnapshotResource) ModelObject() interface{} {
	return &ElasticSANSnapshotResourceModel{}
}

type ElasticSANSnapshotResourceModel struct {
	Name                string                           `tfschema:"name"`
	SourceVolumeSizeGiB int64                            `tfschema:"source_volume_size_gib"`
	CreateSource        []ElasticSANSnapshotCreateSource `tfschema:"creation_source"`
	VolumeGroupId       string                           `tfschema:"volume_group_id"`
	VolumeName          string                           `tfschema:"volume_name"`
}

type ElasticSANSnapshotCreateSource struct {
	SourceId string `tfschema:"source_id"`
}

func (r ElasticSANSnapshotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return snapshots.ValidateSnapshotID
}

func (r ElasticSANSnapshotResource) ResourceType() string {
	return "azurerm_elastic_san_snapshot"
}

func (r ElasticSANSnapshotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ElasticSanSnapshotName,
		},
		"volume_group_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"creation_source": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: volumes.ValidateVolumeID,
					},
				},
			},
		},
	}
}

func (r ElasticSANSnapshotResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"source_volume_size_gib": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},
		"volume_name": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r ElasticSANSnapshotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Snapshots

			var config ElasticSANSnapshotResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			volumeGroupId, err := snapshots.ParseVolumeGroupID(config.VolumeGroupId)
			if err != nil {
				return err
			}

			id := snapshots.NewSnapshotID(subscriptionId, volumeGroupId.ResourceGroupName, volumeGroupId.ElasticSanName, volumeGroupId.VolumeGroupName, config.Name)

			existing, err := client.VolumeSnapshotsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := snapshots.Snapshot{
				Properties: snapshots.SnapshotProperties{
					CreationData: ExpandElasticSANSnapshotCreateSource(config.CreateSource),
				},
			}

			if err := client.VolumeSnapshotsCreateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ElasticSANSnapshotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Snapshots
			schema := ElasticSANSnapshotResourceModel{}

			id, err := snapshots.ParseSnapshotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			volumeGroupId := snapshots.NewVolumeGroupID(id.SubscriptionId, id.ResourceGroupName, id.ElasticSanName, id.VolumeGroupName)

			resp, err := client.VolumeSnapshotsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema.Name = id.SnapshotName
			schema.VolumeGroupId = volumeGroupId.ID()
			if model := resp.Model; model != nil {
				// these properties are not pointer so we don't need to check for nil
				schema.SourceVolumeSizeGiB = pointer.From(model.Properties.SourceVolumeSizeGiB)
				schema.VolumeName = pointer.From(model.Properties.VolumeName)

				createSource, err := FlattenElasticSANSnapshotCreateSource(model.Properties.CreationData)
				if err != nil {
					return err
				}
				schema.CreateSource = *createSource
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ElasticSANSnapshotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Snapshots

			id, err := snapshots.ParseSnapshotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.VolumeSnapshotsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func ExpandElasticSANSnapshotCreateSource(input []ElasticSANSnapshotCreateSource) snapshots.SnapshotCreationData {
	if len(input) == 0 {
		return snapshots.SnapshotCreationData{}
	}

	v := input[0]
	return snapshots.SnapshotCreationData{
		SourceId: v.SourceId,
	}
}

func FlattenElasticSANSnapshotCreateSource(input snapshots.SnapshotCreationData) (*[]ElasticSANSnapshotCreateSource, error) {
	if input.SourceId == "" {
		return &[]ElasticSANSnapshotCreateSource{}, nil
	}

	// for now source ID can only be ElasticSAN Volume ID
	var sourceId string
	parsedSourceId, err := volumes.ParseVolumeIDInsensitively(input.SourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing source ID as ElasticSAN Volume ID: %+v", err)
	}
	sourceId = parsedSourceId.ID()

	return &[]ElasticSANSnapshotCreateSource{
		{
			SourceId: sourceId,
		},
	}, nil
}
