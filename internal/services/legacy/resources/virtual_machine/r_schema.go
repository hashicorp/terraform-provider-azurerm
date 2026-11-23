// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/legacy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualMachineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"zones": {
			// @tombuildsstuff: since this is the legacy VM resource this is intentionally not using commonschema for consistency
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"plan": {
			Type:     pluginsdk.TypeList,
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

		"availability_set_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			StateFunc: func(id any) string {
				return strings.ToLower(id.(string))
			},
			ConflictsWith: []string{"zones"},
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

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				"Windows_Client",
				"Windows_Server",
			}, false),
		},

		"vm_size": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		// lintignore:S018
		"storage_image_reference": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"publisher": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"offer": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// Note: O+C because Azure assigns a version when not explicitly set
						Computed: true,
						ForceNew: true,
					},
				},
			},
			Set: resourceVirtualMachineStorageImageReferenceHash,
		},

		"storage_os_disk": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"os_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForOperatingSystemTypes(), false),
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"vhd_uri": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ConflictsWith: []string{
							"storage_os_disk.0.managed_disk_id",
							"storage_os_disk.0.managed_disk_type",
						},
					},

					"managed_disk_id": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ForceNew:      true,
						Computed:      true, // azignore:AZS007 - pre-existing violation
						ConflictsWith: []string{"storage_os_disk.0.vhd_uri"},
					},

					"managed_disk_type": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						Computed:      true, // azignore:AZS007 - pre-existing violation
						ConflictsWith: []string{"storage_os_disk.0.vhd_uri"},
						ValidateFunc: validation.StringInSlice([]string{
							string(virtualmachines.StorageAccountTypesPremiumLRS),
							string(virtualmachines.StorageAccountTypesStandardLRS),
							string(virtualmachines.StorageAccountTypesStandardSSDLRS),
						}, false),
					},

					"image_uri": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"caching": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
					},

					"create_option": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"disk_size_gb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// Note: O+C because Azure computes disk size when not specified
						Computed:     true,
						ValidateFunc: validation.IntBetween(0, 32767),
					},

					"write_accelerator_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"delete_os_disk_on_termination": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"storage_data_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"vhd_uri": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"managed_disk_id": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true, // azignore:AZS007 - pre-existing violation
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"managed_disk_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice([]string{
							string(virtualmachines.StorageAccountTypesPremiumLRS),
							string(virtualmachines.StorageAccountTypesStandardLRS),
							string(virtualmachines.StorageAccountTypesStandardSSDLRS),
							string(virtualmachines.StorageAccountTypesUltraSSDLRS),
						}, false),
					},

					"create_option": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: suppress.CaseDifference,
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

					"lun": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"write_accelerator_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"delete_data_disks_on_termination": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"boot_diagnostics": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"storage_uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"additional_capabilities": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ultra_ssd_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		// lintignore:S018
		"os_profile": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_name": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Required: true,
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
						ForceNew:  true,
						Optional:  true,
						Computed:  true, // azignore:AZS007 - pre-existing violation
						StateFunc: helpers.UserDataStateFunc,
					},
				},
			},
			Set: resourceVirtualMachineStorageOsProfileHash,
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
						Default:  false,
					},
					"enable_automatic_upgrades": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"timezone": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     validate.VirtualMachineTimeZoneCaseInsensitive(),
					},
					"winrm": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"protocol": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"HTTP",
										"HTTPS",
									}, false),
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
									ValidateFunc: validation.StringInSlice([]string{
										"oobeSystem",
									}, false),
								},
								"component": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Microsoft-Windows-Shell-Setup",
									}, false),
								},
								"setting_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"AutoLogon",
										"FirstLogonCommands",
									}, false),
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
			Set:           resourceVirtualMachineStorageOsProfileWindowsConfigHash,
			ConflictsWith: []string{"os_profile_linux_config"},
		},

		// lintignore:S018
		"os_profile_linux_config": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disable_password_authentication": {
						Type:     pluginsdk.TypeBool,
						Required: true,
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
									Required: true,
								},
							},
						},
					},
				},
			},
			Set:           resourceVirtualMachineStorageOsProfileLinuxConfigHash,
			ConflictsWith: []string{"os_profile_windows_config"},
		},

		"os_profile_secrets": {
			Type:     pluginsdk.TypeList,
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

		"network_interface_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"primary_network_interface_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}
