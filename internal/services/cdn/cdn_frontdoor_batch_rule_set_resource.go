// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cdn_frontdoor_batch_rule_set -properties "name" -compare-values "subscription_id:cdn_frontdoor_profile_id,resource_group_name:cdn_frontdoor_profile_id,profile_name:cdn_frontdoor_profile_id" -test-name "basicAttachedRoute"

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithCustomImporter = CdnFrontDoorBatchRuleSetResource{}
	_ sdk.ResourceWithCustomizeDiff  = CdnFrontDoorBatchRuleSetResource{}
	_ sdk.ResourceWithIdentity       = CdnFrontDoorBatchRuleSetResource{}
	_ sdk.ResourceWithUpdate         = CdnFrontDoorBatchRuleSetResource{}
)

type CdnFrontDoorBatchRuleSetResource struct{}

func (r CdnFrontDoorBatchRuleSetResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Cdn.FrontDoorRuleSetsClient

		id, err := rulesets.ParseRuleSetID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil`", id)
		}

		if resp.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: `properties` was nil`", id)
		}

		if !pointer.From(resp.Model.Properties.BatchMode) {
			return fmt.Errorf("%s was not provisioned using batch mode and cannot be managed by this resource, use `azurerm_cdn_frontdoor_rule_set` instead", id)
		}

		return nil
	}
}

func (CdnFrontDoorBatchRuleSetResource) Identity() resourceids.ResourceId {
	return &rulesets.RuleSetId{}
}

func (CdnFrontDoorBatchRuleSetResource) ResourceType() string {
	return "azurerm_cdn_frontdoor_batch_rule_set"
}

func (CdnFrontDoorBatchRuleSetResource) ModelObject() interface{} {
	return &CdnFrontDoorBatchRuleSetModel{}
}

func (CdnFrontDoorBatchRuleSetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return rulesets.ValidateRuleSetID
}

func (CdnFrontDoorBatchRuleSetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorRuleSetName,
		},
		"cdn_frontdoor_profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: profiles.ValidateProfileID,
		},
		"rule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 100,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.CdnFrontDoorRuleName,
					},
					"behaviour_on_match": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(rulesets.MatchProcessingBehaviorContinue),
						ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForMatchProcessingBehavior(), false),
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
												ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForRedirectType(), false),
											},
											"redirect_protocol": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      string(rulesets.DestinationProtocolMatchRequest),
												ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForDestinationProtocol(), false),
											},
											"destination_path": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringStartsWithOneOf("/"),
											},
											"destination_host_name": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 2048),
											},
											"destination_fragment": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												ValidateFunc: validation.All(
													validation.StringLenBetween(1, 1024),
													validation.StringDoesNotStartWithOneOf("#"),
												),
											},
											"query_string": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validate.CdnFrontDoorUrlRedirectActionQueryString,
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
								"modify_request_header":  cdnFrontDoorBatchRuleSetActionModifyHeaderSchema(),
								"modify_response_header": cdnFrontDoorBatchRuleSetActionModifyHeaderSchema(),
								"route_configuration_override": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
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
														ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForForwardingProtocol(), false),
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
														ValidateFunc: validate.CdnFrontDoorCacheDuration,
													},
													"compression_enabled": {
														Type:     pluginsdk.TypeBool,
														Optional: true,
														Default:  false,
													},
													"query_string_behaviour": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForRuleQueryStringCachingBehavior(), false),
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
									}},
								},
							},
						},
					},

					"conditions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
							"remote_address": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"operator": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues([]string{string(rulesets.RemoteAddressOperatorGeoMatch), string(rulesets.RemoteAddressOperatorIPMatch)}), false),
										},
										"values": {
											Type:     pluginsdk.TypeList,
											Required: true,
											MinItems: 1,
											MaxItems: 25,
											Elem: &pluginsdk.Schema{
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestMethodOperator()), false),
										},
										"values": {
											Type:     pluginsdk.TypeSet,
											Required: true,
											MinItems: 1,
											MaxItems: 7,
											Elem: &pluginsdk.Schema{
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForRequestMethodMatchValue(), false),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForQueryStringOperator()), false),
										},
										"values":     cdnFrontDoorBatchRuleSetConditionValuesSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
									},
								},
							},
							"post_argument": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"name": cdnFrontDoorBatchRuleSetConditionSelectorSchema(),
										"operator": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForPostArgsOperator()), false),
										},
										"values":     cdnFrontDoorBatchRuleSetConditionValuesSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestUriOperator()), false),
										},
										"values":     cdnFrontDoorBatchRuleSetConditionValuesSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
									},
								},
							},
							"request_header": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"name": cdnFrontDoorBatchRuleSetConditionSelectorSchema(),
										"operator": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestHeaderOperator()), false),
										},
										"values":     cdnFrontDoorBatchRuleSetConditionValuesSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestBodyOperator()), false),
										},
										"values":     cdnFrontDoorBatchRuleSetConditionValuesSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestSchemeMatchConditionParametersOperator()), false),
										},
										"values": {
											Type:     pluginsdk.TypeList,
											Required: true,
											MinItems: 1,
											// This uses a List instead of a String to stay consistent with the other conditions, even though only 1 item can be defined
											MaxItems: 1,
											Elem: &pluginsdk.Schema{
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForRequestSchemeMatchValue(), false),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForURLPathOperator()), false),
										},
										"values": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											MaxItems: 25,
											Elem: &pluginsdk.Schema{
												Type: pluginsdk.TypeString,
												ValidateFunc: validation.All(
													validation.StringIsNotEmpty,
													validation.StringDoesNotStartWithOneOf("/"),
												),
											},
										},
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForURLFileExtensionOperator()), false),
										},
										"values": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											MaxItems: 25,
											Elem: &pluginsdk.Schema{
												Type: pluginsdk.TypeString,
												ValidateFunc: validation.All(
													validation.StringIsNotEmpty,
													validation.StringDoesNotStartWithOneOf("."),
												),
											},
										},
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForURLFileNameOperator()), false),
										},
										"values":     cdnFrontDoorBatchRuleSetConditionValuesSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForHTTPVersionOperator()), false),
										},
										"values": {
											Type:     pluginsdk.TypeSet,
											Required: true,
											MinItems: 1,
											MaxItems: 4,
											Elem: &pluginsdk.Schema{
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice([]string{"2.0", "1.1", "1.0", "0.9"}, false),
											},
										},
									},
								},
							},
							"request_cookies": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"name": cdnFrontDoorBatchRuleSetConditionSelectorSchema(),
										"operator": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForCookiesOperator()), false),
										},
										"values":     cdnFrontDoorBatchRuleSetConditionValuesSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForIsDeviceOperator()), false),
										},
										"values": {
											Type:     pluginsdk.TypeList,
											Required: true,
											MinItems: 1,
											// This uses a List instead of a String to stay consistent with the other conditions, even though only 1 item can be defined
											MaxItems: 1,
											Elem: &pluginsdk.Schema{
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForIsDeviceMatchValue(), false),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues([]string{string(rulesets.SocketAddrOperatorIPMatch)}), false),
										},
										"values": {
											Type:     pluginsdk.TypeList,
											Required: true,
											MinItems: 1,
											MaxItems: 25,
											Elem: &pluginsdk.Schema{
												Type:         pluginsdk.TypeString,
												ValidateFunc: validate.FrontDoorRuleCidrIsValid,
											},
										},
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForClientPortOperator()), false),
										},
										"values": cdnFrontDoorBatchRuleSetConditionValuesSchema(),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForQueryStringOperator()), false),
										},
										"values": {
											Type:     pluginsdk.TypeSet,
											Optional: true,
											MaxItems: 2,
											Elem: &pluginsdk.Schema{
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice([]string{"80", "443"}, false),
											},
										},
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForHostNameOperator()), false),
										},
										"values":     cdnFrontDoorBatchRuleSetConditionValuesSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
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
											ValidateFunc: validation.StringInSlice(cdnFrontDoorBatchRuleSetConditionOperatorPossibleValues(rulesets.PossibleValuesForSslProtocolOperator()), false),
										},
										"values": {
											Type:     pluginsdk.TypeSet,
											Required: true,
											MinItems: 1,
											MaxItems: 3,
											Elem: &pluginsdk.Schema{
												Type:         pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForSslProtocol(), false),
											},
										},
									},
								},
							},
						}},
					},
				},
			},
		},
	}
}

func (CdnFrontDoorBatchRuleSetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r CdnFrontDoorBatchRuleSetResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CdnFrontDoorBatchRuleSetModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding diff: %+v", err)
			}

			if err := validateCdnFrontDoorBatchRules(model.Rule); err != nil {
				return err
			}

			if err := validateCdnFrontDoorBatchRuleDiffQuota(metadata.ResourceDiff); err != nil {
				return err
			}

			return nil
		},
	}
}

func validateCdnFrontDoorBatchRuleDiffQuota(diff *pluginsdk.ResourceDiff) error {
	hashFunc := pluginsdk.HashResource(CdnFrontDoorBatchRuleSetResource{}.Arguments()["rule"].Elem.(*pluginsdk.Resource))

	oldRaw, newRaw := diff.GetChange("rule")
	oldRules := cdnFrontDoorBatchRuleSetRules(oldRaw, hashFunc)
	newRules := cdnFrontDoorBatchRuleSetRules(newRaw, hashFunc)

	names := make(map[string]struct{}, len(oldRules)+len(newRules))
	for name := range oldRules {
		names[name] = struct{}{}
	}
	for name := range newRules {
		names[name] = struct{}{}
	}

	changedRules, cacheOperations := 0, 0
	for name := range names {
		oldRule, oldExists := oldRules[name]
		newRule, newExists := newRules[name]

		if oldExists == newExists && oldRule.hash == newRule.hash {
			continue
		}

		changedRules++
		if oldRule.usesCache || newRule.usesCache {
			cacheOperations++
		}
	}

	// NOTE: Front Door Standard/Premium allows up to 100 rules per ruleset. For batch
	// rulesets, internal service guidance states that a rule with a cache-enabled
	// `route_configuration_override_action` consumes two effective rule slots during
	// replacement. This means a batch update can fail even when the final desired
	// ruleset is within the documented 100 rule limit, because the service evaluates the effective
	// change set for the PATCH and returns a 400 Bad Request once that replacement crosses
	// that quota. We validate `changedRules + cacheOperations` here so the plan
	// fails with the same service-side constraint before sending the batch ruleset update.
	if effectiveDiff := changedRules + cacheOperations; effectiveDiff > 100 {
		return fmt.Errorf("the number of changed `rule` blocks exceeds the service-side quota: got changed `%d` rules, (`%d` of which include caching, which counts for 2 rule slots), total `%d` rule slots used, but the maximum allowed is `100`", changedRules, cacheOperations, effectiveDiff)
	}

	return nil
}

type cdnFrontDoorBatchRuleSetRule struct {
	hash      int
	usesCache bool
}

func cdnFrontDoorBatchRuleSetRules(input interface{}, hash pluginsdk.SchemaSetFunc) map[string]cdnFrontDoorBatchRuleSetRule {
	results := make(map[string]cdnFrontDoorBatchRuleSetRule)
	if input == nil {
		return results
	}

	rulesList, ok := input.([]interface{})
	if !ok {
		return results
	}

	for _, rule := range rulesList {
		r, ok := rule.(map[string]interface{})
		if !ok {
			continue
		}

		results[r["name"].(string)] = cdnFrontDoorBatchRuleSetRule{
			hash:      hash(r),
			usesCache: cdnFrontDoorBatchRuleSetRuleUsesCache(r),
		}
	}

	return results
}

func cdnFrontDoorBatchRuleSetRuleUsesCache(input map[string]interface{}) bool {
	actionsRaw, ok := input["actions"].([]interface{})
	if !ok || len(actionsRaw) == 0 || actionsRaw[0] == nil {
		return false
	}

	actions, ok := actionsRaw[0].(map[string]interface{})
	if !ok {
		return false
	}

	routeOverrides, ok := actions["route_configuration_override"].([]interface{})
	if !ok || len(routeOverrides) == 0 || routeOverrides[0] == nil {
		return false
	}

	routeOverride, ok := routeOverrides[0].(map[string]interface{})
	if !ok {
		return false
	}

	cachingList, ok := routeOverride["caching"].([]interface{})
	if !ok {
		return false
	}

	caching, ok := cachingList[0].(map[string]interface{})
	if !ok {
		return false
	}

	cacheBehaviour, ok := caching["behaviour"].(string)
	if !ok {
		return false
	}

	return cacheBehaviour != RuleCacheBehaviorDisabled
}

func (r CdnFrontDoorBatchRuleSetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient

			var model CdnFrontDoorBatchRuleSetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			profileId, err := profiles.ParseProfileID(model.CdnFrontdoorProfileID)
			if err != nil {
				return err
			}

			id := rulesets.NewRuleSetID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, model.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("retrieving %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload, err := expandCdnFrontDoorBatchRuleSetPayload(model)
			if err != nil {
				return err
			}

			if err := client.CreateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r CdnFrontDoorBatchRuleSetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient

			id, err := rulesets.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r CdnFrontDoorBatchRuleSetResource) flatten(metadata sdk.ResourceMetaData, id *rulesets.RuleSetId, model *rulesets.RuleSet) error {
	state := CdnFrontDoorBatchRuleSetModel{
		ID:                    id.ID(),
		Name:                  id.RuleSetName,
		CdnFrontdoorProfileID: profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID(),
	}

	if model != nil && model.Properties != nil {
		rulesetsState, err := flattenCdnFrontDoorBatchRuleSetRules(model.Properties.Rules)
		if err != nil {
			return fmt.Errorf("flattening `rules`: %+v", err)
		}
		state.Rule = rulesetsState
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r CdnFrontDoorBatchRuleSetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient

			id, err := rulesets.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CdnFrontDoorBatchRuleSetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload, err := expandCdnFrontDoorBatchRuleSetPayload(model)
			if err != nil {
				return err
			}

			if err := client.CreateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CdnFrontDoorBatchRuleSetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient

			id, err := rulesets.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
