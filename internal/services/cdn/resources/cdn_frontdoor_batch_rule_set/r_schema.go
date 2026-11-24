// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_batch_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Front Door Batch Rule Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorRuleSetName,
		},
		"cdn_frontdoor_profile_id": {
			Description:  "The resource ID of the Front Door Profile where this Front Door Batch Rule Set should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: profiles.ValidateProfileID,
		},
		"rule": {
			Description: "The configured blocks represent the complete set of rules managed for this Front Door Batch Rule Set. The final rule ordering is determined by each block's `order` value.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    100,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name which should be used for this Front Door Batch Rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: helpers.ValidateCdnFrontDoorRuleName,
					},
					"behaviour_on_match": {
						Description:  "The behaviour on a condition match.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(rulesets.MatchProcessingBehaviorContinue),
						ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForMatchProcessingBehavior(), false),
					},
					"order": {
						Description:  "The order in which this rule will be applied for the Front Door Endpoint. Rules with a lesser `order` value are applied before rules with a greater `order` value.",
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
												Description:  "The response type to return to the requestor.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForRedirectType(), false),
											},
											"redirect_protocol": {
												Description:  "The protocol the request is redirected as.",
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      string(rulesets.DestinationProtocolMatchRequest),
												ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForDestinationProtocol(), false),
											},
											"destination_path": {
												Description:  "The path to use in the redirect. Leave this unset to preserve the incoming path.",
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringStartsWithOneOf("/"),
											},
											"destination_host_name": {
												Description:  "The host name you want the request to be redirected to. Leave this unset to preserve the incoming host.",
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 2048),
											},
											"destination_fragment": {
												Description: "The fragment to use in the redirect. Leave this unset to preserve the incoming fragment.",
												Type:        pluginsdk.TypeString,
												Optional:    true,
												ValidateFunc: validation.All(
													validation.StringLenBetween(1, 1024),
													validation.StringDoesNotStartWithOneOf("#"),
												),
											},
											"query_string": {
												Description:  "The query string used in the redirect URL. Leave this unset to preserve the incoming query string. The maximum allowed length for this field is `2048` characters.",
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
												Description:  "The source pattern in the URL path to replace.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringStartsWithOneOf("/"),
											},
											"destination_path": {
												Description:  "The destination path to use in the rewrite.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringStartsWithOneOf("/"),
											},
											"preserve_unmatched_path_enabled": {
												Description: "Whether to append the remaining path after the source pattern to the new destination path.",
												Type:        pluginsdk.TypeBool,
												Optional:    true,
												Default:     false,
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
															Description:  "The Front Door Origin Group resource ID that the request should be routed to.",
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: afdorigingroups.ValidateOriginGroupID,
														},
														"forwarding_protocol": {
															Description:  "The forwarding protocol the request is redirected as.",
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
															Description:  "Controls how Front Door handles cache behaviour for the response.",
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringInSlice(helpers.PossibleValuesForRuleCacheBehavior(), false),
														},
														"duration": {
															Description:  "When `behaviour` is set to `OverrideAlways` or `OverrideIfOriginMissing`, this field specifies the cache duration to use and is required. The maximum. If the desired maximum cache duration is less than `1` day, specify it in the `HH:MM:SS` format, for example `23:59:59`.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: helpers.ValidateCdnFrontDoorCacheDuration,
														},
														"compression_enabled": {
															Description: "Whether compression is enabled.",
															Type:        pluginsdk.TypeBool,
															Optional:    true,
															Default:     false,
														},
														"query_string_behaviour": {
															Description:  "Controls how query strings contribute to the cache key.",
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForRuleQueryStringCachingBehavior(), false),
														},
														"query_string_parameters": {
															Description: "A list of query string parameter names. A maximum of `100` parameters may be defined.",
															Type:        pluginsdk.TypeList,
															Optional:    true,
															MaxItems:    100,
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
											Description:  "The type of remote address to match.",
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
											Description:  "A condition operator.",
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestMethodOperator()), false),
										},
										"values": {
											Description: "A list of one or more HTTP methods. If multiple values are specified, they are evaluated using `OR` logic.",
											Type:        pluginsdk.TypeSet,
											Required:    true,
											MinItems:    1,
											MaxItems:    7,
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
											Description:  "A condition operator.",
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
										"name": helpers.CdnFrontDoorRuleConditionNameSchema("A string value representing the name of the `POST` argument."),
										"operator": {
											Description:  "A condition operator.",
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
											Description:  "A condition operator.",
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
										"name": helpers.CdnFrontDoorRuleConditionNameSchema("The name of the request header."),
										"operator": {
											Description:  "A condition operator.",
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
											Description:  "A condition operator.",
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
											Description:  "A condition operator.",
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForRequestSchemeMatchConditionParametersOperator()), false),
										},
										"values": {
											Description: "The request protocol to match.",
											Type:        pluginsdk.TypeList,
											Required:    true,
											MinItems:    1,
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
											Description:  "A condition operator.",
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForURLPathOperator()), false),
										},
										"values": {
											Description: "One or more values representing the request path to match. If multiple values are specified, they are evaluated using `OR` logic.",
											Type:        pluginsdk.TypeList,
											Optional:    true,
											MaxItems:    25,
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
											Description:  "A condition operator.",
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForURLFileExtensionOperator()), false),
										},
										"values": {
											Description: "One or more values representing the request file extension to match. If multiple values are specified, they are evaluated using `OR` logic.",
											Type:        pluginsdk.TypeList,
											Optional:    true,
											MaxItems:    25,
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
											Description:  "A condition operator.",
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
											Description:  "A condition operator.",
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
										"name": helpers.CdnFrontDoorRuleConditionNameSchema("The name of the cookie."),
										"operator": {
											Description:  "A condition operator.",
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
											Description:  "A condition operator.",
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForIsDeviceOperator()), false),
										},
										"values": {
											Description: "The device type to match.",
											Type:        pluginsdk.TypeList,
											Required:    true,
											MinItems:    1,
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
											Description:  "The type of match.",
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
											Description:  "A condition operator.",
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
											Description:  "A condition operator.",
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
											Description:  "A condition operator.",
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
											Description:  "A condition operator.",
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rulesets.PossibleValuesForSslProtocolOperator()), false),
										},
										"values": {
											Description: "A list of one or more SSL protocol values.",
											Type:        pluginsdk.TypeSet,
											Required:    true,
											MinItems:    1,
											MaxItems:    3,
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
			Description: "The ID of the Front Door Batch Rule Set.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
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
