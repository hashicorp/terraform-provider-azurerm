package deviceregistry

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	resourceParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	resourceValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	AssetExtendedLocationTypeCustomLocation = "CustomLocation"
)

var _ sdk.Resource = AssetResource{}

type AssetResource struct{}

type AssetResourceModel struct {
	Name                          string                 `tfschema:"name"`
	ResourceGroupId               string                 `tfschema:"resource_group_id"`
	Location                      string                 `tfschema:"location"`
	Tags                          map[string]string      `tfschema:"tags"`
	ExtendedLocationId            string                 `tfschema:"extended_location_id"`
	Enabled                       bool                   `tfschema:"enabled"`
	ExternalAssetId               string                 `tfschema:"external_asset_id"`
	DisplayName                   string                 `tfschema:"display_name"`
	Description                   string                 `tfschema:"description"`
	AssetEndpointProfileReference string                 `tfschema:"asset_endpoint_profile_reference"`
	Manufacturer                  string                 `tfschema:"manufacturer"`
	ManufacturerUri               string                 `tfschema:"manufacturer_uri"`
	Model                         string                 `tfschema:"model"`
	ProductCode                   string                 `tfschema:"product_code"`
	HardwareRevision              string                 `tfschema:"hardware_revision"`
	SoftwareRevision              string                 `tfschema:"software_revision"`
	DocumentationUri              string                 `tfschema:"documentation_uri"`
	SerialNumber                  string                 `tfschema:"serial_number"`
	Attributes                    map[string]interface{} `tfschema:"attributes"`
	DiscoveredAssetReferences     []string               `tfschema:"discovered_asset_references"`
	DefaultDatasetsConfiguration  string                 `tfschema:"default_datasets_configuration"`
	DefaultEventsConfiguration    string                 `tfschema:"default_events_configuration"`
	DefaultTopic                  []TopicModel           `tfschema:"default_topic"`
	Datasets                      []Dataset              `tfschema:"dataset"`
	Events                        []Event                `tfschema:"event"`
}

type Dataset struct {
	Name                 string       `tfschema:"name"`
	DatasetConfiguration string       `tfschema:"dataset_configuration"`
	Topic                []TopicModel `tfschema:"topic"`
	DataPoints           []DataPoint  `tfschema:"data_point"`
}

type DataPoint struct {
	Name                   string `tfschema:"name"`
	DataSource             string `tfschema:"data_source"`
	ObservabilityMode      string `tfschema:"observability_mode"`
	DataPointConfiguration string `tfschema:"data_point_configuration"`
}

type Event struct {
	Name               string       `tfschema:"name"`
	EventNotifier      string       `tfschema:"event_notifier"`
	ObservabilityMode  string       `tfschema:"observability_mode"`
	EventConfiguration string       `tfschema:"event_configuration"`
	Topic              []TopicModel `tfschema:"topic"`
}

type TopicModel struct {
	Path   string `tfschema:"path"`
	Retain string `tfschema:"retain"`
}

func (AssetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: resourceValidate.ResourceGroupID,
		},
		"location": commonschema.Location(),
		"tags":     commonschema.Tags(),
		"extended_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"external_asset_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"asset_endpoint_profile_reference": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"manufacturer": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"manufacturer_uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"model": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"product_code": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"hardware_revision": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"software_revision": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"documentation_uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"serial_number": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"attributes": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"discovered_asset_references": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"default_datasets_configuration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"default_events_configuration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"default_topic": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"path": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"retain": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(assets.PossibleValuesForTopicRetainType(), false),
						Default:      string(assets.TopicRetainTypeNever),
					},
				},
			},
		},
		"dataset": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"dataset_configuration": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"topic": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"path": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"retain": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(assets.PossibleValuesForTopicRetainType(), false),
									Default:      string(assets.TopicRetainTypeNever),
								},
							},
						},
					},
					"data_point": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"data_source": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"observability_mode": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(assets.DataPointObservabilityModeNone),
									ValidateFunc: validation.StringInSlice(assets.PossibleValuesForDataPointObservabilityMode(), false),
								},
								"data_point_configuration": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
		"event": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"event_notifier": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"observability_mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(assets.EventObservabilityModeNone),
						ValidateFunc: validation.StringInSlice(assets.PossibleValuesForEventObservabilityMode(), false),
					},
					"event_configuration": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"topic": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"path": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"retain": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(assets.PossibleValuesForTopicRetainType(), false),
									Default:      string(assets.TopicRetainTypeNever),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (AssetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (AssetResource) ModelObject() interface{} {
	return &AssetResourceModel{}
}

func (AssetResource) ResourceType() string {
	return "azurerm_device_registry_asset"
}

func (r AssetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DeviceRegistry.AssetsClient

			var config AssetResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resourceGroupId, err := resourceParse.ResourceGroupID(config.ResourceGroupId)
			if err != nil {
				return fmt.Errorf("parsing resource group id: %+v", err)
			}

			id := assets.NewAssetID(resourceGroupId.SubscriptionId, resourceGroupId.ResourceGroup, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// Convert the TF model to the ARM model
			// Optional ARM resource properties are pointers.
			param := assets.Asset{
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
				ExtendedLocation: assets.ExtendedLocation{
					Name: config.ExtendedLocationId,
					Type: AssetExtendedLocationTypeCustomLocation,
				},
				Properties: &assets.AssetProperties{
					AssetEndpointProfileRef: config.AssetEndpointProfileReference,
				},
			}

			// Enabled in config will be default set to false if not explicitly set in Terraform.
			// We must check the raw value if it is null so we don't just send false if user didn't specify
			if !pluginsdk.IsExplicitlyNullInConfig(metadata.ResourceData, "enabled") {
				param.Properties.Enabled = pointer.To(config.Enabled)
			}

			if config.ExternalAssetId != "" {
				param.Properties.ExternalAssetId = pointer.To(config.ExternalAssetId)
			}

			if config.DisplayName != "" {
				param.Properties.DisplayName = pointer.To(config.DisplayName)
			}

			if config.Description != "" {
				param.Properties.Description = pointer.To(config.Description)
			}

			if config.Manufacturer != "" {
				param.Properties.Manufacturer = pointer.To(config.Manufacturer)
			}

			if config.ManufacturerUri != "" {
				param.Properties.ManufacturerUri = pointer.To(config.ManufacturerUri)
			}

			if config.Model != "" {
				param.Properties.Model = pointer.To(config.Model)
			}

			if config.ProductCode != "" {
				param.Properties.ProductCode = pointer.To(config.ProductCode)
			}

			if config.HardwareRevision != "" {
				param.Properties.HardwareRevision = pointer.To(config.HardwareRevision)
			}

			if config.SoftwareRevision != "" {
				param.Properties.SoftwareRevision = pointer.To(config.SoftwareRevision)
			}

			if config.DocumentationUri != "" {
				param.Properties.DocumentationUri = pointer.To(config.DocumentationUri)
			}

			if config.SerialNumber != "" {
				param.Properties.SerialNumber = pointer.To(config.SerialNumber)
			}

			if config.Attributes != nil {
				param.Properties.Attributes = pointer.To(config.Attributes)
			}

			if config.DiscoveredAssetReferences != nil {
				param.Properties.DiscoveredAssetRefs = pointer.To(config.DiscoveredAssetReferences)
			}

			if config.DefaultDatasetsConfiguration != "" {
				param.Properties.DefaultDatasetsConfiguration = pointer.To(config.DefaultDatasetsConfiguration)
			}

			if config.DefaultEventsConfiguration != "" {
				param.Properties.DefaultEventsConfiguration = pointer.To(config.DefaultEventsConfiguration)
			}

			param.Properties.DefaultTopic = expandTopic(config.DefaultTopic)

			param.Properties.Datasets = expandDatasets(config.Datasets)

			param.Properties.Events = expandEvents(config.Events)

			if err := client.CreateOrReplaceThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AssetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DeviceRegistry.AssetsClient

			id, err := assets.ParseAssetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config AssetResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Change the properties that can be updated
			param := assets.AssetUpdate{
				Properties: &assets.AssetUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("tags") {
				param.Tags = pointer.To(config.Tags)
			}

			if metadata.ResourceData.HasChange("attributes") {
				param.Properties.Attributes = pointer.To(config.Attributes)
			}

			if metadata.ResourceData.HasChange("dataset") {
				param.Properties.Datasets = expandDatasets(config.Datasets)
			}

			if metadata.ResourceData.HasChange("default_datasets_configuration") {
				param.Properties.DefaultDatasetsConfiguration = pointer.To(config.DefaultDatasetsConfiguration)
			}

			if metadata.ResourceData.HasChange("default_events_configuration") {
				param.Properties.DefaultEventsConfiguration = pointer.To(config.DefaultEventsConfiguration)
			}

			if metadata.ResourceData.HasChange("default_topic") {
				topic := &assets.TopicUpdate{}
				param.Properties.DefaultTopic = topic
				if len(config.DefaultTopic) > 0 {
					topic.Path = pointer.To(config.DefaultTopic[0].Path)
					// Bug with `go-azure-sdk` library: you can't set retain to null because empty string will cause
					// ARM to throw validation error (retain must be one of the property's possible enum values),
					// and go-azure-sdk library will ignore the retain field if it's set to nil, even if explicitly set.
					if config.DefaultTopic[0].Retain != "" {
						topic.Retain = pointer.To(assets.TopicRetainType(config.DefaultTopic[0].Retain))
					}
				}
			}

			if metadata.ResourceData.HasChange("description") {
				param.Properties.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("display_name") {
				param.Properties.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("documentation_uri") {
				param.Properties.DocumentationUri = pointer.To(config.DocumentationUri)
			}

			if metadata.ResourceData.HasChange("enabled") {
				param.Properties.Enabled = pointer.To(config.Enabled)
			}

			if metadata.ResourceData.HasChange("event") {
				param.Properties.Events = expandEvents(config.Events)
			}

			if metadata.ResourceData.HasChange("hardware_revision") {
				param.Properties.HardwareRevision = pointer.To(config.HardwareRevision)
			}

			if metadata.ResourceData.HasChange("manufacturer") {
				param.Properties.Manufacturer = pointer.To(config.Manufacturer)
			}

			if metadata.ResourceData.HasChange("manufacturer_uri") {
				param.Properties.ManufacturerUri = pointer.To(config.ManufacturerUri)
			}

			if metadata.ResourceData.HasChange("model") {
				param.Properties.Model = pointer.To(config.Model)
			}

			if metadata.ResourceData.HasChange("product_code") {
				param.Properties.ProductCode = pointer.To(config.ProductCode)
			}

			if metadata.ResourceData.HasChange("serial_number") {
				param.Properties.SerialNumber = pointer.To(config.SerialNumber)
			}

			if metadata.ResourceData.HasChange("software_revision") {
				param.Properties.SoftwareRevision = pointer.To(config.SoftwareRevision)
			}

			if err := client.UpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (AssetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DeviceRegistry.AssetsClient

			id, err := assets.ParseAssetID(metadata.ResourceData.Id())
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

			resourceGroupId := resourceParse.NewResourceGroupID(id.SubscriptionId, id.ResourceGroupName)

			// Convert the ARM model to the TF model
			state := AssetResourceModel{
				Name:            id.AssetName,
				ResourceGroupId: resourceGroupId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.ExtendedLocationId = model.ExtendedLocation.Name
				if props := model.Properties; props != nil {
					state.AssetEndpointProfileReference = props.AssetEndpointProfileRef
					state.Enabled = pointer.From(props.Enabled)
					state.ExternalAssetId = pointer.From(props.ExternalAssetId)
					state.DisplayName = pointer.From(props.DisplayName)
					state.Description = pointer.From(props.Description)
					state.Manufacturer = pointer.From(props.Manufacturer)
					state.ManufacturerUri = pointer.From(props.ManufacturerUri)
					state.Model = pointer.From(props.Model)
					state.ProductCode = pointer.From(props.ProductCode)
					state.HardwareRevision = pointer.From(props.HardwareRevision)
					state.SoftwareRevision = pointer.From(props.SoftwareRevision)
					state.DocumentationUri = pointer.From(props.DocumentationUri)
					state.SerialNumber = pointer.From(props.SerialNumber)
					state.Attributes = pointer.From(props.Attributes)
					state.DiscoveredAssetReferences = pointer.From(props.DiscoveredAssetRefs)
					state.DefaultDatasetsConfiguration = pointer.From(props.DefaultDatasetsConfiguration)
					state.DefaultEventsConfiguration = pointer.From(props.DefaultEventsConfiguration)

					if defaultTopic := props.DefaultTopic; defaultTopic != nil {
						state.DefaultTopic = flattenTopic(props.DefaultTopic)
					}

					if datasets := props.Datasets; datasets != nil {
						state.Datasets = flattenDatasets(datasets)
					}

					if events := props.Events; events != nil {
						state.Events = flattenEvents(events)
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (AssetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DeviceRegistry.AssetsClient

			id, err := assets.ParseAssetID(metadata.ResourceData.Id())
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

func (AssetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return assets.ValidateAssetID
}

func expandDatasets(datasets []Dataset) *[]assets.Dataset {
	if datasets == nil {
		return nil
	}

	azureDatasets := make([]assets.Dataset, len(datasets))
	for i, dataset := range datasets {
		azureDatasets[i] = assets.Dataset{
			Name:                 dataset.Name,
			DatasetConfiguration: pointer.To(dataset.DatasetConfiguration),
			Topic:                expandTopic(dataset.Topic),
			DataPoints:           expandDataPoints(dataset.DataPoints),
		}
	}

	return &azureDatasets
}

func expandDataPoints(dataPoints []DataPoint) *[]assets.DataPoint {
	if dataPoints == nil {
		return nil
	}

	azureDataPoints := make([]assets.DataPoint, len(dataPoints))
	for i, dataPoint := range dataPoints {
		azureDataPoints[i] = assets.DataPoint{
			Name:                   dataPoint.Name,
			DataSource:             dataPoint.DataSource,
			ObservabilityMode:      pointer.To(assets.DataPointObservabilityMode(dataPoint.ObservabilityMode)),
			DataPointConfiguration: pointer.To(dataPoint.DataPointConfiguration),
		}
	}

	return &azureDataPoints
}

func expandEvents(events []Event) *[]assets.Event {
	if events == nil {
		return nil
	}

	azureEvents := make([]assets.Event, len(events))
	for i, event := range events {
		azureEvents[i] = assets.Event{
			Name:               event.Name,
			EventNotifier:      event.EventNotifier,
			EventConfiguration: pointer.To(event.EventConfiguration),
			ObservabilityMode:  pointer.To(assets.EventObservabilityMode(event.ObservabilityMode)),
			Topic:              expandTopic(event.Topic),
		}
	}

	return &azureEvents
}

func expandTopic(topic []TopicModel) *assets.Topic {
	if len(topic) == 0 {
		return nil
	}

	azureTopic := assets.Topic{
		Path: topic[0].Path,
	}

	// Topic retain is optional, but if it's set, it must be one of the possible values
	if topic[0].Retain != "" {
		azureTopic.Retain = pointer.To(assets.TopicRetainType(topic[0].Retain))
	}

	return &azureTopic
}

func flattenDatasets(datasets *[]assets.Dataset) []Dataset {
	if datasets == nil {
		return nil
	}

	tfDatasets := make([]Dataset, len(*datasets))
	for i, dataset := range *datasets {
		tfDatasets[i] = Dataset{
			Name:                 dataset.Name,
			DatasetConfiguration: pointer.From(dataset.DatasetConfiguration),
			DataPoints:           flattenDataPoints(dataset.DataPoints),
			Topic:                flattenTopic(dataset.Topic),
		}
	}

	return tfDatasets
}

func flattenDataPoints(dataPoints *[]assets.DataPoint) []DataPoint {
	if dataPoints == nil {
		return nil
	}

	tfDataPoints := make([]DataPoint, len(*dataPoints))
	for i, dataPoint := range *dataPoints {
		tfDataPoints[i] = DataPoint{
			Name:                   dataPoint.Name,
			DataSource:             dataPoint.DataSource,
			ObservabilityMode:      string(pointer.From(dataPoint.ObservabilityMode)),
			DataPointConfiguration: pointer.From(dataPoint.DataPointConfiguration),
		}
	}

	return tfDataPoints
}

func flattenEvents(events *[]assets.Event) []Event {
	if events == nil {
		return nil
	}

	tfEvents := make([]Event, len(*events))
	for i, event := range *events {
		tfEvents[i] = Event{
			Name:               event.Name,
			EventNotifier:      event.EventNotifier,
			ObservabilityMode:  string(pointer.From(event.ObservabilityMode)),
			EventConfiguration: pointer.From(event.EventConfiguration),
			Topic:              flattenTopic(event.Topic),
		}
	}

	return tfEvents
}

func flattenTopic(topic *assets.Topic) []TopicModel {
	if topic == nil {
		return nil
	}

	return []TopicModel{
		{
			Path:   topic.Path,
			Retain: string(pointer.From(topic.Retain)),
		},
	}
}
