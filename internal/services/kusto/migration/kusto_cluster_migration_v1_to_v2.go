// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoAttachedClusterV1ToV2 struct{}

func (s KustoAttachedClusterV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
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

		"identity": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"identity_ids": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
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

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"allowed_fqdns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"allowed_ip_ranges": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"trusted_external_tenants": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"optimized_auto_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"minimum_instances": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"maximum_instances": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"virtual_network_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"engine_public_ip_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"data_management_public_ip_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"language_extensions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"data_ingestion_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_ip_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"outbound_network_access_restricted": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"double_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"auto_stop_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"disk_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"streaming_ingestion_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"purge_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"zones": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
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

func (s KustoAttachedClusterV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// This migration fixes #27580, which prevented provider upgrades from <4.0.0 to >=4.0.0.
		// If the current state file contains the `language_extensions` argument as a list, we'll migrate it to the expected
		// block format. Otherwise, do nothing.
		if extensionsRaw, ok := rawState["language_extensions"]; ok {
			if extensions, ok := extensionsRaw.([]interface{}); ok && len(extensions) > 0 {
				if _, ok := extensions[0].(map[string]interface{}); ok {
					return rawState, nil
				}

				log.Print("[DEBUG] Migrating `language_extensions` to the block format")
				newExtensions := make([]map[string]interface{}, 0)

				for _, v := range extensions {
					switch v {
					case "R":
						newExtensions = append(newExtensions, map[string]interface{}{
							"name":  "R",
							"image": "R",
						})
					case "PYTHON":
						newExtensions = append(newExtensions, map[string]interface{}{
							"name":  "PYTHON",
							"image": "Python3_6_5",
						})
					case "PYTHON_3.10.8":
						newExtensions = append(newExtensions, map[string]interface{}{
							"name":  "PYTHON",
							"image": "Python3_10_8",
						})
					}
				}

				rawState["language_extensions"] = newExtensions
			}
		}

		return rawState, nil
	}
}
