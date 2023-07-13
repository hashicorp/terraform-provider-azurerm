// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutoscaleSettingUpgradeV1ToV2 struct{}

func (s AutoscaleSettingUpgradeV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"target_resource_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 20,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"capacity": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"minimum": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
								"maximum": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
								"default": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
							},
						},
					},
					"rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 10,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"metric_trigger": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"metric_name": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"metric_resource_id": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"time_grain": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"statistic": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"time_window": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"time_aggregation": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"operator": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"threshold": {
												Type:     pluginsdk.TypeFloat,
												Required: true,
											},

											"metric_namespace": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},

											"divide_by_instance_count": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},

											"dimensions": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"name": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},

														"operator": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},

														"values": {
															Type:     pluginsdk.TypeList,
															Required: true,
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
								"scale_action": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"direction": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"type": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"value": {
												Type:     pluginsdk.TypeInt,
												Required: true,
											},
											"cooldown": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
										},
									},
								},
							},
						},
					},
					"fixed_date": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"timezone": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"start": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"end": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"recurrence": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"timezone": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"days": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"hours": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeInt,
									},
								},
								"minutes": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeInt,
									},
								},
							},
						},
					},
				},
			},
		},

		"notification": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"send_to_subscription_administrator": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},
								"send_to_subscription_co_administrator": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},
								"custom_emails": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"webhook": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_uri": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"properties": {
									Type:     pluginsdk.TypeMap,
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

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (s AutoscaleSettingUpgradeV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := autoscalesettings.ParseAutoScaleSettingIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
