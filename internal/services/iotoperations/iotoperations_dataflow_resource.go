package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataflowResource struct{}

var _ sdk.ResourceWithUpdate = DataflowResource{}

type DataflowModel struct {
	Name                 string                   `tfschema:"name"`
	ResourceGroupName    string                   `tfschema:"resource_group_name"`
	InstanceName         string                   `tfschema:"instance_name"`
	DataflowProfileName  string                   `tfschema:"dataflow_profile_name"`
	Mode                 *string                  `tfschema:"mode"`
	Operations           []DataflowOperationModel `tfschema:"operations"`
	ExtendedLocationName *string                  `tfschema:"extended_location_name"`
	ExtendedLocationType *string                  `tfschema:"extended_location_type"`
	ProvisioningState    *string                  `tfschema:"provisioning_state"`
}

type DataflowOperationModel struct {
	Name                          *string                                     `tfschema:"name"`
	OperationType                 string                                      `tfschema:"operation_type"`
	SourceSettings                *DataflowSourceOperationSettingsModel       `tfschema:"source_settings"`
	DestinationSettings           *DataflowDestinationOperationSettingsModel  `tfschema:"destination_settings"`
	BuiltInTransformationSettings *DataflowBuiltInTransformationSettingsModel `tfschema:"built_in_transformation_settings"`
}

type DataflowSourceOperationSettingsModel struct {
	DataSources         []string `tfschema:"data_sources"`
	AssetRef            *string  `tfschema:"asset_ref"`
	EndpointRef         string   `tfschema:"endpoint_ref"`
	SchemaRef           *string  `tfschema:"schema_ref"`
	SerializationFormat *string  `tfschema:"serialization_format"`
}

type DataflowDestinationOperationSettingsModel struct {
	DataDestination string `tfschema:"data_destination"`
	EndpointRef     string `tfschema:"endpoint_ref"`
}

type DataflowBuiltInTransformationSettingsModel struct {
	Datasets            []DataflowBuiltInTransformationDatasetModel `tfschema:"datasets"`
	Filter              []DataflowBuiltInTransformationFilterModel  `tfschema:"filter"`
	Map                 []DataflowBuiltInTransformationMapModel     `tfschema:"map"`
	SchemaRef           *string                                     `tfschema:"schema_ref"`
	SerializationFormat *string                                     `tfschema:"serialization_format"`
}

type DataflowBuiltInTransformationFilterModel struct {
	Description *string  `tfschema:"description"`
	Expression  string   `tfschema:"expression"`
	Inputs      []string `tfschema:"inputs"`
	Type        *string  `tfschema:"type"`
}

type DataflowBuiltInTransformationMapModel struct {
	Description *string  `tfschema:"description"`
	Expression  *string  `tfschema:"expression"`
	Inputs      []string `tfschema:"inputs"`
	Output      string   `tfschema:"output"`
	Type        *string  `tfschema:"type"`
}

type DataflowBuiltInTransformationDatasetModel struct {
	Key         string   `tfschema:"key"`
	Description *string  `tfschema:"description"`
	Expression  *string  `tfschema:"expression"`
	Inputs      []string `tfschema:"inputs"`
	SchemaRef   *string  `tfschema:"schema_ref"`
}

func (r DataflowResource) ModelObject() interface{} {
	return &DataflowModel{}
}

func (r DataflowResource) ResourceType() string {
	return "azurerm_iotoperations_dataflow"
}

func (r DataflowResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataflow.ValidateDataflowID
}

func (r DataflowResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 90),
		},
		"instance_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"dataflow_profile_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "Enabled",
			ValidateFunc: validation.StringInSlice([]string{
				"Enabled",
				"Disabled",
			}, false),
		},
		"operations": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 63),
					},
					"operation_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Source",
							"Destination",
							"BuiltInTransformation",
						}, false),
					},
					"source_settings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_sources": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringLenBetween(1, 253),
									},
								},
								"asset_ref": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
								"endpoint_ref": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
								"schema_ref": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
								"serialization_format": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Json",
									}, false),
								},
							},
						},
					},
					"destination_settings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_destination": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
								"endpoint_ref": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
							},
						},
					},
					"built_in_transformation_settings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"datasets": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"key": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 253),
											},
											"description": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 500),
											},
											"expression": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 1000),
											},
											"inputs": {
												Type:     pluginsdk.TypeList,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringLenBetween(1, 253),
												},
											},
											"schema_ref": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 253),
											},
										},
									},
								},
								"filter": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"description": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 500),
											},
											"expression": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 1000),
											},
											"inputs": {
												Type:     pluginsdk.TypeList,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringLenBetween(1, 253),
												},
											},
											"type": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												ValidateFunc: validation.StringInSlice([]string{
													"Filter",
												}, false),
											},
										},
									},
								},
								"map": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"description": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 500),
											},
											"expression": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 1000),
											},
											"inputs": {
												Type:     pluginsdk.TypeList,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringLenBetween(1, 253),
												},
											},
											"output": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 253),
											},
											"type": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												ValidateFunc: validation.StringInSlice([]string{
													"BuiltInFunction",
													"Compute",
													"NewProperties",
													"PassThrough",
													"Rename",
												}, false),
											},
										},
									},
								},
								"schema_ref": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
								"serialization_format": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Delta",
										"Json",
										"Parquet",
									}, false),
								},
							},
						},
					},
				},
			},
		},
		"extended_location": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"CustomLocation",
						}, false),
					},
				},
			},
		},
	}
}

func (r DataflowResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataflowResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowClient

			var model DataflowModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := dataflow.NewDataflowID(subscriptionId, model.ResourceGroupName, model.InstanceName, model.DataflowProfileName, model.Name)

			// Build payload
			extendedLocationName := ""
			if model.ExtendedLocationName != nil {
				extendedLocationName = *model.ExtendedLocationName
			}
			extendedLocationType := dataflow.ExtendedLocationTypeCustomLocation
			if model.ExtendedLocationType != nil {
				extendedLocationType = dataflow.ExtendedLocationType(*model.ExtendedLocationType)
			}
			// no new variables on left side of :=
			extendedLocationType = dataflow.ExtendedLocationTypeCustomLocation
			if model.ExtendedLocationType != nil {
				extendedLocationType = dataflow.ExtendedLocationType(*model.ExtendedLocationType)
			}
			payload := dataflow.DataflowResource{
				ExtendedLocation: dataflow.ExtendedLocation{
					Name: extendedLocationName,
					Type: extendedLocationType,
				},
				Properties: expandDataflowProperties(model),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataflowResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowClient

			id, err := dataflow.ParseDataflowID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := DataflowModel{
				Name:                id.DataflowName,
				ResourceGroupName:   id.ResourceGroupName,
				InstanceName:        id.InstanceName,
				DataflowProfileName: id.DataflowProfileName,
			}

			if respModel := resp.Model; respModel != nil {
				model.ExtendedLocationName = &respModel.ExtendedLocation.Name
				extendedLocationType := string(respModel.ExtendedLocation.Type)
				model.ExtendedLocationType = &extendedLocationType

				if respModel.Properties != nil {
					flattenDataflowProperties(respModel.Properties, &model)

					if respModel.Properties.ProvisioningState != nil {
						provisioningState := string(*respModel.Properties.ProvisioningState)
						model.ProvisioningState = &provisioningState
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r DataflowResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowClient

			id, err := dataflow.ParseDataflowID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DataflowModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// For dataflow, we use CreateOrUpdate for updates since there's no dedicated Update method
			extendedLocationName := ""
			if model.ExtendedLocationName != nil {
				extendedLocationName = *model.ExtendedLocationName
			}
			extendedLocationType := dataflow.ExtendedLocationTypeCustomLocation
			if model.ExtendedLocationType != nil {
				extendedLocationType = dataflow.ExtendedLocationType(*model.ExtendedLocationType)
			}

			payload := dataflow.DataflowResource{
				ExtendedLocation: dataflow.ExtendedLocation{
					Name: extendedLocationName,
					Type: extendedLocationType,
				},
				Properties: expandDataflowProperties(model),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DataflowResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowClient

			id, err := dataflow.ParseDataflowID(metadata.ResourceData.Id())
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

// Helper functions for expand/flatten operations
// expandDataflowExtendedLocation and flattenDataflowExtendedLocation removed; now handled inline with separate fields

func expandDataflowProperties(model DataflowModel) *dataflow.DataflowProperties {
	props := &dataflow.DataflowProperties{
		Operations: expandDataflowOperations(model.Operations),
	}

	if model.Mode != nil {
		mode := dataflow.OperationalMode(*model.Mode)
		props.Mode = &mode
	}

	return props
}

func expandDataflowOperations(operations []DataflowOperationModel) []dataflow.DataflowOperation {
	result := make([]dataflow.DataflowOperation, 0, len(operations))

	for _, op := range operations {
		operation := dataflow.DataflowOperation{
			OperationType: dataflow.OperationType(op.OperationType),
		}

		if op.Name != nil {
			operation.Name = op.Name
		}

		if op.SourceSettings != nil {
			operation.SourceSettings = expandDataflowSourceOperationSettings(*op.SourceSettings)
		}

		if op.DestinationSettings != nil {
			operation.DestinationSettings = expandDataflowDestinationOperationSettings(*op.DestinationSettings)
		}

		if op.BuiltInTransformationSettings != nil {
			operation.BuiltInTransformationSettings = expandDataflowBuiltInTransformationSettings(*op.BuiltInTransformationSettings)
		}

		result = append(result, operation)
	}

	return result
}

func expandDataflowSourceOperationSettings(source DataflowSourceOperationSettingsModel) *dataflow.DataflowSourceOperationSettings {
	result := &dataflow.DataflowSourceOperationSettings{
		DataSources: source.DataSources,
		EndpointRef: source.EndpointRef,
	}

	if source.AssetRef != nil {
		result.AssetRef = source.AssetRef
	}

	if source.SchemaRef != nil {
		result.SchemaRef = source.SchemaRef
	}

	if source.SerializationFormat != nil {
		format := dataflow.SourceSerializationFormat(*source.SerializationFormat)
		result.SerializationFormat = &format
	}

	return result
}

func expandDataflowDestinationOperationSettings(destination DataflowDestinationOperationSettingsModel) *dataflow.DataflowDestinationOperationSettings {
	return &dataflow.DataflowDestinationOperationSettings{
		DataDestination: destination.DataDestination,
		EndpointRef:     destination.EndpointRef,
	}
}

func expandDataflowBuiltInTransformationSettings(settings DataflowBuiltInTransformationSettingsModel) *dataflow.DataflowBuiltInTransformationSettings {
	result := &dataflow.DataflowBuiltInTransformationSettings{}

	if len(settings.Datasets) > 0 {
		result.Datasets = expandDataflowBuiltInTransformationDatasets(settings.Datasets)
	}

	if len(settings.Filter) > 0 {
		result.Filter = expandDataflowBuiltInTransformationFilters(settings.Filter)
	}

	if len(settings.Map) > 0 {
		result.Map = expandDataflowBuiltInTransformationMaps(settings.Map)
	}

	if settings.SchemaRef != nil {
		result.SchemaRef = settings.SchemaRef
	}

	if settings.SerializationFormat != nil {
		format := dataflow.TransformationSerializationFormat(*settings.SerializationFormat)
		result.SerializationFormat = &format
	}

	return result
}

func expandDataflowBuiltInTransformationDatasets(datasets []DataflowBuiltInTransformationDatasetModel) *[]dataflow.DataflowBuiltInTransformationDataset {
	result := make([]dataflow.DataflowBuiltInTransformationDataset, 0, len(datasets))

	for _, dataset := range datasets {
		datasetItem := dataflow.DataflowBuiltInTransformationDataset{
			Key:    dataset.Key,
			Inputs: dataset.Inputs,
		}

		if dataset.Description != nil {
			datasetItem.Description = dataset.Description
		}

		if dataset.Expression != nil {
			datasetItem.Expression = dataset.Expression
		}

		if dataset.SchemaRef != nil {
			datasetItem.SchemaRef = dataset.SchemaRef
		}

		result = append(result, datasetItem)
	}

	return &result
}

func expandDataflowBuiltInTransformationFilters(filters []DataflowBuiltInTransformationFilterModel) *[]dataflow.DataflowBuiltInTransformationFilter {
	result := make([]dataflow.DataflowBuiltInTransformationFilter, 0, len(filters))

	for _, filter := range filters {
		filterItem := dataflow.DataflowBuiltInTransformationFilter{
			Expression: filter.Expression,
			Inputs:     filter.Inputs,
		}

		if filter.Description != nil {
			filterItem.Description = filter.Description
		}

		if filter.Type != nil {
			filterType := dataflow.FilterType(*filter.Type)
			filterItem.Type = &filterType
		}

		result = append(result, filterItem)
	}

	return &result
}

func expandDataflowBuiltInTransformationMaps(maps []DataflowBuiltInTransformationMapModel) *[]dataflow.DataflowBuiltInTransformationMap {
	result := make([]dataflow.DataflowBuiltInTransformationMap, 0, len(maps))

	for _, mapItem := range maps {
		mapTransform := dataflow.DataflowBuiltInTransformationMap{
			Inputs: mapItem.Inputs,
			Output: mapItem.Output,
		}

		if mapItem.Description != nil {
			mapTransform.Description = mapItem.Description
		}

		if mapItem.Expression != nil {
			mapTransform.Expression = mapItem.Expression
		}

		if mapItem.Type != nil {
			mapType := dataflow.DataflowMappingType(*mapItem.Type)
			mapTransform.Type = &mapType
		}

		result = append(result, mapTransform)
	}

	return &result
}

func flattenDataflowProperties(props *dataflow.DataflowProperties, model *DataflowModel) {
	if props == nil {
		return
	}

	if props.Mode != nil {
		mode := string(*props.Mode)
		model.Mode = &mode
	}

	if len(props.Operations) > 0 {
		model.Operations = flattenDataflowOperations(props.Operations)
	}
}

func flattenDataflowOperations(operations []dataflow.DataflowOperation) []DataflowOperationModel {
	result := make([]DataflowOperationModel, 0, len(operations))

	for _, op := range operations {
		operation := DataflowOperationModel{
			OperationType: string(op.OperationType),
		}

		if op.Name != nil {
			operation.Name = op.Name
		}

		if op.SourceSettings != nil {
			operation.SourceSettings = flattenDataflowSourceOperationSettings(*op.SourceSettings)
		}

		if op.DestinationSettings != nil {
			operation.DestinationSettings = flattenDataflowDestinationOperationSettings(*op.DestinationSettings)
		}

		if op.BuiltInTransformationSettings != nil {
			operation.BuiltInTransformationSettings = flattenDataflowBuiltInTransformationSettings(*op.BuiltInTransformationSettings)
		}

		result = append(result, operation)
	}

	return result
}

func flattenDataflowSourceOperationSettings(source dataflow.DataflowSourceOperationSettings) *DataflowSourceOperationSettingsModel {
	result := &DataflowSourceOperationSettingsModel{
		DataSources: source.DataSources,
		EndpointRef: source.EndpointRef,
	}

	if source.AssetRef != nil {
		result.AssetRef = source.AssetRef
	}

	if source.SchemaRef != nil {
		result.SchemaRef = source.SchemaRef
	}

	if source.SerializationFormat != nil {
		format := string(*source.SerializationFormat)
		result.SerializationFormat = &format
	}

	return result
}

func flattenDataflowDestinationOperationSettings(destination dataflow.DataflowDestinationOperationSettings) *DataflowDestinationOperationSettingsModel {
	return &DataflowDestinationOperationSettingsModel{
		DataDestination: destination.DataDestination,
		EndpointRef:     destination.EndpointRef,
	}
}

func flattenDataflowBuiltInTransformationSettings(settings dataflow.DataflowBuiltInTransformationSettings) *DataflowBuiltInTransformationSettingsModel {
	result := &DataflowBuiltInTransformationSettingsModel{}

	if settings.Datasets != nil {
		result.Datasets = flattenDataflowBuiltInTransformationDatasets(*settings.Datasets)
	}

	if settings.Filter != nil {
		result.Filter = flattenDataflowBuiltInTransformationFilters(*settings.Filter)
	}

	if settings.Map != nil {
		result.Map = flattenDataflowBuiltInTransformationMaps(*settings.Map)
	}

	if settings.SchemaRef != nil {
		result.SchemaRef = settings.SchemaRef
	}

	if settings.SerializationFormat != nil {
		format := string(*settings.SerializationFormat)
		result.SerializationFormat = &format
	}

	return result
}

func flattenDataflowBuiltInTransformationDatasets(datasets []dataflow.DataflowBuiltInTransformationDataset) []DataflowBuiltInTransformationDatasetModel {
	result := make([]DataflowBuiltInTransformationDatasetModel, 0, len(datasets))

	for _, dataset := range datasets {
		datasetModel := DataflowBuiltInTransformationDatasetModel{
			Key:    dataset.Key,
			Inputs: dataset.Inputs,
		}

		if dataset.Description != nil {
			datasetModel.Description = dataset.Description
		}

		if dataset.Expression != nil {
			datasetModel.Expression = dataset.Expression
		}

		if dataset.SchemaRef != nil {
			datasetModel.SchemaRef = dataset.SchemaRef
		}

		result = append(result, datasetModel)
	}

	return result
}

func flattenDataflowBuiltInTransformationFilters(filters []dataflow.DataflowBuiltInTransformationFilter) []DataflowBuiltInTransformationFilterModel {
	result := make([]DataflowBuiltInTransformationFilterModel, 0, len(filters))

	for _, filter := range filters {
		filterModel := DataflowBuiltInTransformationFilterModel{
			Expression: filter.Expression,
			Inputs:     filter.Inputs,
		}

		if filter.Description != nil {
			filterModel.Description = filter.Description
		}

		if filter.Type != nil {
			filterType := string(*filter.Type)
			filterModel.Type = &filterType
		}

		result = append(result, filterModel)
	}

	return result
}

func flattenDataflowBuiltInTransformationMaps(maps []dataflow.DataflowBuiltInTransformationMap) []DataflowBuiltInTransformationMapModel {
	result := make([]DataflowBuiltInTransformationMapModel, 0, len(maps))

	for _, mapItem := range maps {
		mapModel := DataflowBuiltInTransformationMapModel{
			Inputs: mapItem.Inputs,
			Output: mapItem.Output,
		}

		if mapItem.Description != nil {
			mapModel.Description = mapItem.Description
		}

		if mapItem.Expression != nil {
			mapModel.Expression = mapItem.Expression
		}

		if mapItem.Type != nil {
			mapType := string(*mapItem.Type)
			mapModel.Type = &mapType
		}

		result = append(result, mapModel)
	}

	return result
}
