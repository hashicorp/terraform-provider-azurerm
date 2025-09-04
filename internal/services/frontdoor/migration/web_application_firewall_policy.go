// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WebApplicationFirewallPolicyV0ToV1{}

type WebApplicationFirewallPolicyV0ToV1 struct{}

func (WebApplicationFirewallPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"redirect_url": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"custom_block_response_status_code": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"custom_block_response_body": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"custom_rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"priority": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"rate_limit_duration_in_minutes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"rate_limit_threshold": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"match_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 100,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"match_variable": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"match_values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"selector": {
									Type:     pluginsdk.TypeString,
									Optional: true,
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
										Type: pluginsdk.TypeString,
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
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"exclusion": {
						Type:     pluginsdk.TypeList,
						MaxItems: 100,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"match_variable": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"selector": {
									Type:     pluginsdk.TypeString,
									Required: true,
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
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"exclusion": {
									Type:     pluginsdk.TypeList,
									MaxItems: 100,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"match_variable": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"operator": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"selector": {
												Type:     pluginsdk.TypeString,
												Required: true,
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
												Type:     pluginsdk.TypeString,
												Required: true,
											},

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},

											"exclusion": {
												Type:     pluginsdk.TypeList,
												MaxItems: 100,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"match_variable": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},
														"operator": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},
														"selector": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},
													},
												},
											},

											"action": {
												Type:     pluginsdk.TypeString,
												Required: true,
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

		"tags": tags.Schema(),
	}
}

func (WebApplicationFirewallPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontdoorwebapplicationfirewallpolicies/{policyName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/{policyName}
		oldId := rawState["id"].(string)
		id, err := parse.WebApplicationFirewallPolicyIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
