// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var (
	_ sdk.ResourceWithUpdate        = DataCollectionRuleResource{}
	_ sdk.ResourceWithCustomizeDiff = DataCollectionRuleResource{}
)

type DataCollectionRule struct {
	DataCollectionEndpointId string                 `tfschema:"data_collection_endpoint_id"`
	DataFlows                []DataFlow             `tfschema:"data_flow"`
	DataSources              []DataSource           `tfschema:"data_sources"`
	Description              string                 `tfschema:"description"`
	Destinations             []Destination          `tfschema:"destinations"`
	ImmutableId              string                 `tfschema:"immutable_id"`
	Kind                     string                 `tfschema:"kind"`
	Name                     string                 `tfschema:"name"`
	Location                 string                 `tfschema:"location"`
	ResourceGroupName        string                 `tfschema:"resource_group_name"`
	StreamDeclaration        []StreamDeclaration    `tfschema:"stream_declaration"`
	Tags                     map[string]interface{} `tfschema:"tags"`
}

type DataFlow struct {
	BuiltInTransform string   `tfschema:"built_in_transform"`
	Destinations     []string `tfschema:"destinations"`
	OutputStream     string   `tfschema:"output_stream"`
	Streams          []string `tfschema:"streams"`
	TransformKql     string   `tfschema:"transform_kql"`
}

type DataSource struct {
	DataImport          []DataImport          `tfschema:"data_import"`
	Extensions          []Extension           `tfschema:"extension"`
	IisLog              []IisLog              `tfschema:"iis_log"`
	LogFile             []LogFile             `tfschema:"log_file"`
	PerformanceCounters []PerfCounter         `tfschema:"performance_counter"`
	PlatformTelemetry   []PlatformTelemetry   `tfschema:"platform_telemetry"`
	PrometheusForwarder []PrometheusForwarder `tfschema:"prometheus_forwarder"`
	Syslog              []Syslog              `tfschema:"syslog"`
	WindowsEventLogs    []WindowsEventLog     `tfschema:"windows_event_log"`
	WindowsFirewallLog  []WindowsFirewallLog  `tfschema:"windows_firewall_log"`
}

type Destination struct {
	AzureMonitorMetrics []AzureMonitorMetric `tfschema:"azure_monitor_metrics"`
	EventHubDirect      []EventHub           `tfschema:"event_hub_direct"`
	EventHub            []EventHub           `tfschema:"event_hub"`
	LogAnalytics        []LogAnalytic        `tfschema:"log_analytics"`
	MonitorAccount      []MonitorAccount     `tfschema:"monitor_account"`
	StorageBlob         []StorageBlob        `tfschema:"storage_blob"`
	StorageBlobDirect   []StorageBlob        `tfschema:"storage_blob_direct"`
	StorageTableDirect  []StorageTableDirect `tfschema:"storage_table_direct"`
}

type DataImport struct {
	EventHubDataSource []EventHubDataSource `tfschema:"event_hub_data_source"`
}

type EventHubDataSource struct {
	ConsumerGroup string `tfschema:"consumer_group"`
	Name          string `tfschema:"name"`
	Stream        string `tfschema:"stream"`
}

type Extension struct {
	ExtensionName     string   `tfschema:"extension_name"`
	ExtensionSettings string   `tfschema:"extension_json"`
	InputDataSources  []string `tfschema:"input_data_sources"`
	Name              string   `tfschema:"name"`
	Streams           []string `tfschema:"streams"`
}

type IisLog struct {
	Name           string   `tfschema:"name"`
	Streams        []string `tfschema:"streams"`
	LogDirectories []string `tfschema:"log_directories"`
}

type LogFile struct {
	Name         string           `tfschema:"name"`
	Streams      []string         `tfschema:"streams"`
	FilePatterns []string         `tfschema:"file_patterns"`
	Format       string           `tfschema:"format"`
	Settings     []LogFileSetting `tfschema:"settings"`
}

type PerfCounter struct {
	CounterSpecifiers          []string `tfschema:"counter_specifiers"`
	Name                       string   `tfschema:"name"`
	SamplingFrequencyInSeconds int64    `tfschema:"sampling_frequency_in_seconds"`
	Streams                    []string `tfschema:"streams"`
}

type PlatformTelemetry struct {
	Name    string   `tfschema:"name"`
	Streams []string `tfschema:"streams"`
}

type PrometheusForwarder struct {
	Name               string               `tfschema:"name"`
	Streams            []string             `tfschema:"streams"`
	LabelIncludeFilter []LabelIncludeFilter `tfschema:"label_include_filter"`
}

type LabelIncludeFilter struct {
	Label string `tfschema:"label"`
	Value string `tfschema:"value"`
}

type Syslog struct {
	FacilityNames []string `tfschema:"facility_names"`
	LogLevels     []string `tfschema:"log_levels"`
	Name          string   `tfschema:"name"`
	Streams       []string `tfschema:"streams"`
}

type WindowsEventLog struct {
	Name         string   `tfschema:"name"`
	Streams      []string `tfschema:"streams"`
	XPathQueries []string `tfschema:"x_path_queries"`
}

type WindowsFirewallLog struct {
	Name    string   `tfschema:"name"`
	Streams []string `tfschema:"streams"`
}

type AzureMonitorMetric struct {
	Name string `tfschema:"name"`
}

type EventHub struct {
	EventHubResourceId string `tfschema:"event_hub_id"`
	Name               string `tfschema:"name"`
}

type LogAnalytic struct {
	Name                string `tfschema:"name"`
	WorkspaceResourceId string `tfschema:"workspace_resource_id"`
}

type MonitorAccount struct {
	AccountId string `tfschema:"monitor_account_id"`
	Name      string `tfschema:"name"`
}

type StorageBlob struct {
	ContainerName    string `tfschema:"container_name"`
	Name             string `tfschema:"name"`
	StorageAccountId string `tfschema:"storage_account_id"`
}

type StorageTableDirect struct {
	Name             string `tfschema:"name"`
	StorageAccountId string `tfschema:"storage_account_id"`
	TableName        string `tfschema:"table_name"`
}

type LogFileSetting struct {
	Text []TextSetting `tfschema:"text"`
}

type TextSetting struct {
	RecordStartTimestampFormat string `tfschema:"record_start_timestamp_format"`
}

type StreamDeclaration struct {
	StreamName string                    `tfschema:"stream_name"`
	Column     []StreamDeclarationColumn `tfschema:"column"`
}

type StreamDeclarationColumn struct {
	Name string `tfschema:"name"`
	Type string `tfschema:"type"`
}

type DataCollectionRuleResource struct{}

func (r DataCollectionRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"data_collection_endpoint_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: datacollectionendpoints.ValidateDataCollectionEndpointID,
		},

		"data_flow": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destinations": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"streams": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"output_stream": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"transform_kql": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"built_in_transform": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"destinations": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"azure_monitor_metrics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
					},
					"event_hub": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"event_hub_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: eventhubs.ValidateEventhubID,
								},
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
					},
					"event_hub_direct": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"event_hub_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: eventhubs.ValidateEventhubID,
								},
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
					},
					"log_analytics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"workspace_resource_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: workspaces.ValidateWorkspaceID,
								},
							},
						},
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
					},
					"monitor_account": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"monitor_account_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: azure.ValidateResourceID,
								},
							},
						},
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
					},
					"storage_blob": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"container_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"storage_account_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateStorageAccountID,
								},
							},
						},
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
					},
					"storage_blob_direct": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"container_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"storage_account_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateStorageAccountID,
								},
							},
						},
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
					},
					"storage_table_direct": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"table_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"storage_account_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateStorageAccountID,
								},
							},
						},
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
					},
				},
			},
		},

		"data_sources": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_import": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"event_hub_data_source": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"stream": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"consumer_group": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},
							},
						},
					},
					"extension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"extension_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"extension_json": {
									Type:             pluginsdk.TypeString,
									Optional:         true,
									ValidateFunc:     validation.StringIsJSON,
									DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
								},
								"input_data_sources": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"iis_log": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"log_directories": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"log_file": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"file_patterns": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"format": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(datacollectionrules.PossibleValuesForKnownLogFilesDataSourceFormat(), false),
								},
								"settings": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"text": {
												Type:     pluginsdk.TypeList,
												Required: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"record_start_timestamp_format": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringInSlice(datacollectionrules.PossibleValuesForKnownLogFileTextSettingsRecordStartTimestampFormat(), false),
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
					"performance_counter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"sampling_frequency_in_seconds": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 1800),
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"counter_specifiers": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"platform_telemetry": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"prometheus_forwarder": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(datacollectionrules.PossibleValuesForKnownPrometheusForwarderDataSourceStreams(), false),
									},
								},
								"label_include_filter": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"label": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													"microsoft_metrics_include_label",
												}, false),
											},
											"value": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},
							},
						},
					},
					"syslog": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"facility_names": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											datacollectionrules.PossibleValuesForKnownSyslogDataSourceFacilityNames(),
											false),
									},
								},
								"log_levels": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											datacollectionrules.PossibleValuesForKnownSyslogDataSourceLogLevels(), false),
									},
								},
								// lintignore:S013
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"windows_event_log": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"x_path_queries": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"windows_firewall_log": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					"Linux",
					"Windows",
					"AgentDirectToStore",
					"WorkspaceTransforms",
				},
				false),
		},

		"stream_declaration": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"stream_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"column": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(datacollectionrules.PossibleValuesForKnownColumnDefinitionType(), false),
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r DataCollectionRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"immutable_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataCollectionRuleResource) ResourceType() string {
	return "azurerm_monitor_data_collection_rule"
}

func (r DataCollectionRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datacollectionrules.ValidateDataCollectionRuleID
}

func (r DataCollectionRuleResource) ModelObject() interface{} {
	return &DataCollectionRule{}
}

func (r DataCollectionRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state DataCollectionRule
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := datacollectionrules.NewDataCollectionRuleID(subscriptionId, state.ResourceGroupName, state.Name)
			metadata.Logger.Infof("creating %s", id)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			dataSources, err := expandDataCollectionRuleDataSources(state.DataSources)
			if err != nil {
				return err
			}

			identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			input := datacollectionrules.DataCollectionRuleResource{
				Identity: identityValue,
				Kind:     expandDataCollectionRuleKind(state.Kind),
				Location: azure.NormalizeLocation(state.Location),
				Name:     utils.String(state.Name),
				Properties: &datacollectionrules.DataCollectionRule{
					DataFlows:          expandDataCollectionRuleDataFlows(state.DataFlows),
					DataSources:        dataSources,
					Description:        utils.String(state.Description),
					Destinations:       expandDataCollectionRuleDestinations(state.Destinations),
					StreamDeclarations: expandDataCollectionRuleStreamDeclarations(state.StreamDeclaration),
				},
				Tags: tags.Expand(state.Tags),
			}

			if state.DataCollectionEndpointId != "" {
				input.Properties.DataCollectionEndpointId = utils.String(state.DataCollectionEndpointId)
			}

			if _, err := client.Create(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r DataCollectionRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRulesClient
			id, err := datacollectionrules.ParseDataCollectionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var dataCollectionEndpointId, description, immutableId, kind, location string
			var tag map[string]interface{}
			var dataFlows []DataFlow
			var dataSources []DataSource
			var destinations []Destination
			var streamDeclaration []StreamDeclaration

			if model := resp.Model; model != nil {
				kind = flattenDataCollectionRuleKind(model.Kind)
				location = azure.NormalizeLocation(model.Location)
				tag = tags.Flatten(model.Tags)

				identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				if prop := model.Properties; prop != nil {
					dataCollectionEndpointId = flattenStringPtr(prop.DataCollectionEndpointId)
					description = flattenStringPtr(prop.Description)
					dataFlows = flattenDataCollectionRuleDataFlows(prop.DataFlows)
					dataSources = flattenDataCollectionRuleDataSources(prop.DataSources)
					destinations = flattenDataCollectionRuleDestinations(prop.Destinations)
					immutableId = flattenStringPtr(prop.ImmutableId)
					streamDeclaration = flattenDataCollectionRuleStreamDeclarations(prop.StreamDeclarations)
				}
			}

			return metadata.Encode(&DataCollectionRule{
				Name:                     id.DataCollectionRuleName,
				ResourceGroupName:        id.ResourceGroupName,
				DataCollectionEndpointId: dataCollectionEndpointId,
				DataFlows:                dataFlows,
				DataSources:              dataSources,
				Description:              description,
				Destinations:             destinations,
				ImmutableId:              immutableId,
				Kind:                     kind,
				Location:                 location,
				StreamDeclaration:        streamDeclaration,
				Tags:                     tag,
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r DataCollectionRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := datacollectionrules.ParseDataCollectionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Monitor.DataCollectionRulesClient
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}
			existing := resp.Model
			if existing.Properties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}

			var state DataCollectionRule
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("kind") {
				existing.Kind = expandDataCollectionRuleKind(state.Kind)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.Expand(state.Tags)
			}

			if metadata.ResourceData.HasChange("data_flow") {
				existing.Properties.DataFlows = expandDataCollectionRuleDataFlows(state.DataFlows)
			}

			if metadata.ResourceData.HasChange("data_sources") {
				dataSource, err := expandDataCollectionRuleDataSources(state.DataSources)
				if err != nil {
					return err
				}
				existing.Properties.DataSources = dataSource
			}

			if metadata.ResourceData.HasChange("data_collection_endpoint_id") {
				if state.DataCollectionEndpointId != "" {
					existing.Properties.DataCollectionEndpointId = utils.String(state.DataCollectionEndpointId)
				} else {
					existing.Properties.DataCollectionEndpointId = nil
				}
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Properties.Description = utils.String(state.Description)
			}

			if metadata.ResourceData.HasChange("destinations") {
				existing.Properties.Destinations = expandDataCollectionRuleDestinations(state.Destinations)
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("stream_declaration") {
				existing.Properties.StreamDeclarations = expandDataCollectionRuleStreamDeclarations(state.StreamDeclaration)
			}

			if _, err := client.Create(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r DataCollectionRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRulesClient
			id, err := datacollectionrules.ParseDataCollectionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			resp, err := client.Delete(ctx, *id)
			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandDataCollectionRuleKind(input string) *datacollectionrules.KnownDataCollectionRuleResourceKind {
	if input == "" {
		return nil
	}

	result := datacollectionrules.KnownDataCollectionRuleResourceKind(input)
	return &result
}

func expandDataCollectionRuleDataFlows(input []DataFlow) *[]datacollectionrules.DataFlow {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.DataFlow, 0)
	for _, v := range input {
		dataFlow := datacollectionrules.DataFlow{
			Destinations: pointer.To(v.Destinations),
			Streams:      expandDataCollectionRuleDataFlowStreams(v.Streams),
		}

		if v.BuiltInTransform != "" {
			dataFlow.BuiltInTransform = utils.String(v.BuiltInTransform)
		}

		if v.OutputStream != "" {
			dataFlow.OutputStream = utils.String(v.OutputStream)
		}

		if v.TransformKql != "" {
			dataFlow.TransformKql = utils.String(v.TransformKql)
		}

		result = append(result, dataFlow)
	}
	return &result
}

func expandDataCollectionRuleDataFlowStreams(input []string) *[]datacollectionrules.KnownDataFlowStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownDataFlowStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownDataFlowStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSources(input []DataSource) (*datacollectionrules.DataSourcesSpec, error) {
	if len(input) == 0 {
		return nil, nil
	}

	extension, err := expandDataCollectionRuleDataSourceExtensions(input[0].Extensions)
	if err != nil {
		return nil, err
	}

	return &datacollectionrules.DataSourcesSpec{
		DataImports:         expandDataCollectionRuleDataSourceDataImports(input[0].DataImport),
		Extensions:          extension,
		IisLogs:             expandDataCollectionRuleDataSourceIisLogs(input[0].IisLog),
		LogFiles:            expandDataCollectionRuleDataSourceLogFiles(input[0].LogFile),
		PerformanceCounters: expandDataCollectionRuleDataSourcePerfCounters(input[0].PerformanceCounters),
		PlatformTelemetry:   expandDataCollectionRuleDataSourcePlatformTelemetry(input[0].PlatformTelemetry),
		PrometheusForwarder: expandDataCollectionRuleDataSourcePrometheusForwarder(input[0].PrometheusForwarder),
		Syslog:              expandDataCollectionRuleDataSourceSyslog(input[0].Syslog),
		WindowsEventLogs:    expandDataCollectionRuleDataSourceWindowsEventLogs(input[0].WindowsEventLogs),
		WindowsFirewallLogs: expandDataCollectionRuleDataSourceWindowsFirewallLogs(input[0].WindowsFirewallLog),
	}, nil
}

func expandDataCollectionRuleDataSourceDataImports(input []DataImport) *datacollectionrules.DataImportSources {
	if len(input) == 0 || len(input[0].EventHubDataSource) == 0 {
		return nil
	}

	result := &datacollectionrules.DataImportSources{
		EventHub: &datacollectionrules.EventHubDataSource{
			Name:   utils.String(input[0].EventHubDataSource[0].Name),
			Stream: utils.String(input[0].EventHubDataSource[0].Stream),
		},
	}

	if consumerGroup := input[0].EventHubDataSource[0].ConsumerGroup; consumerGroup != "" {
		result.EventHub.ConsumerGroup = utils.String(consumerGroup)
	}

	return result
}

func expandDataCollectionRuleDataSourceExtensions(input []Extension) (*[]datacollectionrules.ExtensionDataSource, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]datacollectionrules.ExtensionDataSource, 0)
	for _, v := range input {
		var extensionSettings interface{}
		if v.ExtensionSettings != "" {
			settings, err := pluginsdk.ExpandJsonFromString(v.ExtensionSettings)
			if err != nil {
				return nil, err
			}
			extensionSettings = settings
		}

		result = append(result, datacollectionrules.ExtensionDataSource{
			ExtensionName:     v.ExtensionName,
			ExtensionSettings: &extensionSettings,
			InputDataSources:  pointer.To(v.InputDataSources),
			Name:              utils.String(v.Name),
			Streams:           expandDataCollectionRuleDataSourceExtensionStreams(v.Streams),
		})
	}
	return &result, nil
}

func expandDataCollectionRuleDataSourceExtensionStreams(input []string) *[]datacollectionrules.KnownExtensionDataSourceStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownExtensionDataSourceStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownExtensionDataSourceStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceIisLogs(input []IisLog) *[]datacollectionrules.IisLogsDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.IisLogsDataSource, 0)
	for _, v := range input {
		iisLog := datacollectionrules.IisLogsDataSource{
			Name:    utils.String(v.Name),
			Streams: v.Streams,
		}

		if len(v.LogDirectories) != 0 {
			iisLog.LogDirectories = pointer.To(v.LogDirectories)
		}

		result = append(result, iisLog)
	}

	return &result
}

func expandDataCollectionRuleDataSourceLogFiles(input []LogFile) *[]datacollectionrules.LogFilesDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.LogFilesDataSource, 0)
	for _, v := range input {
		logFile := datacollectionrules.LogFilesDataSource{
			Name:         utils.String(v.Name),
			Streams:      v.Streams,
			FilePatterns: v.FilePatterns,
			Format:       datacollectionrules.KnownLogFilesDataSourceFormat(v.Format),
		}

		if len(v.Settings) != 0 && len(v.Settings[0].Text) != 0 {
			logFile.Settings = &datacollectionrules.LogFileSettings{
				Text: &datacollectionrules.LogFileTextSettings{
					RecordStartTimestampFormat: datacollectionrules.KnownLogFileTextSettingsRecordStartTimestampFormat(v.Settings[0].Text[0].RecordStartTimestampFormat),
				},
			}
		}

		result = append(result, logFile)
	}

	return &result
}

func expandDataCollectionRuleDataSourcePerfCounters(input []PerfCounter) *[]datacollectionrules.PerfCounterDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.PerfCounterDataSource, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.PerfCounterDataSource{
			CounterSpecifiers:          pointer.To(v.CounterSpecifiers),
			Name:                       utils.String(v.Name),
			SamplingFrequencyInSeconds: utils.Int64(v.SamplingFrequencyInSeconds),
			Streams:                    expandDataCollectionRuleDataSourcePerfCounterStreams(v.Streams),
		})
	}

	return &result
}

func expandDataCollectionRuleDataSourcePerfCounterStreams(input []string) *[]datacollectionrules.KnownPerfCounterDataSourceStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownPerfCounterDataSourceStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownPerfCounterDataSourceStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourcePlatformTelemetry(input []PlatformTelemetry) *[]datacollectionrules.PlatformTelemetryDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.PlatformTelemetryDataSource, 0)
	for _, v := range input {
		platformTelemetry := datacollectionrules.PlatformTelemetryDataSource{
			Name:    utils.String(v.Name),
			Streams: v.Streams,
		}

		result = append(result, platformTelemetry)
	}

	return &result
}

func expandDataCollectionRuleDataSourcePrometheusForwarder(input []PrometheusForwarder) *[]datacollectionrules.PrometheusForwarderDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.PrometheusForwarderDataSource, 0)
	for _, v := range input {
		streams := make([]datacollectionrules.KnownPrometheusForwarderDataSourceStreams, 0)
		for _, stream := range v.Streams {
			streams = append(streams, datacollectionrules.KnownPrometheusForwarderDataSourceStreams(stream))
		}

		prometheusForwarder := datacollectionrules.PrometheusForwarderDataSource{
			Name:    utils.String(v.Name),
			Streams: &streams,
		}

		if len(v.LabelIncludeFilter) != 0 {
			labelIncludeFilter := make(map[string]string, 0)
			for _, filter := range v.LabelIncludeFilter {
				labelIncludeFilter[filter.Label] = filter.Value
			}

			prometheusForwarder.LabelIncludeFilter = &labelIncludeFilter
		}

		result = append(result, prometheusForwarder)
	}

	return &result
}

func expandDataCollectionRuleDataSourceSyslog(input []Syslog) *[]datacollectionrules.SyslogDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.SyslogDataSource, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.SyslogDataSource{
			FacilityNames: expandDataCollectionRuleDataSourceSyslogFacilityNames(v.FacilityNames),
			LogLevels:     expandDataCollectionRuleDataSourceSyslogLogLevels(v.LogLevels),
			Name:          utils.String(v.Name),
			Streams:       expandDataCollectionRuleDataSourceSyslogStreams(v.Streams),
		})
	}
	return &result
}

func expandDataCollectionRuleDataSourceSyslogStreams(input []string) *[]datacollectionrules.KnownSyslogDataSourceStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownSyslogDataSourceStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownSyslogDataSourceStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceSyslogFacilityNames(input []string) *[]datacollectionrules.KnownSyslogDataSourceFacilityNames {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownSyslogDataSourceFacilityNames, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownSyslogDataSourceFacilityNames(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceSyslogLogLevels(input []string) *[]datacollectionrules.KnownSyslogDataSourceLogLevels {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownSyslogDataSourceLogLevels, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownSyslogDataSourceLogLevels(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceWindowsEventLogs(input []WindowsEventLog) *[]datacollectionrules.WindowsEventLogDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.WindowsEventLogDataSource, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.WindowsEventLogDataSource{
			Name:         utils.String(v.Name),
			Streams:      expandDataCollectionRuleDataSourceWindowsEventLogsStreams(v.Streams),
			XPathQueries: pointer.To(v.XPathQueries),
		})
	}
	return &result
}

func expandDataCollectionRuleDataSourceWindowsEventLogsStreams(input []string) *[]datacollectionrules.KnownWindowsEventLogDataSourceStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownWindowsEventLogDataSourceStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownWindowsEventLogDataSourceStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceWindowsFirewallLogs(input []WindowsFirewallLog) *[]datacollectionrules.WindowsFirewallLogsDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.WindowsFirewallLogsDataSource, 0)
	for _, v := range input {
		windowsFirewallLog := datacollectionrules.WindowsFirewallLogsDataSource{
			Name:    utils.String(v.Name),
			Streams: v.Streams,
		}

		result = append(result, windowsFirewallLog)
	}

	return &result
}

func expandDataCollectionRuleDestinations(input []Destination) *datacollectionrules.DestinationsSpec {
	if len(input) == 0 {
		return nil
	}

	return &datacollectionrules.DestinationsSpec{
		AzureMonitorMetrics: expandDataCollectionRuleDestinationMetrics(input[0].AzureMonitorMetrics),
		EventHubsDirect:     expandDataCollectionRuleDestinationEventHubsDirect(input[0].EventHubDirect),
		EventHubs:           expandDataCollectionRuleDestinationEventHubs(input[0].EventHub),
		LogAnalytics:        expandDataCollectionRuleDestinationLogAnalytics(input[0].LogAnalytics),
		MonitoringAccounts:  expandDataCollectionRuleDestinationMonitoringAccounts(input[0].MonitorAccount),
		StorageAccounts:     expandDataCollectionRuleDestinationStorageBlobs(input[0].StorageBlob),
		StorageBlobsDirect:  expandDataCollectionRuleDestinationStorageBlobs(input[0].StorageBlobDirect),
		StorageTablesDirect: expandDataCollectionRuleDestinationStorageTableDirect(input[0].StorageTableDirect),
	}
}

func expandDataCollectionRuleDestinationMetrics(input []AzureMonitorMetric) *datacollectionrules.AzureMonitorMetricsDestination {
	if len(input) == 0 {
		return nil
	}

	return &datacollectionrules.AzureMonitorMetricsDestination{
		Name: utils.String(input[0].Name),
	}
}

func expandDataCollectionRuleDestinationEventHubs(input []EventHub) *[]datacollectionrules.EventHubDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.EventHubDestination, 0)
	for _, v := range input {
		eventhub := datacollectionrules.EventHubDestination{
			Name:               utils.String(v.Name),
			EventHubResourceId: utils.String(v.EventHubResourceId),
		}

		result = append(result, eventhub)
	}

	return &result
}

func expandDataCollectionRuleDestinationEventHubsDirect(input []EventHub) *[]datacollectionrules.EventHubDirectDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.EventHubDirectDestination, 0)
	for _, v := range input {
		eventhub := datacollectionrules.EventHubDirectDestination{
			Name:               utils.String(v.Name),
			EventHubResourceId: utils.String(v.EventHubResourceId),
		}

		result = append(result, eventhub)
	}

	return &result
}

func expandDataCollectionRuleDestinationLogAnalytics(input []LogAnalytic) *[]datacollectionrules.LogAnalyticsDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.LogAnalyticsDestination, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.LogAnalyticsDestination{
			Name:                utils.String(v.Name),
			WorkspaceResourceId: utils.String(v.WorkspaceResourceId),
		})
	}
	return &result
}

func expandDataCollectionRuleDestinationMonitoringAccounts(input []MonitorAccount) *[]datacollectionrules.MonitoringAccountDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.MonitoringAccountDestination, 0)
	for _, v := range input {
		monitorAccount := datacollectionrules.MonitoringAccountDestination{
			Name:              utils.String(v.Name),
			AccountResourceId: utils.String(v.AccountId),
		}

		result = append(result, monitorAccount)
	}

	return &result
}

func expandDataCollectionRuleDestinationStorageBlobs(input []StorageBlob) *[]datacollectionrules.StorageBlobDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.StorageBlobDestination, 0)
	for _, v := range input {
		monitorAccount := datacollectionrules.StorageBlobDestination{
			Name:                     utils.String(v.Name),
			StorageAccountResourceId: utils.String(v.StorageAccountId),
			ContainerName:            utils.String(v.ContainerName),
		}

		result = append(result, monitorAccount)
	}

	return &result
}

func expandDataCollectionRuleDestinationStorageTableDirect(input []StorageTableDirect) *[]datacollectionrules.StorageTableDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.StorageTableDestination, 0)
	for _, v := range input {
		monitorAccount := datacollectionrules.StorageTableDestination{
			Name:                     utils.String(v.Name),
			StorageAccountResourceId: utils.String(v.StorageAccountId),
			TableName:                utils.String(v.TableName),
		}

		result = append(result, monitorAccount)
	}

	return &result
}

func expandDataCollectionRuleStreamDeclarations(input []StreamDeclaration) *map[string]datacollectionrules.StreamDeclaration {
	if len(input) == 0 {
		return nil
	}

	result := make(map[string]datacollectionrules.StreamDeclaration, 0)
	for _, v := range input {
		columns := make([]datacollectionrules.ColumnDefinition, 0)
		for _, column := range v.Column {
			columnType := datacollectionrules.KnownColumnDefinitionType(column.Type)
			columns = append(columns, datacollectionrules.ColumnDefinition{
				Name: utils.String(column.Name),
				Type: &columnType,
			})
		}

		result[v.StreamName] = datacollectionrules.StreamDeclaration{
			Columns: &columns,
		}
	}

	return &result
}

func flattenDataCollectionRuleKind(input *datacollectionrules.KnownDataCollectionRuleResourceKind) string {
	if input == nil {
		return ""
	}
	return string(*input)
}

func flattenStringPtr(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}

func flattenStringSlicePtr(input *[]string) []string {
	if input == nil {
		return make([]string, 0)
	}
	return *input
}

func flattenDataCollectionRuleDataFlows(input *[]datacollectionrules.DataFlow) []DataFlow {
	if input == nil {
		return make([]DataFlow, 0)
	}

	result := make([]DataFlow, 0)
	for _, v := range *input {
		result = append(result, DataFlow{
			BuiltInTransform: flattenStringPtr(v.BuiltInTransform),
			Destinations:     flattenStringSlicePtr(v.Destinations),
			OutputStream:     flattenStringPtr(v.OutputStream),
			Streams:          flattenDataCollectionRuleDataFlowStreams(v.Streams),
			TransformKql:     flattenStringPtr(v.TransformKql),
		})
	}
	return result
}

func flattenDataCollectionRuleDataFlowStreams(input *[]datacollectionrules.KnownDataFlowStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleDataSources(input *datacollectionrules.DataSourcesSpec) []DataSource {
	if input == nil {
		return make([]DataSource, 0)
	}

	return []DataSource{{
		DataImport:          flattenDataCollectionRuleDataSourceDataImports(input.DataImports),
		Extensions:          flattenDataCollectionRuleDataSourceExtensions(input.Extensions),
		IisLog:              flattenDataCollectionRuleDataSourceIisLog(input.IisLogs),
		LogFile:             flattenDataCollectionRuleDataSourceLogFiles(input.LogFiles),
		PerformanceCounters: flattenDataCollectionRuleDataSourcePerfCounters(input.PerformanceCounters),
		PlatformTelemetry:   flattenDataCollectionRuleDataSourcePlatformTelemetry(input.PlatformTelemetry),
		PrometheusForwarder: flattenDataCollectionRuleDataSourcePrometheusForwarder(input.PrometheusForwarder),
		Syslog:              flattenDataCollectionRuleDataSourceSyslog(input.Syslog),
		WindowsEventLogs:    flattenDataCollectionRuleWindowsEventLogs(input.WindowsEventLogs),
		WindowsFirewallLog:  flattenDataCollectionRuleWindowsFirewallLog(input.WindowsFirewallLogs),
	}}
}

func flattenDataCollectionRuleDataSourceDataImports(input *datacollectionrules.DataImportSources) []DataImport {
	if input == nil || input.EventHub == nil {
		return make([]DataImport, 0)
	}

	return []DataImport{
		{
			EventHubDataSource: []EventHubDataSource{
				{
					ConsumerGroup: flattenStringPtr(input.EventHub.ConsumerGroup),
					Name:          flattenStringPtr(input.EventHub.Name),
					Stream:        flattenStringPtr(input.EventHub.Stream),
				},
			},
		},
	}
}

func flattenDataCollectionRuleDataSourceExtensions(input *[]datacollectionrules.ExtensionDataSource) []Extension {
	if input == nil {
		return make([]Extension, 0)
	}

	result := make([]Extension, 0)
	for _, v := range *input {
		extensionSettings := ""
		if v.ExtensionSettings != nil {
			settingString, _ := pluginsdk.FlattenJsonToString((*v.ExtensionSettings).(map[string]interface{}))
			extensionSettings = settingString
		}
		result = append(result, Extension{
			ExtensionName:     v.ExtensionName,
			Name:              flattenStringPtr(v.Name),
			ExtensionSettings: extensionSettings,
			InputDataSources:  flattenStringSlicePtr(v.InputDataSources),
			Streams:           flattenDataCollectionRuleDataSourceExtensionStreams(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourceExtensionStreams(input *[]datacollectionrules.KnownExtensionDataSourceStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleDataSourceIisLog(input *[]datacollectionrules.IisLogsDataSource) []IisLog {
	if input == nil {
		return make([]IisLog, 0)
	}

	result := make([]IisLog, 0)
	for _, v := range *input {
		result = append(result, IisLog{
			Name:           flattenStringPtr(v.Name),
			LogDirectories: flattenStringSlicePtr(v.LogDirectories),
			Streams:        v.Streams,
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourceLogFiles(input *[]datacollectionrules.LogFilesDataSource) []LogFile {
	if input == nil {
		return make([]LogFile, 0)
	}

	result := make([]LogFile, 0)
	for _, v := range *input {
		setting := make([]LogFileSetting, 0)
		if v.Settings != nil && v.Settings.Text != nil {
			setting = append(setting, LogFileSetting{
				Text: []TextSetting{
					{
						RecordStartTimestampFormat: string(v.Settings.Text.RecordStartTimestampFormat),
					},
				},
			})
		}

		result = append(result, LogFile{
			Name:         flattenStringPtr(v.Name),
			Format:       string(v.Format),
			FilePatterns: v.FilePatterns,
			Streams:      v.Streams,
			Settings:     setting,
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourcePerfCounters(input *[]datacollectionrules.PerfCounterDataSource) []PerfCounter {
	if input == nil {
		return make([]PerfCounter, 0)
	}

	result := make([]PerfCounter, 0)
	for _, v := range *input {
		result = append(result, PerfCounter{
			Name:                       flattenStringPtr(v.Name),
			CounterSpecifiers:          flattenStringSlicePtr(v.CounterSpecifiers),
			SamplingFrequencyInSeconds: utils.NormaliseNilableInt64(v.SamplingFrequencyInSeconds),
			Streams:                    flattenDataCollectionRuleDataSourcePerfCounterStreams(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourcePerfCounterStreams(input *[]datacollectionrules.KnownPerfCounterDataSourceStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleDataSourcePlatformTelemetry(input *[]datacollectionrules.PlatformTelemetryDataSource) []PlatformTelemetry {
	if input == nil {
		return make([]PlatformTelemetry, 0)
	}

	result := make([]PlatformTelemetry, 0)
	for _, v := range *input {
		result = append(result, PlatformTelemetry{
			Name:    flattenStringPtr(v.Name),
			Streams: v.Streams,
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourcePrometheusForwarder(input *[]datacollectionrules.PrometheusForwarderDataSource) []PrometheusForwarder {
	if input == nil {
		return make([]PrometheusForwarder, 0)
	}

	result := make([]PrometheusForwarder, 0)
	for _, v := range *input {
		labelIncludeFilter := make([]LabelIncludeFilter, 0)
		if v.LabelIncludeFilter != nil {
			for label, value := range *v.LabelIncludeFilter {
				labelIncludeFilter = append(labelIncludeFilter, LabelIncludeFilter{
					Label: label,
					Value: value,
				})
			}
		}

		result = append(result, PrometheusForwarder{
			Name:               flattenStringPtr(v.Name),
			Streams:            flattenDataCollectionRuleDataSourcePrometheusForwarderStreams(v.Streams),
			LabelIncludeFilter: labelIncludeFilter,
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourcePrometheusForwarderStreams(input *[]datacollectionrules.KnownPrometheusForwarderDataSourceStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleDataSourceSyslog(input *[]datacollectionrules.SyslogDataSource) []Syslog {
	if input == nil {
		return make([]Syslog, 0)
	}

	result := make([]Syslog, 0)
	for _, v := range *input {
		result = append(result, Syslog{
			Name:          flattenStringPtr(v.Name),
			FacilityNames: flattenDataCollectionRuleDataSourceSyslogFacilityNames(v.FacilityNames),
			LogLevels:     flattenDataCollectionRuleDataSourceSyslogLogLevels(v.LogLevels),
			Streams:       flattenDataCollectionRuleSyslogStreams(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourceSyslogFacilityNames(input *[]datacollectionrules.KnownSyslogDataSourceFacilityNames) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleDataSourceSyslogLogLevels(input *[]datacollectionrules.KnownSyslogDataSourceLogLevels) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleSyslogStreams(input *[]datacollectionrules.KnownSyslogDataSourceStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleWindowsEventLogs(input *[]datacollectionrules.WindowsEventLogDataSource) []WindowsEventLog {
	if input == nil {
		return make([]WindowsEventLog, 0)
	}

	result := make([]WindowsEventLog, 0)
	for _, v := range *input {
		result = append(result, WindowsEventLog{
			Name:         flattenStringPtr(v.Name),
			XPathQueries: flattenStringSlicePtr(v.XPathQueries),
			Streams:      flattenDataCollectionRuleWindowsEventLogStreams(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleWindowsEventLogStreams(input *[]datacollectionrules.KnownWindowsEventLogDataSourceStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleWindowsFirewallLog(input *[]datacollectionrules.WindowsFirewallLogsDataSource) []WindowsFirewallLog {
	if input == nil {
		return make([]WindowsFirewallLog, 0)
	}

	result := make([]WindowsFirewallLog, 0)
	for _, v := range *input {
		result = append(result, WindowsFirewallLog{
			Name:    flattenStringPtr(v.Name),
			Streams: v.Streams,
		})
	}
	return result
}

func flattenDataCollectionRuleDestinations(input *datacollectionrules.DestinationsSpec) []Destination {
	if input == nil {
		return make([]Destination, 0)
	}

	return []Destination{{
		AzureMonitorMetrics: flattenDataCollectionRuleDestinationMetrics(input.AzureMonitorMetrics),
		EventHub:            flattenDataCollectionRuleDestinationEventHubs(input.EventHubs),
		EventHubDirect:      flattenDataCollectionRuleDestinationEventHubDirect(input.EventHubsDirect),
		LogAnalytics:        flattenDataCollectionRuleDestinationLogAnalytics(input.LogAnalytics),
		MonitorAccount:      flattenDataCollectionRuleDestinationMonitorAccount(input.MonitoringAccounts),
		StorageBlob:         flattenDataCollectionRuleDestinationStorageBlob(input.StorageAccounts),
		StorageBlobDirect:   flattenDataCollectionRuleDestinationStorageBlob(input.StorageBlobsDirect),
		StorageTableDirect:  flattenDataCollectionRuleDestinationStorageTableDirect(input.StorageTablesDirect),
	}}
}

func flattenDataCollectionRuleDestinationMetrics(input *datacollectionrules.AzureMonitorMetricsDestination) []AzureMonitorMetric {
	if input == nil {
		return make([]AzureMonitorMetric, 0)
	}

	return []AzureMonitorMetric{{
		Name: flattenStringPtr(input.Name),
	}}
}

func flattenDataCollectionRuleDestinationEventHubs(input *[]datacollectionrules.EventHubDestination) []EventHub {
	if input == nil {
		return make([]EventHub, 0)
	}

	result := make([]EventHub, 0)
	for _, v := range *input {
		result = append(result, EventHub{
			Name:               flattenStringPtr(v.Name),
			EventHubResourceId: flattenStringPtr(v.EventHubResourceId),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationEventHubDirect(input *[]datacollectionrules.EventHubDirectDestination) []EventHub {
	if input == nil {
		return make([]EventHub, 0)
	}

	result := make([]EventHub, 0)
	for _, v := range *input {
		result = append(result, EventHub{
			Name:               flattenStringPtr(v.Name),
			EventHubResourceId: flattenStringPtr(v.EventHubResourceId),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationLogAnalytics(input *[]datacollectionrules.LogAnalyticsDestination) []LogAnalytic {
	if input == nil {
		return make([]LogAnalytic, 0)
	}

	result := make([]LogAnalytic, 0)
	for _, v := range *input {
		result = append(result, LogAnalytic{
			Name:                flattenStringPtr(v.Name),
			WorkspaceResourceId: flattenStringPtr(v.WorkspaceResourceId),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationMonitorAccount(input *[]datacollectionrules.MonitoringAccountDestination) []MonitorAccount {
	if input == nil {
		return make([]MonitorAccount, 0)
	}

	result := make([]MonitorAccount, 0)
	for _, v := range *input {
		result = append(result, MonitorAccount{
			Name:      flattenStringPtr(v.Name),
			AccountId: flattenStringPtr(v.AccountResourceId),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationStorageBlob(input *[]datacollectionrules.StorageBlobDestination) []StorageBlob {
	if input == nil {
		return make([]StorageBlob, 0)
	}

	result := make([]StorageBlob, 0)
	for _, v := range *input {
		result = append(result, StorageBlob{
			Name:             flattenStringPtr(v.Name),
			StorageAccountId: flattenStringPtr(v.StorageAccountResourceId),
			ContainerName:    flattenStringPtr(v.ContainerName),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationStorageTableDirect(input *[]datacollectionrules.StorageTableDestination) []StorageTableDirect {
	if input == nil {
		return make([]StorageTableDirect, 0)
	}

	result := make([]StorageTableDirect, 0)
	for _, v := range *input {
		result = append(result, StorageTableDirect{
			Name:             flattenStringPtr(v.Name),
			StorageAccountId: flattenStringPtr(v.StorageAccountResourceId),
			TableName:        flattenStringPtr(v.TableName),
		})
	}
	return result
}

func flattenDataCollectionRuleStreamDeclarations(input *map[string]datacollectionrules.StreamDeclaration) []StreamDeclaration {
	if input == nil {
		return make([]StreamDeclaration, 0)
	}

	result := make([]StreamDeclaration, 0)
	for name, stream := range *input {
		if stream.Columns == nil {
			continue
		}

		columns := make([]StreamDeclarationColumn, 0)
		for _, column := range *stream.Columns {
			if column.Name == nil || column.Type == nil {
				continue
			}
			columns = append(columns, StreamDeclarationColumn{
				Name: *column.Name,
				Type: string(*column.Type),
			})
		}

		streamDeclaration := StreamDeclaration{
			StreamName: name,
			Column:     columns,
		}
		result = append(result, streamDeclaration)
	}

	return result
}

func (r DataCollectionRuleResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldValue, newValue := metadata.ResourceDiff.GetChange("kind"); oldValue.(string) != newValue.(string) && oldValue.(string) != "" {
				if err := metadata.ResourceDiff.ForceNew("kind"); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
