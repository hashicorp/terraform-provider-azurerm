// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/legacy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualMachineScaleSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"zones": {
			// @tombuildsstuff: since this is the legacy VMSS resource this is intentionally not using commonschema for consistency
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"tier": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true, // azignore:AZS007 - pre-existing violation
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				"Windows_Client",
				"Windows_Server",
			}, false),
		},

		"upgrade_policy_mode": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForUpgradeMode(), false),
		},

		"health_probe_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
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
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      20,
						ValidateFunc: validation.IntBetween(5, 100),
					},

					"max_unhealthy_instance_percent": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      20,
						ValidateFunc: validation.IntBetween(5, 100),
					},

					"max_unhealthy_upgraded_instance_percent": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      20,
						ValidateFunc: validation.IntBetween(5, 100),
					},

					"pause_time_between_batches": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "PT0S",
						ValidateFunc: validation.ISO8601Duration,
					},
				},
			},
			DiffSuppressFunc: azureRmVirtualMachineScaleSetSuppressRollingUpgradePolicyDiff,
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
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualmachinescalesets.VirtualMachinePriorityTypesLow),
				string(virtualmachinescalesets.VirtualMachinePriorityTypesRegular),
			}, false),
		},

		"eviction_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualmachinescalesets.PossibleValuesForVirtualMachineEvictionPolicyTypes(), false),
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"admin_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"custom_data": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						StateFunc:        helpers.UserDataStateFunc,
						DiffSuppressFunc: userDataDiffSuppressFunc,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
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
			Set: resourceVirtualMachineScaleSetOsProfileWindowsConfigHash,
		},

		// lintignore:S018
		"os_profile_linux_config": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
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
			Set: resourceVirtualMachineScaleSetOsProfileLinuxConfigHash,
		},

		// lintignore:S018
		"network_profile": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID,
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
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"subnet_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: azure.ValidateResourceID,
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
										Type:         pluginsdk.TypeString,
										ValidateFunc: azure.ValidateResourceID,
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
									Computed: true, // azignore:AZS007 - pre-existing violation
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
												Type:         pluginsdk.TypeInt,
												Required:     true,
												ValidateFunc: validation.IntBetween(4, 32),
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
			Set: resourceVirtualMachineScaleSetNetworkConfigurationHash,
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice([]string{
							string(virtualmachinescalesets.StorageAccountTypesPremiumLRS),
							string(virtualmachinescalesets.StorageAccountTypesStandardLRS),
							string(virtualmachinescalesets.StorageAccountTypesStandardSSDLRS),
						}, false),
					},

					"caching": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
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
			Set: resourceVirtualMachineScaleSetStorageProfileOsDiskHash,
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
						Computed: true, // azignore:AZS007 - pre-existing violation
					},

					"disk_size_gb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// Note: O+C because Azure computes disk size when not specified
						Computed:     true,
						ValidateFunc: validation.IntBetween(0, 32767),
					},

					"managed_disk_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice([]string{
							string(virtualmachinescalesets.StorageAccountTypesPremiumLRS),
							string(virtualmachinescalesets.StorageAccountTypesStandardLRS),
							string(virtualmachinescalesets.StorageAccountTypesStandardSSDLRS),
						}, false),
					},
				},
			},
		},

		// lintignore:S018
		"storage_profile_image_reference": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
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
			Set: resourceVirtualMachineScaleSetStorageProfileImageReferenceHash,
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

		// lintignore:S018
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

					"provision_after_extensions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						Set: pluginsdk.HashString,
					},

					"settings": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
					},

					"protected_settings": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Sensitive:        true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
					},
				},
			},
			Set: resourceVirtualMachineScaleSetExtensionHash,
		},

		"proximity_placement_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,

			// We have to ignore case due to incorrect capitalisation of resource group name in
			// proximity placement group ID in the response we get from the API request
			//
			// todo can be removed when https://github.com/Azure/azure-sdk-for-go/issues/5699 is fixed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"tags": commonschema.Tags(),
	}
}
