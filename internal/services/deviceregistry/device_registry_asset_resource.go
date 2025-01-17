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
	Type                         string                 `tfschema:"type"`
	Tags                         map[string]string      `tfschema:"tags"`
	ExtendedLocationName         string                 `tfschema:"extended_location_name"`
	ExtendedLocationType         string                 `tfschema:"extended_location_type"`
	ProvisioningState            string                 `tfschema:"provisioning_state"`
	Uuid                         string                 `tfschema:"uuid"`
	Enabled                      bool                   `tfschema:"enabled"`
	ExternalAssetId              string                 `tfschema:"external_asset_id"`
	DisplayName                  string                 `tfschema:"display_name"`
	Description                  string                 `tfschema:"description"`
	AssetEndpointProfileRef      string                 `tfschema:"asset_endpoint_profile_ref"`
	Version                      int64                  `tfschema:"version"`
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
	DefaultTopic                 Topic                  `tfschema:"default_topic"`
	Datasets                     []Dataset              `tfschema:"datasets"`
	Events                       []Event                `tfschema:"events"`
	Status                       AssetStatus            `tfschema:"status"`
}

type Topic struct {
	Path   string `tfschema:"path"`
	Retain string `tfschema:"retain"`
}

type Dataset struct {
	Name                 string      `tfschema:"name"`
	DatasetConfiguration string      `tfschema:"dataset_configuration"`
	Topic                Topic       `tfschema:"topic"`
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
	Topic              Topic  `tfschema:"topic"`
}

type AssetStatus struct {
	Errors   []ErrorStatus   `tfschema:"errors"`
	Version  int64           `tfschema:"version"`
	Datasets []DatasetStatus `tfschema:"datasets"`
	Events   []EventStatus   `tfschema:"events"`
}

type ErrorStatus struct {
	Code    string `tfschema:"code"`
	Message string `tfschema:"message"`
}

type DatasetStatus struct {
	Name                   string                 `tfschema:"name"`
	MessageSchemaReference MessageSchemaReference `tfschema:"message_schema_reference"`
}

type EventStatus struct {
	Name                   string                 `tfschema:"name"`
	MessageSchemaReference MessageSchemaReference `tfschema:"message_schema_reference"`
}

type MessageSchemaReference struct {
	SchemaRegistryNamespace string `tfschema:"schema_registry_namespace"`
	SchemaName              string `tfschema:"schema_name"`
	SchemaVersion           string `tfschema:"schema_version"`
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
		"default_topic": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: map[string]*pluginsdk.Schema{
				"path": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"retain": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
		"datasets": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeMap,
				Elem: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"dataset_configuration": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"topic": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: map[string]*pluginsdk.Schema{
							"path": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
							"retain": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},
						},
					},
					"data_points": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeMap,
							Elem: map[string]*pluginsdk.Schema{
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
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeMap,
				Elem: map[string]*pluginsdk.Schema{
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
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: map[string]*pluginsdk.Schema{
							"path": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
							"retain": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}

func (AssetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"uuid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"version": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"status": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: map[string]*pluginsdk.Schema{
				"errors": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeMap,
						Elem: map[string]*pluginsdk.Schema{
							"code": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"message": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
				"version": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
				"datasets": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeMap,
						Elem: map[string]*pluginsdk.Schema{
							"name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"message_schema_reference": {
								Type:     pluginsdk.TypeMap,
								Computed: true,
								Elem: map[string]*pluginsdk.Schema{
									"schema_registry_namespace": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"schema_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"schema_version": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
				"events": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeMap,
						Elem: map[string]*pluginsdk.Schema{
							"name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"message_schema_reference": {
								Type:     pluginsdk.TypeMap,
								Computed: true,
								Elem: map[string]*pluginsdk.Schema{
									"schema_registry_namespace": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"schema_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"schema_version": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
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
					AssetEndpointProfileRef:      config.AssetEndpointProfileRef,
					Enabled:                      pointer.To(config.Enabled),
					ExternalAssetId:              pointer.To(config.ExternalAssetId),
					DisplayName:                  pointer.To(config.DisplayName),
					Description:                  pointer.To(config.Description),
					Manufacturer:                 pointer.To(config.Manufacturer),
					ManufacturerUri:              pointer.To(config.ManufacturerUri),
					Model:                        pointer.To(config.Model),
					ProductCode:                  pointer.To(config.ProductCode),
					HardwareRevision:             pointer.To(config.HardwareRevision),
					SoftwareRevision:             pointer.To(config.SoftwareRevision),
					DocumentationUri:             pointer.To(config.DocumentationUri),
					SerialNumber:                 pointer.To(config.SerialNumber),
					Attributes:                   pointer.To(config.Attributes),
					DiscoveredAssetRefs:          pointer.To(config.DiscoveredAssetRefs),
					DefaultDatasetsConfiguration: pointer.To(config.DefaultDatasetsConfiguration),
					DefaultEventsConfiguration:   pointer.To(config.DefaultEventsConfiguration),
				},
			}

			if config.DefaultTopic.Path != "" {
				param.Properties.DefaultTopic = toAzureTopic(config.DefaultTopic)
			}

			if len(config.Datasets) > 0 {
				param.Properties.Datasets = toAzureDatasets(config.Datasets)
			}

			if len(config.Events) > 0 {
				param.Properties.Events = toAzureEvents(config.Events)
			}

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

			id, err := assets.ParseAssetID(metadata.ResourceData.Get("id").(string))
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

			if metadata.ResourceData.HasChange("default_topic") {
				param.Properties.DefaultTopic = &assets.TopicUpdate{
					Path:   pointer.To(config.DefaultTopic.Path),
					Retain: pointer.To(assets.TopicRetainType(config.DefaultTopic.Retain)),
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
					state.Type = pointer.From(model.Type)
					state.ProvisioningState = string(pointer.From(props.ProvisioningState))
					state.Uuid = pointer.From(props.Uuid)
					state.Enabled = pointer.From(props.Enabled)
					state.ExternalAssetId = pointer.From(props.ExternalAssetId)
					state.DisplayName = pointer.From(props.DisplayName)
					state.Description = pointer.From(props.Description)
					state.Version = pointer.From(props.Version)
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
					state.DefaultTopic = toTFTopic(props.DefaultTopic)

					if datasets := props.Datasets; datasets != nil {
						state.Datasets = toTFDatasets(datasets)
					}

					if events := props.Events; events != nil {
						state.Events = toTFEvents(events)
					}

					if status := props.Status; status != nil {
						state.Status.Version = pointer.From(status.Version)
						state.Status.Errors = toTFAssetErrorStatuses(status.Errors)
						state.Status.Datasets = toTFDatasetStatuses(status.Datasets)
						state.Status.Events = toTFEventStatuses(status.Events)
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
	if len(datasets) == 0 {
		return nil
	}

	azureDatasets := make([]assets.Dataset, len(datasets))
	for i, dataset := range datasets {
		azureDatasets[i] = assets.Dataset{
			Name:                 dataset.Name,
			DatasetConfiguration: pointer.To(dataset.DatasetConfiguration),
			Topic:                toAzureTopic(dataset.Topic),
			DataPoints:           toAzureDataPoints(dataset.DataPoints),
		}
	}

	return &azureDatasets
}

func toAzureDataPoints(dataPoints []DataPoint) *[]assets.DataPoint {
	if len(dataPoints) == 0 {
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
	if len(events) == 0 {
		return nil
	}

	azureEvents := make([]assets.Event, len(events))
	for i, event := range events {
		azureEvents[i] = assets.Event{
			Name:               event.Name,
			EventNotifier:      event.EventNotifier,
			EventConfiguration: pointer.To(event.EventConfiguration),
			ObservabilityMode:  pointer.To(assets.EventObservabilityMode(event.ObservabilityMode)),
			Topic:              toAzureTopic(event.Topic),
		}
	}

	return &azureEvents
}

func toAzureTopic(topic Topic) *assets.Topic {
	if topic.Path == "" {
		return nil
	}

	azureTopic := assets.Topic{
		Path: topic.Path,
	}

	if topic.Retain != "" {
		azureTopic.Retain = pointer.To(assets.TopicRetainType(topic.Retain))
	}

	return &azureTopic
}

func toTFDatasets(datasets *[]assets.Dataset) []Dataset {
	if datasets == nil {
		return nil
	}

	tfDatasets := make([]Dataset, len(*datasets))
	for i, dataset := range *datasets {
		tfDatasets[i] = Dataset{
			Name:                 dataset.Name,
			DatasetConfiguration: pointer.From(dataset.DatasetConfiguration),
			Topic:                toTFTopic(dataset.Topic),
			DataPoints:           toTFDataPoints(dataset.DataPoints),
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
		tfEvents[i] = Event{
			Name:               event.Name,
			EventNotifier:      event.EventNotifier,
			ObservabilityMode:  string(pointer.From(event.ObservabilityMode)),
			EventConfiguration: pointer.From(event.EventConfiguration),
			Topic:              toTFTopic(event.Topic),
		}
	}

	return tfEvents
}

func toTFTopic(topic *assets.Topic) Topic {
	if topic == nil {
		return Topic{}
	}

	return Topic{
		Path:   topic.Path,
		Retain: string(pointer.From(topic.Retain)),
	}
}

func toTFAssetErrorStatuses(errorStatuses *[]assets.AssetStatusError) []ErrorStatus {
	if errorStatuses == nil {
		return nil
	}

	tfErrorStatuses := make([]ErrorStatus, len(*errorStatuses))
	for i, errorStatus := range *errorStatuses {
		tfErrorStatuses[i] = ErrorStatus{
			Code:    string(pointer.From(errorStatus.Code)),
			Message: pointer.From(errorStatus.Message),
		}
	}

	return tfErrorStatuses
}

func toTFDatasetStatuses(datasetStatuses *[]assets.AssetStatusDataset) []DatasetStatus {
	if datasetStatuses == nil {
		return nil
	}

	tfDatasetStatuses := make([]DatasetStatus, len(*datasetStatuses))
	for i, datasetStatus := range *datasetStatuses {
		tfDatasetStatuses[i] = DatasetStatus{
			Name:                   datasetStatus.Name,
			MessageSchemaReference: toTFMessageSchemaReference(datasetStatus.MessageSchemaReference),
		}
	}

	return tfDatasetStatuses
}

func toTFEventStatuses(eventStatuses *[]assets.AssetStatusEvent) []EventStatus {
	if eventStatuses == nil {
		return nil
	}

	tfEventStatuses := make([]EventStatus, len(*eventStatuses))
	for i, eventStatus := range *eventStatuses {
		tfEventStatuses[i] = EventStatus{
			Name:                   eventStatus.Name,
			MessageSchemaReference: toTFMessageSchemaReference(eventStatus.MessageSchemaReference),
		}
	}

	return tfEventStatuses
}

func toTFMessageSchemaReference(messageSchemaReference *assets.MessageSchemaReference) MessageSchemaReference {
	if messageSchemaReference == nil {
		return MessageSchemaReference{}
	}

	return MessageSchemaReference{
		SchemaRegistryNamespace: messageSchemaReference.SchemaRegistryNamespace,
		SchemaName:              messageSchemaReference.SchemaName,
		SchemaVersion:           messageSchemaReference.SchemaVersion,
	}
}
