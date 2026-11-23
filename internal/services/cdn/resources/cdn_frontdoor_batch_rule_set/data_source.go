// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_batch_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

type CdnFrontDoorBatchRuleSetDataSourceModel struct {
	Name                  string                              `tfschema:"name"`
	ProfileName           string                              `tfschema:"profile_name"`
	ResourceGroupName     string                              `tfschema:"resource_group_name"`
	CdnFrontDoorProfileID string                              `tfschema:"cdn_frontdoor_profile_id"`
	Rules                 []CdnFrontDoorBatchRuleSetRuleModel `tfschema:"rules"`
}

func (DataSource) ResourceType() string {
	return "azurerm_cdn_frontdoor_batch_rule_set"
}

func (DataSource) Description() string {
	return "Gets information about an existing Front Door (standard/premium) Batch Rule Set."
}

func (DataSource) ModelObject() any {
	return &CdnFrontDoorBatchRuleSetDataSourceModel{}
}

func (DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Front Door Batch Rule Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.FrontDoorRuleSetName,
		},
		"profile_name": {
			Description:  "The name of the Front Door Profile.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.FrontDoorName,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cdn_frontdoor_profile_id": {
			Description: "The ID of the Front Door Profile associated with this Front Door Batch Rule Set.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"rules": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the Front Door Batch Rule.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"behaviour_on_match": {
						Description: "Whether the rules engine continues processing after this rule matches.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"order": {
						Description: "The order in which this rule is applied.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
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
											Description: "The redirect type for the redirect action.",
											Type:        pluginsdk.TypeString,
											Computed:    true,
										},
										"redirect_protocol": {
											Description: "The redirect protocol for the redirect action.",
											Type:        pluginsdk.TypeString,
											Computed:    true,
										},
										"destination_path": {
											Description: "The destination path for the redirect action.",
											Type:        pluginsdk.TypeString,
											Computed:    true,
										},
										"destination_host_name": {
											Description: "The destination host name for the redirect action.",
											Type:        pluginsdk.TypeString,
											Computed:    true,
										},
										"destination_fragment": {
											Description: "The destination fragment for the redirect action.",
											Type:        pluginsdk.TypeString,
											Computed:    true,
										},
										"query_string": {
											Description: "The query string for the redirect action.",
											Type:        pluginsdk.TypeString,
											Computed:    true,
										},
									}},
								},
								"url_rewrite": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
										"source_pattern": {
											Description: "The source pattern for the rewrite action.",
											Type:        pluginsdk.TypeString,
											Computed:    true,
										},
										"destination_path": {
											Description: "The destination path for the rewrite action.",
											Type:        pluginsdk.TypeString,
											Computed:    true,
										},
										"preserve_unmatched_path_enabled": {
											Description: "Whether to preserve the unmatched part of the path.",
											Type:        pluginsdk.TypeBool,
											Computed:    true,
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
													Description: "The ID of the Front Door Origin Group associated with this action.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
												"forwarding_protocol": {
													Description: "The forwarding protocol applied to this action.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
											}},
										},
										"caching": {
											Type:     pluginsdk.TypeList,
											Computed: true,
											Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
												"behaviour": {
													Description: "The cache behaviour applied to this action.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
												"duration": {
													Description: "The cache duration applied to this action.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
												"compression_enabled": {
													Description: "Whether compression is enabled.",
													Type:        pluginsdk.TypeBool,
													Computed:    true,
												},
												"query_string_behaviour": {
													Description: "The query string caching behaviour applied to this action.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
												"query_string_parameters": cdnFrontDoorBatchRuleSetConditionStringListSchema("The query string parameters associated with this action."),
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
							"remote_address":         cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema("The remote address values associated with this condition."),
							"request_method":         cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema("The request methods associated with this condition."),
							"query_string":           cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema("The query string values associated with this condition."),
							"post_argument":          cdnFrontDoorBatchRuleSetConditionWithNameAndTransformsComputedSchema("The POST argument values associated with this condition."),
							"request_url":            cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema("The request URL values associated with this condition."),
							"request_header":         cdnFrontDoorBatchRuleSetConditionWithNameAndTransformsComputedSchema("The request header values associated with this condition."),
							"request_body":           cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema("The request body values associated with this condition."),
							"request_scheme":         cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema("The request schemes associated with this condition."),
							"request_path":           cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema("The request path values associated with this condition."),
							"request_file_extension": cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema("The request file extension values associated with this condition."),
							"request_filename":       cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema("The request file name values associated with this condition."),
							"http_version":           cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema("The HTTP versions associated with this condition."),
							"request_cookies":        cdnFrontDoorBatchRuleSetConditionWithNameAndTransformsComputedSchema("The cookie values associated with this condition."),
							"device_type":            cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema("The device types associated with this condition."),
							"socket_address":         cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema("The socket address values associated with this condition."),
							"client_port":            cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema("The client port values associated with this condition."),
							"server_port":            cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema("The server port values associated with this condition."),
							"host_name":              cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema("The host names associated with this condition."),
							"ssl_protocol":           cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema("The SSL protocol values associated with this condition."),
						}},
					},
				},
			},
		},
	}
}

func (DataSource) Read() sdk.ResourceFunc {
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
				Description: "The action applied to the header.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"header_name": {
				Description: "The header name.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"header_value": {
				Description: "The header value.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionBaseListValuesComputedSchema(valuesDescription string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Description: "The operator for this condition.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"values": cdnFrontDoorBatchRuleSetConditionStringListSchema(valuesDescription),
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionBaseSetValuesComputedSchema(valuesDescription string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Description: "The operator for this condition.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"values": cdnFrontDoorBatchRuleSetConditionStringSetSchema(valuesDescription),
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionWithTransformsComputedSchema(valuesDescription string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Description: "The operator for this condition.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"values":     cdnFrontDoorBatchRuleSetConditionStringListSchema(valuesDescription),
			"transforms": cdnFrontDoorBatchRuleSetConditionStringSetSchema("The transforms associated with this condition."),
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionWithNameAndTransformsComputedSchema(valuesDescription string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"operator": {
				Description: "The operator for this condition.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"values":     cdnFrontDoorBatchRuleSetConditionStringListSchema(valuesDescription),
			"transforms": cdnFrontDoorBatchRuleSetConditionStringSetSchema("The transforms associated with this condition."),
		}},
	}
}

func cdnFrontDoorBatchRuleSetConditionStringListSchema(description string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description: description,
		Type:        pluginsdk.TypeList,
		Computed:    true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func cdnFrontDoorBatchRuleSetConditionStringSetSchema(description string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description: description,
		Type:        pluginsdk.TypeSet,
		Computed:    true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}
