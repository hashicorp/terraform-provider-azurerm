package azurestackhci

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/marketplacegalleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = StackHCIMarketplaceGalleryImageResource{}
	_ sdk.ResourceWithUpdate = StackHCIMarketplaceGalleryImageResource{}
)

type StackHCIMarketplaceGalleryImageResource struct{}

func (StackHCIMarketplaceGalleryImageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return marketplacegalleryimages.ValidateMarketplaceGalleryImageID
}

func (StackHCIMarketplaceGalleryImageResource) ResourceType() string {
	return "azurerm_stack_hci_marketplace_gallery_image"
}

func (StackHCIMarketplaceGalleryImageResource) ModelObject() interface{} {
	return &StackHCIMarketplaceGalleryImageResourceModel{}
}

type StackHCIMarketplaceGalleryImageResourceModel struct {
	Name              string                                      `tfschema:"name"`
	ResourceGroupName string                                      `tfschema:"resource_group_name"`
	Location          string                                      `tfschema:"location"`
	CustomLocationId  string                                      `tfschema:"custom_location_id"`
	HypervGeneration  string                                      `tfschema:"hyperv_generation"`
	Identifier        []StackHCIMarketplaceGalleryImageIdentifier `tfschema:"identifier"`
	OsType            string                                      `tfschema:"os_type"`
	Version           string                                      `tfschema:"version"`
	StoragePathId     string                                      `tfschema:"storage_path_id"`
	Tags              map[string]interface{}                      `tfschema:"tags"`
}

type StackHCIMarketplaceGalleryImageIdentifier struct {
	Offer     string `tfschema:"offer"`
	Publisher string `tfschema:"publisher"`
	Sku       string `tfschema:"sku"`
}

func (StackHCIMarketplaceGalleryImageResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,78}[a-zA-Z0-9]$`),
				"name must be between 2 and 80 characters and can only contain alphanumberic characters, hyphen, dot and underline",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},

		"hyperv_generation": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(marketplacegalleryimages.HyperVGenerationVOne),
				string(marketplacegalleryimages.HyperVGenerationVTwo),
			}, false),
		},

		"identifier": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"publisher": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"offer": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sku": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(marketplacegalleryimages.OperatingSystemTypesLinux),
				string(marketplacegalleryimages.OperatingSystemTypesWindows),
			}, false),
		},

		"version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_path_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: storagecontainers.ValidateStorageContainerID,
		},

		"tags": commonschema.Tags(),
	}
}

func (StackHCIMarketplaceGalleryImageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCIMarketplaceGalleryImageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.MarketplaceGalleryImages

			var config StackHCIMarketplaceGalleryImageResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := marketplacegalleryimages.NewMarketplaceGalleryImageID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := marketplacegalleryimages.MarketplaceGalleryImages{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				ExtendedLocation: &marketplacegalleryimages.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(marketplacegalleryimages.ExtendedLocationTypesCustomLocation),
				},
				Properties: &marketplacegalleryimages.MarketplaceGalleryImageProperties{
					Identifier: expandStackHCIMarketplaceGalleryImageIdentifier(config.Identifier),
					OsType:     marketplacegalleryimages.OperatingSystemTypes(config.OsType),
				},
			}

			if config.StoragePathId != "" {
				payload.Properties.ContainerId = pointer.To(config.StoragePathId)
			}

			if config.HypervGeneration != "" {
				payload.Properties.HyperVGeneration = pointer.To(marketplacegalleryimages.HyperVGeneration(config.HypervGeneration))
			}

			if config.Version != "" {
				payload.Properties.Version = &marketplacegalleryimages.GalleryImageVersion{
					Name: pointer.To(config.Version),
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("performing create %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r StackHCIMarketplaceGalleryImageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.MarketplaceGalleryImages

			id, err := marketplacegalleryimages.ParseMarketplaceGalleryImageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schema := StackHCIMarketplaceGalleryImageResourceModel{
				Name:              id.MarketplaceGalleryImageName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				if model.ExtendedLocation != nil && model.ExtendedLocation.Name != nil {
					customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(*model.ExtendedLocation.Name)
					if err != nil {
						return err
					}

					schema.CustomLocationId = customLocationId.ID()
				}

				if props := model.Properties; props != nil {
					schema.StoragePathId = pointer.From(props.ContainerId)
					schema.OsType = string(props.OsType)
					schema.HypervGeneration = string(pointer.From(props.HyperVGeneration))
					schema.Identifier = flattenStackHCIMarketplaceGalleryImageIdentifier(props.Identifier)

					if props.Version != nil {
						schema.Version = pointer.From(props.Version.Name)
					}
				}

			}
			return metadata.Encode(&schema)
		},
	}
}

func (r StackHCIMarketplaceGalleryImageResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.MarketplaceGalleryImages

			id, err := marketplacegalleryimages.ParseMarketplaceGalleryImageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StackHCIMarketplaceGalleryImageResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &marketplacegalleryimages.MarketplaceGalleryImagesUpdateRequest{}
			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.Expand(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r StackHCIMarketplaceGalleryImageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 1 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.MarketplaceGalleryImages

			id, err := marketplacegalleryimages.ParseMarketplaceGalleryImageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandStackHCIMarketplaceGalleryImageIdentifier(input []StackHCIMarketplaceGalleryImageIdentifier) *marketplacegalleryimages.GalleryImageIdentifier {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &marketplacegalleryimages.GalleryImageIdentifier{
		Offer:     v.Offer,
		Publisher: v.Publisher,
		Sku:       v.Sku,
	}
}

func flattenStackHCIMarketplaceGalleryImageIdentifier(input *marketplacegalleryimages.GalleryImageIdentifier) []StackHCIMarketplaceGalleryImageIdentifier {
	if input == nil {
		return make([]StackHCIMarketplaceGalleryImageIdentifier, 0)
	}

	return []StackHCIMarketplaceGalleryImageIdentifier{
		{
			Offer:     input.Offer,
			Publisher: input.Publisher,
			Sku:       input.Sku,
		},
	}
}
