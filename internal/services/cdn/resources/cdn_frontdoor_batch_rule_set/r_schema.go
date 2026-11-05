// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_batch_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
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
						ValidateFunc: helpers.ValidateCdnFrontDoorRuleName,
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
												ValidateFunc: helpers.ValidateCdnFrontDoorUrlRedirectActionQueryString,
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
															ValidateFunc: validation.StringInSlice(helpers.PossibleValuesForRuleCacheBehavior(), false),
														},
														"duration": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: helpers.ValidateCdnFrontDoorCacheDuration,
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
						Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
							"remote_address": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"operator": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues([]string{string(rulesets.RemoteAddressOperatorGeoMatch), string(rulesets.RemoteAddressOperatorIPMatch)}), false),
										},
										"values": helpers.CdnFrontDoorRuleValuesRequiredSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestMethodOperator()), false),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForQueryStringOperator()), false),
										},
										"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
									},
								},
							},
							"post_argument": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"name": helpers.CdnFrontDoorRuleConditionNameSchema(),
										"operator": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForPostArgsOperator()), false),
										},
										"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestUriOperator()), false),
										},
										"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
										"transforms": cdnFrontDoorBatchRuleSetConditionTransformsSchema(),
									},
								},
							},
							"request_header": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"name": helpers.CdnFrontDoorRuleConditionNameSchema(),
										"operator": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestHeaderOperator()), false),
										},
										"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestBodyOperator()), false),
										},
										"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestSchemeMatchConditionParametersOperator()), false),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForURLPathOperator()), false),
										},
										"values": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											MaxItems: 25,
											Elem: &pluginsdk.Schema{
												// Azure accepts `/` as the root path, while non-root values must omit the leading slash.
												Type: pluginsdk.TypeString,
												ValidateFunc: validation.All(
													validation.StringIsNotEmpty,
													validation.Any(
														validation.StringDoesNotStartWithOneOf("/"),
														validation.StringInSlice([]string{"/"}, false),
													),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForURLFileExtensionOperator()), false),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForURLFileNameOperator()), false),
										},
										"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForHTTPVersionOperator()), false),
										},
										"values": helpers.CdnFrontDoorRuleHttpVersionValuesSchema(),
									},
								},
							},
							"request_cookies": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"name": helpers.CdnFrontDoorRuleConditionNameSchema(),
										"operator": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForCookiesOperator()), false),
										},
										"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForIsDeviceOperator()), false),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues([]string{string(rulesets.SocketAddrOperatorIPMatch)}), false),
										},
										"values": helpers.CdnFrontDoorSocketAddressValuesSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForClientPortOperator()), false),
										},
										"values": helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForQueryStringOperator()), false),
										},
										"values": helpers.CdnFrontDoorServerPortValuesSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForHostNameOperator()), false),
										},
										"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForSslProtocolOperator()), false),
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

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
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
