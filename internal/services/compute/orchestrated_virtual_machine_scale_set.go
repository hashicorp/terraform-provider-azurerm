// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func OrchestratedVirtualMachineScaleSetOSProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"custom_data": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsBase64,
				},
				"windows_configuration": OrchestratedVirtualMachineScaleSetWindowsConfigurationSchema(),
				"linux_configuration":   OrchestratedVirtualMachineScaleSetLinuxConfigurationSchema(),
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetWindowsConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"admin_username": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validateAdminUsernameWindows,
				},

				"admin_password": {
					Type:             pluginsdk.TypeString,
					Required:         true,
					ForceNew:         true,
					Sensitive:        true,
					DiffSuppressFunc: adminPasswordDiffSuppressFunc,
					ValidateFunc:     validatePasswordComplexityWindows,
				},

				"computer_name_prefix": computerPrefixWindowsSchema(),

				"additional_unattend_content": additionalUnattendContentSchema(),

				// TODO 4.0: change this from enable_* to *_enabled
				"enable_automatic_updates": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"hotpatching_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"provision_vm_agent": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
					ForceNew: true,
				},

				"patch_assessment_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(virtualmachinescalesets.WindowsPatchAssessmentModeImageDefault),
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform),
						string(virtualmachinescalesets.WindowsPatchAssessmentModeImageDefault),
					}, false),
				},

				"patch_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByOS),
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByOS),
						string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform),
						string(virtualmachinescalesets.WindowsVMGuestPatchModeManual),
					}, false),
				},

				"secret": windowsSecretSchema(),

				"timezone": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.VirtualMachineTimeZone(),
				},

				"winrm_listener": winRmListenerSchema(),
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetLinuxConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"admin_username": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validateAdminUsernameLinux,
				},

				"admin_password": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					ForceNew:         true,
					Sensitive:        true,
					DiffSuppressFunc: adminPasswordDiffSuppressFunc,
					ValidateFunc:     validatePasswordComplexityLinux,
				},

				"admin_ssh_key":        SSHKeysSchema(false),
				"computer_name_prefix": computerPrefixLinuxSchema(),

				"disable_password_authentication": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"provision_vm_agent": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
					ForceNew: true,
				},

				"patch_assessment_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(virtualmachinescalesets.LinuxPatchAssessmentModeImageDefault),
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform),
						string(virtualmachinescalesets.LinuxPatchAssessmentModeImageDefault),
					}, false),
				},

				"patch_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(virtualmachinescalesets.LinuxVMGuestPatchModeImageDefault),
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.LinuxVMGuestPatchModeImageDefault),
						string(virtualmachinescalesets.LinuxVMGuestPatchModeAutomaticByPlatform),
					}, false),
				},

				"secret": linuxSecretSchema(),
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetExtensionsSchema() *pluginsdk.Schema {
	schema := &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type_handler_version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"auto_upgrade_minor_version_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				// Only supported in Orchestrated mode
				"failure_suppression_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"force_extension_execution_on_change": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"protected_settings": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsJSON,
				},

				// Need to check `protected_settings_from_key_vault` conflicting with `protected_settings` in iteration
				"protected_settings_from_key_vault": protectedSettingsFromKeyVaultSchema(false),

				"extensions_to_provision_after_vm_creation": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"settings": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					ValidateFunc:     validation.StringIsJSON,
					DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				},
			},
		},
	}

	if !features.FourPointOhBeta() {
		schema.Elem.(*pluginsdk.Resource).Schema["failure_suppression_enabled"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Optional: true,
		}
	}

	return schema
}

func OrchestratedVirtualMachineScaleSetNetworkInterfaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"ip_configuration": orchestratedVirtualMachineScaleSetIPConfigurationSchema(),

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				// TODO 4.0: change this from enable_* to *_enabled
				"enable_accelerated_networking": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				// TODO 4.0: change this from enable_* to *_enabled
				"enable_ip_forwarding": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"network_security_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func orchestratedVirtualMachineScaleSetIPConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				// Optional
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
						ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
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

				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"public_ip_address": orchestratedVirtualMachineScaleSetPublicIPAddressSchema(),

				"subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(virtualmachinescalesets.IPVersionIPvFour),
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.IPVersionIPvFour),
						string(virtualmachinescalesets.IPVersionIPvSix),
					}, false),
				},
			},
		},
	}
}

func orchestratedVirtualMachineScaleSetPublicIPAddressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				// Optional
				"domain_name_label": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.OrchestratedDomainNameLabel,
				},
				"idle_timeout_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(4, 32),
				},
				"ip_tag": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"tag": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"type": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				// TODO: preview feature
				// $ az feature register --namespace Microsoft.Network --name AllowBringYourOwnPublicIpAddress
				// $ az provider register -n Microsoft.Network
				"public_ip_prefix_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
				},

				"sku_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validate.OrchestratedVirtualMachineScaleSetPublicIPSku,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					Default:  string(virtualmachinescalesets.IPVersionIPvFour),
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.IPVersionIPvFour),
						string(virtualmachinescalesets.IPVersionIPvSix),
					}, false),
				},
			},
		},
	}
}

func computerPrefixWindowsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,

		// Computed since we reuse the VM name if one's not specified
		Computed:     true,
		ForceNew:     true,
		ValidateFunc: validate.WindowsComputerNamePrefix,
	}
}

func computerPrefixLinuxSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,

		// Computed since we reuse the VM name if one's not specified
		Computed:     true,
		ForceNew:     true,
		ValidateFunc: validate.LinuxComputerNamePrefix,
	}
}

func OrchestratedVirtualMachineScaleSetDataDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"caching": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.CachingTypesNone),
						string(virtualmachinescalesets.CachingTypesReadOnly),
						string(virtualmachinescalesets.CachingTypesReadWrite),
					}, false),
				},

				"create_option": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.DiskCreateOptionTypesEmpty),
						string(virtualmachinescalesets.DiskCreateOptionTypesFromImage),
					}, false),
					Default: string(virtualmachinescalesets.DiskCreateOptionTypesEmpty),
				},

				"disk_encryption_set_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// whilst the API allows updating this value, it's never actually set at Azure's end
					// presumably this'll take effect once key rotation is supported a few months post-GA?
					// however for now let's make this ForceNew since it can't be (successfully) updated
					ForceNew:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"disk_size_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"lun": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 2000), // TODO: confirm upper bounds
				},

				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.StorageAccountTypesPremiumLRS),
						string(virtualmachinescalesets.StorageAccountTypesPremiumVTwoLRS),
						string(virtualmachinescalesets.StorageAccountTypesPremiumZRS),
						string(virtualmachinescalesets.StorageAccountTypesStandardLRS),
						string(virtualmachinescalesets.StorageAccountTypesStandardSSDLRS),
						string(virtualmachinescalesets.StorageAccountTypesStandardSSDZRS),
						string(virtualmachinescalesets.StorageAccountTypesUltraSSDLRS),
					}, false),
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				// TODO rename `ultra_ssd_disk_iops_read_write` to `disk_iops_read_write` in 4.0
				"ultra_ssd_disk_iops_read_write": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
					Computed:     true,
				},

				// TODO rename `ultra_ssd_disk_mbps_read_write` to `disk_mbps_read_write` in 4.0
				"ultra_ssd_disk_mbps_read_write": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
					Computed:     true,
				},
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetAdditionalCapabilitiesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ultra_ssd_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
					ForceNew: true,
				},
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetOSDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"caching": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachinescalesets.CachingTypesNone),
						string(virtualmachinescalesets.CachingTypesReadOnly),
						string(virtualmachinescalesets.CachingTypesReadWrite),
					}, false),
				},
				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					// whilst this appears in the Update block the API returns this when changing:
					// Changing property 'osDisk.managedDisk.storageAccountType' is not allowed
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: OS Disks don't support Ultra SSDs or PremiumV2_LRS
						string(virtualmachinescalesets.StorageAccountTypesPremiumLRS),
						string(virtualmachinescalesets.StorageAccountTypesPremiumZRS),
						string(virtualmachinescalesets.StorageAccountTypesStandardLRS),
						string(virtualmachinescalesets.StorageAccountTypesStandardSSDLRS),
						string(virtualmachinescalesets.StorageAccountTypesStandardSSDZRS),
					}, false),
				},

				"diff_disk_settings": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"option": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(virtualmachinescalesets.DiffDiskOptionsLocal),
								}, false),
							},
							"placement": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ForceNew: true,
								Default:  string(virtualmachinescalesets.DiffDiskPlacementCacheDisk),
								ValidateFunc: validation.StringInSlice([]string{
									string(virtualmachinescalesets.DiffDiskPlacementCacheDisk),
									string(virtualmachinescalesets.DiffDiskPlacementResourceDisk),
								}, false),
							}},
					},
				},

				"disk_encryption_set_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// whilst the API allows updating this value, it's never actually set at Azure's end
					// presumably this'll take effect once key rotation is supported a few months post-GA?
					// however for now let's make this ForceNew since it can't be (successfully) updated
					ForceNew:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"disk_size_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 4095),
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetTerminationNotificationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				"timeout": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
					Default:      "PT5M",
				},
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetPriorityMixPolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"base_regular_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 1000),
					Default:      0,
				},
				"regular_percentage_above_base": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 100),
					Default:      0,
				},
			},
		},
	}
}

func FlattenOrchestratedVirtualMachineScaleSetOSProfile(input *virtualmachinescalesets.VirtualMachineScaleSetOSProfile, d *pluginsdk.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if input.CustomData != nil {
		output["custom_data"] = *input.CustomData
	} else {
		// look up the current custom data
		output["custom_data"] = utils.Base64EncodeIfNot(d.Get("os_profile.0.custom_data").(string))
	}

	if winConfig := input.WindowsConfiguration; winConfig != nil {
		output["windows_configuration"] = flattenOrchestratedVirtualMachineScaleSetWindowsConfiguration(input, d)
	}

	if linConfig := input.LinuxConfiguration; linConfig != nil {
		output["linux_configuration"] = flattenOrchestratedVirtualMachineScaleSetLinuxConfiguration(input, d)
	}

	return []interface{}{output}
}

func validateAdminUsernameWindows(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	// **Disallowed values:**
	invalidUserNames := []string{
		" ", "administrator", "admin", "user", "user1", "test", "user2", "test1", "user3", "admin1", "1", "123", "a",
		"actuser", "adm", "admin2", "aspnet", "backup", "console", "david", "guest", "john", "owner", "root", "server",
		"sql", "support", "support_388945a0", "sys", "test2", "test3", "user4", "user5",
	}

	for _, str := range invalidUserNames {
		if strings.EqualFold(v, str) {
			errors = append(errors, fmt.Errorf("%q can not be one of %v, got %q", key, invalidUserNames, v))
			return warnings, errors
		}
	}

	// Cannot end in "."
	if strings.HasSuffix(input.(string), ".") {
		errors = append(errors, fmt.Errorf("%q can not end with a '.', got %q", key, v))
		return warnings, errors
	}

	if len(v) < 1 || len(v) > 20 {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 20 characters in length, got %q(%d characters)", key, v, len(v)))
		return warnings, errors
	}

	return
}

func validateAdminUsernameLinux(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	// **Disallowed values:**
	invalidUserNames := []string{
		" ", "abrt", "adm", "admin", "audio", "backup", "bin", "cdrom", "cgred", "console", "crontab", "daemon", "dbus", "dialout", "dip",
		"disk", "fax", "floppy", "ftp", "fuse", "games", "gnats", "gopher", "haldaemon", "halt", "irc", "kmem", "landscape", "libuuid", "list",
		"lock", "lp", "mail", "maildrop", "man", "mem", "messagebus", "mlocate", "modem", "netdev", "news", "nfsnobody", "nobody", "nogroup",
		"ntp", "operator", "oprofile", "plugdev", "polkituser", "postdrop", "postfix", "proxy", "public", "qpidd", "root", "rpc", "rpcuser",
		"sasl", "saslauth", "shadow", "shutdown", "slocate", "src", "ssh", "sshd", "staff", "stapdev", "stapusr", "sudo", "sync", "sys", "syslog",
		"tape", "tcpdump", "test", "trusted", "tty", "users", "utempter", "utmp", "uucp", "uuidd", "vcsa", "video", "voice", "wheel", "whoopsie",
		"www", "www-data", "wwwrun", "xok",
	}

	for _, str := range invalidUserNames {
		if strings.EqualFold(v, str) {
			errors = append(errors, fmt.Errorf("%q can not be one of %s, got %q", key, azure.QuotedStringSlice(invalidUserNames), v))
			return warnings, errors
		}
	}

	if len(v) < 1 || len(v) > 64 {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length, got %q(%d characters)", key, v, len(v)))
		return warnings, errors
	}

	return
}

func validatePasswordComplexityWindows(input interface{}, key string) (warnings []string, errors []error) {
	return validatePasswordComplexity(input, key, 8, 123)
}

func validatePasswordComplexityLinux(input interface{}, key string) (warnings []string, errors []error) {
	return validatePasswordComplexity(input, key, 6, 72)
}

func validatePasswordComplexity(input interface{}, key string, min int, max int) (warnings []string, errors []error) {
	password, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return warnings, errors
	}

	complexityMatch := 0
	re := regexp.MustCompile(`[a-z]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	re = regexp.MustCompile(`[A-Z]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	re = regexp.MustCompile(`[0-9]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	re = regexp.MustCompile(`[\W_]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	if complexityMatch < 3 {
		errors = append(errors, fmt.Errorf("%q did not meet minimum password complexity requirements. A password must contain at least 3 of the 4 following conditions: a lower case character, a upper case character, a digit and/or a special character. Got %q", key, password))
		return warnings, errors
	}

	if len(password) < min || len(password) > max {
		errors = append(errors, fmt.Errorf("%q must be at least 6 characters long and less than 72 characters long. Got %q(%d characters)", key, password, len(password)))
		return warnings, errors
	}

	// NOTE: I realize that some of these will not pass the above complexity checks, but they are in the API so I am checking
	// the same values that the API is...
	disallowedValues := []string{
		"abc@123", "P@$$w0rd", "P@ssw0rd", "P@ssword123", "Pa$$word", "pass@word1", "Password!", "Password1", "Password22", "iloveyou!",
	}

	for _, str := range disallowedValues {
		if password == str {
			errors = append(errors, fmt.Errorf("%q can not be one of %s, got %q", key, azure.QuotedStringSlice(disallowedValues), password))
			return warnings, errors
		}
	}

	return warnings, errors
}

func ExpandOrchestratedVirtualMachineScaleSetAdditionalCapabilities(input []interface{}) *virtualmachinescalesets.AdditionalCapabilities {
	capabilities := virtualmachinescalesets.AdditionalCapabilities{}

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})

		capabilities.UltraSSDEnabled = pointer.To(raw["ultra_ssd_enabled"].(bool))
	}

	return &capabilities
}

func FlattenOrchestratedVirtualMachineScaleSetAdditionalCapabilities(input *virtualmachinescalesets.AdditionalCapabilities) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ultraSsdEnabled := false

	if input.UltraSSDEnabled != nil {
		ultraSsdEnabled = *input.UltraSSDEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"ultra_ssd_enabled": ultraSsdEnabled,
		},
	}
}

func expandOrchestratedVirtualMachineScaleSetOsProfileWithWindowsConfiguration(input map[string]interface{}, customData string) *virtualmachinescalesets.VirtualMachineScaleSetOSProfile {
	osProfile := virtualmachinescalesets.VirtualMachineScaleSetOSProfile{}
	winConfig := virtualmachinescalesets.WindowsConfiguration{}
	patchSettings := virtualmachinescalesets.PatchSettings{}

	if len(input) > 0 {
		osProfile.CustomData = pointer.To(customData)
		osProfile.AdminUsername = pointer.To(input["admin_username"].(string))
		osProfile.AdminPassword = pointer.To(input["admin_password"].(string))

		if computerPrefix := input["computer_name_prefix"].(string); computerPrefix != "" {
			osProfile.ComputerNamePrefix = pointer.To(computerPrefix)
		}

		if secrets := input["secret"].([]interface{}); len(secrets) > 0 {
			osProfile.Secrets = expandWindowsSecretsVMSS(secrets)
		}

		if additionalUnattendContents := input["additional_unattend_content"].([]interface{}); len(additionalUnattendContents) > 0 {
			winConfig.AdditionalUnattendContent = expandWindowsConfigurationAdditionalUnattendContent(input["additional_unattend_content"].([]interface{}))
		}
		winConfig.EnableAutomaticUpdates = pointer.To(input["enable_automatic_updates"].(bool))
		winConfig.ProvisionVMAgent = pointer.To(input["provision_vm_agent"].(bool))
		winRmListenersRaw := input["winrm_listener"].(*pluginsdk.Set).List()
		winConfig.WinRM = expandWinRMListenerVMSS(winRmListenersRaw)

		// Automatic VM Guest Patching and Hotpatching settings
		patchSettings.AssessmentMode = pointer.To(virtualmachinescalesets.WindowsPatchAssessmentMode(input["patch_assessment_mode"].(string)))
		patchSettings.PatchMode = pointer.To(virtualmachinescalesets.WindowsVMGuestPatchMode(input["patch_mode"].(string)))
		patchSettings.EnableHotpatching = pointer.To(input["hotpatching_enabled"].(bool))
		winConfig.PatchSettings = &patchSettings

		// due to a change in RP behavor, it will now throw and error if we pass an empty
		// string add check to only include it if it is actually defined in the config file
		timeZone := input["timezone"].(string)
		if timeZone != "" {
			winConfig.TimeZone = pointer.To(timeZone)
		}
	}

	osProfile.WindowsConfiguration = &winConfig

	return &osProfile
}

func expandOrchestratedVirtualMachineScaleSetOsProfileWithLinuxConfiguration(input map[string]interface{}, customData string) *virtualmachinescalesets.VirtualMachineScaleSetOSProfile {
	osProfile := virtualmachinescalesets.VirtualMachineScaleSetOSProfile{}
	linConfig := virtualmachinescalesets.LinuxConfiguration{}
	patchSettings := virtualmachinescalesets.LinuxPatchSettings{}

	if len(input) > 0 {
		osProfile.CustomData = pointer.To(customData)
		osProfile.AdminUsername = pointer.To(input["admin_username"].(string))

		if adminPassword := input["admin_password"].(string); adminPassword != "" {
			osProfile.AdminPassword = pointer.To(adminPassword)
		}

		if computerPrefix := input["computer_name_prefix"].(string); computerPrefix != "" {
			osProfile.ComputerNamePrefix = pointer.To(computerPrefix)
		}

		if secrets := input["secret"].([]interface{}); len(secrets) > 0 {
			osProfile.Secrets = expandLinuxSecretsVMSS(secrets)
		}

		if sshPublicKeys := expandSSHKeysVMSS(input["admin_ssh_key"].(*pluginsdk.Set).List()); len(sshPublicKeys) > 0 {
			if linConfig.Ssh == nil {
				linConfig.Ssh = &virtualmachinescalesets.SshConfiguration{}
			}
			linConfig.Ssh.PublicKeys = &sshPublicKeys
		}

		linConfig.DisablePasswordAuthentication = pointer.To(input["disable_password_authentication"].(bool))
		linConfig.ProvisionVMAgent = pointer.To(input["provision_vm_agent"].(bool))

		// Automatic VM Guest Patching
		patchSettings.AssessmentMode = pointer.To(virtualmachinescalesets.LinuxPatchAssessmentMode(input["patch_assessment_mode"].(string)))
		patchSettings.PatchMode = pointer.To(virtualmachinescalesets.LinuxVMGuestPatchMode(input["patch_mode"].(string)))
		linConfig.PatchSettings = &patchSettings
	}

	osProfile.LinuxConfiguration = &linConfig

	return &osProfile
}

func expandWindowsConfigurationAdditionalUnattendContent(input []interface{}) *[]virtualmachinescalesets.AdditionalUnattendContent {
	output := make([]virtualmachinescalesets.AdditionalUnattendContent, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		output = append(output, virtualmachinescalesets.AdditionalUnattendContent{
			SettingName: pointer.To(virtualmachinescalesets.SettingNames(raw["setting"].(string))),
			Content:     pointer.To(raw["content"].(string)),

			// no other possible values
			PassName:      pointer.To(virtualmachinescalesets.PassNamesOobeSystem),
			ComponentName: pointer.To(virtualmachinescalesets.ComponentNamesMicrosoftNegativeWindowsNegativeShellNegativeSetup),
		})
	}

	return &output
}

func ExpandOrchestratedVirtualMachineScaleSetNetworkInterface(input []interface{}) (*[]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration, error) {
	output := make([]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		dnsServers := utils.ExpandStringSlice(raw["dns_servers"].([]interface{}))

		ipConfigurations := make([]virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration, 0)
		ipConfigurationsRaw := raw["ip_configuration"].([]interface{})
		for _, configV := range ipConfigurationsRaw {
			configRaw := configV.(map[string]interface{})
			ipConfiguration, err := expandOrchestratedVirtualMachineScaleSetIPConfiguration(configRaw)
			if err != nil {
				return nil, err
			}

			ipConfigurations = append(ipConfigurations, *ipConfiguration)
		}

		config := virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration{
			Name: raw["name"].(string),
			Properties: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationProperties{
				DnsSettings: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
					DnsServers: dnsServers,
				},
				EnableAcceleratedNetworking: pointer.To(raw["enable_accelerated_networking"].(bool)),
				EnableIPForwarding:          pointer.To(raw["enable_ip_forwarding"].(bool)),
				IPConfigurations:            ipConfigurations,
				Primary:                     pointer.To(raw["primary"].(bool)),
			},
		}

		if nsgId := raw["network_security_group_id"].(string); nsgId != "" {
			config.Properties.NetworkSecurityGroup = &virtualmachinescalesets.SubResource{
				Id: pointer.To(nsgId),
			}
		}

		output = append(output, config)
	}

	return &output, nil
}

func expandOrchestratedVirtualMachineScaleSetIPConfiguration(raw map[string]interface{}) (*virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration, error) {
	applicationGatewayBackendAddressPoolIdsRaw := raw["application_gateway_backend_address_pool_ids"].(*pluginsdk.Set).List()
	applicationGatewayBackendAddressPoolIds := expandIDsToSubResources(applicationGatewayBackendAddressPoolIdsRaw)

	applicationSecurityGroupIdsRaw := raw["application_security_group_ids"].(*pluginsdk.Set).List()
	applicationSecurityGroupIds := expandIDsToSubResources(applicationSecurityGroupIdsRaw)

	loadBalancerBackendAddressPoolIdsRaw := raw["load_balancer_backend_address_pool_ids"].(*pluginsdk.Set).List()
	loadBalancerBackendAddressPoolIds := expandIDsToSubResources(loadBalancerBackendAddressPoolIdsRaw)

	primary := raw["primary"].(bool)
	version := virtualmachinescalesets.IPVersion(raw["version"].(string))
	if primary && version == virtualmachinescalesets.IPVersionIPvSix {
		return nil, fmt.Errorf("an IPv6 Primary IP Configuration is unsupported - instead add a IPv4 IP Configuration as the Primary and make the IPv6 IP Configuration the secondary")
	}

	ipConfiguration := virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration{
		Name: raw["name"].(string),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetIPConfigurationProperties{
			Primary:                               pointer.To(primary),
			PrivateIPAddressVersion:               pointer.To(version),
			ApplicationGatewayBackendAddressPools: applicationGatewayBackendAddressPoolIds,
			ApplicationSecurityGroups:             applicationSecurityGroupIds,
			LoadBalancerBackendAddressPools:       loadBalancerBackendAddressPoolIds,
			// LoadBalancerInboundNatPools removed per service team this attribute will never be used in VMSS Flex
		},
	}

	if subnetId := raw["subnet_id"].(string); subnetId != "" {
		ipConfiguration.Properties.Subnet = &virtualmachinescalesets.ApiEntityReference{
			Id: pointer.To(subnetId),
		}
	}

	publicIPConfigsRaw := raw["public_ip_address"].([]interface{})
	if len(publicIPConfigsRaw) > 0 {
		publicIPConfigRaw := publicIPConfigsRaw[0].(map[string]interface{})
		publicIPAddressConfig := expandOrchestratedVirtualMachineScaleSetPublicIPAddress(publicIPConfigRaw)
		ipConfiguration.Properties.PublicIPAddressConfiguration = publicIPAddressConfig
	}

	return &ipConfiguration, nil
}

func expandOrchestratedVirtualMachineScaleSetPublicIPAddress(raw map[string]interface{}) *virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration {
	ipTagsRaw := raw["ip_tag"].([]interface{})
	ipTags := make([]virtualmachinescalesets.VirtualMachineScaleSetIPTag, 0)
	for _, ipTagV := range ipTagsRaw {
		ipTagRaw := ipTagV.(map[string]interface{})
		ipTags = append(ipTags, virtualmachinescalesets.VirtualMachineScaleSetIPTag{
			Tag:       pointer.To(ipTagRaw["tag"].(string)),
			IPTagType: pointer.To(ipTagRaw["type"].(string)),
		})
	}

	publicIPAddressConfig := virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration{
		Name: raw["name"].(string),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
			IPTags: &ipTags,
		},
	}

	if domainNameLabel := raw["domain_name_label"].(string); domainNameLabel != "" {
		dns := &virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
			DomainNameLabel: domainNameLabel,
		}
		publicIPAddressConfig.Properties.DnsSettings = dns
	}

	if idleTimeout := raw["idle_timeout_in_minutes"].(int); idleTimeout > 0 {
		publicIPAddressConfig.Properties.IdleTimeoutInMinutes = pointer.To(int64(raw["idle_timeout_in_minutes"].(int)))
	}

	if publicIPPrefixID := raw["public_ip_prefix_id"].(string); publicIPPrefixID != "" {
		publicIPAddressConfig.Properties.PublicIPPrefix = &virtualmachinescalesets.SubResource{
			Id: pointer.To(publicIPPrefixID),
		}
	}

	if sku := raw["sku_name"].(string); sku != "" {
		v := expandOrchestratedVirtualMachineScaleSetPublicIPSku(sku)
		publicIPAddressConfig.Sku = v
	}

	if version := raw["version"].(string); version != "" {
		publicIPAddressConfig.Properties.PublicIPAddressVersion = pointer.To(virtualmachinescalesets.IPVersion(version))
	}

	return &publicIPAddressConfig
}

func ExpandOrchestratedVirtualMachineScaleSetNetworkInterfaceUpdate(input []interface{}) (*[]virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfiguration, error) {
	output := make([]virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		dnsServers := utils.ExpandStringSlice(raw["dns_servers"].([]interface{}))

		ipConfigurations := make([]virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfiguration, 0)
		ipConfigurationsRaw := raw["ip_configuration"].([]interface{})
		for _, configV := range ipConfigurationsRaw {
			configRaw := configV.(map[string]interface{})
			ipConfiguration, err := expandOrchestratedVirtualMachineScaleSetIPConfigurationUpdate(configRaw)
			if err != nil {
				return nil, err
			}

			ipConfigurations = append(ipConfigurations, *ipConfiguration)
		}

		config := virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfiguration{
			Name: pointer.To(raw["name"].(string)),
			Properties: &virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfigurationProperties{
				DnsSettings: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
					DnsServers: dnsServers,
				},
				EnableAcceleratedNetworking: pointer.To(raw["enable_accelerated_networking"].(bool)),
				EnableIPForwarding:          pointer.To(raw["enable_ip_forwarding"].(bool)),
				IPConfigurations:            &ipConfigurations,
				Primary:                     pointer.To(raw["primary"].(bool)),
			},
		}

		if nsgId := raw["network_security_group_id"].(string); nsgId != "" {
			config.Properties.NetworkSecurityGroup = &virtualmachinescalesets.SubResource{
				Id: pointer.To(nsgId),
			}
		}

		output = append(output, config)
	}

	return &output, nil
}

func expandOrchestratedVirtualMachineScaleSetIPConfigurationUpdate(raw map[string]interface{}) (*virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfiguration, error) {
	applicationGatewayBackendAddressPoolIdsRaw := raw["application_gateway_backend_address_pool_ids"].(*pluginsdk.Set).List()
	applicationGatewayBackendAddressPoolIds := expandIDsToSubResources(applicationGatewayBackendAddressPoolIdsRaw)

	applicationSecurityGroupIdsRaw := raw["application_security_group_ids"].(*pluginsdk.Set).List()
	applicationSecurityGroupIds := expandIDsToSubResources(applicationSecurityGroupIdsRaw)

	loadBalancerBackendAddressPoolIdsRaw := raw["load_balancer_backend_address_pool_ids"].(*pluginsdk.Set).List()
	loadBalancerBackendAddressPoolIds := expandIDsToSubResources(loadBalancerBackendAddressPoolIdsRaw)

	primary := raw["primary"].(bool)
	version := virtualmachinescalesets.IPVersion(raw["version"].(string))

	if primary && version == virtualmachinescalesets.IPVersionIPvSix {
		return nil, fmt.Errorf("an IPv6 Primary IP Configuration is unsupported - instead add a IPv4 IP Configuration as the Primary and make the IPv6 IP Configuration the secondary")
	}

	ipConfiguration := virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfiguration{
		Name: pointer.To(raw["name"].(string)),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfigurationProperties{
			Primary:                               pointer.To(primary),
			PrivateIPAddressVersion:               pointer.To(version),
			ApplicationGatewayBackendAddressPools: applicationGatewayBackendAddressPoolIds,
			ApplicationSecurityGroups:             applicationSecurityGroupIds,
			LoadBalancerBackendAddressPools:       loadBalancerBackendAddressPoolIds,
			// LoadBalancerInboundNatPools removed per service team this attribute will never be used in VMSS Flex
		},
	}

	if subnetId := raw["subnet_id"].(string); subnetId != "" {
		ipConfiguration.Properties.Subnet = &virtualmachinescalesets.ApiEntityReference{
			Id: pointer.To(subnetId),
		}
	}

	publicIPConfigsRaw := raw["public_ip_address"].([]interface{})
	if len(publicIPConfigsRaw) > 0 {
		publicIPConfigRaw := publicIPConfigsRaw[0].(map[string]interface{})
		publicIPAddressConfig := expandOrchestratedVirtualMachineScaleSetPublicIPAddressUpdate(publicIPConfigRaw)
		ipConfiguration.Properties.PublicIPAddressConfiguration = publicIPAddressConfig
	}

	return &ipConfiguration, nil
}

func expandOrchestratedVirtualMachineScaleSetPublicIPAddressUpdate(raw map[string]interface{}) *virtualmachinescalesets.VirtualMachineScaleSetUpdatePublicIPAddressConfiguration {
	publicIPAddressConfig := virtualmachinescalesets.VirtualMachineScaleSetUpdatePublicIPAddressConfiguration{
		Name:       pointer.To(raw["name"].(string)),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetUpdatePublicIPAddressConfigurationProperties{},
	}

	if domainNameLabel := raw["domain_name_label"].(string); domainNameLabel != "" {
		dns := &virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
			DomainNameLabel: domainNameLabel,
		}
		publicIPAddressConfig.Properties.DnsSettings = dns
	}

	if idleTimeout := raw["idle_timeout_in_minutes"].(int); idleTimeout > 0 {
		publicIPAddressConfig.Properties.IdleTimeoutInMinutes = pointer.To(int64(raw["idle_timeout_in_minutes"].(int)))
	}

	return &publicIPAddressConfig
}

func ExpandOrchestratedVirtualMachineScaleSetDataDisk(input []interface{}, ultraSSDEnabled bool) (*[]virtualmachinescalesets.VirtualMachineScaleSetDataDisk, error) {
	disks := make([]virtualmachinescalesets.VirtualMachineScaleSetDataDisk, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		storageAccountType := virtualmachinescalesets.StorageAccountTypes(raw["storage_account_type"].(string))
		disk := virtualmachinescalesets.VirtualMachineScaleSetDataDisk{
			Caching: pointer.To(virtualmachinescalesets.CachingTypes(raw["caching"].(string))),
			ManagedDisk: &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{
				StorageAccountType: pointer.To(storageAccountType),
			},
			WriteAcceleratorEnabled: pointer.To(raw["write_accelerator_enabled"].(bool)),
			CreateOption:            virtualmachinescalesets.DiskCreateOptionTypes(raw["create_option"].(string)),
		}

		if dataDiskSize := raw["disk_size_gb"].(int); dataDiskSize > 0 {
			disk.DiskSizeGB = pointer.To(int64(dataDiskSize))
		}

		if lun := raw["lun"].(int); lun >= 0 {
			disk.Lun = int64(lun)
		}

		if id := raw["disk_encryption_set_id"].(string); id != "" {
			disk.ManagedDisk.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
				Id: pointer.To(id),
			}
		}

		var iops int
		if diskIops, ok := raw["ultra_ssd_disk_iops_read_write"]; ok && diskIops.(int) > 0 {
			iops = diskIops.(int)
		}

		if iops > 0 && !ultraSSDEnabled && storageAccountType != virtualmachinescalesets.StorageAccountTypesPremiumVTwoLRS {
			return nil, fmt.Errorf("`ultra_ssd_disk_iops_read_write` can only be set when `storage_account_type` is set to `PremiumV2_LRS` or `UltraSSD_LRS`")
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if iops > 0 {
			disk.DiskIOPSReadWrite = pointer.To(int64(iops))
		}

		var mbps int
		if diskMbps, ok := raw["ultra_ssd_disk_mbps_read_write"]; ok && diskMbps.(int) > 0 {
			mbps = diskMbps.(int)
		}

		if mbps > 0 && !ultraSSDEnabled && storageAccountType != virtualmachinescalesets.StorageAccountTypesPremiumVTwoLRS {
			return nil, fmt.Errorf("`ultra_ssd_disk_mbps_read_write` can only be set when `storage_account_type` is set to `PremiumV2_LRS` or `UltraSSD_LRS`")
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if mbps > 0 {
			disk.DiskMBpsReadWrite = pointer.To(int64(mbps))
		}

		disks = append(disks, disk)
	}

	return &disks, nil
}

func ExpandOrchestratedVirtualMachineScaleSetOSDisk(input []interface{}, osType virtualmachinescalesets.OperatingSystemTypes) *virtualmachinescalesets.VirtualMachineScaleSetOSDisk {
	raw := input[0].(map[string]interface{})
	disk := virtualmachinescalesets.VirtualMachineScaleSetOSDisk{
		Caching: pointer.To(virtualmachinescalesets.CachingTypes(raw["caching"].(string))),
		ManagedDisk: &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: pointer.To(virtualmachinescalesets.StorageAccountTypes(raw["storage_account_type"].(string))),
		},
		WriteAcceleratorEnabled: pointer.To(raw["write_accelerator_enabled"].(bool)),

		// these have to be hard-coded so there's no point exposing them
		CreateOption: virtualmachinescalesets.DiskCreateOptionTypesFromImage,
		OsType:       pointer.To(osType),
	}

	if diskEncryptionSetId := raw["disk_encryption_set_id"].(string); diskEncryptionSetId != "" {
		disk.ManagedDisk.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
			Id: pointer.To(diskEncryptionSetId),
		}
	}

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = pointer.To(int64(osDiskSize))
	}

	if diffDiskSettingsRaw := raw["diff_disk_settings"].([]interface{}); len(diffDiskSettingsRaw) > 0 {
		diffDiskRaw := diffDiskSettingsRaw[0].(map[string]interface{})
		disk.DiffDiskSettings = &virtualmachinescalesets.DiffDiskSettings{
			Option:    pointer.To(virtualmachinescalesets.DiffDiskOptions(diffDiskRaw["option"].(string))),
			Placement: pointer.To(virtualmachinescalesets.DiffDiskPlacement(diffDiskRaw["placement"].(string))),
		}
	}

	return &disk
}

func ExpandOrchestratedVirtualMachineScaleSetOSDiskUpdate(input []interface{}) *virtualmachinescalesets.VirtualMachineScaleSetUpdateOSDisk {
	raw := input[0].(map[string]interface{})
	disk := virtualmachinescalesets.VirtualMachineScaleSetUpdateOSDisk{
		Caching: pointer.To(virtualmachinescalesets.CachingTypes(raw["caching"].(string))),
		ManagedDisk: &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: pointer.To(virtualmachinescalesets.StorageAccountTypes(raw["storage_account_type"].(string))),
		},
		WriteAcceleratorEnabled: pointer.To(raw["write_accelerator_enabled"].(bool)),
	}

	if diskEncryptionSetId := raw["disk_encryption_set_id"].(string); diskEncryptionSetId != "" {
		disk.ManagedDisk.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
			Id: pointer.To(diskEncryptionSetId),
		}
	}

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = pointer.To(int64(osDiskSize))
	}

	return &disk
}

func ExpandOrchestratedVirtualMachineScaleSetScheduledEventsProfile(input []interface{}) *virtualmachinescalesets.ScheduledEventsProfile {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	enabled := raw["enabled"].(bool)
	timeout := raw["timeout"].(string)

	return &virtualmachinescalesets.ScheduledEventsProfile{
		TerminateNotificationProfile: &virtualmachinescalesets.TerminateNotificationProfile{
			Enable:           &enabled,
			NotBeforeTimeout: &timeout,
		},
	}
}

func expandOrchestratedVirtualMachineScaleSetExtensions(input []interface{}) (extensionProfile *virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile, hasHealthExtension bool, err error) {
	extensionProfile = &virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile{}
	if len(input) == 0 {
		return nil, false, nil
	}

	extensions := make([]virtualmachinescalesets.VirtualMachineScaleSetExtension, 0)
	for _, v := range input {
		extensionRaw := v.(map[string]interface{})
		extension := virtualmachinescalesets.VirtualMachineScaleSetExtension{
			Name: pointer.To(extensionRaw["name"].(string)),
		}
		extensionType := extensionRaw["type"].(string)

		autoUpgradeMinorVersion, _ := extensionRaw["auto_upgrade_minor_version_enabled"].(bool)

		extensionProps := virtualmachinescalesets.VirtualMachineScaleSetExtensionProperties{
			Publisher:                pointer.To(extensionRaw["publisher"].(string)),
			Type:                     &extensionType,
			TypeHandlerVersion:       pointer.To(extensionRaw["type_handler_version"].(string)),
			AutoUpgradeMinorVersion:  pointer.To(autoUpgradeMinorVersion),
			ProvisionAfterExtensions: utils.ExpandStringSlice(extensionRaw["extensions_to_provision_after_vm_creation"].([]interface{})),
		}

		if extensionType == "ApplicationHealthLinux" || extensionType == "ApplicationHealthWindows" {
			hasHealthExtension = true
		}

		if val, ok := extensionRaw["failure_suppression_enabled"]; ok {
			extensionProps.SuppressFailures = pointer.To(val.(bool))
		}

		if forceUpdateTag := extensionRaw["force_extension_execution_on_change"]; forceUpdateTag != nil {
			extensionProps.ForceUpdateTag = pointer.To(forceUpdateTag.(string))
		}

		if val, ok := extensionRaw["settings"]; ok && val.(string) != "" {
			var result interface{}
			err := json.Unmarshal([]byte(val.(string)), &result)
			if err != nil {
				return nil, false, fmt.Errorf("unmarshaling `settings`: %+v", err)
			}
			extensionProps.Settings = pointer.To(result)
		}

		protectedSettingsFromKeyVault := expandProtectedSettingsFromKeyVaultVMSS(extensionRaw["protected_settings_from_key_vault"].([]interface{}))
		extensionProps.ProtectedSettingsFromKeyVault = (protectedSettingsFromKeyVault)

		if val, ok := extensionRaw["protected_settings"]; ok && val.(string) != "" {
			if protectedSettingsFromKeyVault != nil {
				return nil, false, fmt.Errorf("`protected_settings_from_key_vault` cannot be used with `protected_settings`")
			}

			var result interface{}
			err := json.Unmarshal([]byte(val.(string)), &result)
			if err != nil {
				return nil, false, fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
			}
			extensionProps.ProtectedSettings = pointer.To(result)
		}

		extension.Properties = &extensionProps
		extensions = append(extensions, extension)
	}
	extensionProfile.Extensions = &extensions

	return extensionProfile, hasHealthExtension, nil
}

func ExpandOrchestratedVirtualMachineScaleSetPriorityMixPolicy(input []interface{}) *virtualmachinescalesets.PriorityMixPolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &virtualmachinescalesets.PriorityMixPolicy{
		BaseRegularPriorityCount:           pointer.To(int64(raw["base_regular_count"].(int))),
		RegularPriorityPercentageAboveBase: pointer.To(int64(raw["regular_percentage_above_base"].(int))),
	}
}

func flattenOrchestratedVirtualMachineScaleSetExtensions(input *virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile, d *pluginsdk.ResourceData) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	if input == nil || input.Extensions == nil {
		return result, nil
	}

	// extensionsFromState holds the "extension" block, which is used to retrieve the "protected_settings" to fill it back the state,
	// since it is not returned from the API.
	extensionsFromState := map[string]map[string]interface{}{}
	if extSet, ok := d.GetOk("extension"); ok && extSet != nil {
		extensions := extSet.(*pluginsdk.Set).List()
		for _, ext := range extensions {
			if ext == nil {
				continue
			}
			ext := ext.(map[string]interface{})
			extensionsFromState[ext["name"].(string)] = ext
		}
	}

	for _, v := range *input.Extensions {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		autoUpgradeMinorVersion := false
		suppressFailures := false
		forceUpdateTag := ""
		provisionAfterExtension := make([]interface{}, 0)
		protectedSettings := ""
		var protectedSettingsFromKeyVault *virtualmachinescalesets.KeyVaultSecretReference
		extPublisher := ""
		extSettings := ""
		extType := ""
		extTypeVersion := ""

		if props := v.Properties; props != nil {
			if props.Publisher != nil {
				extPublisher = *props.Publisher
			}

			if props.Type != nil {
				extType = *props.Type
			}

			if props.TypeHandlerVersion != nil {
				extTypeVersion = *props.TypeHandlerVersion
			}

			if props.AutoUpgradeMinorVersion != nil {
				autoUpgradeMinorVersion = *props.AutoUpgradeMinorVersion
			}

			if props.ForceUpdateTag != nil {
				forceUpdateTag = *props.ForceUpdateTag
			}

			if props.SuppressFailures != nil {
				suppressFailures = *props.SuppressFailures
			}

			if props.ProvisionAfterExtensions != nil {
				provisionAfterExtension = utils.FlattenStringSlice(props.ProvisionAfterExtensions)
			}

			if props.Settings != nil {
				settings, err := json.Marshal(props.Settings)
				if err != nil {
					return nil, fmt.Errorf("unmarshaling `settings`: %+v", err)
				}

				extSettings = string(settings)
			}

			protectedSettingsFromKeyVault = props.ProtectedSettingsFromKeyVault
		}
		// protected_settings isn't returned, so we attempt to get it from state otherwise set to empty string
		if ext, ok := extensionsFromState[name]; ok {
			if protectedSettingsFromState, ok := ext["protected_settings"]; ok {
				if protectedSettingsFromState.(string) != "" && protectedSettingsFromState.(string) != "{}" {
					protectedSettings = protectedSettingsFromState.(string)
				}
			}
		}

		result = append(result, map[string]interface{}{
			"name":                                      name,
			"auto_upgrade_minor_version_enabled":        autoUpgradeMinorVersion,
			"force_extension_execution_on_change":       forceUpdateTag,
			"failure_suppression_enabled":               suppressFailures,
			"extensions_to_provision_after_vm_creation": provisionAfterExtension,
			"protected_settings":                        protectedSettings,
			"protected_settings_from_key_vault":         flattenProtectedSettingsFromKeyVaultVMSS(protectedSettingsFromKeyVault),
			"publisher":                                 extPublisher,
			"settings":                                  extSettings,
			"type":                                      extType,
			"type_handler_version":                      extTypeVersion,
		})
	}
	return result, nil
}

func FlattenOrchestratedVirtualMachineScaleSetIPConfiguration(input virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration) map[string]interface{} {
	var subnetId, version string
	var primary bool
	var publicIPAddresses, applicationGatewayBackendAddressPoolIds, applicationSecurityGroupIds, loadBalancerBackendAddressPoolIds []interface{}

	if props := input.Properties; props != nil {
		if props.Subnet != nil && props.Subnet.Id != nil {
			subnetId = *props.Subnet.Id
		}

		if props.Primary != nil {
			primary = *props.Primary
		}

		if props.PublicIPAddressConfiguration != nil {
			publicIPAddresses = append(publicIPAddresses, FlattenOrchestratedVirtualMachineScaleSetPublicIPAddress(*props.PublicIPAddressConfiguration))
		}

		version = string(pointer.From(props.PrivateIPAddressVersion))

		applicationGatewayBackendAddressPoolIds = flattenSubResourcesToIDs(props.ApplicationGatewayBackendAddressPools)
		applicationSecurityGroupIds = flattenSubResourcesToIDs(props.ApplicationSecurityGroups)
		loadBalancerBackendAddressPoolIds = flattenSubResourcesToIDs(props.LoadBalancerBackendAddressPools)

	}

	return map[string]interface{}{
		"name":              input.Name,
		"primary":           primary,
		"public_ip_address": publicIPAddresses,
		"subnet_id":         subnetId,
		"version":           version,
		"application_gateway_backend_address_pool_ids": applicationGatewayBackendAddressPoolIds,
		"application_security_group_ids":               applicationSecurityGroupIds,
		"load_balancer_backend_address_pool_ids":       loadBalancerBackendAddressPoolIds,
		// load_balancer_inbound_nat_rules_ids removed per service team this attribute will never be used in VMSS Flex
	}
}

func FlattenOrchestratedVirtualMachineScaleSetPublicIPAddress(input virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration) map[string]interface{} {
	ipTags := make([]interface{}, 0)
	var domainNameLabel, publicIPPrefixId, version, sku string
	var idleTimeoutInMinutes int

	if props := input.Properties; props != nil {
		if props.IPTags != nil {
			for _, rawTag := range *props.IPTags {
				var tag, tagType string

				if rawTag.IPTagType != nil {
					tagType = *rawTag.IPTagType
				}

				if rawTag.Tag != nil {
					tag = *rawTag.Tag
				}

				ipTags = append(ipTags, map[string]interface{}{
					"tag":  tag,
					"type": tagType,
				})
			}
		}

		if props.DnsSettings != nil {
			domainNameLabel = props.DnsSettings.DomainNameLabel
		}

		if props.PublicIPPrefix != nil && props.PublicIPPrefix.Id != nil {
			publicIPPrefixId = *props.PublicIPPrefix.Id
		}

		if props.IdleTimeoutInMinutes != nil {
			idleTimeoutInMinutes = int(*props.IdleTimeoutInMinutes)
		}

		if props.PublicIPAddressVersion != nil {
			version = string(pointer.From(props.PublicIPAddressVersion))
		}

		if input.Sku != nil && input.Sku.Name != nil && input.Sku.Tier != nil {
			sku = flattenOrchestratedVirtualMachineScaleSetPublicIPSku(input.Sku)
		}

	}

	return map[string]interface{}{
		"name":                    input.Name,
		"domain_name_label":       domainNameLabel,
		"idle_timeout_in_minutes": idleTimeoutInMinutes,
		"ip_tag":                  ipTags,
		"public_ip_prefix_id":     publicIPPrefixId,
		"sku_name":                sku,
		"version":                 version,
	}
}

func flattenOrchestratedVirtualMachineScaleSetWindowsConfiguration(input *virtualmachinescalesets.VirtualMachineScaleSetOSProfile, d *pluginsdk.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})
	winConfig := input.WindowsConfiguration
	patchSettings := winConfig.PatchSettings

	if v := input.AdminUsername; v != nil {
		output["admin_username"] = *v
	}

	if v := d.Get("os_profile").([]interface{}); len(v) > 0 {
		osProfile := v[0].(map[string]interface{})
		if winConfigRaw := osProfile["windows_configuration"].([]interface{}); len(winConfigRaw) > 0 {
			winCfg := winConfigRaw[0].(map[string]interface{})
			output["admin_password"] = winCfg["admin_password"].(string)
		}
	}

	if v := input.ComputerNamePrefix; v != nil {
		output["computer_name_prefix"] = *v
	}

	if v := winConfig.AdditionalUnattendContent; v != nil {
		output["additional_unattend_content"] = flattenWindowsConfigurationAdditionalUnattendContent(winConfig, d)
	}

	if v := winConfig.EnableAutomaticUpdates; v != nil {
		output["enable_automatic_updates"] = *v
	}

	if v := winConfig.ProvisionVMAgent; v != nil {
		output["provision_vm_agent"] = *v
	}

	if v := input.Secrets; v != nil {
		output["secret"] = flattenWindowsSecretsVMSS(v)
	}

	if v := winConfig.WinRM; v != nil {
		output["winrm_listener"] = flattenWinRMListenerVMSS(winConfig.WinRM)
	}

	if v := winConfig.TimeZone; v != nil {
		output["timezone"] = v
	}

	output["patch_mode"] = string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByOS)
	output["patch_assessment_mode"] = string(virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform)
	output["hotpatching_enabled"] = false

	if patchSettings != nil {
		output["patch_mode"] = pointer.From(patchSettings.PatchMode)
		output["patch_assessment_mode"] = pointer.From(patchSettings.AssessmentMode)

		if v := patchSettings.EnableHotpatching; v != nil {
			output["hotpatching_enabled"] = *v
		}
	}

	return []interface{}{output}
}

func flattenOrchestratedVirtualMachineScaleSetLinuxConfiguration(input *virtualmachinescalesets.VirtualMachineScaleSetOSProfile, d *pluginsdk.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})
	linConfig := input.LinuxConfiguration

	if v := input.AdminUsername; v != nil {
		output["admin_username"] = *v
	}

	if v := d.Get("os_profile").([]interface{}); len(v) > 0 {
		osProfile := v[0].(map[string]interface{})
		if linConfigRaw := osProfile["linux_configuration"].([]interface{}); len(linConfigRaw) > 0 {
			linCfg := linConfigRaw[0].(map[string]interface{})
			output["admin_password"] = linCfg["admin_password"].(string)
		}
	}

	if v := input.AdminPassword; v != nil {
		output["admin_password"] = *v
	}

	if v := linConfig.Ssh; v != nil {
		if sshKeys, _ := flattenSSHKeysVMSS(v); sshKeys != nil {
			output["admin_ssh_key"] = pluginsdk.NewSet(SSHKeySchemaHash, *sshKeys)
		}
	}

	if v := input.ComputerNamePrefix; v != nil {
		output["computer_name_prefix"] = *v
	}

	if v := linConfig.DisablePasswordAuthentication; v != nil {
		output["disable_password_authentication"] = *v
	}

	if v := linConfig.PatchSettings; v != nil {
		output["patch_mode"] = v.PatchMode
		output["patch_assessment_mode"] = v.AssessmentMode
	}

	if v := linConfig.ProvisionVMAgent; v != nil {
		output["provision_vm_agent"] = *v
	}

	if v := input.Secrets; v != nil {
		output["secret"] = flattenLinuxSecretsVMSS(v)
	}

	return []interface{}{output}
}

func FlattenOrchestratedVirtualMachineScaleSetNetworkInterface(input *[]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		var networkSecurityGroupId string
		var enableAcceleratedNetworking, enableIPForwarding, primary bool
		var dnsServers []interface{}
		var ipConfigurations []interface{}
		if props := v.Properties; props != nil {
			if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.Id != nil {
				networkSecurityGroupId = *props.NetworkSecurityGroup.Id
			}

			if props.EnableAcceleratedNetworking != nil {
				enableAcceleratedNetworking = *props.EnableAcceleratedNetworking
			}
			if props.EnableIPForwarding != nil {
				enableIPForwarding = *props.EnableIPForwarding
			}
			if props.Primary != nil {
				primary = *props.Primary
			}

			if settings := props.DnsSettings; settings != nil {
				dnsServers = utils.FlattenStringSlice(props.DnsSettings.DnsServers)
			}

			if len(props.IPConfigurations) != 0 {
				for _, configRaw := range props.IPConfigurations {
					config := FlattenOrchestratedVirtualMachineScaleSetIPConfiguration(configRaw)
					ipConfigurations = append(ipConfigurations, config)
				}
			}

		}

		results = append(results, map[string]interface{}{
			"name":                          v.Name,
			"dns_servers":                   dnsServers,
			"enable_accelerated_networking": enableAcceleratedNetworking,
			"enable_ip_forwarding":          enableIPForwarding,
			"ip_configuration":              ipConfigurations,
			"network_security_group_id":     networkSecurityGroupId,
			"primary":                       primary,
		})
	}

	return results
}

func FlattenOrchestratedVirtualMachineScaleSetDataDisk(input *[]virtualmachinescalesets.VirtualMachineScaleSetDataDisk) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		diskSizeGb := 0
		if v.DiskSizeGB != nil && *v.DiskSizeGB != 0 {
			diskSizeGb = int(*v.DiskSizeGB)
		}

		storageAccountType := ""
		diskEncryptionSetId := ""
		if v.ManagedDisk != nil {
			storageAccountType = string(pointer.From(v.ManagedDisk.StorageAccountType))
			if v.ManagedDisk.DiskEncryptionSet != nil && v.ManagedDisk.DiskEncryptionSet.Id != nil {
				diskEncryptionSetId = *v.ManagedDisk.DiskEncryptionSet.Id
			}
		}

		writeAcceleratorEnabled := false
		if v.WriteAcceleratorEnabled != nil {
			writeAcceleratorEnabled = *v.WriteAcceleratorEnabled
		}

		iops := 0
		if v.DiskIOPSReadWrite != nil {
			iops = int(*v.DiskIOPSReadWrite)
		}

		mbps := 0
		if v.DiskMBpsReadWrite != nil {
			mbps = int(*v.DiskMBpsReadWrite)
		}

		output = append(output, map[string]interface{}{
			"caching":                        pointer.From(v.Caching),
			"create_option":                  string(v.CreateOption),
			"lun":                            v.Lun,
			"disk_encryption_set_id":         diskEncryptionSetId,
			"disk_size_gb":                   diskSizeGb,
			"storage_account_type":           storageAccountType,
			"write_accelerator_enabled":      writeAcceleratorEnabled,
			"ultra_ssd_disk_iops_read_write": iops,
			"ultra_ssd_disk_mbps_read_write": mbps,
		})
	}

	return output
}

func flattenWindowsConfigurationAdditionalUnattendContent(input *virtualmachinescalesets.WindowsConfiguration, d *pluginsdk.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	existing := make([]interface{}, 0)
	if v, ok := d.GetOk("os_profile.0.windows_configuration.0.additional_unattend_content"); ok {
		existing = v.([]interface{})
	}

	output := make([]interface{}, 0)
	for i, v := range *input.AdditionalUnattendContent {
		// content isn't returned by the API since it's sensitive data so we need to pull from the state file.
		content := ""
		if len(existing) > i {
			existingVal := existing[i]
			existingRaw, ok := existingVal.(map[string]interface{})
			if ok {
				contentRaw, ok := existingRaw["content"]
				if ok {
					content = contentRaw.(string)
				}
			}
		}
		output = append(output, map[string]interface{}{
			"content": content,
			"setting": pointer.From(v.SettingName),
		})
	}

	return output
}

func FlattenOrchestratedVirtualMachineScaleSetOSDisk(input *virtualmachinescalesets.VirtualMachineScaleSetOSDisk) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	diffDiskSettings := make([]interface{}, 0)
	if input.DiffDiskSettings != nil {
		diffDiskSettings = append(diffDiskSettings, map[string]interface{}{
			"option":    pointer.From(input.DiffDiskSettings.Option),
			"placement": pointer.From(input.DiffDiskSettings.Placement),
		})
	}

	diskSizeGb := 0
	if input.DiskSizeGB != nil && *input.DiskSizeGB != 0 {
		diskSizeGb = int(*input.DiskSizeGB)
	}

	storageAccountType := ""
	diskEncryptionSetId := ""
	if input.ManagedDisk != nil {
		storageAccountType = string(pointer.From(input.ManagedDisk.StorageAccountType))
		if input.ManagedDisk.DiskEncryptionSet != nil && input.ManagedDisk.DiskEncryptionSet.Id != nil {
			diskEncryptionSetId = *input.ManagedDisk.DiskEncryptionSet.Id
		}
	}

	writeAcceleratorEnabled := false
	if input.WriteAcceleratorEnabled != nil {
		writeAcceleratorEnabled = *input.WriteAcceleratorEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"caching":                   pointer.From(input.Caching),
			"disk_size_gb":              diskSizeGb,
			"diff_disk_settings":        diffDiskSettings,
			"storage_account_type":      storageAccountType,
			"write_accelerator_enabled": writeAcceleratorEnabled,
			"disk_encryption_set_id":    diskEncryptionSetId,
		},
	}
}

func FlattenOrchestratedVirtualMachineScaleSetScheduledEventsProfile(input *virtualmachinescalesets.ScheduledEventsProfile) []interface{} {
	// if enabled is set to false, there will be no ScheduledEventsProfile in response, to avoid plan non empty when
	// a user explicitly set enabled to false, we need to assign a default block to this field

	enabled := false
	if input != nil && input.TerminateNotificationProfile != nil && input.TerminateNotificationProfile.Enable != nil {
		enabled = *input.TerminateNotificationProfile.Enable
	}

	timeout := "PT5M"
	if input != nil && input.TerminateNotificationProfile != nil && input.TerminateNotificationProfile.NotBeforeTimeout != nil {
		timeout = *input.TerminateNotificationProfile.NotBeforeTimeout
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": enabled,
			"timeout": timeout,
		},
	}
}

func FlattenOrchestratedVirtualMachineScaleSetPriorityMixPolicy(input *virtualmachinescalesets.PriorityMixPolicy) []interface{} {

	baseRegularPriorityCount := int64(0)
	if input != nil && input.BaseRegularPriorityCount != nil {
		baseRegularPriorityCount = *input.BaseRegularPriorityCount
	}

	regularPriorityPercentageAboveBase := int64(0)
	if input != nil && input.RegularPriorityPercentageAboveBase != nil {
		regularPriorityPercentageAboveBase = *input.RegularPriorityPercentageAboveBase
	}

	return []interface{}{
		map[string]interface{}{
			"base_regular_count":            baseRegularPriorityCount,
			"regular_percentage_above_base": regularPriorityPercentageAboveBase,
		},
	}
}
