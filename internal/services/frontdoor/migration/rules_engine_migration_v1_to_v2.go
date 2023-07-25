// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RulesEngineV1ToV2 struct{}

func (s RulesEngineV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"frontdoor_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"location": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"match_condition": {
						Type:     pluginsdk.TypeList,
						MaxItems: 100,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"variable": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"selector": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"transform": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 6,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"value": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 25,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"action": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"request_header": {
									Type:     pluginsdk.TypeList,
									MaxItems: 100,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"header_action_type": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},

											"header_name": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},

											"value": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
										},
									},
								},

								"response_header": {
									Type:     pluginsdk.TypeList,
									MaxItems: 100,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"header_action_type": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},

											"header_name": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},

											"value": {
												Type:     pluginsdk.TypeString,
												Optional: true,
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
	}
}

func (s RulesEngineV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.RulesEngineIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
