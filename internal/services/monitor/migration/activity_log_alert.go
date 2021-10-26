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

var _ pluginsdk.StateUpgrade = ActivityLogAlertUpgradeV0ToV1{}

type ActivityLogAlertUpgradeV0ToV1 struct{}

func (ActivityLogAlertUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return activityLogAlertSchemaForV0AndV1()
}

func (ActivityLogAlertUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/activityLogAlerts/alert1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/activityLogAlerts/alert1
		oldId, err := azure.ParseAzureResourceID(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		alertName := ""
		for key, value := range oldId.Path {
			if strings.EqualFold(key, "activityLogAlerts") {
				alertName = value
				break
			}
		}

		if alertName == "" {
			return rawState, fmt.Errorf("couldn't find the `activityLogAlerts` segment in the old resource id %q", oldId)
		}

		newId := parse.NewActivityLogAlertID(oldId.SubscriptionID, oldId.ResourceGroup, alertName)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId.ID())

		rawState["id"] = newId.ID()

		return rawState, nil
	}
}

func activityLogAlertSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"scopes": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"criteria": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"operation_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"caller": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"level": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"resource_provider": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"resource_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"resource_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"resource_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"status": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"sub_status": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"recommendation_category": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"recommendation_impact": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"recommendation_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					//lintignore:XS003
					"service_health": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"events": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"locations": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"services": {
									Type:     pluginsdk.TypeSet,
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

		"action": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action_group_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"webhook_properties": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": tags.Schema(),
	}
}
