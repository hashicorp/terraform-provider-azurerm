package elasticsan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = ElasticSANVolumeResource{}
	_ sdk.ResourceWithUpdate = ElasticSANVolumeResource{}
)

type ElasticSANVolumeResource struct{}

func (r ElasticSANVolumeResource) ModelObject() interface{} {
	return &ElasticSANVolumeResourceModel{}
}

type ElasticSANVolumeResourceModel struct {
	CreateSourceId       string `tfschema:"create_source_id"`
	Name                 string `tfschema:"name"`
	SizeInGiB            int64  `tfschema:"size_in_gib"`
	TargetIqn            string `tfschema:"target_iqn"`
	TargetPortalHostname string `tfschema:"target_portal_hostname"`
	TargetPortalPort     int64  `tfschema:"target_portal_port"`
	VolumeGroupId        string `tfschema:"volume_group_id"`
	VolumeId             string `tfschema:"volume_id"`
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

		"volume_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: volumes.ValidateVolumeGroupID,
		},

		"size_in_gib": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65536),
		},

		"create_source_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
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
					SizeGiB: config.SizeInGiB,
				},
			}

			if config.CreateSourceId != "" {
				payload.Properties.CreationData = &volumes.SourceCreationData{
					SourceId: pointer.To(config.CreateSourceId),
				}
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

				if creationData := props.CreationData; creationData != nil {
					schema.CreateSourceId = pointer.From(creationData.SourceId)
				}

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
