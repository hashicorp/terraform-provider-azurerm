package migration

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var _ pluginsdk.StateUpgrade = LegacyVMSSV0ToV1{}

type LegacyVMSSV0ToV1 struct{}

func (LegacyVMSSV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"zones": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
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

					"tier": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"upgrade_policy_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"health_probe_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"automatic_os_upgrade": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"rolling_upgrade_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"max_batch_instance_percent": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  20,
					},

					"max_unhealthy_instance_percent": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  20,
					},

					"max_unhealthy_upgraded_instance_percent": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  20,
					},

					"pause_time_between_batches": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "PT0S",
					},
				},
			},
		},

		"overprovision": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"single_placement_group": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"priority": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"eviction_policy": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"os_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_name_prefix": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"admin_username": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"admin_password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},

					"custom_data": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						StateFunc: userDataStateFunc,
					},
				},
			},
		},

		"os_profile_secrets": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source_vault_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"vault_certificates": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"certificate_url": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"certificate_store": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},

		// lintignore:S018
		"os_profile_windows_config": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"provision_vm_agent": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"enable_automatic_upgrades": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"winrm": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"protocol": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"certificate_url": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					"additional_unattend_config": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"pass": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"component": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"setting_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"content": {
									Type:      pluginsdk.TypeString,
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},
				},
			},
			Set: resourceArmVirtualMachineScaleSetOsProfileWindowsConfigHash,
		},

		// lintignore:S018
		"os_profile_linux_config": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disable_password_authentication": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
						ForceNew: true,
					},
					"ssh_keys": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"path": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"key_data": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
			Set: resourceArmVirtualMachineScaleSetOsProfileLinuxConfigHash,
		},

		"network_profile": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"primary": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"accelerated_networking": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"ip_forwarding": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"network_security_group_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"dns_settings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dns_servers": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"ip_configuration": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"subnet_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"application_gateway_backend_address_pool_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
									Set:      pluginsdk.HashString,
								},

								"application_security_group_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
									Set:      pluginsdk.HashString,
									MaxItems: 20,
								},

								"load_balancer_backend_address_pool_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
									Set:      pluginsdk.HashString,
								},

								"load_balancer_inbound_nat_rules_ids": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Computed: true,
									Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
									Set:      pluginsdk.HashString,
								},

								"primary": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"public_ip_address_configuration": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},

											"idle_timeout": {
												Type:     pluginsdk.TypeInt,
												Required: true,
											},

											"domain_name_label": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Set: resourceArmVirtualMachineScaleSetNetworkConfigurationHash,
		},

		"boot_diagnostics": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"storage_uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		// lintignore:S018
		"storage_profile_os_disk": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"image": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"vhd_containers": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						Set:      pluginsdk.HashString,
					},

					"managed_disk_type": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						Computed:      true,
						ConflictsWith: []string{"storage_profile_os_disk.vhd_containers"},
					},

					"caching": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"os_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"create_option": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
			Set: resourceArmVirtualMachineScaleSetStorageProfileOsDiskHash,
		},

		"storage_profile_data_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"lun": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"create_option": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"caching": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"disk_size_gb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"managed_disk_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		// lintignore:S018
		"storage_profile_image_reference": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"publisher": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"offer": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
			Set: resourceArmVirtualMachineScaleSetStorageProfileImageReferenceHash,
		},

		// lintignore:S018
		"plan": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"publisher": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"product": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"extension": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"publisher": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"type_handler_version": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"auto_upgrade_minor_version": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"settings": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"protected_settings": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},
				},
			},
			Set: resourceArmVirtualMachineScaleSetExtensionHash,
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

func (LegacyVMSSV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// @tombuildsstuff: NOTE, this state migration is essentially pointless
		// however it existed in the legacy migration so even though this is
		//  essentially a noop there's no reason this shouldn't be the same I guess

		client := meta.(*clients.Client).Compute.VMScaleSetClient

		resGroup := rawState["resource_group_name"].(string)
		name := rawState["name"].(string)

		read, err := client.Get(ctx, resGroup, name)
		if err != nil {
			return rawState, err
		}

		rawState["id"] = *read.ID
		return rawState, nil
	}
}

func userDataStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		s = utils.Base64EncodeIfNot(s)
		hash := sha1.Sum([]byte(s))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}

func resourceArmVirtualMachineScaleSetOsProfileWindowsConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["provision_vm_agent"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
		if v, ok := m["enable_automatic_upgrades"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceArmVirtualMachineScaleSetOsProfileLinuxConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%t-", m["disable_password_authentication"].(bool)))
	}

	return pluginsdk.HashString(buf.String())
}

func resourceArmVirtualMachineScaleSetNetworkConfigurationHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%t-", m["primary"].(bool)))
	}

	return pluginsdk.HashString(buf.String())
}

func resourceArmVirtualMachineScaleSetStorageProfileOsDiskHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))

		if v, ok := m["vhd_containers"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(*pluginsdk.Set).List()))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceArmVirtualMachineScaleSetStorageProfileImageReferenceHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["publisher"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["offer"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["sku"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["version"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func resourceArmVirtualMachineScaleSetExtensionHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["publisher"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type_handler_version"].(string)))

		if v, ok := m["auto_upgrade_minor_version"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}

		if v, ok := m["provision_after_extensions"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(*pluginsdk.Set).List()))
		}

		// we need to ensure the whitespace is consistent
		settings := m["settings"].(string)
		if settings != "" {
			expandedSettings, err := pluginsdk.ExpandJsonFromString(settings)
			if err == nil {
				serialisedSettings, err := pluginsdk.FlattenJsonToString(expandedSettings)
				if err == nil {
					buf.WriteString(fmt.Sprintf("%s-", serialisedSettings))
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}
