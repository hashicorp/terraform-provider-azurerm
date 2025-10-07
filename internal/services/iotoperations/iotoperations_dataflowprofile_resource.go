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
	InstanceCount     *int                             `tfschema:"instance_count"`
	Diagnostics       *DataflowProfileDiagnosticsModel `tfschema:"diagnostics"`
	ExtendedLocation  *ExtendedLocationModel           `tfschema:"extended_location"`
	Tags              map[string]string                `tfschema:"tags"`
	ProvisioningState *string                          `tfschema:"provisioning_state"`
}

type DataflowProfileDiagnosticsModel struct {
	Logs    *DataflowProfileDiagnosticsLogsModel    `tfschema:"logs"`
	Metrics *DataflowProfileDiagnosticsMetricsModel `tfschema:"metrics"`
}

type DataflowProfileDiagnosticsLogsModel struct {
	Level           *string `tfschema:"level"`
	OpenTelemetryExportConfig *DataflowProfileDiagnosticsLogsOpenTelemetryExportConfigModel `tfschema:"open_telemetry_export_config"`
}

type DataflowProfileDiagnosticsLogsOpenTelemetryExportConfigModel struct {
	Level    string  `tfschema:"level"`
	OtlpGrpcEndpoint *string `tfschema:"otlp_grpc_endpoint"`
	IntervalSeconds  *int    `tfschema:"interval_seconds"`
}

type DataflowProfileDiagnosticsMetricsModel struct {
	PrometheusPort            *int    `tfschema:"prometheus_port"`
	OpenTelemetryExportConfig *DataflowProfileDiagnosticsMetricsOpenTelemetryExportConfigModel `tfschema:"open_telemetry_export_config"`
}

type DataflowProfileDiagnosticsMetricsOpenTelemetryExportConfigModel struct {
	Level    string  `tfschema:"level"`
	OtlpGrpcEndpoint *string `tfschema:"otlp_grpc_endpoint"`
	IntervalSeconds  *int    `tfschema:"interval_seconds"`
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
		"instance_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 16),
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
										"Debug",
										"Info",
										"Warn",
										"Error",
										"Trace",
									}, false),
								},
								"open_telemetry_export_config": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"level": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													"Debug",
													"Info",
													"Warn",
													"Error",
													"Trace",
												}, false),
											},
											"otlp_grpc_endpoint": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 253),
											},
											"interval_seconds": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 3600),
											},
										},
									},
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
									ValidateFunc: validation.IntBetween(1024, 65535),
								},
								"open_telemetry_export_config": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"level": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													"Debug",
													"Info",
													"Warn",
													"Error",
													"Trace",
												}, false),
											},
											"otlp_grpc_endpoint": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 253),
											},
											"interval_seconds": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 3600),
											},
										},
									},
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

func (r DataflowProfileResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			// NOTE: O+C Azure automatically assigns provisioning state during resource lifecycle
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
				Properties: expandDataflowProfileProperties(model),
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
				if respModel.ExtendedLocation != nil {
					model.ExtendedLocation = flattenExtendedLocation(respModel.ExtendedLocation)
				}

				if respModel.Tags != nil {
					model.Tags = *respModel.Tags
				}

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

			// Check if anything actually changed before making API call
			if !metadata.ResourceData.HasChange("tags") && 
			   !metadata.ResourceData.HasChange("instance_count") &&
			   !metadata.ResourceData.HasChange("diagnostics") {
				return nil
			}

			payload := dataflowprofile.DataflowProfilePatchModel{}
			hasChanges := false

			// Only include tags if they changed
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
				hasChanges = true
			}

			// Only include properties if they changed
			if metadata.ResourceData.HasChange("instance_count") || metadata.ResourceData.HasChange("diagnostics") {
				payload.Properties = &dataflowprofile.DataflowProfilePropertiesPatch{
					InstanceCount: expandInstanceCount(model.InstanceCount),
					Diagnostics:   expandDataflowProfileDiagnosticsPatch(model.Diagnostics),
				}
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
func expandDataflowProfileProperties(model DataflowProfileModel) *dataflowprofile.DataflowProfileProperties {
	props := &dataflowprofile.DataflowProfileProperties{}

	if model.InstanceCount != nil {
		props.InstanceCount = expandInstanceCount(model.InstanceCount)
	}

	if model.Diagnostics != nil {
		props.Diagnostics = expandDataflowProfileDiagnostics(*model.Diagnostics)
	}

	return props
}

func expandInstanceCount(instanceCount *int) *int64 {
	if instanceCount == nil {
		return nil
	}
	count := int64(*instanceCount)
	return &count
}

func expandDataflowProfileDiagnostics(model DataflowProfileDiagnosticsModel) *dataflowprofile.ProfileDiagnostics {
	diagnostics := &dataflowprofile.ProfileDiagnostics{}

	if model.Logs != nil {
		diagnostics.Logs = expandDataflowProfileDiagnosticsLogs(*model.Logs)
	}

	if model.Metrics != nil {
		diagnostics.Metrics = expandDataflowProfileDiagnosticsMetrics(*model.Metrics)
	}

	return diagnostics
}

func expandDataflowProfileDiagnosticsLogs(model DataflowProfileDiagnosticsLogsModel) *dataflowprofile.DiagnosticsLogs {
	logs := &dataflowprofile.DiagnosticsLogs{}

	if model.Level != nil {
		level := dataflowprofile.Level(*model.Level)
		logs.Level = &level
	}

	if model.OpenTelemetryExportConfig != nil {
		logs.OpenTelemetryExportConfig = expandDataflowProfileDiagnosticsLogsOpenTelemetryExportConfig(*model.OpenTelemetryExportConfig)
	}

	return logs
}

func expandDataflowProfileDiagnosticsLogsOpenTelemetryExportConfig(model DataflowProfileDiagnosticsLogsOpenTelemetryExportConfigModel) *dataflowprofile.OpenTelemetryExportConfig {
	config := &dataflowprofile.OpenTelemetryExportConfig{
		Level: dataflowprofile.Level(model.Level),
	}

	if model.OtlpGrpcEndpoint != nil {
		config.OtlpGrpcEndpoint = model.OtlpGrpcEndpoint
	}

	if model.IntervalSeconds != nil {
		intervalSeconds := int64(*model.IntervalSeconds)
		config.IntervalSeconds = &intervalSeconds
	}

	return config
}

func expandDataflowProfileDiagnosticsMetrics(model DataflowProfileDiagnosticsMetricsModel) *dataflowprofile.Metrics {
	metrics := &dataflowprofile.Metrics{}

	if model.PrometheusPort != nil {
		prometheusPort := int64(*model.PrometheusPort)
		metrics.PrometheusPort = &prometheusPort
	}

	if model.OpenTelemetryExportConfig != nil {
		metrics.OpenTelemetryExportConfig = expandDataflowProfileDiagnosticsMetricsOpenTelemetryExportConfig(*model.OpenTelemetryExportConfig)
	}

	return metrics
}

func expandDataflowProfileDiagnosticsMetricsOpenTelemetryExportConfig(model DataflowProfileDiagnosticsMetricsOpenTelemetryExportConfigModel) *dataflowprofile.OpenTelemetryExportConfig {
	config := &dataflowprofile.OpenTelemetryExportConfig{
		Level: dataflowprofile.Level(model.Level),
	}

	if model.OtlpGrpcEndpoint != nil {
		config.OtlpGrpcEndpoint = model.OtlpGrpcEndpoint
	}

	if model.IntervalSeconds != nil {
		intervalSeconds := int64(*model.IntervalSeconds)
		config.IntervalSeconds = &intervalSeconds
	}

	return config
}

func flattenDataflowProfileProperties(props *dataflowprofile.DataflowProfileProperties, model *DataflowProfileModel) {
	if props == nil {
		return
	}

	if props.InstanceCount != nil {
		instanceCount := int(*props.InstanceCount)
		model.InstanceCount = &instanceCount
	}

	if props.Diagnostics != nil {
		model.Diagnostics = flattenDataflowProfileDiagnostics(*props.Diagnostics)
	}
}

func flattenDataflowProfileDiagnostics(diagnostics dataflowprofile.ProfileDiagnostics) *DataflowProfileDiagnosticsModel {
	result := &DataflowProfileDiagnosticsModel{}

	if diagnostics.Logs != nil {
		result.Logs = flattenDataflowProfileDiagnosticsLogs(*diagnostics.Logs)
	}

	if diagnostics.Metrics != nil {
		result.Metrics = flattenDataflowProfileDiagnosticsMetrics(*diagnostics.Metrics)
	}

	return result
}

func flattenDataflowProfileDiagnosticsLogs(logs dataflowprofile.DiagnosticsLogs) *DataflowProfileDiagnosticsLogsModel {
	result := &DataflowProfileDiagnosticsLogsModel{}

	if logs.Level != nil {
		level := string(*logs.Level)
		result.Level = &level
	}

	if logs.OpenTelemetryExportConfig != nil {
		result.OpenTelemetryExportConfig = flattenDataflowProfileDiagnosticsLogsOpenTelemetryExportConfig(*logs.OpenTelemetryExportConfig)
	}

	return result
}

func flattenDataflowProfileDiagnosticsLogsOpenTelemetryExportConfig(config dataflowprofile.OpenTelemetryExportConfig) *DataflowProfileDiagnosticsLogsOpenTelemetryExportConfigModel {
	result := &DataflowProfileDiagnosticsLogsOpenTelemetryExportConfigModel{
		Level: string(config.Level),
	}

	if config.OtlpGrpcEndpoint != nil {
		result.OtlpGrpcEndpoint = config.OtlpGrpcEndpoint
	}

	if config.IntervalSeconds != nil {
		intervalSeconds := int(*config.IntervalSeconds)
		result.IntervalSeconds = &intervalSeconds
	}

	return result
}

func flattenDataflowProfileDiagnosticsMetrics(metrics dataflowprofile.Metrics) *DataflowProfileDiagnosticsMetricsModel {
	result := &DataflowProfileDiagnosticsMetricsModel{}

	if metrics.PrometheusPort != nil {
		prometheusPort := int(*metrics.PrometheusPort)
		result.PrometheusPort = &prometheusPort
	}

	if metrics.OpenTelemetryExportConfig != nil {
		result.OpenTelemetryExportConfig = flattenDataflowProfileDiagnosticsMetricsOpenTelemetryExportConfig(*metrics.OpenTelemetryExportConfig)
	}

	return result
}

func flattenDataflowProfileDiagnosticsMetricsOpenTelemetryExportConfig(config dataflowprofile.OpenTelemetryExportConfig) *DataflowProfileDiagnosticsMetricsOpenTelemetryExportConfigModel {
	result := &DataflowProfileDiagnosticsMetricsOpenTelemetryExportConfigModel{
		Level: string(config.Level),
	}

	if config.OtlpGrpcEndpoint != nil {
		result.OtlpGrpcEndpoint = config.OtlpGrpcEndpoint
	}

	if config.IntervalSeconds != nil {
		intervalSeconds := int(*config.IntervalSeconds)
		result.IntervalSeconds = &intervalSeconds
	}

	return result
}

// Patch functions for update operations
func expandDataflowProfileDiagnosticsPatch(model *DataflowProfileDiagnosticsModel) *dataflowprofile.ProfileDiagnostics {
	if model == nil {
		return nil
	}
	return expandDataflowProfileDiagnostics(*model)
}