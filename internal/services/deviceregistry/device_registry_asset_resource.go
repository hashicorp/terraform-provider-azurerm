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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = AssetResource{}

type AssetResource struct{}

type AssetResourceModel struct {
	Name                         string                 `tfschema:"name"`
	ResourceGroupName            string                 `tfschema:"resource_group_name"`
	Location                     string                 `tfschema:"location"`
	Tags                         map[string]string      `tfschema:"tags"`
	ExtendedLocationName         string                 `tfschema:"extended_location_name"`
	ExtendedLocationType         string                 `tfschema:"extended_location_type"`
	Enabled                      bool                   `tfschema:"enabled"`
	ExternalAssetId              string                 `tfschema:"external_asset_id"`
	DisplayName                  string                 `tfschema:"display_name"`
	Description                  string                 `tfschema:"description"`
	AssetEndpointProfileRef      string                 `tfschema:"asset_endpoint_profile_ref"`
	Manufacturer                 string                 `tfschema:"manufacturer"`
	ManufacturerUri              string                 `tfschema:"manufacturer_uri"`
	Model                        string                 `tfschema:"model"`
	ProductCode                  string                 `tfschema:"product_code"`
	HardwareRevision             string                 `tfschema:"hardware_revision"`
	SoftwareRevision             string                 `tfschema:"software_revision"`
	DocumentationUri             string                 `tfschema:"documentation_uri"`
	SerialNumber                 string                 `tfschema:"serial_number"`
	Attributes                   map[string]interface{} `tfschema:"attributes"`
	DiscoveredAssetRefs          []string               `tfschema:"discovered_asset_refs"`
	DefaultDatasetsConfiguration string                 `tfschema:"default_datasets_configuration"`
	DefaultEventsConfiguration   string                 `tfschema:"default_events_configuration"`
	DefaultTopicPath             string                 `tfschema:"default_topic_path"`
	DefaultTopicRetain           string                 `tfschema:"default_topic_retain"`
	Datasets                     []Dataset              `tfschema:"datasets"`
	Events                       []Event                `tfschema:"events"`
}

type Dataset struct {
	Name                 string      `tfschema:"name"`
	DatasetConfiguration string      `tfschema:"dataset_configuration"`
	TopicPath            string      `tfschema:"topic_path"`
	TopicRetain          string      `tfschema:"topic_retain"`
	DataPoints           []DataPoint `tfschema:"data_points"`
}

type DataPoint struct {
	Name                   string `tfschema:"name"`
	DataSource             string `tfschema:"data_source"`
	ObservabilityMode      string `tfschema:"observability_mode"`
	DataPointConfiguration string `tfschema:"data_point_configuration"`
}

type Event struct {
	Name               string `tfschema:"name"`
	EventNotifier      string `tfschema:"event_notifier"`
	ObservabilityMode  string `tfschema:"observability_mode"`
	EventConfiguration string `tfschema:"event_configuration"`
	TopicPath          string `tfschema:"topic_path"`
	TopicRetain        string `tfschema:"topic_retain"`
}

func (AssetResource) Arguments() map[string]*pluginsdk.Schema {
	// add the other asset properties
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"tags":                commonschema.Tags(),
		"extended_location_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"extended_location_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
		"asset_endpoint_profile_ref": {
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
		"discovered_asset_refs": {
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
		"default_topic_path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"default_topic_retain": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(assets.PossibleValuesForTopicRetainType(), false),
			RequiredWith: []string{"default_topic_path"},
		},
		"datasets": {
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
					"topic_path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"topic_retain": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(assets.PossibleValuesForTopicRetainType(), false),
					},
					"data_points": {
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
		"events": {
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
					"topic_path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"topic_retain": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(assets.PossibleValuesForTopicRetainType(), false),
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
			client := metadata.Client.DeviceRegistry.AssetClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config AssetResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := assets.NewAssetID(subscriptionId, config.ResourceGroupName, config.Name)

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
					Name: config.ExtendedLocationName,
					Type: config.ExtendedLocationType,
				},
				Properties: &assets.AssetProperties{
					AssetEndpointProfileRef: config.AssetEndpointProfileRef,
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

			if config.DiscoveredAssetRefs != nil {
				param.Properties.DiscoveredAssetRefs = pointer.To(config.DiscoveredAssetRefs)
			}

			if config.DefaultDatasetsConfiguration != "" {
				param.Properties.DefaultDatasetsConfiguration = pointer.To(config.DefaultDatasetsConfiguration)
			}

			if config.DefaultEventsConfiguration != "" {
				param.Properties.DefaultEventsConfiguration = pointer.To(config.DefaultEventsConfiguration)
			}

			param.Properties.DefaultTopic = toAzureTopic(config.DefaultTopicPath, config.DefaultTopicRetain)

			param.Properties.Datasets = toAzureDatasets(config.Datasets)

			param.Properties.Events = toAzureEvents(config.Events)

			if _, err := client.CreateOrReplace(ctx, id, param); err != nil {
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
			client := metadata.Client.DeviceRegistry.AssetClient

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

			if metadata.ResourceData.HasChange("datasets") {
				param.Properties.Datasets = toAzureDatasets(config.Datasets)
			}

			if metadata.ResourceData.HasChange("default_datasets_configuration") {
				param.Properties.DefaultDatasetsConfiguration = pointer.To(config.DefaultDatasetsConfiguration)
			}

			if metadata.ResourceData.HasChange("default_events_configuration") {
				param.Properties.DefaultEventsConfiguration = pointer.To(config.DefaultEventsConfiguration)
			}

			defaultTopicPathChanged := metadata.ResourceData.HasChange("default_topic_path")
			defaultTopicRetainChanged := metadata.ResourceData.HasChange("default_topic_retain")
			if defaultTopicPathChanged || defaultTopicRetainChanged {
				param.Properties.DefaultTopic = &assets.TopicUpdate{}
				if defaultTopicPathChanged {
					param.Properties.DefaultTopic.Path = pointer.To(config.DefaultTopicPath)
				}
				if defaultTopicRetainChanged {
					// Bug with `go-azure-sdk` library: you can't set retain to null because empty string will cause
					// ARM to throw validation error (retain must be one of the possible enum values),
					// and go-azure-sdk library will ignore the retain field if it's set to nil, even if explicitly set.
					param.Properties.DefaultTopic.Retain = pointer.To(assets.TopicRetainType(config.DefaultTopicRetain))
				}
			}

			if metadata.ResourceData.HasChange("description") {
				param.Properties.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("display_name") {
				param.Properties.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("enabled") {
				param.Properties.Enabled = pointer.To(config.Enabled)
			}

			if metadata.ResourceData.HasChange("events") {
				param.Properties.Events = toAzureEvents(config.Events)
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

			if _, err := client.Update(ctx, *id, param); err != nil {
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
			client := metadata.Client.DeviceRegistry.AssetClient

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

			// Convert the ARM model to the TF model
			state := AssetResourceModel{
				Name:              id.AssetName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.ExtendedLocationName = model.ExtendedLocation.Name
				state.ExtendedLocationType = model.ExtendedLocation.Type
				if props := model.Properties; props != nil {
					state.AssetEndpointProfileRef = props.AssetEndpointProfileRef
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
					state.DiscoveredAssetRefs = pointer.From(props.DiscoveredAssetRefs)
					state.DefaultDatasetsConfiguration = pointer.From(props.DefaultDatasetsConfiguration)
					state.DefaultEventsConfiguration = pointer.From(props.DefaultEventsConfiguration)

					if defaultTopic := props.DefaultTopic; defaultTopic != nil {
						state.DefaultTopicPath, state.DefaultTopicRetain = toTFTopic(props.DefaultTopic)
					}

					if datasets := props.Datasets; datasets != nil {
						state.Datasets = toTFDatasets(datasets)
					}

					if events := props.Events; events != nil {
						state.Events = toTFEvents(events)
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
			client := metadata.Client.DeviceRegistry.AssetClient

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

func toAzureDatasets(datasets []Dataset) *[]assets.Dataset {
	if datasets == nil {
		return nil
	}

	azureDatasets := make([]assets.Dataset, len(datasets))
	for i, dataset := range datasets {
		azureDatasets[i] = assets.Dataset{
			Name:                 dataset.Name,
			DatasetConfiguration: pointer.To(dataset.DatasetConfiguration),
			Topic:                toAzureTopic(dataset.TopicPath, dataset.TopicRetain),
			DataPoints:           toAzureDataPoints(dataset.DataPoints),
		}
	}

	return &azureDatasets
}

func toAzureDataPoints(dataPoints []DataPoint) *[]assets.DataPoint {
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

func toAzureEvents(events []Event) *[]assets.Event {
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
			Topic:              toAzureTopic(event.TopicPath, event.TopicRetain),
		}
	}

	return &azureEvents
}

func toAzureTopic(topicPath string, topicRetain string) *assets.Topic {
	if topicPath == "" && topicRetain == "" {
		return nil
	}

	azureTopic := assets.Topic{
		Path: topicPath,
	}

	// Topic retain is optional, but if it's set, it must be one of the possible values
	if topicRetain != "" {
		azureTopic.Retain = pointer.To(assets.TopicRetainType(topicRetain))
	}

	return &azureTopic
}

func toTFDatasets(datasets *[]assets.Dataset) []Dataset {
	if datasets == nil {
		return nil
	}

	tfDatasets := make([]Dataset, len(*datasets))
	for i, dataset := range *datasets {
		topicPath, topicRetain := toTFTopic(dataset.Topic)
		tfDatasets[i] = Dataset{
			Name:                 dataset.Name,
			DatasetConfiguration: pointer.From(dataset.DatasetConfiguration),
			DataPoints:           toTFDataPoints(dataset.DataPoints),
			TopicPath:            topicPath,
			TopicRetain:          topicRetain,
		}
	}

	return tfDatasets
}

func toTFDataPoints(dataPoints *[]assets.DataPoint) []DataPoint {
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

func toTFEvents(events *[]assets.Event) []Event {
	if events == nil {
		return nil
	}

	tfEvents := make([]Event, len(*events))
	for i, event := range *events {
		topicPath, topicRetain := toTFTopic(event.Topic)
		tfEvents[i] = Event{
			Name:               event.Name,
			EventNotifier:      event.EventNotifier,
			ObservabilityMode:  string(pointer.From(event.ObservabilityMode)),
			EventConfiguration: pointer.From(event.EventConfiguration),
			TopicPath:          topicPath,
			TopicRetain:        topicRetain,
		}
	}

	return tfEvents
}

func toTFTopic(topic *assets.Topic) (string, string) {
	if topic == nil {
		return "", ""
	}

	if topic.Retain == nil {
		return topic.Path, ""
	} else {
		return topic.Path, string(pointer.From(topic.Retain))
	}
}
