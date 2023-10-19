package elasticsan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ElasticSANResource{}
var _ sdk.ResourceWithUpdate = ElasticSANResource{}

type ElasticSANResource struct{}

func (r ElasticSANResource) ModelObject() interface{} {
	return &ElasticSANResourceModel{}
}

type ElasticSANResourceModel struct {
	BaseSizeInTiB             int64                        `tfschema:"base_size_in_tib"`
	ExtendedCapacitySizeInTiB int64                        `tfschema:"extended_capacity_size_in_tib"`
	Location                  string                       `tfschema:"location"`
	Name                      string                       `tfschema:"name"`
	ResourceGroupName         string                       `tfschema:"resource_group_name"`
	Sku                       []ElasticSANResourceSkuModel `tfschema:"sku"`
	Tags                      map[string]interface{}       `tfschema:"tags"`
	TotalIops                 int64                        `tfschema:"total_iops"`
	TotalMBps                 int64                        `tfschema:"total_mbps"`
	TotalSizeInTiB            int64                        `tfschema:"total_size_in_tib"`
	TotalVolumeSizeInGiB      int64                        `tfschema:"total_volume_size_in_gib"`
	VolumeGroupCount          int64                        `tfschema:"volume_group_count"`
	Zones                     []string                     `tfschema:"zones"`
}

type ElasticSANResourceSkuModel struct {
	Name string `tfschema:"name"`
	Tier string `tfschema:"tier"`
}

func (r ElasticSANResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return elasticsans.ValidateElasticSanID
}

func (r ElasticSANResource) ResourceType() string {
	return "azurerm_elastic_san"
}

func (r ElasticSANResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Required: true,
			ForceNew: true,
			Type:     pluginsdk.TypeString,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"base_size_in_tib": {
			Required:     true,
			Type:         pluginsdk.TypeInt,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"extended_capacity_size_in_tib": {
			Optional:     true,
			Type:         pluginsdk.TypeInt,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"sku": {
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Type:     pluginsdk.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Required: true,
						ForceNew: true,
						Type:     pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							elasticsans.PossibleValuesForSkuName(),
							false,
						),
					},

					"tier": {
						Optional: true,
						Type:     pluginsdk.TypeString,
						Default:  string(elasticsans.SkuTierPremium),
						ValidateFunc: validation.StringInSlice(
							elasticsans.PossibleValuesForSkuTier(),
							false,
						),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r ElasticSANResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"total_iops": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"total_mbps": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"total_size_in_tib": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"total_volume_size_in_gib": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"volume_group_count": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},
	}
}

func (r ElasticSANResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Client.ElasticSans

			var config ElasticSANResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := elasticsans.NewElasticSanID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := elasticsans.ElasticSan{
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				Properties: elasticsans.ElasticSanProperties{
					BaseSizeTiB:             config.BaseSizeInTiB,
					ExtendedCapacitySizeTiB: config.ExtendedCapacitySizeInTiB,
					Sku:                     r.ExpandSku(config.Sku),
				},
			}

			if config.Zones != nil {
				payload.Properties.AvailabilityZones = pointer.To(zones.Expand(config.Zones))
			}

			if err := client.CreateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ElasticSANResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Client.ElasticSans
			schema := ElasticSANResourceModel{}

			id, err := elasticsans.ParseElasticSanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema.Name = id.ElasticSanName
			schema.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				prop := model.Properties
				schema.Sku = r.FlattenSku(prop.Sku)
				schema.Zones = zones.Flatten(prop.AvailabilityZones)
				schema.BaseSizeInTiB = prop.BaseSizeTiB
				schema.ExtendedCapacitySizeInTiB = prop.ExtendedCapacitySizeTiB
				schema.TotalIops = pointer.From(prop.TotalIops)
				schema.TotalMBps = pointer.From(prop.TotalMBps)
				schema.TotalSizeInTiB = pointer.From(prop.TotalSizeTiB)
				schema.TotalVolumeSizeInGiB = pointer.From(prop.TotalVolumeSizeGiB)
				schema.VolumeGroupCount = pointer.From(prop.VolumeGroupCount)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ElasticSANResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Client.ElasticSans

			id, err := elasticsans.ParseElasticSanID(metadata.ResourceData.Id())
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

func (r ElasticSANResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Client.ElasticSans

			id, err := elasticsans.ParseElasticSanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ElasticSANResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := elasticsans.ElasticSanUpdate{}

			if metadata.ResourceData.HasChange("base_size_in_tib") {
				if payload.Properties == nil {
					payload.Properties = &elasticsans.ElasticSanUpdateProperties{}
				}
				payload.Properties.BaseSizeTiB = pointer.FromInt64(config.BaseSizeInTiB)
			}

			if metadata.ResourceData.HasChange("extended_capacity_size_in_tib") {
				if payload.Properties == nil {
					payload.Properties = &elasticsans.ElasticSanUpdateProperties{}
				}
				payload.Properties.ExtendedCapacitySizeTiB = pointer.FromInt64(config.ExtendedCapacitySizeInTiB)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ElasticSANResource) ExpandSku(input []ElasticSANResourceSkuModel) elasticsans.Sku {
	output := elasticsans.Sku{
		Name: elasticsans.SkuName(input[0].Name),
	}

	if input[0].Tier != "" {
		output.Tier = pointer.To(elasticsans.SkuTier(input[0].Tier))
	}

	return output
}

func (r ElasticSANResource) FlattenSku(input elasticsans.Sku) []ElasticSANResourceSkuModel {
	return []ElasticSANResourceSkuModel{
		{
			Name: string(input.Name),
			Tier: string(pointer.From(input.Tier)),
		},
	}
}
