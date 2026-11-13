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
												ValidateFunc: validation.StringLenBetween(1, 32),
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
									ValidateFunc: validation.StringLenBetween(1, 32),
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
									ValidateFunc: validation.StringLenBetween(1, 32),
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
									ValidateFunc: validation.StringLenBetween(1, 32),
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
									ValidateFunc: validation.StringLenBetween(1, 32),
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
									ValidateFunc: validation.StringLenBetween(1, 32),
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
									ValidateFunc: validation.StringLenBetween(1, 32),
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
									ValidateFunc: validation.StringLenBetween(1, 32),
								},
								"facility_names": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											datacollectionrules.PossibleValuesForKnownSyslogDataSourceFacilityNames(),
											false,
										),
									},
								},
								"log_levels": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											datacollectionrules.PossibleValuesForKnownSyslogDataSourceLogLevels(), false,
										),
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
									ValidateFunc: validation.StringLenBetween(1, 32),
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
									ValidateFunc: validation.StringLenBetween(1, 32),
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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"immutable_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"logs_ingestion_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"metrics_ingestion_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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
