package migration

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var _ pluginsdk.StateUpgrade = LegacyVMSSV0ToV1{}

type LegacyVMSSV0ToV1 struct{}

func (LegacyVMSSV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"zones": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
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
					"identity_ids": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"principal_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"sku": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},

					"tier": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},

					"capacity": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},

		"license_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"upgrade_policy_mode": {
			Type:     schema.TypeString,
			Required: true,
		},

		"health_probe_id": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"automatic_os_upgrade": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"rolling_upgrade_policy": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"max_batch_instance_percent": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  20,
					},

					"max_unhealthy_instance_percent": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  20,
					},

					"max_unhealthy_upgraded_instance_percent": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  20,
					},

					"pause_time_between_batches": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  "PT0S",
					},
				},
			},
		},

		"overprovision": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},

		"single_placement_group": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"priority": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"eviction_policy": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"os_profile": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"computer_name_prefix": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},

					"admin_username": {
						Type:     schema.TypeString,
						Required: true,
					},

					"admin_password": {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
					},

					"custom_data": {
						Type:      schema.TypeString,
						Optional:  true,
						StateFunc: userDataStateFunc,
					},
				},
			},
		},

		"os_profile_secrets": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"source_vault_id": {
						Type:     schema.TypeString,
						Required: true,
					},

					"vault_certificates": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"certificate_url": {
									Type:     schema.TypeString,
									Required: true,
								},
								"certificate_store": {
									Type:     schema.TypeString,
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
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"provision_vm_agent": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"enable_automatic_upgrades": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"winrm": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"protocol": {
									Type:     schema.TypeString,
									Required: true,
								},
								"certificate_url": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"additional_unattend_config": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"pass": {
									Type:     schema.TypeString,
									Required: true,
								},
								"component": {
									Type:     schema.TypeString,
									Required: true,
								},
								"setting_name": {
									Type:     schema.TypeString,
									Required: true,
								},
								"content": {
									Type:      schema.TypeString,
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
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disable_password_authentication": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
						ForceNew: true,
					},
					"ssh_keys": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"path": {
									Type:     schema.TypeString,
									Required: true,
								},
								"key_data": {
									Type:     schema.TypeString,
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
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},

					"primary": {
						Type:     schema.TypeBool,
						Required: true,
					},

					"accelerated_networking": {
						Type:     schema.TypeBool,
						Optional: true,
					},

					"ip_forwarding": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},

					"network_security_group_id": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"dns_settings": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dns_servers": {
									Type:     schema.TypeList,
									Required: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},

					"ip_configuration": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Required: true,
								},

								"subnet_id": {
									Type:     schema.TypeString,
									Required: true,
								},

								"application_gateway_backend_address_pool_ids": {
									Type:     schema.TypeSet,
									Optional: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
									Set:      schema.HashString,
								},

								"application_security_group_ids": {
									Type:     schema.TypeSet,
									Optional: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
									Set:      schema.HashString,
									MaxItems: 20,
								},

								"load_balancer_backend_address_pool_ids": {
									Type:     schema.TypeSet,
									Optional: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
									Set:      schema.HashString,
								},

								"load_balancer_inbound_nat_rules_ids": {
									Type:     schema.TypeSet,
									Optional: true,
									Computed: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
									Set:      schema.HashString,
								},

								"primary": {
									Type:     schema.TypeBool,
									Required: true,
								},

								"public_ip_address_configuration": {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"name": {
												Type:     schema.TypeString,
												Required: true,
											},

											"idle_timeout": {
												Type:     schema.TypeInt,
												Required: true,
											},

											"domain_name_label": {
												Type:     schema.TypeString,
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
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},

					"storage_uri": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},

		// lintignore:S018
		"storage_profile_os_disk": {
			Type:     schema.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"image": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"vhd_containers": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
						Set:      schema.HashString,
					},

					"managed_disk_type": {
						Type:          schema.TypeString,
						Optional:      true,
						Computed:      true,
						ConflictsWith: []string{"storage_profile_os_disk.vhd_containers"},
					},

					"caching": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},

					"os_type": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"create_option": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			Set: resourceArmVirtualMachineScaleSetStorageProfileOsDiskHash,
		},

		"storage_profile_data_disk": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"lun": {
						Type:     schema.TypeInt,
						Required: true,
					},

					"create_option": {
						Type:     schema.TypeString,
						Required: true,
					},

					"caching": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},

					"disk_size_gb": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},

					"managed_disk_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		// lintignore:S018
		"storage_profile_image_reference": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"publisher": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"offer": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"sku": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"version": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			Set: resourceArmVirtualMachineScaleSetStorageProfileImageReferenceHash,
		},

		// lintignore:S018
		"plan": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},

					"publisher": {
						Type:     schema.TypeString,
						Required: true,
					},

					"product": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},

		"extension": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},

					"publisher": {
						Type:     schema.TypeString,
						Required: true,
					},

					"type": {
						Type:     schema.TypeString,
						Required: true,
					},

					"type_handler_version": {
						Type:     schema.TypeString,
						Required: true,
					},

					"auto_upgrade_minor_version": {
						Type:     schema.TypeBool,
						Optional: true,
					},

					"settings": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"protected_settings": {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
					},
				},
			},
			Set: resourceArmVirtualMachineScaleSetExtensionHash,
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

	return schema.HashString(buf.String())
}

func resourceArmVirtualMachineScaleSetOsProfileLinuxConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%t-", m["disable_password_authentication"].(bool)))
	}

	return schema.HashString(buf.String())
}

func resourceArmVirtualMachineScaleSetNetworkConfigurationHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%t-", m["primary"].(bool)))
	}

	return schema.HashString(buf.String())
}

func resourceArmVirtualMachineScaleSetStorageProfileOsDiskHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))

		if v, ok := m["vhd_containers"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(*schema.Set).List()))
		}
	}

	return schema.HashString(buf.String())
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

	return schema.HashString(buf.String())
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
			buf.WriteString(fmt.Sprintf("%s-", v.(*schema.Set).List()))
		}

		// we need to ensure the whitespace is consistent
		settings := m["settings"].(string)
		if settings != "" {
			expandedSettings, err := structure.ExpandJsonFromString(settings)
			if err == nil {
				serialisedSettings, err := structure.FlattenJsonToString(expandedSettings)
				if err == nil {
					buf.WriteString(fmt.Sprintf("%s-", serialisedSettings))
				}
			}
		}
	}

	return schema.HashString(buf.String())
}
