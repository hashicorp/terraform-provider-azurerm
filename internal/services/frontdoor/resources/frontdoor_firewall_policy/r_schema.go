// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_firewall_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func frontdoorFirewallPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorWAFName,
		},

		"location": commonschema.LocationComputed(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"enabled": {
			Description: "Whether the policy is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"mode": {
			Description:  "The Firewall Policy mode.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForPolicyMode(), false),
			Default:      string(webapplicationfirewallpolicies.PolicyModePrevention),
		},

		"redirect_url": {
			Description:  "If action type is redirect, this field represents redirect URL for the client.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
		},

		"custom_block_response_status_code": {
			Description: "If a `custom_rule` block's action type is `block`, this is the response status code.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			ValidateFunc: validation.IntInSlice([]int{
				200,
				403,
				405,
				406,
				429,
			}),
		},

		"custom_block_response_body": {
			Description:  "If a `custom_rule` block's action type is `block`, this is the response body. The body must be specified in base64 encoding.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.CustomBlockResponseBody,
		},

		"custom_rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the resource that is unique within a policy. This name can be used to access the resource.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"enabled": {
						Description: "Whether the rule is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"priority": {
						Description: "The priority of the rule. Rules with a lower value will be evaluated before rules with a higher value.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     1,
					},

					"type": {
						Description:  "The type of rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForRuleType(), false),
					},

					"rate_limit_duration_in_minutes": {
						Description: "The rate limit duration in minutes.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     1,
					},

					"rate_limit_threshold": {
						Description: "The rate limit threshold.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     10,
					},

					"action": {
						Description:  "The action to perform when the rule is matched.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForActionType(), false),
					},

					"match_condition": {
						Description: "The conditions a request must match for this custom rule to apply.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    10,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"match_variable": {
									Description:  "The request variable to compare with.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForMatchVariable(), false),
								},

								"match_values": {
									Description: "Up to `600` possible values to match. Limit is in total across all `match_condition` blocks and `match_values` arguments.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MaxItems:    600,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringLenBetween(1, 256),
									},
								},

								"operator": {
									Description:  "Comparison type to use for matching with the variable value.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForOperator(), false),
								},

								"selector": {
									Description:  "Match against a specific key if the `match_variable` is `QueryString`, `PostArgs`, `RequestHeader` or `Cookies`.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"negation_condition": {
									Description: "Whether the result of the condition should be negated.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},

								"transforms": {
									Description: "Up to `5` transforms to apply.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									MaxItems:    5,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForTransformType(), false),
									},
								},
							},
						},
					},
				},
			},
		},

		"managed_rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Description:  "The name of the managed rule to use with this resource.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Description:  "The version on the managed rule to use with this resource.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"exclusion": {
						Type:     pluginsdk.TypeList,
						MaxItems: 100,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"match_variable": {
									Description:  "The variable type to be excluded.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionMatchVariable(), false),
								},
								"operator": {
									Description:  "Comparison operator to apply to the selector when specifying which elements in the collection this exclusion applies to.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionSelectorMatchOperator(), false),
								},
								"selector": {
									Description:  "Selector for the value in the `match_variable` attribute this exclusion applies to.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"override": {
						Type:     pluginsdk.TypeList,
						MaxItems: 100,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule_group_name": {
									Description:  "The managed rule group to override.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"exclusion": {
									Type:     pluginsdk.TypeList,
									MaxItems: 100,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"match_variable": {
												Description:  "The variable type to be excluded.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionMatchVariable(), false),
											},
											"operator": {
												Description:  "Comparison operator to apply to the selector when specifying which elements in the collection this exclusion applies to.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionSelectorMatchOperator(), false),
											},
											"selector": {
												Description:  "Selector for the value in the `match_variable` attribute this exclusion applies to.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"rule": {
									Description: "If none are specified, all of the rules in the group will be disabled.",
									Type:        pluginsdk.TypeList,
									MaxItems:    1000,
									Optional:    true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"rule_id": {
												Description:  "Identifier for the managed rule.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"enabled": {
												Description: "Whether the managed rule override is enabled.",
												Type:        pluginsdk.TypeBool,
												Optional:    true,
												Default:     false,
											},

											"exclusion": {
												Type:     pluginsdk.TypeList,
												MaxItems: 100,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"match_variable": {
															Description:  "The variable type to be excluded.",
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionMatchVariable(), false),
														},
														"operator": {
															Description:  "Comparison operator to apply to the selector when specifying which elements in the collection this exclusion applies to.",
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionSelectorMatchOperator(), false),
														},
														"selector": {
															Description:  "Selector for the value in the `match_variable` attribute this exclusion applies to.",
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"action": {
												Description:  "The action to be applied when the rule matches.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForActionType(), false),
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

		"frontend_endpoint_ids": {
			Description: "The Frontend Endpoints associated with this Front Door Web Application Firewall Policy.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
