// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Data Collection Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"data_collection_endpoint_id": {
			Description:  "The resource ID of the Data Collection Endpoint that this rule can be used with.",
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
						Description: "A list of destination names. An `azure_monitor_metrics` data source only allows for stream of kind `Microsoft-InsightsMetrics`.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MinItems:    1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"streams": {
						Description: "A list of streams.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MinItems:    1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"output_stream": {
						Description:  "The output stream of the transform. Only required if the data flow changes data to a different stream.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"transform_kql": {
						Description:  "The KQL query to transform stream data.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"built_in_transform": {
						Description:  "The built-in transform to transform stream data.",
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
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"azure_monitor_metrics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description:  "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"event_hub": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"event_hub_id": {
									Description:  "The resource ID of the Event Hub.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: eventhubs.ValidateEventhubID,
								},
								"name": {
									Description:  "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"event_hub_direct": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"event_hub_id": {
									Description:  "The resource ID of the Event Hub.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: eventhubs.ValidateEventhubID,
								},
								"name": {
									Description:  "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"log_analytics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description:  "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"workspace_resource_id": {
									Description:  "The ID of a Log Analytic Workspace resource.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: workspaces.ValidateWorkspaceID,
								},
							},
						},
					},
					"monitor_account": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description:  "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"monitor_account_id": {
									Description:  "The resource ID of the Monitor Account.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: azure.ValidateResourceID,
								},
							},
						},
					},
					"storage_blob": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description:  "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"container_name": {
									Description:  "The Storage Container name.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"storage_account_id": {
									Description:  "The resource ID of the Storage Account.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateStorageAccountID,
								},
							},
						},
					},
					"storage_blob_direct": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description:  "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"container_name": {
									Description:  "The Storage Container name.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"storage_account_id": {
									Description:  "The resource ID of the Storage Account.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateStorageAccountID,
								},
							},
						},
					},
					"storage_table_direct": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						AtLeastOneOf: []string{
							"destinations.0.azure_monitor_metrics", "destinations.0.event_hub",
							"destinations.0.event_hub_direct", "destinations.0.log_analytics",
							"destinations.0.monitor_account", "destinations.0.storage_blob",
							"destinations.0.storage_blob_direct", "destinations.0.storage_table_direct",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description:  "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"table_name": {
									Description:  "The Storage Table name.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"storage_account_id": {
									Description:  "The resource ID of the Storage Account.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateStorageAccountID,
								},
							},
						},
					},
				},
			},
		},

		"data_sources": {
			Description: "This property is optional and can be omitted if the rule is meant to be used via direct calls to the provisioned endpoint.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
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
												Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringLenBetween(1, 32),
											},
											"stream": {
												Description:  "The stream to collect from Event Hub.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"consumer_group": {
												Description:  "The Event Hub consumer group name.",
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
									Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"extension_name": {
									Description:  "The name of the VM extension.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"extension_json": {
									Description:      "A JSON String which specifies the extension setting.",
									Type:             pluginsdk.TypeString,
									Optional:         true,
									ValidateFunc:     validation.StringIsJSON,
									DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
								},
								"input_data_sources": {
									Description: "A list of data sources this extension needs data from. An item should be a name of a supported data source which produces only one stream. Supported data sources type: `performance_counter`, `windows_event_log`, and `syslog`.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
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
									Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"log_directories": {
									Description: "A list of absolute paths where the log files are located.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									MinItems:    1,
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
									Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"file_patterns": {
									Description: "A list of file patterns where the log files are located. For example, `C:\\\\JavaLogs\\\\*.log`.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"format": {
									Description:  "The data format of the log files.",
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
															Description:  "The timestamp format of the text log files.",
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
									Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"sampling_frequency_in_seconds": {
									Description:  "The number of seconds between consecutive counter measurements (samples). `sampling_frequency_in_seconds` must be equal to `60` seconds for counters collected with `Microsoft-InsightsMetrics` stream.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 1800),
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"counter_specifiers": {
									Description: "A list of specifier names of the performance counters you want to collect. To get a list of performance counters on Windows, run the command `typeperf`. See [this document](https://learn.microsoft.com/azure/azure-monitor/agents/data-sources-performance-counters#configure-performance-counters) for more information.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
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
									Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
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
									Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
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
												Description: "The label of the filter. This label should be unique across all `label_include_filter` block.",
												Type:        pluginsdk.TypeString,
												Required:    true,
												ValidateFunc: validation.StringInSlice([]string{
													"microsoft_metrics_include_label",
												}, false),
											},
											"value": {
												Description:  "The value of the filter.",
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
									Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"facility_names": {
									Description: "A list of facility names. Use a wildcard `*` to collect logs for all facility names.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											datacollectionrules.PossibleValuesForKnownSyslogDataSourceFacilityNames(),
											false,
										),
									},
								},
								"log_levels": {
									Description: "A list of log levels. Use a wildcard `*` to collect logs for all log levels.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											datacollectionrules.PossibleValuesForKnownSyslogDataSourceLogLevels(), false,
										),
									},
								},
								// lintignore:S013
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
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
									Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"x_path_queries": {
									Description: "A list of Windows Event Log queries in XPath expression. See [this document](https://learn.microsoft.com/azure/azure-monitor/agents/data-collection-rule-azure-monitor-agent?tabs=cli#filter-events-using-xpath-queries) for more information.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
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
									Description:  "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
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
			Description: "The description of the Data Collection Rule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"kind": {
			Description: "The kind of the Data Collection Rule. A rule of kind `Linux` does not allow for `windows_event_log` data sources, and a rule of kind `Windows` does not allow for `syslog` data sources. If kind is not specified, all kinds of data sources are allowed. A rule of kind `Direct` must use a `data_collection_endpoint_id` if you're sending data to a Log Analytics Workspace configured for Private Link.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					"Linux",
					"Windows",
					"AgentDirectToStore",
					"WorkspaceTransforms",
					"Direct",
				},
				false,
			),
		},

		"stream_declaration": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"stream_name": {
						Description:  "The name of the custom stream. This name should be unique across all `stream_declaration` blocks and must begin with a prefix of `Custom-`.",
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
									Description:  "The name of the column.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"type": {
									Description:  "The type of the column data.",
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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"immutable_id": {
			Description: "The immutable ID of the Data Collection Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"logs_ingestion_endpoint": {
			Description: "The endpoint used for ingesting logs.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"metrics_ingestion_endpoint": {
			Description: "The endpoint used for ingesting metrics.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldValue, newValue := metadata.ResourceDiff.GetChange("kind"); oldValue.(string) != newValue.(string) && oldValue.(string) != "" {
				if err := metadata.ResourceDiff.ForceNew("kind"); err != nil {
					return err
				}
			}

			if kind := metadata.ResourceDiff.Get("kind").(string); kind == "Direct" {
				if v := metadata.ResourceDiff.Get("data_sources").([]any); len(v) > 0 {
					return fmt.Errorf("`data_sources` cannot be used when `kind` is set to `Direct`")
				}

				if v := metadata.ResourceDiff.Get("destinations.0.event_hub_direct").([]any); len(v) > 0 {
					return fmt.Errorf("`event_hub_direct`, `storage_blob_direct`, and `storage_table_direct` destinations cannot be used when `kind` is set to `Direct`")
				}

				if v := metadata.ResourceDiff.Get("destinations.0.storage_blob_direct").([]any); len(v) > 0 {
					return fmt.Errorf("`event_hub_direct`, `storage_blob_direct`, and `storage_table_direct` destinations cannot be used when `kind` is set to `Direct`")
				}

				if v := metadata.ResourceDiff.Get("destinations.0.storage_table_direct").([]any); len(v) > 0 {
					return fmt.Errorf("`event_hub_direct`, `storage_blob_direct`, and `storage_table_direct` destinations cannot be used when `kind` is set to `Direct`")
				}
			}

			return nil
		},
	}
}
