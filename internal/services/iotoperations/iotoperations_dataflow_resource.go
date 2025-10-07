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
	Name              string                    `tfschema:"name"`
	ResourceGroupName string                    `tfschema:"resource_group_name"`
	InstanceName      string                    `tfschema:"instance_name"`
	DataflowProfileName string                  `tfschema:"dataflow_profile_name"`
	Mode              *string                   `tfschema:"mode"`
	Operations        []DataflowOperationModel  `tfschema:"operations"`
	ExtendedLocation  *ExtendedLocationModel    `tfschema:"extended_location"`
	Tags              map[string]string         `tfschema:"tags"`
	ProvisioningState *string                   `tfschema:"provisioning_state"`
}

type DataflowOperationModel struct {
	Name            string                             `tfschema:"name"`
	OperationType   string                             `tfschema:"operation_type"`
	Source          *DataflowOperationSourceModel      `tfschema:"source"`
	Destination     *DataflowOperationDestinationModel `tfschema:"destination"`
	BuiltInTransformations []DataflowBuiltInTransformationModel `tfschema:"built_in_transformations"`
}

type DataflowOperationSourceModel struct {
	DataSource   string  `tfschema:"data_source"`
	AssetRef     *string `tfschema:"asset_ref"`
	EndpointRef  string  `tfschema:"endpoint_ref"`
	SchemaRef    *string `tfschema:"schema_ref"`
	SerializationFormat *string `tfschema:"serialization_format"`
}

type DataflowOperationDestinationModel struct {
	DataDestination string  `tfschema:"data_destination"`
	EndpointRef     string  `tfschema:"endpoint_ref"`
	SchemaRef       *string `tfschema:"schema_ref"`
	SerializationFormat *string `tfschema:"serialization_format"`
}

type DataflowBuiltInTransformationModel struct {
	Filter         []DataflowFilterModel         `tfschema:"filter"`
	Map            []DataflowMapModel            `tfschema:"map"`
	Datasets       []DataflowDatasetModel        `tfschema:"datasets"`
	SerializationFormat *string                  `tfschema:"serialization_format"`
	SchemaRef      *string                       `tfschema:"schema_ref"`
}

type DataflowFilterModel struct {
	Type        string                        `tfschema:"type"`
	Description *string                       `tfschema:"description"`
	Inputs      []string                      `tfschema:"inputs"`
	Expression  string                        `tfschema:"expression"`
}

type DataflowMapModel struct {
	Type        string                        `tfschema:"type"`
	Description *string                       `tfschema:"description"`
	Inputs      []string                      `tfschema:"inputs"`
	Output      string                        `tfschema:"output"`
	Expression  string                        `tfschema:"expression"`
}

type DataflowDatasetModel struct {
	Key         string                        `tfschema:"key"`
	Description *string                       `tfschema:"description"`
	Inputs      []string                      `tfschema:"inputs"`
	Expression  string                        `tfschema:"expression"`
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
						Required:     true,
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
					"source": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_source": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
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
										"JSON",
										"Parquet",
										"Delta",
									}, false),
								},
							},
						},
					},
					"destination": {
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
								"schema_ref": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
								"serialization_format": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										"JSON",
										"Parquet",
										"Delta",
									}, false),
								},
							},
						},
					},
					"built_in_transformations": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"filter": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"type": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 63),
											},
											"description": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 500),
											},
											"inputs": {
												Type:     pluginsdk.TypeList,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringLenBetween(1, 253),
												},
											},
											"expression": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 1000),
											},
										},
									},
								},
								"map": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"type": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 63),
											},
											"description": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 500),
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
											"expression": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 1000),
											},
										},
									},
								},
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
											"inputs": {
												Type:     pluginsdk.TypeList,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringLenBetween(1, 253),
												},
											},
											"expression": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 1000),
											},
										},
									},
								},
								"serialization_format": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										"JSON",
										"Parquet",
										"Delta",
									}, false),
								},
								"schema_ref": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringLenBetween(1, 253),
								},
							},
						},
					},
				},
			},
		},
		"extended_location": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
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
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r DataflowResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			// NOTE: O+C Azure automatically assigns provisioning state during resource lifecycle
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
			payload := dataflow.DataflowResource{
				Properties: expandDataflowProperties(model),
			}

			if model.ExtendedLocation != nil {
				payload.ExtendedLocation = expandExtendedLocation(model.ExtendedLocation)
			}

			if len(model.Tags) > 0 {
				payload.Tags = &model.Tags
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
				if respModel.ExtendedLocation != nil {
					model.ExtendedLocation = flattenExtendedLocation(respModel.ExtendedLocation)
				}

				if respModel.Tags != nil {
					model.Tags = *respModel.Tags
				}

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

			// Check if anything actually changed before making API call
			if !metadata.ResourceData.HasChange("tags") && 
			   !metadata.ResourceData.HasChange("mode") &&
			   !metadata.ResourceData.HasChange("operations") {
				return nil
			}

			payload := dataflow.DataflowPatchModel{}
			hasChanges := false

			// Only include tags if they changed
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
				hasChanges = true
			}

			// Only include properties if they changed
			if metadata.ResourceData.HasChange("mode") ||
			   metadata.ResourceData.HasChange("operations") {
				patchProps := &dataflow.DataflowPropertiesPatch{}
				
				if metadata.ResourceData.HasChange("mode") {
					if model.Mode != nil {
						mode := dataflow.OperationalMode(*model.Mode)
						patchProps.Mode = &mode
					}
				}

				if metadata.ResourceData.HasChange("operations") {
					patchProps.Operations = expandDataflowOperations(model.Operations)
				}

				payload.Properties = patchProps
				hasChanges = true
			}

			// Only make API call if something actually changed
			if !hasChanges {
				return nil
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
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

func expandDataflowOperations(operations []DataflowOperationModel) *[]dataflow.DataflowOperation {
	if len(operations) == 0 {
		return nil
	}

	result := make([]dataflow.DataflowOperation, 0, len(operations))

	for _, op := range operations {
		operation := dataflow.DataflowOperation{
			Name:          op.Name,
			OperationType: dataflow.OperationType(op.OperationType),
		}

		if op.Source != nil {
			operation.Source = expandDataflowOperationSource(*op.Source)
		}

		if op.Destination != nil {
			operation.Destination = expandDataflowOperationDestination(*op.Destination)
		}

		if len(op.BuiltInTransformations) > 0 {
			operation.BuiltInTransformations = expandDataflowBuiltInTransformations(op.BuiltInTransformations)
		}

		result = append(result, operation)
	}

	return &result
}

func expandDataflowOperationSource(source DataflowOperationSourceModel) *dataflow.DataflowSourceOperation {
	result := &dataflow.DataflowSourceOperation{
		DataSource:  source.DataSource,
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

func expandDataflowOperationDestination(destination DataflowOperationDestinationModel) *dataflow.DataflowDestinationOperation {
	result := &dataflow.DataflowDestinationOperation{
		DataDestination: destination.DataDestination,
		EndpointRef:     destination.EndpointRef,
	}

	if destination.SchemaRef != nil {
		result.SchemaRef = destination.SchemaRef
	}

	if destination.SerializationFormat != nil {
		format := dataflow.DestinationSerializationFormat(*destination.SerializationFormat)
		result.SerializationFormat = &format
	}

	return result
}

func expandDataflowBuiltInTransformations(transformations []DataflowBuiltInTransformationModel) *[]dataflow.DataflowBuiltInTransformation {
	if len(transformations) == 0 {
		return nil
	}

	result := make([]dataflow.DataflowBuiltInTransformation, 0, len(transformations))

	for _, transform := range transformations {
		transformation := dataflow.DataflowBuiltInTransformation{}

		if len(transform.Filter) > 0 {
			transformation.Filter = expandDataflowFilters(transform.Filter)
		}

		if len(transform.Map) > 0 {
			transformation.Map = expandDataflowMaps(transform.Map)
		}

		if len(transform.Datasets) > 0 {
			transformation.Datasets = expandDataflowDatasets(transform.Datasets)
		}

		if transform.SerializationFormat != nil {
			format := dataflow.TransformationSerializationFormat(*transform.SerializationFormat)
			transformation.SerializationFormat = &format
		}

		if transform.SchemaRef != nil {
			transformation.SchemaRef = transform.SchemaRef
		}

		result = append(result, transformation)
	}

	return &result
}

func expandDataflowFilters(filters []DataflowFilterModel) *[]dataflow.DataflowBuiltInTransformationFilter {
	result := make([]dataflow.DataflowBuiltInTransformationFilter, 0, len(filters))

	for _, filter := range filters {
		filterItem := dataflow.DataflowBuiltInTransformationFilter{
			Type:       filter.Type,
			Inputs:     filter.Inputs,
			Expression: filter.Expression,
		}

		if filter.Description != nil {
			filterItem.Description = filter.Description
		}

		result = append(result, filterItem)
	}

	return &result
}

func expandDataflowMaps(maps []DataflowMapModel) *[]dataflow.DataflowBuiltInTransformationMap {
	result := make([]dataflow.DataflowBuiltInTransformationMap, 0, len(maps))

	for _, mapItem := range maps {
		mapTransform := dataflow.DataflowBuiltInTransformationMap{
			Type:       mapItem.Type,
			Inputs:     mapItem.Inputs,
			Output:     mapItem.Output,
			Expression: mapItem.Expression,
		}

		if mapItem.Description != nil {
			mapTransform.Description = mapItem.Description
		}

		result = append(result, mapTransform)
	}

	return &result
}

func expandDataflowDatasets(datasets []DataflowDatasetModel) *[]dataflow.DataflowBuiltInTransformationDataset {
	result := make([]dataflow.DataflowBuiltInTransformationDataset, 0, len(datasets))

	for _, dataset := range datasets {
		datasetItem := dataflow.DataflowBuiltInTransformationDataset{
			Key:        dataset.Key,
			Inputs:     dataset.Inputs,
			Expression: dataset.Expression,
		}

		if dataset.Description != nil {
			datasetItem.Description = dataset.Description
		}

		result = append(result, datasetItem)
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

	if props.Operations != nil {
		model.Operations = flattenDataflowOperations(*props.Operations)
	}
}

func flattenDataflowOperations(operations []dataflow.DataflowOperation) []DataflowOperationModel {
	result := make([]DataflowOperationModel, 0, len(operations))

	for _, op := range operations {
		operation := DataflowOperationModel{
			Name:          op.Name,
			OperationType: string(op.OperationType),
		}

		if op.Source != nil {
			operation.Source = flattenDataflowOperationSource(*op.Source)
		}

		if op.Destination != nil {
			operation.Destination = flattenDataflowOperationDestination(*op.Destination)
		}

		if op.BuiltInTransformations != nil {
			operation.BuiltInTransformations = flattenDataflowBuiltInTransformations(*op.BuiltInTransformations)
		}

		result = append(result, operation)
	}

	return result
}

func flattenDataflowOperationSource(source dataflow.DataflowSourceOperation) *DataflowOperationSourceModel {
	result := &DataflowOperationSourceModel{
		DataSource:  source.DataSource,
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

func flattenDataflowOperationDestination(destination dataflow.DataflowDestinationOperation) *DataflowOperationDestinationModel {
	result := &DataflowOperationDestinationModel{
		DataDestination: destination.DataDestination,
		EndpointRef:     destination.EndpointRef,
	}

	if destination.SchemaRef != nil {
		result.SchemaRef = destination.SchemaRef
	}

	if destination.SerializationFormat != nil {
		format := string(*destination.SerializationFormat)
		result.SerializationFormat = &format
	}

	return result
}

func flattenDataflowBuiltInTransformations(transformations []dataflow.DataflowBuiltInTransformation) []DataflowBuiltInTransformationModel {
	result := make([]DataflowBuiltInTransformationModel, 0, len(transformations))

	for _, transform := range transformations {
		transformation := DataflowBuiltInTransformationModel{}

		if transform.Filter != nil {
			transformation.Filter = flattenDataflowFilters(*transform.Filter)
		}

		if transform.Map != nil {
			transformation.Map = flattenDataflowMaps(*transform.Map)
		}

		if transform.Datasets != nil {
			transformation.Datasets = flattenDataflowDatasets(*transform.Datasets)
		}

		if transform.SerializationFormat != nil {
			format := string(*transform.SerializationFormat)
			transformation.SerializationFormat = &format
		}

		if transform.SchemaRef != nil {
			transformation.SchemaRef = transform.SchemaRef
		}

		result = append(result, transformation)
	}

	return result
}

func flattenDataflowFilters(filters []dataflow.DataflowBuiltInTransformationFilter) []DataflowFilterModel {
	result := make([]DataflowFilterModel, 0, len(filters))

	for _, filter := range filters {
		filterModel := DataflowFilterModel{
			Type:       filter.Type,
			Inputs:     filter.Inputs,
			Expression: filter.Expression,
		}

		if filter.Description != nil {
			filterModel.Description = filter.Description
		}

		result = append(result, filterModel)
	}

	return result
}

func flattenDataflowMaps(maps []dataflow.DataflowBuiltInTransformationMap) []DataflowMapModel {
	result := make([]DataflowMapModel, 0, len(maps))

	for _, mapItem := range maps {
		mapModel := DataflowMapModel{
			Type:       mapItem.Type,
			Inputs:     mapItem.Inputs,
			Output:     mapItem.Output,
			Expression: mapItem.Expression,
		}

		if mapItem.Description != nil {
			mapModel.Description = mapItem.Description
		}

		result = append(result, mapModel)
	}

	return result
}

func flattenDataflowDatasets(datasets []dataflow.DataflowBuiltInTransformationDataset) []DataflowDatasetModel {
	result := make([]DataflowDatasetModel, 0, len(datasets))

	for _, dataset := range datasets {
		datasetModel := DataflowDatasetModel{
			Key:        dataset.Key,
			Inputs:     dataset.Inputs,
			Expression: dataset.Expression,
		}

		if dataset.Description != nil {
			datasetModel.Description = dataset.Description
		}

		result = append(result, datasetModel)
	}

	return result
}