// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationsecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
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

				// I am only commenting this out as this is going to be supported in the next release of the API in October 2021
				// "additional_unattend_content": additionalUnattendContentSchema(),

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
					Default:  string(compute.WindowsPatchAssessmentModeImageDefault),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.WindowsPatchAssessmentModeAutomaticByPlatform),
						string(compute.WindowsPatchAssessmentModeImageDefault),
					}, false),
				},

				"patch_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(compute.WindowsVMGuestPatchModeAutomaticByOS),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.WindowsVMGuestPatchModeAutomaticByOS),
						string(compute.WindowsVMGuestPatchModeAutomaticByPlatform),
						string(compute.WindowsVMGuestPatchModeManual),
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
					Default:  string(compute.LinuxPatchAssessmentModeImageDefault),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.LinuxPatchAssessmentModeAutomaticByPlatform),
						string(compute.LinuxPatchAssessmentModeImageDefault),
					}, false),
				},

				"patch_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(compute.LinuxVMGuestPatchModeImageDefault),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.LinuxVMGuestPatchModeImageDefault),
						string(compute.LinuxVMGuestPatchModeAutomaticByPlatform),
					}, false),
				},

				"secret": linuxSecretSchema(),
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetExtensionsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
					ValidateFunc: networkValidate.NetworkSecurityGroupID,
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
					Default:  string(compute.IPVersionIPv4),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.IPVersionIPv4),
						string(compute.IPVersionIPv6),
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
					ValidateFunc: networkValidate.PublicIpPrefixID,
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
					Default:  string(compute.IPVersionIPv4),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.IPVersionIPv4),
						string(compute.IPVersionIPv6),
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
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
					}, false),
				},

				"create_option": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.DiskCreateOptionTypesEmpty),
						string(compute.DiskCreateOptionTypesFromImage),
					}, false),
					Default: string(compute.DiskCreateOptionTypesEmpty),
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
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"lun": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 2000), // TODO: confirm upper bounds
				},

				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesPremiumV2LRS),
						string(compute.StorageAccountTypesPremiumZRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
						string(compute.StorageAccountTypesStandardSSDZRS),
						string(compute.StorageAccountTypesUltraSSDLRS),
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
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
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
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesPremiumZRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
						string(compute.StorageAccountTypesStandardSSDZRS),
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
									string(compute.DiffDiskOptionsLocal),
								}, false),
							},
							"placement": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ForceNew: true,
								Default:  string(compute.DiffDiskPlacementCacheDisk),
								ValidateFunc: validation.StringInSlice([]string{
									string(compute.DiffDiskPlacementCacheDisk),
									string(compute.DiffDiskPlacementResourceDisk),
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

func FlattenOrchestratedVirtualMachineScaleSetOSProfile(input *compute.VirtualMachineScaleSetOSProfile, d *pluginsdk.ResourceData) []interface{} {
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

func ExpandOrchestratedVirtualMachineScaleSetAdditionalCapabilities(input []interface{}) *compute.AdditionalCapabilities {
	capabilities := compute.AdditionalCapabilities{}

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})

		capabilities.UltraSSDEnabled = utils.Bool(raw["ultra_ssd_enabled"].(bool))
	}

	return &capabilities
}

func FlattenOrchestratedVirtualMachineScaleSetAdditionalCapabilities(input *compute.AdditionalCapabilities) []interface{} {
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

func expandOrchestratedVirtualMachineScaleSetOsProfileWithWindowsConfiguration(input map[string]interface{}, customData string) *compute.VirtualMachineScaleSetOSProfile {
	osProfile := compute.VirtualMachineScaleSetOSProfile{}
	winConfig := compute.WindowsConfiguration{}
	patchSettings := compute.PatchSettings{}

	if len(input) > 0 {
		osProfile.CustomData = utils.String(customData)
		osProfile.AdminUsername = utils.String(input["admin_username"].(string))
		osProfile.AdminPassword = utils.String(input["admin_password"].(string))

		if computerPrefix := input["computer_name_prefix"].(string); computerPrefix != "" {
			osProfile.ComputerNamePrefix = utils.String(computerPrefix)
		}

		if secrets := input["secret"].([]interface{}); len(secrets) > 0 {
			osProfile.Secrets = expandWindowsSecrets(secrets)
		}

		// I am only commenting this out as this is going to be supported in the next release of the API in October 2021
		// winConfig.AdditionalUnattendContent = expandWindowsConfigurationAdditionalUnattendContent(input["additional_unattend_content"].([]interface{}))
		winConfig.EnableAutomaticUpdates = utils.Bool(input["enable_automatic_updates"].(bool))
		winConfig.ProvisionVMAgent = utils.Bool(input["provision_vm_agent"].(bool))
		winRmListenersRaw := input["winrm_listener"].(*pluginsdk.Set).List()
		winConfig.WinRM = expandWinRMListener(winRmListenersRaw)

		// Automatic VM Guest Patching and Hotpatching settings
		patchSettings.AssessmentMode = compute.WindowsPatchAssessmentMode(input["patch_assessment_mode"].(string))
		patchSettings.PatchMode = compute.WindowsVMGuestPatchMode(input["patch_mode"].(string))
		patchSettings.EnableHotpatching = utils.Bool(input["hotpatching_enabled"].(bool))
		winConfig.PatchSettings = &patchSettings

		// due to a change in RP behavor, it will now throw and error if we pass an empty
		// string add check to only include it if it is actually defined in the config file
		timeZone := input["timezone"].(string)
		if timeZone != "" {
			winConfig.TimeZone = utils.String(timeZone)
		}
	}

	osProfile.WindowsConfiguration = &winConfig

	return &osProfile
}

func expandOrchestratedVirtualMachineScaleSetOsProfileWithLinuxConfiguration(input map[string]interface{}, customData string) *compute.VirtualMachineScaleSetOSProfile {
	osProfile := compute.VirtualMachineScaleSetOSProfile{}
	linConfig := compute.LinuxConfiguration{}
	patchSettings := compute.LinuxPatchSettings{}

	if len(input) > 0 {
		osProfile.CustomData = utils.String(customData)
		osProfile.AdminUsername = utils.String(input["admin_username"].(string))

		if adminPassword := input["admin_password"].(string); adminPassword != "" {
			osProfile.AdminPassword = utils.String(adminPassword)
		}

		if computerPrefix := input["computer_name_prefix"].(string); computerPrefix != "" {
			osProfile.ComputerNamePrefix = utils.String(computerPrefix)
		}

		if secrets := input["secret"].([]interface{}); len(secrets) > 0 {
			osProfile.Secrets = expandLinuxSecrets(secrets)
		}

		if sshPublicKeys := ExpandSSHKeys(input["admin_ssh_key"].(*pluginsdk.Set).List()); len(sshPublicKeys) > 0 {
			if linConfig.SSH == nil {
				linConfig.SSH = &compute.SSHConfiguration{}
			}
			linConfig.SSH.PublicKeys = &sshPublicKeys
		}

		linConfig.DisablePasswordAuthentication = utils.Bool(input["disable_password_authentication"].(bool))
		linConfig.ProvisionVMAgent = utils.Bool(input["provision_vm_agent"].(bool))

		// Automatic VM Guest Patching
		patchSettings.AssessmentMode = compute.LinuxPatchAssessmentMode(input["patch_assessment_mode"].(string))
		patchSettings.PatchMode = compute.LinuxVMGuestPatchMode(input["patch_mode"].(string))
		linConfig.PatchSettings = &patchSettings
	}

	osProfile.LinuxConfiguration = &linConfig

	return &osProfile
}

// I am only commenting this out as this is going to be supported in the next release of the API version 2021-10-01
// func expandWindowsConfigurationAdditionalUnattendContent(input []interface{}) *[]compute.AdditionalUnattendContent {
// 	output := make([]compute.AdditionalUnattendContent, 0)

// 	for _, v := range input {
// 		raw := v.(map[string]interface{})

// 		output = append(output, compute.AdditionalUnattendContent{
// 			SettingName: compute.SettingNames(raw["setting"].(string)),
// 			Content:     utils.String(raw["content"].(string)),

// 			// no other possible values
// 			PassName:      compute.PassNamesOobeSystem,
// 			ComponentName: compute.ComponentNamesMicrosoftWindowsShellSetup,
// 		})
// 	}

// 	return &output
// }

func ExpandOrchestratedVirtualMachineScaleSetNetworkInterface(input []interface{}) (*[]compute.VirtualMachineScaleSetNetworkConfiguration, error) {
	output := make([]compute.VirtualMachineScaleSetNetworkConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		dnsServers := utils.ExpandStringSlice(raw["dns_servers"].([]interface{}))

		ipConfigurations := make([]compute.VirtualMachineScaleSetIPConfiguration, 0)
		ipConfigurationsRaw := raw["ip_configuration"].([]interface{})
		for _, configV := range ipConfigurationsRaw {
			configRaw := configV.(map[string]interface{})
			ipConfiguration, err := expandOrchestratedVirtualMachineScaleSetIPConfiguration(configRaw)
			if err != nil {
				return nil, err
			}

			ipConfigurations = append(ipConfigurations, *ipConfiguration)
		}

		config := compute.VirtualMachineScaleSetNetworkConfiguration{
			Name: utils.String(raw["name"].(string)),
			VirtualMachineScaleSetNetworkConfigurationProperties: &compute.VirtualMachineScaleSetNetworkConfigurationProperties{
				DNSSettings: &compute.VirtualMachineScaleSetNetworkConfigurationDNSSettings{
					DNSServers: dnsServers,
				},
				EnableAcceleratedNetworking: utils.Bool(raw["enable_accelerated_networking"].(bool)),
				EnableIPForwarding:          utils.Bool(raw["enable_ip_forwarding"].(bool)),
				IPConfigurations:            &ipConfigurations,
				Primary:                     utils.Bool(raw["primary"].(bool)),
			},
		}

		if nsgId := raw["network_security_group_id"].(string); nsgId != "" {
			config.VirtualMachineScaleSetNetworkConfigurationProperties.NetworkSecurityGroup = &compute.SubResource{
				ID: utils.String(nsgId),
			}
		}

		output = append(output, config)
	}

	return &output, nil
}

func expandOrchestratedVirtualMachineScaleSetIPConfiguration(raw map[string]interface{}) (*compute.VirtualMachineScaleSetIPConfiguration, error) {
	applicationGatewayBackendAddressPoolIdsRaw := raw["application_gateway_backend_address_pool_ids"].(*pluginsdk.Set).List()
	applicationGatewayBackendAddressPoolIds := expandIDsToSubResources(applicationGatewayBackendAddressPoolIdsRaw)

	applicationSecurityGroupIdsRaw := raw["application_security_group_ids"].(*pluginsdk.Set).List()
	applicationSecurityGroupIds := expandIDsToSubResources(applicationSecurityGroupIdsRaw)

	loadBalancerBackendAddressPoolIdsRaw := raw["load_balancer_backend_address_pool_ids"].(*pluginsdk.Set).List()
	loadBalancerBackendAddressPoolIds := expandIDsToSubResources(loadBalancerBackendAddressPoolIdsRaw)

	primary := raw["primary"].(bool)
	version := compute.IPVersion(raw["version"].(string))
	if primary && version == compute.IPVersionIPv6 {
		return nil, fmt.Errorf("an IPv6 Primary IP Configuration is unsupported - instead add a IPv4 IP Configuration as the Primary and make the IPv6 IP Configuration the secondary")
	}

	ipConfiguration := compute.VirtualMachineScaleSetIPConfiguration{
		Name: utils.String(raw["name"].(string)),
		VirtualMachineScaleSetIPConfigurationProperties: &compute.VirtualMachineScaleSetIPConfigurationProperties{
			Primary:                               utils.Bool(primary),
			PrivateIPAddressVersion:               version,
			ApplicationGatewayBackendAddressPools: applicationGatewayBackendAddressPoolIds,
			ApplicationSecurityGroups:             applicationSecurityGroupIds,
			LoadBalancerBackendAddressPools:       loadBalancerBackendAddressPoolIds,
			// LoadBalancerInboundNatPools removed per service team this attribute will never be used in VMSS Flex
		},
	}

	if subnetId := raw["subnet_id"].(string); subnetId != "" {
		ipConfiguration.VirtualMachineScaleSetIPConfigurationProperties.Subnet = &compute.APIEntityReference{
			ID: utils.String(subnetId),
		}
	}

	publicIPConfigsRaw := raw["public_ip_address"].([]interface{})
	if len(publicIPConfigsRaw) > 0 {
		publicIPConfigRaw := publicIPConfigsRaw[0].(map[string]interface{})
		publicIPAddressConfig := expandOrchestratedVirtualMachineScaleSetPublicIPAddress(publicIPConfigRaw)
		ipConfiguration.VirtualMachineScaleSetIPConfigurationProperties.PublicIPAddressConfiguration = publicIPAddressConfig
	}

	return &ipConfiguration, nil
}

func expandOrchestratedVirtualMachineScaleSetPublicIPAddress(raw map[string]interface{}) *compute.VirtualMachineScaleSetPublicIPAddressConfiguration {
	ipTagsRaw := raw["ip_tag"].([]interface{})
	ipTags := make([]compute.VirtualMachineScaleSetIPTag, 0)
	for _, ipTagV := range ipTagsRaw {
		ipTagRaw := ipTagV.(map[string]interface{})
		ipTags = append(ipTags, compute.VirtualMachineScaleSetIPTag{
			Tag:       utils.String(ipTagRaw["tag"].(string)),
			IPTagType: utils.String(ipTagRaw["type"].(string)),
		})
	}

	publicIPAddressConfig := compute.VirtualMachineScaleSetPublicIPAddressConfiguration{
		Name: utils.String(raw["name"].(string)),
		VirtualMachineScaleSetPublicIPAddressConfigurationProperties: &compute.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
			IPTags: &ipTags,
		},
	}

	if domainNameLabel := raw["domain_name_label"].(string); domainNameLabel != "" {
		dns := &compute.VirtualMachineScaleSetPublicIPAddressConfigurationDNSSettings{
			DomainNameLabel: utils.String(domainNameLabel),
		}
		publicIPAddressConfig.VirtualMachineScaleSetPublicIPAddressConfigurationProperties.DNSSettings = dns
	}

	if idleTimeout := raw["idle_timeout_in_minutes"].(int); idleTimeout > 0 {
		publicIPAddressConfig.VirtualMachineScaleSetPublicIPAddressConfigurationProperties.IdleTimeoutInMinutes = utils.Int32(int32(raw["idle_timeout_in_minutes"].(int)))
	}

	if publicIPPrefixID := raw["public_ip_prefix_id"].(string); publicIPPrefixID != "" {
		publicIPAddressConfig.VirtualMachineScaleSetPublicIPAddressConfigurationProperties.PublicIPPrefix = &compute.SubResource{
			ID: utils.String(publicIPPrefixID),
		}
	}

	if sku := raw["sku_name"].(string); sku != "" {
		v := expandOrchestratedVirtualMachineScaleSetPublicIPSku(sku)
		publicIPAddressConfig.Sku = v
	}

	if version := raw["version"].(string); version != "" {
		publicIPAddressConfig.PublicIPAddressVersion = compute.IPVersion(version)
	}

	return &publicIPAddressConfig
}

func ExpandOrchestratedVirtualMachineScaleSetNetworkInterfaceUpdate(input []interface{}) (*[]compute.VirtualMachineScaleSetUpdateNetworkConfiguration, error) {
	output := make([]compute.VirtualMachineScaleSetUpdateNetworkConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		dnsServers := utils.ExpandStringSlice(raw["dns_servers"].([]interface{}))

		ipConfigurations := make([]compute.VirtualMachineScaleSetUpdateIPConfiguration, 0)
		ipConfigurationsRaw := raw["ip_configuration"].([]interface{})
		for _, configV := range ipConfigurationsRaw {
			configRaw := configV.(map[string]interface{})
			ipConfiguration, err := expandOrchestratedVirtualMachineScaleSetIPConfigurationUpdate(configRaw)
			if err != nil {
				return nil, err
			}

			ipConfigurations = append(ipConfigurations, *ipConfiguration)
		}

		config := compute.VirtualMachineScaleSetUpdateNetworkConfiguration{
			Name: utils.String(raw["name"].(string)),
			VirtualMachineScaleSetUpdateNetworkConfigurationProperties: &compute.VirtualMachineScaleSetUpdateNetworkConfigurationProperties{
				DNSSettings: &compute.VirtualMachineScaleSetNetworkConfigurationDNSSettings{
					DNSServers: dnsServers,
				},
				EnableAcceleratedNetworking: utils.Bool(raw["enable_accelerated_networking"].(bool)),
				EnableIPForwarding:          utils.Bool(raw["enable_ip_forwarding"].(bool)),
				IPConfigurations:            &ipConfigurations,
				Primary:                     utils.Bool(raw["primary"].(bool)),
			},
		}

		if nsgId := raw["network_security_group_id"].(string); nsgId != "" {
			config.VirtualMachineScaleSetUpdateNetworkConfigurationProperties.NetworkSecurityGroup = &compute.SubResource{
				ID: utils.String(nsgId),
			}
		}

		output = append(output, config)
	}

	return &output, nil
}

func expandOrchestratedVirtualMachineScaleSetIPConfigurationUpdate(raw map[string]interface{}) (*compute.VirtualMachineScaleSetUpdateIPConfiguration, error) {
	applicationGatewayBackendAddressPoolIdsRaw := raw["application_gateway_backend_address_pool_ids"].(*pluginsdk.Set).List()
	applicationGatewayBackendAddressPoolIds := expandIDsToSubResources(applicationGatewayBackendAddressPoolIdsRaw)

	applicationSecurityGroupIdsRaw := raw["application_security_group_ids"].(*pluginsdk.Set).List()
	applicationSecurityGroupIds := expandIDsToSubResources(applicationSecurityGroupIdsRaw)

	loadBalancerBackendAddressPoolIdsRaw := raw["load_balancer_backend_address_pool_ids"].(*pluginsdk.Set).List()
	loadBalancerBackendAddressPoolIds := expandIDsToSubResources(loadBalancerBackendAddressPoolIdsRaw)

	primary := raw["primary"].(bool)
	version := compute.IPVersion(raw["version"].(string))

	if primary && version == compute.IPVersionIPv6 {
		return nil, fmt.Errorf("an IPv6 Primary IP Configuration is unsupported - instead add a IPv4 IP Configuration as the Primary and make the IPv6 IP Configuration the secondary")
	}

	ipConfiguration := compute.VirtualMachineScaleSetUpdateIPConfiguration{
		Name: utils.String(raw["name"].(string)),
		VirtualMachineScaleSetUpdateIPConfigurationProperties: &compute.VirtualMachineScaleSetUpdateIPConfigurationProperties{
			Primary:                               utils.Bool(primary),
			PrivateIPAddressVersion:               version,
			ApplicationGatewayBackendAddressPools: applicationGatewayBackendAddressPoolIds,
			ApplicationSecurityGroups:             applicationSecurityGroupIds,
			LoadBalancerBackendAddressPools:       loadBalancerBackendAddressPoolIds,
			// LoadBalancerInboundNatPools removed per service team this attribute will never be used in VMSS Flex
		},
	}

	if subnetId := raw["subnet_id"].(string); subnetId != "" {
		ipConfiguration.VirtualMachineScaleSetUpdateIPConfigurationProperties.Subnet = &compute.APIEntityReference{
			ID: utils.String(subnetId),
		}
	}

	publicIPConfigsRaw := raw["public_ip_address"].([]interface{})
	if len(publicIPConfigsRaw) > 0 {
		publicIPConfigRaw := publicIPConfigsRaw[0].(map[string]interface{})
		publicIPAddressConfig := expandOrchestratedVirtualMachineScaleSetPublicIPAddressUpdate(publicIPConfigRaw)
		ipConfiguration.VirtualMachineScaleSetUpdateIPConfigurationProperties.PublicIPAddressConfiguration = publicIPAddressConfig
	}

	return &ipConfiguration, nil
}

func expandOrchestratedVirtualMachineScaleSetPublicIPAddressUpdate(raw map[string]interface{}) *compute.VirtualMachineScaleSetUpdatePublicIPAddressConfiguration {
	publicIPAddressConfig := compute.VirtualMachineScaleSetUpdatePublicIPAddressConfiguration{
		Name: utils.String(raw["name"].(string)),
		VirtualMachineScaleSetUpdatePublicIPAddressConfigurationProperties: &compute.VirtualMachineScaleSetUpdatePublicIPAddressConfigurationProperties{},
	}

	if domainNameLabel := raw["domain_name_label"].(string); domainNameLabel != "" {
		dns := &compute.VirtualMachineScaleSetPublicIPAddressConfigurationDNSSettings{
			DomainNameLabel: utils.String(domainNameLabel),
		}
		publicIPAddressConfig.VirtualMachineScaleSetUpdatePublicIPAddressConfigurationProperties.DNSSettings = dns
	}

	if idleTimeout := raw["idle_timeout_in_minutes"].(int); idleTimeout > 0 {
		publicIPAddressConfig.VirtualMachineScaleSetUpdatePublicIPAddressConfigurationProperties.IdleTimeoutInMinutes = utils.Int32(int32(raw["idle_timeout_in_minutes"].(int)))
	}

	return &publicIPAddressConfig
}

func ExpandOrchestratedVirtualMachineScaleSetDataDisk(input []interface{}, ultraSSDEnabled bool) (*[]compute.VirtualMachineScaleSetDataDisk, error) {
	disks := make([]compute.VirtualMachineScaleSetDataDisk, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		storageAccountType := compute.StorageAccountTypes(raw["storage_account_type"].(string))
		disk := compute.VirtualMachineScaleSetDataDisk{
			Caching:    compute.CachingTypes(raw["caching"].(string)),
			DiskSizeGB: utils.Int32(int32(raw["disk_size_gb"].(int))),
			Lun:        utils.Int32(int32(raw["lun"].(int))),
			ManagedDisk: &compute.VirtualMachineScaleSetManagedDiskParameters{
				StorageAccountType: storageAccountType,
			},
			WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),
			CreateOption:            compute.DiskCreateOptionTypes(raw["create_option"].(string)),
		}

		if id := raw["disk_encryption_set_id"].(string); id != "" {
			disk.ManagedDisk.DiskEncryptionSet = &compute.DiskEncryptionSetParameters{
				ID: utils.String(id),
			}
		}

		var iops int
		if diskIops, ok := raw["ultra_ssd_disk_iops_read_write"]; ok && diskIops.(int) > 0 {
			iops = diskIops.(int)
		}

		if iops > 0 && !ultraSSDEnabled && storageAccountType != compute.StorageAccountTypesPremiumV2LRS {
			return nil, fmt.Errorf("`ultra_ssd_disk_iops_read_write` can only be set when `storage_account_type` is set to `PremiumV2_LRS` or `UltraSSD_LRS`")
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if iops > 0 {
			disk.DiskIOPSReadWrite = utils.Int64(int64(iops))
		}

		var mbps int
		if diskMbps, ok := raw["ultra_ssd_disk_mbps_read_write"]; ok && diskMbps.(int) > 0 {
			mbps = diskMbps.(int)
		}

		if mbps > 0 && !ultraSSDEnabled && storageAccountType != compute.StorageAccountTypesPremiumV2LRS {
			return nil, fmt.Errorf("`ultra_ssd_disk_mbps_read_write` can only be set when `storage_account_type` is set to `PremiumV2_LRS` or `UltraSSD_LRS`")
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if mbps > 0 {
			disk.DiskMBpsReadWrite = utils.Int64(int64(mbps))
		}

		disks = append(disks, disk)
	}

	return &disks, nil
}

func ExpandOrchestratedVirtualMachineScaleSetOSDisk(input []interface{}, osType compute.OperatingSystemTypes) *compute.VirtualMachineScaleSetOSDisk {
	raw := input[0].(map[string]interface{})
	disk := compute.VirtualMachineScaleSetOSDisk{
		Caching: compute.CachingTypes(raw["caching"].(string)),
		ManagedDisk: &compute.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: compute.StorageAccountTypes(raw["storage_account_type"].(string)),
		},
		WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),

		// these have to be hard-coded so there's no point exposing them
		CreateOption: compute.DiskCreateOptionTypesFromImage,
		OsType:       osType,
	}

	if diskEncryptionSetId := raw["disk_encryption_set_id"].(string); diskEncryptionSetId != "" {
		disk.ManagedDisk.DiskEncryptionSet = &compute.DiskEncryptionSetParameters{
			ID: utils.String(diskEncryptionSetId),
		}
	}

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = utils.Int32(int32(osDiskSize))
	}

	if diffDiskSettingsRaw := raw["diff_disk_settings"].([]interface{}); len(diffDiskSettingsRaw) > 0 {
		diffDiskRaw := diffDiskSettingsRaw[0].(map[string]interface{})
		disk.DiffDiskSettings = &compute.DiffDiskSettings{
			Option:    compute.DiffDiskOptions(diffDiskRaw["option"].(string)),
			Placement: compute.DiffDiskPlacement(diffDiskRaw["placement"].(string)),
		}
	}

	return &disk
}

func ExpandOrchestratedVirtualMachineScaleSetOSDiskUpdate(input []interface{}) *compute.VirtualMachineScaleSetUpdateOSDisk {
	raw := input[0].(map[string]interface{})
	disk := compute.VirtualMachineScaleSetUpdateOSDisk{
		Caching: compute.CachingTypes(raw["caching"].(string)),
		ManagedDisk: &compute.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: compute.StorageAccountTypes(raw["storage_account_type"].(string)),
		},
		WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),
	}

	if diskEncryptionSetId := raw["disk_encryption_set_id"].(string); diskEncryptionSetId != "" {
		disk.ManagedDisk.DiskEncryptionSet = &compute.DiskEncryptionSetParameters{
			ID: utils.String(diskEncryptionSetId),
		}
	}

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = utils.Int32(int32(osDiskSize))
	}

	return &disk
}

func ExpandOrchestratedVirtualMachineScaleSetScheduledEventsProfile(input []interface{}) *compute.ScheduledEventsProfile {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	enabled := raw["enabled"].(bool)
	timeout := raw["timeout"].(string)

	return &compute.ScheduledEventsProfile{
		TerminateNotificationProfile: &compute.TerminateNotificationProfile{
			Enable:           &enabled,
			NotBeforeTimeout: &timeout,
		},
	}
}

func expandOrchestratedVirtualMachineScaleSetExtensions(input []interface{}) (extensionProfile *compute.VirtualMachineScaleSetExtensionProfile, hasHealthExtension bool, err error) {
	extensionProfile = &compute.VirtualMachineScaleSetExtensionProfile{}
	if len(input) == 0 {
		return nil, false, nil
	}

	extensions := make([]compute.VirtualMachineScaleSetExtension, 0)
	for _, v := range input {
		extensionRaw := v.(map[string]interface{})
		extension := compute.VirtualMachineScaleSetExtension{
			Name: utils.String(extensionRaw["name"].(string)),
		}
		extensionType := extensionRaw["type"].(string)

		extensionProps := compute.VirtualMachineScaleSetExtensionProperties{
			Publisher:                utils.String(extensionRaw["publisher"].(string)),
			Type:                     &extensionType,
			TypeHandlerVersion:       utils.String(extensionRaw["type_handler_version"].(string)),
			AutoUpgradeMinorVersion:  utils.Bool(extensionRaw["auto_upgrade_minor_version_enabled"].(bool)),
			ProvisionAfterExtensions: utils.ExpandStringSlice(extensionRaw["extensions_to_provision_after_vm_creation"].([]interface{})),
		}

		if extensionType == "ApplicationHealthLinux" || extensionType == "ApplicationHealthWindows" {
			hasHealthExtension = true
		}

		if val, ok := extensionRaw["failure_suppression_enabled"]; ok {
			extensionProps.SuppressFailures = utils.Bool(val.(bool))
		}

		if forceUpdateTag := extensionRaw["force_extension_execution_on_change"]; forceUpdateTag != nil {
			extensionProps.ForceUpdateTag = utils.String(forceUpdateTag.(string))
		}

		if val, ok := extensionRaw["settings"]; ok && val.(string) != "" {
			settings, err := pluginsdk.ExpandJsonFromString(val.(string))
			if err != nil {
				return nil, false, fmt.Errorf("failed to parse JSON from `settings`: %+v", err)
			}
			extensionProps.Settings = settings
		}

		protectedSettingsFromKeyVault := expandProtectedSettingsFromKeyVault(extensionRaw["protected_settings_from_key_vault"].([]interface{}))
		extensionProps.ProtectedSettingsFromKeyVault = protectedSettingsFromKeyVault

		if val, ok := extensionRaw["protected_settings"]; ok && val.(string) != "" {
			if protectedSettingsFromKeyVault != nil {
				return nil, false, fmt.Errorf("`protected_settings_from_key_vault` cannot be used with `protected_settings`")
			}

			protectedSettings, err := pluginsdk.ExpandJsonFromString(val.(string))
			if err != nil {
				return nil, false, fmt.Errorf("failed to parse JSON for `protected_settings`: %+v", err)
			}
			extensionProps.ProtectedSettings = protectedSettings
		}

		extension.VirtualMachineScaleSetExtensionProperties = &extensionProps
		extensions = append(extensions, extension)
	}
	extensionProfile.Extensions = &extensions

	return extensionProfile, hasHealthExtension, nil
}

func ExpandOrchestratedVirtualMachineScaleSetPriorityMixPolicy(input []interface{}) *compute.PriorityMixPolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &compute.PriorityMixPolicy{
		BaseRegularPriorityCount:           utils.Int32(int32(raw["base_regular_count"].(int))),
		RegularPriorityPercentageAboveBase: utils.Int32(int32(raw["regular_percentage_above_base"].(int))),
	}
}

func flattenOrchestratedVirtualMachineScaleSetExtensions(input *compute.VirtualMachineScaleSetExtensionProfile, d *pluginsdk.ResourceData) ([]map[string]interface{}, error) {
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
		// Automatic Upgrade is not yet supported in VMSS Flex
		// enableAutomaticUpgrade := false
		forceUpdateTag := ""
		provisionAfterExtension := make([]interface{}, 0)
		protectedSettings := ""
		var protectedSettingsFromKeyVault *compute.KeyVaultSecretReference
		extPublisher := ""
		extSettings := ""
		extType := ""
		extTypeVersion := ""

		if props := v.VirtualMachineScaleSetExtensionProperties; props != nil {
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

			// if props.EnableAutomaticUpgrade != nil {
			// 	enableAutomaticUpgrade = *props.EnableAutomaticUpgrade
			// }

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
				extSettingsRaw, err := pluginsdk.FlattenJsonToString(props.Settings.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				extSettings = extSettingsRaw
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
			"name":                               name,
			"auto_upgrade_minor_version_enabled": autoUpgradeMinorVersion,
			// "automatic_upgrade_enabled":  enableAutomaticUpgrade,
			"force_extension_execution_on_change":       forceUpdateTag,
			"failure_suppression_enabled":               suppressFailures,
			"extensions_to_provision_after_vm_creation": provisionAfterExtension,
			"protected_settings":                        protectedSettings,
			"protected_settings_from_key_vault":         flattenProtectedSettingsFromKeyVault(protectedSettingsFromKeyVault),
			"publisher":                                 extPublisher,
			"settings":                                  extSettings,
			"type":                                      extType,
			"type_handler_version":                      extTypeVersion,
		})
	}
	return result, nil
}

func FlattenOrchestratedVirtualMachineScaleSetIPConfiguration(input compute.VirtualMachineScaleSetIPConfiguration) map[string]interface{} {
	var name, subnetId string
	if input.Name != nil {
		name = *input.Name
	}
	if input.Subnet != nil && input.Subnet.ID != nil {
		subnetId = *input.Subnet.ID
	}

	var primary bool
	if input.Primary != nil {
		primary = *input.Primary
	}

	var publicIPAddresses []interface{}
	if input.PublicIPAddressConfiguration != nil {
		publicIPAddresses = append(publicIPAddresses, FlattenOrchestratedVirtualMachineScaleSetPublicIPAddress(*input.PublicIPAddressConfiguration))
	}

	applicationGatewayBackendAddressPoolIds := flattenSubResourcesToIDs(input.ApplicationGatewayBackendAddressPools)
	applicationSecurityGroupIds := flattenSubResourcesToIDs(input.ApplicationSecurityGroups)
	loadBalancerBackendAddressPoolIds := flattenSubResourcesToIDs(input.LoadBalancerBackendAddressPools)

	return map[string]interface{}{
		"name":              name,
		"primary":           primary,
		"public_ip_address": publicIPAddresses,
		"subnet_id":         subnetId,
		"version":           string(input.PrivateIPAddressVersion),
		"application_gateway_backend_address_pool_ids": applicationGatewayBackendAddressPoolIds,
		"application_security_group_ids":               applicationSecurityGroupIds,
		"load_balancer_backend_address_pool_ids":       loadBalancerBackendAddressPoolIds,
		// load_balancer_inbound_nat_rules_ids removed per service team this attribute will never be used in VMSS Flex
	}
}

func FlattenOrchestratedVirtualMachineScaleSetPublicIPAddress(input compute.VirtualMachineScaleSetPublicIPAddressConfiguration) map[string]interface{} {
	ipTags := make([]interface{}, 0)
	if input.IPTags != nil {
		for _, rawTag := range *input.IPTags {
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

	var domainNameLabel, name, publicIPPrefixId string
	if input.DNSSettings != nil && input.DNSSettings.DomainNameLabel != nil {
		domainNameLabel = *input.DNSSettings.DomainNameLabel
	}
	if input.Name != nil {
		name = *input.Name
	}
	if input.PublicIPPrefix != nil && input.PublicIPPrefix.ID != nil {
		publicIPPrefixId = *input.PublicIPPrefix.ID
	}

	var idleTimeoutInMinutes int
	if input.IdleTimeoutInMinutes != nil {
		idleTimeoutInMinutes = int(*input.IdleTimeoutInMinutes)
	}

	var version string
	if input.PublicIPAddressVersion != "" {
		version = string(input.PublicIPAddressVersion)
	}

	var sku string
	if input.Sku != nil && input.Sku.Name != "" && input.Sku.Tier != "" {
		sku = flattenOrchestratedVirtualMachineScaleSetPublicIPSku(input.Sku)
	}

	return map[string]interface{}{
		"name":                    name,
		"domain_name_label":       domainNameLabel,
		"idle_timeout_in_minutes": idleTimeoutInMinutes,
		"ip_tag":                  ipTags,
		"public_ip_prefix_id":     publicIPPrefixId,
		"sku_name":                sku,
		"version":                 version,
	}
}

func flattenOrchestratedVirtualMachineScaleSetWindowsConfiguration(input *compute.VirtualMachineScaleSetOSProfile, d *pluginsdk.ResourceData) []interface{} {
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

	// I am only commenting this out as this is going to be supported in the next release of the API in October 2021
	// if v := winConfig.AdditionalUnattendContent; v != nil {
	// 	output["additional_unattend_content"] = flattenWindowsConfigurationAdditionalUnattendContent(winConfig)
	// }

	if v := winConfig.EnableAutomaticUpdates; v != nil {
		output["enable_automatic_updates"] = *v
	}

	if v := winConfig.ProvisionVMAgent; v != nil {
		output["provision_vm_agent"] = *v
	}

	if v := input.Secrets; v != nil {
		output["secret"] = flattenWindowsSecrets(v)
	}

	if v := winConfig.WinRM; v != nil {
		output["winrm_listener"] = flattenWinRMListener(winConfig.WinRM)
	}

	if v := winConfig.TimeZone; v != nil {
		output["timezone"] = v
	}

	output["patch_mode"] = string(compute.WindowsVMGuestPatchModeAutomaticByOS)
	output["patch_assessment_mode"] = string(compute.WindowsPatchAssessmentModeAutomaticByPlatform)
	output["hotpatching_enabled"] = false

	if patchSettings != nil {
		output["patch_mode"] = string(patchSettings.PatchMode)
		output["patch_assessment_mode"] = string(patchSettings.AssessmentMode)

		if v := patchSettings.EnableHotpatching; v != nil {
			output["hotpatching_enabled"] = *v
		}
	}

	return []interface{}{output}
}

func flattenOrchestratedVirtualMachineScaleSetLinuxConfiguration(input *compute.VirtualMachineScaleSetOSProfile, d *pluginsdk.ResourceData) []interface{} {
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

	if v := linConfig.SSH; v != nil {
		if sshKeys, _ := FlattenSSHKeys(v); sshKeys != nil {
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
		output["secret"] = flattenLinuxSecrets(v)
	}

	return []interface{}{output}
}

func FlattenOrchestratedVirtualMachineScaleSetNetworkInterface(input *[]compute.VirtualMachineScaleSetNetworkConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		var name, networkSecurityGroupId string
		if v.Name != nil {
			name = *v.Name
		}
		if v.NetworkSecurityGroup != nil && v.NetworkSecurityGroup.ID != nil {
			networkSecurityGroupId = *v.NetworkSecurityGroup.ID
		}

		var enableAcceleratedNetworking, enableIPForwarding, primary bool
		if v.EnableAcceleratedNetworking != nil {
			enableAcceleratedNetworking = *v.EnableAcceleratedNetworking
		}
		if v.EnableIPForwarding != nil {
			enableIPForwarding = *v.EnableIPForwarding
		}
		if v.Primary != nil {
			primary = *v.Primary
		}

		var dnsServers []interface{}
		if settings := v.DNSSettings; settings != nil {
			dnsServers = utils.FlattenStringSlice(v.DNSSettings.DNSServers)
		}

		var ipConfigurations []interface{}
		if v.IPConfigurations != nil {
			for _, configRaw := range *v.IPConfigurations {
				config := FlattenOrchestratedVirtualMachineScaleSetIPConfiguration(configRaw)
				ipConfigurations = append(ipConfigurations, config)
			}
		}

		results = append(results, map[string]interface{}{
			"name":                          name,
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

func FlattenOrchestratedVirtualMachineScaleSetDataDisk(input *[]compute.VirtualMachineScaleSetDataDisk) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		diskSizeGb := 0
		if v.DiskSizeGB != nil && *v.DiskSizeGB != 0 {
			diskSizeGb = int(*v.DiskSizeGB)
		}

		lun := 0
		if v.Lun != nil {
			lun = int(*v.Lun)
		}

		storageAccountType := ""
		diskEncryptionSetId := ""
		if v.ManagedDisk != nil {
			storageAccountType = string(v.ManagedDisk.StorageAccountType)
			if v.ManagedDisk.DiskEncryptionSet != nil && v.ManagedDisk.DiskEncryptionSet.ID != nil {
				diskEncryptionSetId = *v.ManagedDisk.DiskEncryptionSet.ID
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
			"caching":                        string(v.Caching),
			"create_option":                  string(v.CreateOption),
			"lun":                            lun,
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

// I am only commenting this out as this is going to be supported in the next release of the API in October 2021
// func flattenWindowsConfigurationAdditionalUnattendContent(input *compute.WindowsConfiguration) []interface{} {
// 	if input == nil {
// 		return []interface{}{}
// 	}

// 	output := make([]interface{}, 0)
// 	for _, v := range *input.AdditionalUnattendContent {
// 		// content isn't returned by the API since it's sensitive data so we need to look it up later
// 		// where we can pull it out of the state file.
// 		output = append(output, map[string]interface{}{
// 			"content": "",
// 			"setting": string(v.SettingName),
// 		})
// 	}

// 	return output
// }

func FlattenOrchestratedVirtualMachineScaleSetOSDisk(input *compute.VirtualMachineScaleSetOSDisk) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	diffDiskSettings := make([]interface{}, 0)
	if input.DiffDiskSettings != nil {
		diffDiskSettings = append(diffDiskSettings, map[string]interface{}{
			"option":    string(input.DiffDiskSettings.Option),
			"placement": string(input.DiffDiskSettings.Placement),
		})
	}

	diskSizeGb := 0
	if input.DiskSizeGB != nil && *input.DiskSizeGB != 0 {
		diskSizeGb = int(*input.DiskSizeGB)
	}

	storageAccountType := ""
	diskEncryptionSetId := ""
	if input.ManagedDisk != nil {
		storageAccountType = string(input.ManagedDisk.StorageAccountType)
		if input.ManagedDisk.DiskEncryptionSet != nil && input.ManagedDisk.DiskEncryptionSet.ID != nil {
			diskEncryptionSetId = *input.ManagedDisk.DiskEncryptionSet.ID
		}
	}

	writeAcceleratorEnabled := false
	if input.WriteAcceleratorEnabled != nil {
		writeAcceleratorEnabled = *input.WriteAcceleratorEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"caching":                   string(input.Caching),
			"disk_size_gb":              diskSizeGb,
			"diff_disk_settings":        diffDiskSettings,
			"storage_account_type":      storageAccountType,
			"write_accelerator_enabled": writeAcceleratorEnabled,
			"disk_encryption_set_id":    diskEncryptionSetId,
		},
	}
}

func FlattenOrchestratedVirtualMachineScaleSetScheduledEventsProfile(input *compute.ScheduledEventsProfile) []interface{} {
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

func FlattenOrchestratedVirtualMachineScaleSetPriorityMixPolicy(input *compute.PriorityMixPolicy) []interface{} {

	baseRegularPriorityCount := int32(0)
	if input != nil && input.BaseRegularPriorityCount != nil {
		baseRegularPriorityCount = *input.BaseRegularPriorityCount
	}

	regularPriorityPercentageAboveBase := int32(0)
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
