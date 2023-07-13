// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataCollectionRuleDataSource struct{}

var _ sdk.DataSource = DataCollectionRuleDataSource{}

func (d DataCollectionRuleDataSource) ModelObject() interface{} {
	return &DataCollectionRule{}
}

func (d DataCollectionRuleDataSource) ResourceType() string {
	return "azurerm_monitor_data_collection_rule"
}

func (d DataCollectionRuleDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d DataCollectionRuleDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"location": commonschema.LocationComputed(),

		"data_collection_endpoint_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"data_flow": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destinations": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"streams": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"output_stream": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"transform_kql": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"built_in_transform": {
						Type:     pluginsdk.TypeString,
						Computed: true,
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
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"event_hub_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
					"event_hub_direct": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"event_hub_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"workspace_resource_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"monitor_account_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"container_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"storage_account_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"container_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"storage_account_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"table_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"storage_account_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"data_sources": {
			Type:     pluginsdk.TypeList,
			Computed: true,
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
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"stream": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"consumer_group": {
												Type:     pluginsdk.TypeString,
												Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"extension_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"extension_json": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"input_data_sources": {
									Type:     pluginsdk.TypeList,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"log_directories": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"file_patterns": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"format": {
									Type:     pluginsdk.TypeString,
									Computed: true,
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
						},
					},
					"performance_counter": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"sampling_frequency_in_seconds": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"counter_specifiers": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Computed: true,
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
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"value": {
												Type:     pluginsdk.TypeString,
												Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"facility_names": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"log_levels": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"x_path_queries": {
									Type:     pluginsdk.TypeList,
									Computed: true,
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Computed: true,
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
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"immutable_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"stream_declaration": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"stream_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"column": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
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

func (d DataCollectionRuleDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DataCollectionRule
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := datacollectionrules.NewDataCollectionRuleID(subscriptionId, state.ResourceGroupName, state.Name)
			metadata.Logger.Infof("retrieving %s", id)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
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
				Location:                 location,
				StreamDeclaration:        streamDeclaration,
				Tags:                     tag,
			})
		},
		Timeout: 5 * time.Minute,
	}
}
