// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppV0ToV1 struct{}

func envSchemaV0() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MinItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"value": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"secret_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func secretSchemaV0() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:      pluginsdk.TypeSet,
		Optional:  true,
		Sensitive: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"identity": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"key_vault_secret_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"value": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
				},
			},
		},
	}
}

// probeHeaderSchemaV0 returns the inlined header schema shared by all probe types.
func probeHeaderSchemaV0() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"value": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			},
		},
	}
}

// livenessProbeSchemaV0 returns the inlined liveness_probe schema at v0.
func livenessProbeSchemaV0() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"transport": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"port": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"host": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},
				"header": probeHeaderSchemaV0(),
				"initial_delay": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  1,
				},
				"interval_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  10,
				},
				"timeout": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  1,
				},
				"failure_count_threshold": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  3,
				},
				"termination_grace_period_seconds": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

// readinessProbeSchemaV0 returns the inlined readiness_probe schema at v0.
func readinessProbeSchemaV0() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"transport": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"port": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"host": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},
				"header": probeHeaderSchemaV0(),
				"initial_delay": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  0,
				},
				"interval_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  10,
				},
				"timeout": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  1,
				},
				"failure_count_threshold": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  3,
				},
				"success_count_threshold": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  3,
				},
			},
		},
	}
}

// startupProbeSchemaV0 returns the inlined startup_probe schema at v0.
func startupProbeSchemaV0() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"transport": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"port": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"host": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},
				"header": probeHeaderSchemaV0(),
				"initial_delay": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  0,
				},
				"interval_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  10,
				},
				"timeout": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  1,
				},
				"failure_count_threshold": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  3,
				},
				"termination_grace_period_seconds": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

// volumeMountSchemaV0 returns the inlined volume_mounts schema at v0.
func volumeMountSchemaV0() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"path": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"sub_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

// scaleRuleAuthSchemaV0 returns the inlined authentication schema shared by scale rules.
func scaleRuleAuthSchemaV0() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"secret_name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"trigger_parameter": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			},
		},
	}
}

// httpScaleRuleAuthSchemaV0 is the same as scaleRuleAuthSchemaV0 but trigger_parameter is Optional.
func httpScaleRuleAuthSchemaV0() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"secret_name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"trigger_parameter": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func (ContainerAppV0ToV1) Schema() map[string]*pluginsdk.Schema {
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
			Computed: true,
		},

		"container_app_environment_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"workload_profile_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"template": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"image": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"cpu": {
									Type:     pluginsdk.TypeFloat,
									Required: true,
								},
								"memory": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"ephemeral_storage": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"env": envSchemaV0(),
								"args": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"command": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"liveness_probe":  livenessProbeSchemaV0(),
								"readiness_probe": readinessProbeSchemaV0(),
								"startup_probe":   startupProbeSchemaV0(),
								"volume_mounts":   volumeMountSchemaV0(),
							},
						},
					},

					"init_container": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"image": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"cpu": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
								},
								"memory": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"ephemeral_storage": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"env": envSchemaV0(),
								"args": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"command": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"volume_mounts": volumeMountSchemaV0(),
							},
						},
					},

					"min_replicas": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
					},
					"max_replicas": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  10,
					},
					"cooldown_period_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  300,
					},
					"polling_interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  30,
					},

					"azure_queue_scale_rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"queue_length": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
								"queue_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"authentication": scaleRuleAuthSchemaV0(),
							},
						},
					},

					"custom_scale_rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"metadata": {
									Type:     pluginsdk.TypeMap,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"custom_rule_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"authentication": scaleRuleAuthSchemaV0(),
								"identity_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"http_scale_rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"concurrent_requests": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"authentication": httpScaleRuleAuthSchemaV0(),
							},
						},
					},

					"tcp_scale_rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"concurrent_requests": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"authentication": httpScaleRuleAuthSchemaV0(),
							},
						},
					},

					"volume": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"storage_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "EmptyDir",
								},
								"storage_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"mount_options": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"revision_suffix": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"termination_grace_period_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
					},
				},
			},
		},

		"secret": secretSchemaV0(),

		"registry": {
			Type:     pluginsdk.TypeList,
			MinItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"server": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"username": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"password_secret_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"identity": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
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

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"revision_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"ingress": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allow_insecure_connections": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"external_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"custom_domain": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"certificate_binding_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"certificate_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
					"cors": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allowed_origins": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"allow_credentials_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"allowed_headers": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"allowed_methods": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"exposed_headers": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"max_age_in_seconds": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
					},
					"fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"ip_security_restriction": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"description": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"ip_address_range": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"target_port": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"exposed_port": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"traffic_weight": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"label": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"revision_suffix": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"latest_revision": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"percentage": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
							},
						},
					},
					"transport": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "auto",
					},
					"client_certificate_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"dapr": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"app_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"app_port": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"app_protocol": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "http",
					},
				},
			},
		},

		"max_inactive_revisions": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"latest_revision_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"latest_revision_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"custom_domain_verification_id": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (ContainerAppV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		normalizeEnvVars(rawState)
		normalizeSecrets(rawState)
		return rawState, nil
	}
}

func normalizeEnvVars(rawState map[string]interface{}) {
	templates, _ := rawState["template"].([]interface{})
	for _, tmpl := range templates {
		tmplMap, ok := tmpl.(map[string]interface{})
		if !ok {
			continue
		}
		for _, key := range []string{"container", "init_container"} {
			containers, _ := tmplMap[key].([]interface{})
			for _, c := range containers {
				cMap, ok := c.(map[string]interface{})
				if !ok {
					continue
				}
				envs, _ := cMap["env"].([]interface{})
				for _, e := range envs {
					eMap, ok := e.(map[string]interface{})
					if !ok {
						continue
					}
					if eMap["value"] == nil {
						eMap["value"] = ""
					}
					if eMap["secret_name"] == nil {
						eMap["secret_name"] = ""
					}
				}
			}
		}
	}
}

func normalizeSecrets(rawState map[string]interface{}) {
	secrets, _ := rawState["secret"].([]interface{})
	for _, s := range secrets {
		sMap, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		if sMap["identity"] == nil {
			sMap["identity"] = ""
		}
		if sMap["key_vault_secret_id"] == nil {
			sMap["key_vault_secret_id"] = ""
		}
		if sMap["value"] == nil {
			sMap["value"] = ""
		}
	}
}
