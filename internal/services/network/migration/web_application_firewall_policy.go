// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WebApplicationFirewallPolicyV0ToV1{}

type WebApplicationFirewallPolicyV0ToV1 struct{}

func (WebApplicationFirewallPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_rules": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"action": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"match_conditions": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"match_values": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
						"match_variables": {
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"selector": {
									Optional: true,
									Type:     pluginsdk.TypeString,
								},
								"variable_name": {
									Required: true,
									Type:     pluginsdk.TypeString,
								},
							}},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"negation_condition": {
							Optional: true,
							Type:     pluginsdk.TypeBool,
						},
						"operator": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"transforms": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Optional: true,
							Type:     pluginsdk.TypeSet,
						},
					}},
					Required: true,
					Type:     pluginsdk.TypeList,
				},
				"name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"priority": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
				"rule_type": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"http_listener_ids": {
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Type:     pluginsdk.TypeList,
		},
		"location": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"managed_rules": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"exclusion": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"excluded_rule_set": {
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"rule_group": {
									Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
										"excluded_rules": {
											Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
											Optional: true,
											Type:     pluginsdk.TypeList,
										},
										"rule_group_name": {
											Required: true,
											Type:     pluginsdk.TypeString,
										},
									}},
									Optional: true,
									Type:     pluginsdk.TypeList,
								},
								"type": {
									Optional: true,
									Type:     pluginsdk.TypeString,
								},
								"version": {
									Optional: true,
									Type:     pluginsdk.TypeString,
								},
							}},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
						"match_variable": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"selector": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"selector_match_operator": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"managed_rule_set": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"rule_group_override": {
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"disabled_rules": {
									Computed: true,
									Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
									Optional: true,
									Type:     pluginsdk.TypeList,
								},
								"rule": {
									Computed: true,
									Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
										"action": {
											Optional: true,
											Type:     pluginsdk.TypeString,
										},
										"enabled": {
											Optional: true,
											Type:     pluginsdk.TypeBool,
										},
										"id": {
											Required: true,
											Type:     pluginsdk.TypeString,
										},
									}},
									Optional: true,
									Type:     pluginsdk.TypeList,
								},
								"rule_group_name": {
									Required: true,
									Type:     pluginsdk.TypeString,
								},
							}},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
						"type": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"version": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Required: true,
					Type:     pluginsdk.TypeList,
				},
			}},
			Required: true,
			Type:     pluginsdk.TypeList,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"path_based_rule_ids": {
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Type:     pluginsdk.TypeList,
		},
		"policy_settings": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"file_upload_limit_in_mb": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"max_request_body_size_in_kb": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"mode": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"request_body_check": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"resource_group_name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"tags": {
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Optional: true,
			Type:     pluginsdk.TypeMap,
		},
	}
}

func (WebApplicationFirewallPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[Debug] start upgrade web application firewall policy id")
		oldID := rawState["id"].(string)
		if newID, err := normalizeWebAppFirewallPolicyID(oldID); err != nil {
			return nil, err
		} else if newID != nil {
			rawState["id"] = *newID
		}
		return rawState, nil
	}
}

func normalizeWebAppFirewallPolicyID(id string) (*string, error) {
	if id == "" {
		return nil, nil
	}
	parseID, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(id)
	if err != nil {
		return nil, fmt.Errorf("prase id: %v", err)
	}
	normalizedID := parseID.ID()
	return &normalizedID, nil
}
