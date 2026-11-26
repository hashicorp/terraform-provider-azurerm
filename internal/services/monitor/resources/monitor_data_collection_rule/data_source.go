// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

func (d DataSource) ModelObject() any {
	return &DataCollectionRule{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_monitor_data_collection_rule"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Data Collection Rule."
}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Collection Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"data_collection_endpoint_id": {
			Description: "The resource ID of the Data Collection Endpoint that this rule can be used with.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"data_flow": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destinations": {
						Description: "A list of destination names. An `azure_monitor_metrics` data source only allows for stream of kind `Microsoft-InsightsMetrics`.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"streams": {
						Description: "A list of streams.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"output_stream": {
						Description: "The output stream of the transform. Only required if the data flow changes data to a different stream.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"transform_kql": {
						Description: "The KQL query to transform stream data.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"built_in_transform": {
						Description: "The built-in transform to transform stream data.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"destinations": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"event_hub": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"event_hub_id": {
									Description: "The resource ID of the Event Hub.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"name": {
									Description: "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"event_hub_direct": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"event_hub_id": {
									Description: "The resource ID of the Event Hub.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"name": {
									Description: "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"azure_monitor_metrics": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"log_analytics": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"workspace_resource_id": {
									Description: "The ID of a Log Analytic Workspace resource.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"monitor_account": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"monitor_account_id": {
									Description: "The resource ID of the Monitor Account.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"storage_blob": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"container_name": {
									Description: "The Storage Container name.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"storage_account_id": {
									Description: "The resource ID of the Storage Account.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"storage_blob_direct": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"container_name": {
									Description: "The Storage Container name.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"storage_account_id": {
									Description: "The resource ID of the Storage Account.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"storage_table_direct": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "The name which should be used for this destination. This name should be unique across all destinations regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"table_name": {
									Description: "The Storage Table name.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"storage_account_id": {
									Description: "The resource ID of the Storage Account.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
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
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_import": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"event_hub_data_source": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},
											"stream": {
												Description: "The stream to collect from Event Hub.",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},
											"consumer_group": {
												Description: "The Event Hub consumer group name.",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},
										},
									},
								},
							},
						},
					},
					"extension": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"extension_name": {
									Description: "The name of the VM extension.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"extension_json": {
									Description: "A JSON String which specifies the extension setting.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"input_data_sources": {
									Description: "A list of data sources this extension needs data from. An item should be a name of a supported data source which produces only one stream. Supported data sources type: `performance_counter`, `windows_event_log`, and `syslog`.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"iis_log": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"log_directories": {
									Description: "A list of absolute paths where the log files are located.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"log_file": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"file_patterns": {
									Description: "A list of file patterns where the log files are located. For example, `C:\\\\JavaLogs\\\\*.log`.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"format": {
									Description: "The data format of the log files.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"settings": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"text": {
												Type:     pluginsdk.TypeList,
												Computed: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"record_start_timestamp_format": {
															Description: "The timestamp format of the text log files.",
															Type:        pluginsdk.TypeString,
															Computed:    true,
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
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"sampling_frequency_in_seconds": {
									Description: "The number of seconds between consecutive counter measurements (samples).",
									Type:        pluginsdk.TypeInt,
									Computed:    true,
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"counter_specifiers": {
									Description: "A list of specifier names of the performance counters you want to collect. Use a wildcard `*` to collect counters for all instances. To get a list of performance counters on Windows, run the command `typeperf`.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"platform_telemetry": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"prometheus_forwarder": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"label_include_filter": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"label": {
												Description: "The label of the filter. This label should be unique across all `label_include_filter` block.",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},
											"value": {
												Description: "The value of the filter.",
												Type:        pluginsdk.TypeString,
												Computed:    true,
											},
										},
									},
								},
							},
						},
					},
					"syslog": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"facility_names": {
									Description: "A list of facility names. Use a wildcard `*` to collect logs for all facility names.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"log_levels": {
									Description: "A list of log levels. Use a wildcard `*` to collect logs for all log levels.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"windows_event_log": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"x_path_queries": {
									Description: "A list of Windows Event Log queries in XPath expression.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"windows_firewall_log": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name which should be used for this data source. This name should be unique across all data sources regardless of type within the Data Collection Rule.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"streams": {
									Description: "A list of streams that this data source will be sent to. A stream indicates what schema will be used for this data and usually what table in Log Analytics the data will be sent to.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
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
			Computed:    true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"immutable_id": {
			Description: "The immutable ID of the Data Collection Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kind": {
			Description: "The kind of the Data Collection Rule. A rule of kind `Linux` does not allow for `windows_event_log` data sources. A rule of kind `Windows` does not allow for `syslog` data sources. If kind is not specified, all kinds of data sources are allowed.",
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

		"stream_declaration": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"stream_name": {
						Description: "The name of the custom stream. This name should be unique across all `stream_declaration` blocks.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"column": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "The name of the column.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"type": {
									Description: "The type of the column data.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DataCollectionRule
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := datacollectionrules.NewDataCollectionRuleID(subscriptionId, state.ResourceGroupName, state.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var dataCollectionEndpointId, description, immutableId, kind, loc string
			var tag map[string]any
			var dataFlows []DataFlow
			var dataSources []DataSourceModel
			var destinations []Destination
			var logsIngestionEndpoint, metricsIngestionEndpoint string
			var streamDeclaration []StreamDeclaration

			if model := resp.Model; model != nil {
				kind = pointer.FromEnum(model.Kind)
				loc = location.Normalize(model.Location)
				tag = tags.Flatten(model.Tags)

				identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				if prop := model.Properties; prop != nil {
					dataCollectionEndpointId = pointer.From(prop.DataCollectionEndpointId)
					description = pointer.From(prop.Description)
					dataFlows = flattenDataCollectionRuleDataFlows(prop.DataFlows)
					dataSources = flattenDataCollectionRuleDataSources(prop.DataSources)
					destinations = flattenDataCollectionRuleDestinations(prop.Destinations)
					logsIngestionEndpoint, metricsIngestionEndpoint = flattenDataCollectionRuleEndpoints(prop.Endpoints)
					immutableId = pointer.From(prop.ImmutableId)
					streamDeclaration = flattenDataCollectionRuleStreamDeclarations(prop.StreamDeclarations)
				}
			}

			metadata.SetID(id)

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
				Location:                 loc,
				LogsIngestionEndpoint:    logsIngestionEndpoint,
				MetricsIngestionEndpoint: metricsIngestionEndpoint,
				StreamDeclaration:        streamDeclaration,
				Tags:                     tag,
			})
		},
		Timeout: 5 * time.Minute,
	}
}
