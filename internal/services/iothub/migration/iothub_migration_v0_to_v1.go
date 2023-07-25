// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type IoTHubV0ToV1 struct{}

func (s IoTHubV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
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

		"sku": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"shared_access_policy": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"primary_key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"secondary_key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"permissions": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"event_hub_partition_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},
		"event_hub_retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"file_upload": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_string": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"container_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"authentication_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"notifications": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"max_delivery_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  10,
					},
					"sas_ttl": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"default_ttl": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"lock_duration": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"endpoint": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"authentication_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"endpoint_uri": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"entity_path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"connection_string": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"batch_frequency_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  300,
					},

					"max_chunk_size_in_bytes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  314572800,
					},

					"container_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"encoding": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"file_name_format": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"resource_group_name": commonschema.ResourceGroupNameOptional(),
				},
			},
		},

		"route": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"source": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"condition": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "true",
					},
					"endpoint_names": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Required: true,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"enrichment": {
			Type: pluginsdk.TypeList,
			// Currently only 10 enrichments is allowed for standard or basic tier, 2 for Free tier.
			MaxItems: 10,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"value": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"endpoint_names": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Required: true,
					},
				},
			},
		},

		"fallback_route": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"condition": {
						// The condition is a string value representing device-to-cloud message routes query expression
						// https://docs.microsoft.com/en-us/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "true",
					},
					"endpoint_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"network_rule_set": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"apply_to_builtin_eventhub_endpoint": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"ip_rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"ip_mask": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"action": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},

		"cloud_to_device": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"max_delivery_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  10,
					},
					"default_ttl": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "PT1H",
					},
					"feedback": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"time_to_live": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "PT1H",
								},
								"max_delivery_count": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  10,
								},
								"lock_duration": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "PT60S",
								},
							},
						},
					},
				},
			},
		},

		"min_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"event_hub_events_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"event_hub_events_namespace": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"event_hub_operations_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"event_hub_events_path": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"event_hub_operations_path": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (s IoTHubV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.IotHubIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
