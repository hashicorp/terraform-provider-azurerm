package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowprofile"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataflowProfileResource struct{}

var _ sdk.ResourceWithUpdate = DataflowProfileResource{}

type DataflowProfileModel struct {
	Name              string                           `tfschema:"name"`
	ResourceGroupName string                           `tfschema:"resource_group_name"`
	InstanceName      string                           `tfschema:"instance_name"`
	InstanceCount     *int64                           `tfschema:"instance_count"`
	Diagnostics       *DataflowProfileDiagnosticsModel `tfschema:"diagnostics"`
	ExtendedLocation  ExtendedLocationModel            `tfschema:"extended_location"`
	ProvisioningState *string                          `tfschema:"provisioning_state"`
}

type DataflowProfileDiagnosticsModel struct {
	Logs    *DataflowProfileDiagnosticsLogsModel    `tfschema:"logs"`
	Metrics *DataflowProfileDiagnosticsMetricsModel `tfschema:"metrics"`
}

type DataflowProfileDiagnosticsLogsModel struct {
	Level *string `tfschema:"level"`
}

type DataflowProfileDiagnosticsMetricsModel struct {
	PrometheusPort *int64 `tfschema:"prometheus_port"`
}

func (r DataflowProfileResource) ModelObject() interface{} {
	return &DataflowProfileModel{}
}

func (r DataflowProfileResource) ResourceType() string {
	return "azurerm_iotoperations_dataflow_profile"
}

func (r DataflowProfileResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataflowprofile.ValidateDataflowProfileID
}

func (r DataflowProfileResource) Arguments() map[string]*pluginsdk.Schema {
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
		"instance_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 1000),
		},
		"diagnostics": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"logs": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"level": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										"trace",
										"debug",
										"info",
										"warn",
										"error",
									}, false),
								},
							},
						},
					},
					"metrics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"prometheus_port": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r DataflowProfileResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataflowProfileResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowProfileClient

			var model DataflowProfileModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := dataflowprofile.NewDataflowProfileID(subscriptionId, model.ResourceGroupName, model.InstanceName, model.Name)

			// Build payload
			payload := dataflowprofile.DataflowProfileResource{
				ExtendedLocation: expandDataflowProfileExtendedLocation(model.ExtendedLocation),
				Properties:       expandDataflowProfileProperties(model),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataflowProfileResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowProfileClient

			id, err := dataflowprofile.ParseDataflowProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := DataflowProfileModel{
				Name:              id.DataflowProfileName,
				ResourceGroupName: id.ResourceGroupName,
				InstanceName:      id.InstanceName,
			}

			if respModel := resp.Model; respModel != nil {
				model.ExtendedLocation = flattenDataflowProfileExtendedLocation(respModel.ExtendedLocation)

				if respModel.Properties != nil {
					flattenDataflowProfileProperties(respModel.Properties, &model)

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

func (r DataflowProfileResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowProfileClient

			id, err := dataflowprofile.ParseDataflowProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DataflowProfileModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// For dataflow profile, we use CreateOrUpdate for updates since there's no dedicated Update method
			payload := dataflowprofile.DataflowProfileResource{
				ExtendedLocation: expandDataflowProfileExtendedLocation(model.ExtendedLocation),
				Properties:       expandDataflowProfileProperties(model),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DataflowProfileResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.DataflowProfileClient

			id, err := dataflowprofile.ParseDataflowProfileID(metadata.ResourceData.Id())
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
func expandDataflowProfileExtendedLocation(input ExtendedLocationModel) dataflowprofile.ExtendedLocation {
	return dataflowprofile.ExtendedLocation{
		Name: *input.Name,
		Type: dataflowprofile.ExtendedLocationType(*input.Type),
	}
}

func flattenDataflowProfileExtendedLocation(input dataflowprofile.ExtendedLocation) ExtendedLocationModel {
	typeStr := string(input.Type)
	return ExtendedLocationModel{
		Name: &input.Name,
		Type: &typeStr,
	}
}

func expandDataflowProfileProperties(model DataflowProfileModel) *dataflowprofile.DataflowProfileProperties {
	props := &dataflowprofile.DataflowProfileProperties{}

	if model.InstanceCount != nil {
		props.InstanceCount = model.InstanceCount
	}

	if model.Diagnostics != nil {
		props.Diagnostics = expandDataflowProfileDiagnostics(*model.Diagnostics)
	}

	return props
}

func expandDataflowProfileDiagnostics(model DataflowProfileDiagnosticsModel) *dataflowprofile.ProfileDiagnostics {
	result := &dataflowprofile.ProfileDiagnostics{}

	if model.Logs != nil {
		result.Logs = &dataflowprofile.DiagnosticsLogs{
			Level: model.Logs.Level,
		}
	}

	if model.Metrics != nil {
		result.Metrics = &dataflowprofile.Metrics{
			PrometheusPort: model.Metrics.PrometheusPort,
		}
	}

	return result
}

func flattenDataflowProfileProperties(props *dataflowprofile.DataflowProfileProperties, model *DataflowProfileModel) {
	if props == nil {
		return
	}

	if props.InstanceCount != nil {
		model.InstanceCount = props.InstanceCount
	}

	if props.Diagnostics != nil {
		model.Diagnostics = flattenDataflowProfileDiagnostics(*props.Diagnostics)
	}
}

func flattenDataflowProfileDiagnostics(diagnostics dataflowprofile.ProfileDiagnostics) *DataflowProfileDiagnosticsModel {
	result := &DataflowProfileDiagnosticsModel{}

	if diagnostics.Logs != nil {
		result.Logs = &DataflowProfileDiagnosticsLogsModel{
			Level: diagnostics.Logs.Level,
		}
	}

	if diagnostics.Metrics != nil {
		result.Metrics = &DataflowProfileDiagnosticsMetricsModel{
			PrometheusPort: diagnostics.Metrics.PrometheusPort,
		}
	}

	return result
}
