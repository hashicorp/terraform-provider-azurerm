package migration

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DataFactoryV0ToV1{}
var _ pluginsdk.StateUpgrade = DataFactoryV1ToV2{}

type DataFactoryV0ToV1 struct{}

func (DataFactoryV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"github_configuration": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"vsts_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"git_url": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"repository_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"vsts_configuration": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"github_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"project_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"repository_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
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

func (DataFactoryV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Updating `public_network_enabled` to %q", datafactory.PublicNetworkAccessEnabled)

		rawState["public_network_enabled"] = true

		return rawState, nil
	}
}

type DataFactoryV1ToV2 struct{}

func (DataFactoryV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"identity_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"github_configuration": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"vsts_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"git_url": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"repository_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"vsts_configuration": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"github_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"project_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"repository_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},
		"public_network_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"customer_managed_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			RequiredWith: []string{"identity.0.identity_ids"},
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

func (DataFactoryV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Updating `id` if resourceName is in upper case")

		oldId := rawState["id"].(string)
		id, err := parse.DataFactoryID(oldId)
		if err != nil {
			return nil, err
		}

		id.FactoryName = rawState["name"].(string)
		rawState["id"] = id.ID()

		return rawState, nil
	}
}
