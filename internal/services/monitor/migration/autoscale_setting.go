package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = AutoscaleSettingUpgradeV0ToV1{}

type AutoscaleSettingUpgradeV0ToV1 struct{}

func (AutoscaleSettingUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return autoscaleSettingSchemaForV0AndV1()
}

func (AutoscaleSettingUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/autoscalesettings/{settingName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/autoscaleSettings/{settingName}
		oldId, err := azure.ParseAzureResourceID(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		settingName := ""
		for key, value := range oldId.Path {
			if strings.EqualFold(key, "autoscaleSettings") {
				settingName = value
				break
			}
		}

		if settingName == "" {
			return rawState, fmt.Errorf("couldn't find the `autoscaleSettings` segment in the old resource id %q", oldId)
		}

		newId := parse.NewAutoscaleSettingID(oldId.SubscriptionID, oldId.ResourceGroup, settingName)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId.ID())

		rawState["id"] = newId.ID()

		return rawState, nil
	}
}

func autoscaleSettingSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

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

		"tags": tags.Schema(),
	}
}
