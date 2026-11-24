// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_firewall_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnFrontdoorFirewallPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorFirewallPolicyName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Description: "The SKU's pricing tier for this Front Door Firewall Policy.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapplicationfirewallpolicies.SkuNameStandardAzureFrontDoor),
				string(webapplicationfirewallpolicies.SkuNamePremiumAzureFrontDoor),
			}, false),
		},

		"mode": {
			Description:  "The Front Door Firewall Policy mode.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForPolicyMode(), false),
		},

		"enabled": {
			Description: "Is the Front Door Firewall Policy enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"js_challenge_cookie_expiration_in_minutes": {
			Description: "The JavaScript challenge cookie lifetime in minutes, after which the user will be revalidated.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because 'js challenge expiration' is always
			// enabled no matter what and cannot be disabled for Premium_AzureFrontDoor
			// and is not supported in Standard_AzureFrontDoor...
			Computed:     true,
			ValidateFunc: validation.IntBetween(5, 1440),
		},

		"captcha_cookie_expiration_in_minutes": {
			Description: "The Captcha cookie lifetime in minutes.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because Azure automatically enables CAPTCHA policy with a default of 30 minutes on the Premium_AzureFrontDoor SKU
			// but it cannot be defined for Standard_AzureFrontDoor
			Computed:     true,
			ValidateFunc: validation.IntBetween(5, 1440),
		},

		"redirect_url": {
			Description:  "If action type is redirect, this field represents redirect URL for the client.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
		},

		"request_body_check_enabled": {
			Description: "Should policy managed rules inspect the request body content?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
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
				990,
				991,
				992,
				993,
				994,
				995,
				996,
				997,
				998,
				999,
			}),
		},

		"custom_block_response_body": {
			Description:  "If a `custom_rule` block's action type is `block`, this is the response body. The body must be specified in base64 encoding.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"custom_rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "Gets name of the resource that is unique within a policy. This name can be used to access the resource.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"enabled": {
						Description: "Is the rule is enabled or disabled?",
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
						Description: "The action to perform when the rule is matched.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(webapplicationfirewallpolicies.ActionTypeAllow),
							string(webapplicationfirewallpolicies.ActionTypeBlock),
							string(webapplicationfirewallpolicies.ActionTypeLog),
							string(webapplicationfirewallpolicies.ActionTypeRedirect),
							string(webapplicationfirewallpolicies.ActionTypeJSChallenge),
							string(webapplicationfirewallpolicies.ActionTypeCAPTCHA),
						}, false),
					},

					"match_condition": {
						Description: "Can support up to `10` `match_condition` blocks?",
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
									Description:  "Match against a specific key if the `match_variable` is `QueryString`, `PostArgs`, `RequestHeader`, or `Cookies`.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"negation_condition": {
									Description: "Should the result of the condition be negated?",
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
						Description:  "The version of the managed rule to use with this resource. For `Microsoft_DefaultRuleSet` the. For `BotProtection` the value must be `preview-0.1` and for `Microsoft_BotManagerRuleSet` the.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"action": {
						Description: "The action to perform for all default rule set rules when the managed rule is matched or when the anomaly score is 5 or greater depending on which version of the default rule set you are using.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(webapplicationfirewallpolicies.ActionTypeAllow),
							string(webapplicationfirewallpolicies.ActionTypeLog),
							string(webapplicationfirewallpolicies.ActionTypeBlock),
							string(webapplicationfirewallpolicies.ActionTypeRedirect),
						}, false),
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
												Description: "Is the managed rule override enabled or disabled?",
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

											// NOTE: 'ActionTypeAnomalyScoring' is only valid with 2.0 and above
											//       'ActionTypeJSChallenge' is only valid with BotManagerRuleSets
											"action": {
												Description: "The action to be applied when the managed rule matches or when the anomaly score is 5 or greater.",
												Type:        pluginsdk.TypeString,
												Required:    true,
												ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForActionType(),
													false),
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

		"log_scrubbing": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Description: "Is log scrubbing enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"scrubbing_rule": {
						Type:     pluginsdk.TypeList,
						MaxItems: 100,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Description: "Is this `scrubbing_rule` enabled?",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     true,
								},

								"match_variable": {
									Description: "The variable to be scrubbed from the logs.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForScrubbingRuleEntryMatchVariable(),
										false),
								},

								"operator": {
									Description: "When the `match_variable` is a collection, operate on the `selector` to specify which elements in the collection this `scrubbing_rule` applies to.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									Default:     string(webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperatorEquals),
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForScrubbingRuleEntryMatchOperator(),
										false),
								},

								"selector": {
									Description:  "When the `match_variable` is a collection, the `operator` is used to specify which elements in the collection this `scrubbing_rule` applies to.",
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

		"frontend_endpoint_ids": {
			Description: "The Front Door Profiles frontend endpoints associated with this Front Door Firewall Policy.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
