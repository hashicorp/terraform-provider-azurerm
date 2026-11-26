// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_application_firewall_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func webApplicationFirewallPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"custom_rules": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(webapplicationfirewallpolicies.WebApplicationFirewallActionAllow),
							string(webapplicationfirewallpolicies.WebApplicationFirewallActionBlock),
							string(webapplicationfirewallpolicies.WebApplicationFirewallActionLog),
							string(webapplicationfirewallpolicies.ActionTypeJSChallenge),
						}, false),
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"match_conditions": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"match_values": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"match_variables": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"variable_name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForWebApplicationFirewallMatchVariable(), false),
											},
											"selector": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
										},
									},
								},
								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForWebApplicationFirewallOperator(), false),
								},
								"negation_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},
								"transforms": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForWebApplicationFirewallTransform(), false),
									},
								},
							},
						},
					},
					"priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"rule_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForWebApplicationFirewallRuleType(), false),
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"rate_limit_duration": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForApplicationGatewayFirewallRateLimitDuration(), false),
					},
					"rate_limit_threshold": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"group_rate_limit_by": {
						// group variables combination not supported yet, use a single variable name
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForApplicationGatewayFirewallUserSessionVariable(), false),
					},
				},
			},
		},

		"managed_rules": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"exception": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MinItems: 1,
						MaxItems: 60, // https://learn.microsoft.com/azure/web-application-firewall/ag/application-gateway-exceptions#limitations
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"exception_rule_set": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"type": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validate.ValidateWebApplicationFirewallPolicyExclusionRuleSetType,
											},
											"version": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validate.ValidateWebApplicationFirewallPolicyExclusionRuleSetVersion,
											},
											"rule_group": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"rule_group_name": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleGroupName,
														},
														"rules": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															Elem: &pluginsdk.Schema{
																Type:         pluginsdk.TypeString,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
											},
										},
									},
								},
								"match_variable": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForExceptionEntryMatchVariable(), false),
								},
								"value_match_operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForExceptionEntryValueMatchOperator(), false),
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									MaxItems: 600, // https://learn.microsoft.com/azure/web-application-firewall/ag/application-gateway-exceptions#limitations
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringLenBetween(1, 256),
									},
								},
								"selector": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"selector_match_operator": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForExceptionEntrySelectorMatchOperator(), false),
								},
							},
						},
					},
					"exclusion": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"match_variable": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForOwaspCrsExclusionEntryMatchVariable(), false),
								},
								"selector": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.NoZeroValues,
								},
								"selector_match_operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForOwaspCrsExclusionEntrySelectorMatchOperator(), false),
								},
								"excluded_rule_set": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"type": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      "OWASP",
												ValidateFunc: validate.ValidateWebApplicationFirewallPolicyExclusionRuleSetType,
											},
											"version": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      "3.2",
												ValidateFunc: validate.ValidateWebApplicationFirewallPolicyExclusionRuleSetVersion,
											},
											"rule_group": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"rule_group_name": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleGroupName,
														},
														"excluded_rules": {
															Type:     pluginsdk.TypeList,
															Optional: true,
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
					"managed_rule_set": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "OWASP",
									ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleSetType,
								},
								"version": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleSetVersion,
								},
								"rule_group_override": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"rule_group_name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleGroupName,
											},
											"rule": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"id": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"enabled": {
															Type:     pluginsdk.TypeBool,
															Optional: true,
															Default:  false,
														},

														"action": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
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
				},
			},
		},

		"policy_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(webapplicationfirewallpolicies.WebApplicationFirewallModePrevention),
						ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForWebApplicationFirewallMode(), false),
					},

					"request_body_check": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"file_upload_limit_in_mb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      100,
						ValidateFunc: validation.IntBetween(1, 4000),
					},

					"request_body_enforcement": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"file_upload_enforcement": {
						Type: pluginsdk.TypeBool,
						/*
							NOTE: O+C: This value defaults to true but is only available under certain conditions (i.e. when version is 3.2)
								managed_rules {
									managed_rule_set {
									 type  = "OWASP"
									 version = "3.2"
									}
								 }
						*/
						Optional: true,
					},

					"max_request_body_size_in_kb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      128,
						ValidateFunc: validation.IntBetween(8, 2000),
					},

					"request_body_inspect_limit_in_kb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      128,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"js_challenge_cookie_expiration_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      30,
						ValidateFunc: validation.IntBetween(5, 1440),
					},

					"log_scrubbing": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},

								"rule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
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
												ValidateFunc: validation.StringInSlice(
													webapplicationfirewallpolicies.PossibleValuesForScrubbingRuleEntryMatchVariable(),
													false,
												),
											},

											"selector_match_operator": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												Default:  "Equals",
												ValidateFunc: validation.StringInSlice(
													webapplicationfirewallpolicies.PossibleValuesForScrubbingRuleEntryMatchOperator(),
													false,
												),
											},

											"selector": {
												Description: "When matchVariable is a collection, operator used to specify which elements in the collection this rule applies to.",
												Type:        pluginsdk.TypeString,
												Optional:    true,
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

		"http_listener_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"path_based_rule_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"tags": commonschema.Tags(),
	}
}
