// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppJobV0ToV1 struct{}

func jobEnvSchemaV0() *pluginsdk.Schema {
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

func jobSecretSchemaV0() *pluginsdk.Schema {
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

// jobProbeHeaderSchemaV0 returns the inlined header schema shared by all probe types.
func jobProbeHeaderSchemaV0() *pluginsdk.Schema {
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

// jobLivenessProbeSchemaV0 returns the inlined liveness_probe schema at v0.
func jobLivenessProbeSchemaV0() *pluginsdk.Schema {
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
				"header": jobProbeHeaderSchemaV0(),
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

// jobReadinessProbeSchemaV0 returns the inlined readiness_probe schema at v0.
func jobReadinessProbeSchemaV0() *pluginsdk.Schema {
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
				"header": jobProbeHeaderSchemaV0(),
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

// jobStartupProbeSchemaV0 returns the inlined startup_probe schema at v0.
func jobStartupProbeSchemaV0() *pluginsdk.Schema {
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
				"header": jobProbeHeaderSchemaV0(),
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

// jobVolumeMountSchemaV0 returns the inlined volume_mounts schema at v0.
func jobVolumeMountSchemaV0() *pluginsdk.Schema {
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

// jobScaleRuleAuthSchemaV0 returns the inlined authentication schema for scale rules.
func jobScaleRuleAuthSchemaV0() *pluginsdk.Schema {
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

func (ContainerAppJobV0ToV1) Schema() map[string]*pluginsdk.Schema {
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
			Required: true,
			MaxItems: 1,
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
								"env": jobEnvSchemaV0(),
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
								"liveness_probe":  jobLivenessProbeSchemaV0(),
								"readiness_probe": jobReadinessProbeSchemaV0(),
								"startup_probe":   jobStartupProbeSchemaV0(),
								"volume_mounts":   jobVolumeMountSchemaV0(),
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
								"env": jobEnvSchemaV0(),
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
								"volume_mounts": jobVolumeMountSchemaV0(),
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
				},
			},
		},

		"secret": jobSecretSchemaV0(),

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

		"replica_timeout_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"replica_retry_limit": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"event_trigger_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"parallelism": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"replica_completion_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"scale": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max_executions": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  100,
								},
								"min_executions": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  0,
								},
								"polling_interval_in_seconds": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  30,
								},
								"rules": {
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
											"authentication": jobScaleRuleAuthSchemaV0(),
											"identity_id": {
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

		"schedule_trigger_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cron_expression": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"parallelism": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"replica_completion_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"manual_trigger_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"parallelism": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"replica_completion_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"event_stream_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ContainerAppJobV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		normalizeJobEnvVars(rawState)
		normalizeJobSecrets(rawState)
		return rawState, nil
	}
}

func normalizeJobEnvVars(rawState map[string]interface{}) {
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

func normalizeJobSecrets(rawState map[string]interface{}) {
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
