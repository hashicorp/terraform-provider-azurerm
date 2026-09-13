// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = CdnFrontDoorBatchRuleSetDataSource{}

type CdnFrontDoorBatchRuleSetDataSource struct{}

type CdnFrontDoorBatchRuleSetDataSourceModel struct {
	Name                  string                              `tfschema:"name"`
	ProfileName           string                              `tfschema:"profile_name"`
	ResourceGroupName     string                              `tfschema:"resource_group_name"`
	CdnFrontDoorProfileID string                              `tfschema:"cdn_frontdoor_profile_id"`
	Rules                 []CdnFrontDoorBatchRuleSetRuleModel `tfschema:"rules"`
}

func (CdnFrontDoorBatchRuleSetDataSource) ResourceType() string {
	return "azurerm_cdn_frontdoor_batch_rule_set"
}

func (CdnFrontDoorBatchRuleSetDataSource) ModelObject() interface{} {
	return &CdnFrontDoorBatchRuleSetDataSourceModel{}
}

func (CdnFrontDoorBatchRuleSetDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.FrontDoorRuleSetName,
		},
		"profile_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.FrontDoorName,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (CdnFrontDoorBatchRuleSetDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cdn_frontdoor_profile_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"rules": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"behaviour_on_match": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"order": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"actions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"url_redirect": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
										"redirect_type": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"redirect_protocol": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"destination_path": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"destination_host_name": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"destination_fragment": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"query_string": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
									}},
								},
								"url_rewrite": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
										"source_pattern": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"destination_path": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"preserve_unmatched_path_enabled": {
											Type:     pluginsdk.TypeBool,
											Computed: true,
										},
									}},
								},
								"modify_request_header":  cdnFrontDoorBatchRuleSetRuleHeaderActionComputedSchema(),
								"modify_response_header": cdnFrontDoorBatchRuleSetRuleHeaderActionComputedSchema(),
								"route_configuration_override": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
										"origin_group": {
											Type:     pluginsdk.TypeList,
											Computed: true,
											Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
												"cdn_frontdoor_origin_group_id": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"forwarding_protocol": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
											}},
										},
										"caching": {
											Type:     pluginsdk.TypeList,
											Computed: true,
											Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
												"behaviour": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"duration": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"compression_enabled": {
													Type:     pluginsdk.TypeBool,
													Computed: true,
												},
												"query_string_behaviour": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"query_string_parameters": cdnFrontDoorBatchRuleSetConditionStringListSchema(),
											}},
										},
									}},
								},
							},
						},
					},
					"conditions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
							"remote_address":         cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema(),
							"request_method":         cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema(),
							"query_string":           cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema(),
							"post_argument":          cdnFrontDoorBatchRuleSetConditionWithNameAndTransformsComputedSchema(),
							"request_url":            cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema(),
							"request_header":         cdnFrontDoorBatchRuleSetConditionWithNameAndTransformsComputedSchema(),
							"request_body":           cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema(),
							"request_scheme":         cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema(),
							"request_path":           cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema(),
							"request_file_extension": cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema(),
							"request_filename":       cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema(),
							"http_version":           cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema(),
							"request_cookies":        cdnFrontDoorBatchRuleSetConditionWithNameAndTransformsComputedSchema(),
							"device_type":            cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema(),
							"socket_address":         cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema(),
							"client_port":            cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema(),
							"server_port":            cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema(),
							"host_name":              cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema(),
							"ssl_protocol":           cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema(),
						}},
					},
				},
			},
		},
	}
}

func (CdnFrontDoorBatchRuleSetDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient

			var state CdnFrontDoorBatchRuleSetDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := rulesets.NewRuleSetID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.ProfileName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model != nil && resp.Model.Properties != nil && !pointer.From(resp.Model.Properties.BatchMode) {
				return fmt.Errorf("%s was not provisioned using batch mode, and cannot be read by this data source, use `azurerm_cdn_frontdoor_rule_set` instead", id)
			}

			metadata.SetID(&id)

			state.Name = id.RuleSetName
			state.ProfileName = id.ProfileName
			state.ResourceGroupName = id.ResourceGroupName
			state.CdnFrontDoorProfileID = profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					rulesState, err := flattenCdnFrontDoorBatchRuleSetRules(props.Rules)
					if err != nil {
						return fmt.Errorf("flattening `rules`: %+v", err)
					}
					state.Rules = rulesState
				}
			}

			return metadata.Encode(&state)
		},
	}
}

// Reusable schema helpers

func cdnFrontDoorBatchRuleSetRuleHeaderActionComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"header_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"header_value": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"values": cdnFrontDoorBatchRuleSetConditionStringListSchema(),
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"values": cdnFrontDoorBatchRuleSetConditionStringSetSchema(),
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"values":     cdnFrontDoorBatchRuleSetConditionStringListSchema(),
			"transforms": cdnFrontDoorBatchRuleSetConditionStringSetSchema(),
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionWithNameAndTransformsComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"operator": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"values":     cdnFrontDoorBatchRuleSetConditionStringListSchema(),
			"transforms": cdnFrontDoorBatchRuleSetConditionStringSetSchema(),
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionStringListSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func cdnFrontDoorBatchRuleSetConditionStringSetSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Computed: true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}
