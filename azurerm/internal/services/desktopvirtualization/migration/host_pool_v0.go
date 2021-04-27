package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = HostPoolV0ToV1{}

type HostPoolV0ToV1 struct{}

func (HostPoolV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": azure.SchemaLocation(),

		"resource_group_name": azure.SchemaResourceGroupName(),

		"type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"load_balancer_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"friendly_name": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"validate_environment": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"personal_desktop_assignment_type": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"maximum_sessions_allowed": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  999999,
		},

		"preferred_app_group_type": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Preferred App Group type to display",
		},

		"registration_info": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"expiration_date": {
						Type:     schema.TypeString,
						Required: true,
					},

					"reset_token": {
						Type:     schema.TypeBool,
						Computed: true,
					},

					"token": {
						Type:      schema.TypeString,
						Sensitive: true,
						Computed:  true,
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}

func (HostPoolV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		id, err := parse.HostPoolIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := id.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
