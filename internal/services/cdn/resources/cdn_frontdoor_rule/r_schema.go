// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (c Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: helpers.ValidateCdnFrontDoorRuleName,
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
									ValidateFunc: helpers.ValidateCdnFrontDoorUrlRedirectActionQueryString,
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues([]string{string(rules.RemoteAddressOperatorGeoMatch), string(rules.RemoteAddressOperatorIPMatch)}), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestMethodOperator()), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForQueryStringOperator()), false),
								},
								"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForPostArgsOperator()), false),
								},
								"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestUriOperator()), false),
								},
								"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
								"transforms": cdnFrontDoorRuleTransformsSchema(),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestHeaderOperator()), false),
								},
								"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestBodyOperator()), false),
								},
								"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForRequestSchemeMatchConditionParametersOperator()), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForURLPathOperator()), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForURLFileExtensionOperator()), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForURLFileNameOperator()), false),
								},
								// making optional for issue #23504
								"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForHTTPVersionOperator()), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForCookiesOperator()), false),
								},
								"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForIsDeviceOperator()), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues([]string{string(rules.SocketAddrOperatorIPMatch)}), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForClientPortOperator()), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForQueryStringOperator()), false),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForHostNameOperator()), false),
								},
								"values":     helpers.CdnFrontDoorRuleValuesOptionalSchema(),
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
									ValidateFunc: validation.StringInSlice(helpers.CdnFrontDoorRuleConditionOperatorPossibleValues(rules.PossibleValuesForSslProtocolOperator()), false),
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

func (c Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cdn_frontdoor_rule_set_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
