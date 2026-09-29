// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithUpdate         = CdnFrontDoorRuleResource{}
	_ sdk.ResourceWithCustomImporter = CdnFrontDoorRuleResource{}
)

type CdnFrontDoorRuleResource struct{}

type CdnFrontDoorRuleResourceModel struct {
	Name                    string                            `tfschema:"name"`
	CdnFrontDoorRuleSetID   string                            `tfschema:"cdn_frontdoor_rule_set_id"`
	BehaviourOnMatch        string                            `tfschema:"behaviour_on_match"`
	Order                   int64                             `tfschema:"order"`
	Actions                 []CdnFrontDoorRuleActionsModel    `tfschema:"actions"`
	Conditions              []CdnFrontDoorRuleConditionsModel `tfschema:"conditions"`
	CdnFrontDoorRuleSetName string                            `tfschema:"cdn_frontdoor_rule_set_name"`
}

func (c CdnFrontDoorRuleResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, rmd sdk.ResourceMetaData) error {
		client := rmd.Client.Cdn.FrontDoorRuleSetsClient

		id, err := rules.ParseRuleID(rmd.ResourceData.Id())
		if err != nil {
			return err
		}
		ruleSetID := rulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)

		resp, err := client.Get(ctx, ruleSetID)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", ruleSetID, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil`", id)
		}

		if resp.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: `properties` was nil`", id)
		}

		if pointer.From(resp.Model.Properties.BatchMode) {
			return fmt.Errorf("the parent ruleset (%s) was provisioned using batch mode, and individual rules for this cannot be managed by this resource, use `azurerm_cdn_frontdoor_batch_rule_set` instead, or create a non-batch Rule Set with `azurerm_cdn_frontdoor_rule_set`", ruleSetID)
		}

		return nil
	}
}

func (c CdnFrontDoorRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateCdnFrontDoorRuleName,
		},

		"cdn_frontdoor_rule_set_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: rules.ValidateRuleSetID,
		},

		"behaviour_on_match": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(rules.MatchProcessingBehaviorContinue),
			ValidateFunc: validation.StringInSlice(rules.PossibleValuesForMatchProcessingBehavior(), false),
		},

		"order": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"actions": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url_redirect": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"redirect_type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRedirectType(), false),
								},

								"redirect_protocol": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(rules.DestinationProtocolMatchRequest),
									ValidateFunc: validation.StringInSlice(rules.PossibleValuesForDestinationProtocol(), false),
								},

								// Omit to preserve the incoming path. Issue #18249
								"destination_path": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringStartsWithOneOf("/"),
								},

								// Omit to preserve the incoming host. Issue #18249
								"destination_host_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringLenBetween(1, 2048),
								},

								// Omit to preserve the incoming query string. Issue #18249 & #19682
								"query_string": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									// Update validation logic to match RP. Issue #19097
									ValidateFunc: validateCdnFrontDoorUrlRedirectActionQueryString,
								},

								// NOTE: it is valid for the destination fragment to be an empty string.
								// Leave blank to preserve the incoming fragment. Issue #18249
								"destination_fragment": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.All(
										validation.StringLenBetween(1, 1024),
										validation.StringDoesNotStartWithOneOf("#"),
									),
								},
							},
						},
					},

					"url_rewrite": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"source_pattern": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringStartsWithOneOf("/"),
								},

								"destination_path": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringStartsWithOneOf("/"),
								},

								"preserve_unmatched_path_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},

					"modify_request_header": cdnFrontDoorRuleActionModifyHeaderSchema(),

					"modify_response_header": cdnFrontDoorRuleActionModifyHeaderSchema(),

					"route_configuration_override": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"origin_group": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"cdn_frontdoor_origin_group_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: afdorigingroups.ValidateOriginGroupID,
											},
											"forwarding_protocol": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(rules.PossibleValuesForForwardingProtocol(), false),
											},
										},
									},
								},

								"caching": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"behaviour": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(PossibleValuesForRuleCacheBehavior(), false),
											},
											"duration": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validateCdnFrontDoorCacheDuration,
											},
											"compression_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  false,
											},
											"query_string_behaviour": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRuleQueryStringCachingBehavior(), false),
											},
											"query_string_parameters": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 100,
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
				},
			},
		},

		"conditions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"remote_address": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues([]string{string(rules.RemoteAddressOperatorGeoMatch), string(rules.RemoteAddressOperatorIPMatch)}), false),
								},
								"values": cdnFrontDoorRuleValuesRequiredSchema(),
							},
						},
					},

					"request_method": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestMethodOperator()), false),
								},
								"values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									MaxItems: 7,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRequestMethodMatchValue(), false),
									},
								},
							},
						},
					},

					"query_string": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForQueryStringOperator()), false),
								},
								"values":     cdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"post_argument": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": cdnFrontDoorRuleConditionNameSchema(),
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForPostArgsOperator()), false),
								},
								"values":     cdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"request_url": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestUriOperator()), false),
								},
								"values":     cdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"request_header": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": cdnFrontDoorRuleConditionNameSchema(),
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestHeaderOperator()), false),
								},
								"values":     cdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"request_body": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestBodyOperator()), false),
								},
								"values":     cdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"request_scheme": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestSchemeMatchConditionParametersOperator()), false),
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									// This uses a List instead of a String to stay consistent with the other conditions, even though only 1 item can be defined
									MaxItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRequestSchemeMatchValue(), false),
									},
								},
							},
						},
					},

					"request_path": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForURLPathOperator()), false),
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 25,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"request_file_extension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForURLFileExtensionOperator()), false),
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 25,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"request_filename": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForURLFileNameOperator()), false),
								},
								// making optional for issue #23504
								"values":     cdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"http_version": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForHTTPVersionOperator()), false),
								},
								"values": cdnFrontDoorRuleHttpVersionValuesSchema(),
							},
						},
					},

					"request_cookies": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": cdnFrontDoorRuleConditionNameSchema(),
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForCookiesOperator()), false),
								},
								"values":     cdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"device_type": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForIsDeviceOperator()), false),
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									// This uses a List instead of a String to stay consistent with the other conditions, even though only 1 item can be defined
									MaxItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(rules.PossibleValuesForIsDeviceMatchValue(), false),
									},
								},
							},
						},
					},

					"socket_address": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues([]string{string(rules.SocketAddrOperatorIPMatch)}), false),
								},
								"values": cdnFrontDoorSocketAddressValuesSchema(),
							},
						},
					},

					"client_port": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForClientPortOperator()), false),
								},
								"values": cdnFrontDoorRuleValuesOptionalSchema(),
							},
						},
					},

					"server_port": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForQueryStringOperator()), false),
								},
								"values": cdnFrontDoorServerPortValuesSchema(),
							},
						},
					},

					"host_name": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForHostNameOperator()), false),
								},
								"values":     cdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
							},
						},
					},

					"ssl_protocol": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(cdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForSslProtocolOperator()), false),
								},
								"values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									MaxItems: 3,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(rules.PossibleValuesForSslProtocol(), false),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c CdnFrontDoorRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cdn_frontdoor_rule_set_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (c CdnFrontDoorRuleResource) ModelObject() interface{} {
	return &CdnFrontDoorRuleResourceModel{}
}

func (c CdnFrontDoorRuleResource) ResourceType() string {
	return "azurerm_cdn_frontdoor_rule"
}

func (c CdnFrontDoorRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.Cdn.FrontDoorRulesClient
			ruleSetsClient := rmd.Client.Cdn.FrontDoorRuleSetsClient

			config := CdnFrontDoorRuleResourceModel{}
			if err := rmd.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			ruleSetId, err := rulesets.ParseRuleSetID(config.CdnFrontDoorRuleSetID)
			if err != nil {
				return err
			}

			id := rules.NewRuleID(ruleSetId.SubscriptionId, ruleSetId.ResourceGroupName, ruleSetId.ProfileName, ruleSetId.RuleSetName, config.Name)

			ruleSet, err := ruleSetsClient.Get(ctx, *ruleSetId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", ruleSetId, err)
			}

			if ruleSet.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", ruleSetId)
			}

			if ruleSet.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", ruleSetId)
			}

			if pointer.From(ruleSet.Model.Properties.BatchMode) {
				return fmt.Errorf("the parent ruleset (%s) was provisioned using batch mode, and individual rules for this cannot be managed by this resource, use `azurerm_cdn_frontdoor_batch_rule_set` instead, or create a non-batch Rule Set with `azurerm_cdn_frontdoor_rule_set`", ruleSetId)
			}

			if !rmd.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				result, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(result.HttpResponse) {
						return fmt.Errorf("checking for existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(result.HttpResponse) {
					return rmd.ResourceRequiresImport("azurerm_cdn_frontdoor_rule", id)
				}
			}

			actions, err := expandCdnFrontDoorRuleActions(config.Actions)
			if err != nil {
				return fmt.Errorf("expanding `actions`: %+v", err)
			}

			conditions, err := expandCdnFrontDoorRuleConditions(config.Conditions)
			if err != nil {
				return fmt.Errorf("expanding `conditions`: %+v", err)
			}

			props := rules.Rule{
				Properties: &rules.RuleProperties{
					Actions:                 &actions,
					Conditions:              &conditions,
					MatchProcessingBehavior: pointer.ToEnum[rules.MatchProcessingBehavior](config.BehaviourOnMatch),
					RuleSetName:             &ruleSetId.RuleSetName,
					Order:                   pointer.To(config.Order),
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, props, rmd.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			rmd.SetID(id)

			return nil
		},
	}
}

func (c CdnFrontDoorRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.Cdn.FrontDoorRulesClient

			id, err := rules.ParseRuleID(rmd.ResourceData.Id())
			if err != nil {
				return err
			}

			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return rmd.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := CdnFrontDoorRuleResourceModel{
				Name:                    id.RuleName,
				CdnFrontDoorRuleSetID:   rulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName).ID(),
				CdnFrontDoorRuleSetName: id.RuleSetName,
			}

			if model := result.Model; model != nil {
				if props := model.Properties; props != nil {
					state.BehaviourOnMatch = pointer.FromEnum(props.MatchProcessingBehavior)
					state.Order = pointer.From(props.Order)

					actions, err := flattenCdnFrontDoorRuleActions(props.Actions)
					if err != nil {
						return fmt.Errorf("flattening `actions`: %+v", err)
					}
					state.Actions = actions

					conditions, err := flattenCdnFrontDoorRuleConditions(props.Conditions)
					if err != nil {
						return fmt.Errorf("flattening `conditions`: %+v", err)
					}
					state.Conditions = conditions
				}
			}

			return rmd.Encode(&state)
		},
	}
}

func (c CdnFrontDoorRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.Cdn.FrontDoorRulesClient

			state := CdnFrontDoorRuleResourceModel{}
			if err := rmd.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			id, err := rules.ParseRuleID(rmd.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}
			props := existing.Model.Properties

			if rmd.ResourceData.HasChange("behaviour_on_match") {
				props.MatchProcessingBehavior = pointer.ToEnum[rules.MatchProcessingBehavior](state.BehaviourOnMatch)
			}

			if rmd.ResourceData.HasChange("order") {
				props.Order = pointer.To(state.Order)
			}

			if rmd.ResourceData.HasChange("actions") {
				actions, err := expandCdnFrontDoorRuleActions(state.Actions)
				if err != nil {
					return fmt.Errorf("expanding `actions`: %+v", err)
				}
				props.Actions = &actions
			}

			if rmd.ResourceData.HasChange("conditions") {
				conditions, err := expandCdnFrontDoorRuleConditions(state.Conditions)
				if err != nil {
					return fmt.Errorf("expanding `conditions`: %+v", err)
				}

				if len(conditions) > 10 {
					return fmt.Errorf("expanding `conditions`: configuration file exceeds the maximum of 10 match conditions, got %d", len(conditions))
				}

				props.Conditions = &conditions
			}

			if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (c CdnFrontDoorRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.Cdn.FrontDoorRulesClient

			id, err := rules.ParseRuleID(rmd.ResourceData.Id())
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

func (c CdnFrontDoorRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return rules.ValidateRuleID
}
