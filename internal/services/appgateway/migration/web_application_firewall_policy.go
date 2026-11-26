// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0
package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WebApplicationFirewallPolicyV0ToV1{}

type WebApplicationFirewallPolicyV0ToV1 struct{}

func (WebApplicationFirewallPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_rules": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"action": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"match_conditions": {
					Type:     pluginsdk.TypeList,
					Required: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"match_values": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"match_variables": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"selector": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"variable_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							}},
						},
						"negation_condition": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"operator": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"transforms": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					}},
				},
				"name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"priority": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"rule_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			}},
		},
		"http_listener_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"managed_rules": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"exclusion": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"excluded_rule_set": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"rule_group": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
										"excluded_rules": {
											Type:     pluginsdk.TypeList,
											Optional: true,
											Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										},
										"rule_group_name": {
											Type:     pluginsdk.TypeString,
											Required: true,
										},
									}},
								},
								"type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							}},
						},
						"match_variable": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"selector": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"selector_match_operator": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					}},
				},
				"managed_rule_set": {
					Type:     pluginsdk.TypeList,
					Required: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"rule_group_override": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"disabled_rules": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Computed: true,
									Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
								},
								"rule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Computed: true,
									Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
										"action": {
											Type:     pluginsdk.TypeString,
											Optional: true,
										},
										"enabled": {
											Type:     pluginsdk.TypeBool,
											Optional: true,
										},
										"id": {
											Type:     pluginsdk.TypeString,
											Required: true,
										},
									}},
								},
								"rule_group_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							}},
						},
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"version": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					}},
				},
			}},
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"path_based_rule_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"policy_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"file_upload_limit_in_mb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
				"max_request_body_size_in_kb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
				"mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"request_body_check": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			}},
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
	}
}

func (WebApplicationFirewallPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		log.Printf("[Debug] start upgrade web application firewall policy id")
		if newID, err := normalizeWebAppFirewallPolicyID(rawState["id"].(string)); err != nil {
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
