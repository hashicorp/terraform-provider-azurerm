// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ConfigurationV0ToV1{}

type ConfigurationV0ToV1 struct {
}

func (ConfigurationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"antimalware": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"real_time_protection_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"scheduled_scan_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"scheduled_scan_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"scheduled_scan_day": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"scheduled_scan_time_in_minutes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"exclusions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"extensions": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"paths": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"processes": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
		"azure_security_baseline": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"assignment_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
		"backup": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"policy_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"time_zone": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"instant_rp_retention_range_in_days": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"schedule_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"schedule_run_frequency": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"schedule_run_times": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"schedule_run_days": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"schedule_policy_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					"retention_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"retention_policy_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"daily_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"retention_times": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},
														"duration_type": {
															Type:     pluginsdk.TypeString,
															Optional: true,
														},
													},
												},
											},
										},
									},
								},
								"weekly_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"retention_times": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},
														"duration_type": {
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
				},
			},
		},
		"automation_account_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"boot_diagnostics_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"defender_for_cloud_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"guest_configuration_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"log_analytics_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"status_change_alert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
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

func (c ConfigurationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		parsed, err := configurationprofiles.ParseConfigurationProfileIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := parsed.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
