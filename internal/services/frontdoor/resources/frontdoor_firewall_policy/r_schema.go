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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorWAFName,
		},

		"location": commonschema.LocationComputed(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForPolicyMode(), false),
			Default:      string(webapplicationfirewallpolicies.PolicyModePrevention),
		},

		"redirect_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
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
			}),
		},

		"custom_block_response_body": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForActionType(), false),
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

											"action": {
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
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
