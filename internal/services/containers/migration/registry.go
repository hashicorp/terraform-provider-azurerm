// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ pluginsdk.StateUpgrade = RegistryV0ToV1{}
	_ pluginsdk.StateUpgrade = RegistryV1ToV2{}
	_ pluginsdk.StateUpgrade = RegistryV2ToV3{}
)

type RegistryV0ToV1 struct{}

func (RegistryV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return registrySchemaForV0AndV1()
}

func (RegistryV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		rawState["sku"] = "Basic"
		return rawState, nil
	}
}

type RegistryV1ToV2 struct{}

func (RegistryV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return registrySchemaForV0AndV1()
}

func (RegistryV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// Basic's been renamed Classic to allow for "ManagedBasic" ¯\_(ツ)_/¯
		rawState["sku"] = "Classic"

		storageAccountId := ""
		if v, ok := rawState["storage_account"]; ok {
			subscriptionId := meta.(*clients.Client).Account.SubscriptionId
			ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, time.Minute*5)
			defer cancel()

			raw := v.(*pluginsdk.Set).List()
			rawVals := raw[0].(map[string]interface{})
			storageAccountName := rawVals["name"].(string)

			account, err := meta.(*clients.Client).Storage.FindAccount(ctx, subscriptionId, storageAccountName)
			if err != nil {
				return nil, fmt.Errorf("finding Storage Account %q: %+v", storageAccountName, err)
			}

			storageAccountId = account.StorageAccountId.ID()
		}

		if storageAccountId == "" {
			return rawState, fmt.Errorf("unable to determine storage account ID")
		}

		return rawState, nil
	}
}

func registrySchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"admin_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// lintignore:S018
		"storage_account": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"access_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"login_server": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_password": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

type RegistryV2ToV3 struct{}

func (r RegistryV2ToV3) Schema() map[string]*pluginsdk.Schema {
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

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"admin_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"georeplications": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAuto,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"zone_redundancy_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"regional_endpoint_enabled": {
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
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"login_server": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_password": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"identity_ids": {
						Type:     pluginsdk.TypeSet,
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

		"encryption": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity_client_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"key_vault_key_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"network_rule_set": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"ip_rule": {
						Type:       pluginsdk.TypeSet,
						Optional:   true,
						Computed:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"ip_range": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"quarantine_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"retention_policy_in_days": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"trust_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"export_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"zone_redundancy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"anonymous_pull_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"data_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"data_endpoint_host_names": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"network_rule_bypass_option": {
			Type:     pluginsdk.TypeString,
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

func (r RegistryV2ToV3) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		delete(rawState, "zone_redundancy_enabled")
		if geos, ok := rawState["georeplications"]; ok {
			for _, geo := range geos.([]any) {
				delete(geo.(map[string]any), "zone_redundancy_enabled")
			}
		}
		return rawState, nil
	}
}
