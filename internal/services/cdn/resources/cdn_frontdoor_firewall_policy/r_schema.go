// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorFirewallPolicyName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapplicationfirewallpolicies.SkuNameStandardAzureFrontDoor),
				string(webapplicationfirewallpolicies.SkuNamePremiumAzureFrontDoor),
			}, false),
		},

		"mode": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForPolicyMode(), false),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"js_challenge_cookie_expiration_in_minutes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// Note: O+C because 'js challenge expiration' is always
			// enabled no matter what and cannot be disabled for Premium_AzureFrontDoor
			// and is not supported in Standard_AzureFrontDoor...
			Computed:     true,
			ValidateFunc: validation.IntBetween(5, 1440),
		},

		"captcha_cookie_expiration_in_minutes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// Note: O+C because Azure automatically enables CAPTCHA policy with a default of 30 minutes on the Premium_AzureFrontDoor SKU
			// but it cannot be defined for Standard_AzureFrontDoor
			Computed:     true,
			ValidateFunc: validation.IntBetween(5, 1440),
		},

		"redirect_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
		},

		"request_body_check_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"custom_block_response_status_code": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"priority": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  1,
					},

					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForRuleType(), false),
					},

					"rate_limit_duration_in_minutes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  1,
					},

					"rate_limit_threshold": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  10,
					},

					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
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
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 10,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"match_variable": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForMatchVariable(), false),
								},

								"match_values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 600,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringLenBetween(1, 256),
									},
								},

								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForOperator(), false),
								},

								"selector": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"negation_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 5,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionMatchVariable(), false),
								},
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionSelectorMatchOperator(), false),
								},
								"selector": {
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
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionMatchVariable(), false),
											},
											"operator": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionSelectorMatchOperator(), false),
											},
											"selector": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"rule": {
									Type:     pluginsdk.TypeList,
									MaxItems: 1000,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"rule_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  false,
											},

											"exclusion": {
												Type:     pluginsdk.TypeList,
												MaxItems: 100,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"match_variable": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionMatchVariable(), false),
														},
														"operator": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForManagedRuleExclusionSelectorMatchOperator(), false),
														},
														"selector": {
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
												Type:     pluginsdk.TypeString,
												Required: true,
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
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"scrubbing_rule": {
						Type:     pluginsdk.TypeList,
						MaxItems: 100,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},

								"match_variable": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForScrubbingRuleEntryMatchVariable(),
										false),
								},

								"operator": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperatorEquals),
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForScrubbingRuleEntryMatchOperator(),
										false),
								},

								"selector": {
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
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
