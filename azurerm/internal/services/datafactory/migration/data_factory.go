package migration

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DataFactoryV0ToV1{}

type DataFactoryV0ToV1 struct{}

func (DataFactoryV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"identity": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"principal_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"github_configuration": {
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"vsts_configuration"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"git_url": {
						Type:     schema.TypeString,
						Required: true,
					},
					"repository_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},

		"vsts_configuration": {
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"github_configuration"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"project_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"repository_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     schema.TypeString,
						Required: true,
					},
					"tenant_id": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (DataFactoryV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Updating `public_network_enabled` to %q", datafactory.PublicNetworkAccessEnabled)

		rawState["public_network_enabled"] = true

		return rawState, nil
	}
}
