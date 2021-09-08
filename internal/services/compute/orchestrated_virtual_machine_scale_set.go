package compute

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	msiparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
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
				"custom_data":           base64.OptionalSchema(false),
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

				"enable_automatic_updates": {
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

				"auto_upgrade_minor_version": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"force_update_tag": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"protected_settings": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsJSON,
				},

				"provision_after_extensions": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
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

func OrchestratedVirtualMachineScaleSetIdentitySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.ResourceIdentityTypeSystemAssigned),
						string(compute.ResourceIdentityTypeUserAssigned),
						string(compute.ResourceIdentityTypeSystemAssignedUserAssigned),
					}, false),
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
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetNetworkInterfaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
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
				"enable_accelerated_networking": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
				"enable_ip_forwarding": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
				"network_security_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateResourceIDOrEmpty,
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

func OrchestratedVirtualMachineScaleSetNetworkInterfaceSchemaForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"ip_configuration": orchestratedVirtualMachineScaleSetIPConfigurationSchemaForDataSource(),

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"enable_accelerated_networking": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
				"enable_ip_forwarding": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
				"network_security_group_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"primary": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
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

				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"public_ip_address": orchestratedVirtualMachineScaleSetPublicIPAddressSchema(),

				"subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateResourceID,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(compute.IPv4),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.IPv4),
						string(compute.IPv6),
					}, false),
				},
			},
		},
	}
}

func orchestratedVirtualMachineScaleSetIPConfigurationSchemaForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"application_gateway_backend_address_pool_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"load_balancer_backend_address_pool_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"public_ip_address": virtualMachineScaleSetPublicIPAddressSchemaForDataSource(),

				"subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
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
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"idle_timeout_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(4, 32),
				},
				"ip_tag": {
					// TODO: does this want to be a Set?
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
					ValidateFunc: azure.ValidateResourceIDOrEmpty,
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
		// TODO: does this want to be a Set?
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
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
						string(compute.StorageAccountTypesUltraSSDLRS),
					}, false),
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				// TODO 3.0 - change this to ultra_ssd_disk_iops_read_write
				"disk_iops_read_write": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Computed: true,
				},

				// TODO 3.0 - change this to ultra_ssd_disk_iops_read_write
				"disk_mbps_read_write": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetOSDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
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
						// note: OS Disks don't support Ultra SSDs
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
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
									string(compute.Local),
								}, false),
							},
						},
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

func OrchestratedVirtualMachineScaleSetTerminateNotificationSchema() *pluginsdk.Schema {
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
					ValidateFunc: azValidate.ISO8601Duration,
					Default:      "PT5M",
				},
			},
		},
	}
}

func OrchestratedVirtualMachineScaleSetAutomaticRepairsPolicySchema() *pluginsdk.Schema {
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
				"grace_period": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "PT30M",
					// this field actually has a range from 30m to 90m, is there a function that can do this validation?
					ValidateFunc: azValidate.ISO8601Duration,
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

	if v := input.CustomData; v != nil {
		output["custom_data"] = *v
	}

	if winConfig := input.WindowsConfiguration; winConfig != nil {
		output["windows_configuration"] = flattenOrchestratedVirtualMachineScaleSetWindowsConfiguration(input, d)
	}

	if linConfig := input.LinuxConfiguration; linConfig != nil {
		output["linux_configuration"] = flattenOrchestratedVirtualMachineScaleSetLinuxConfiguration(input)
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
		errors = append(errors, fmt.Errorf("%q must be at least 6 characters long and shorter than 72 characters long. Got %q(%d characters)", key, password, len(password)))
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

func expandOrchestratedVirtualMachineScaleSetOsProfileWithWindowsConfiguration(input map[string]interface{}) *compute.VirtualMachineScaleSetOSProfile {
	osProfile := compute.VirtualMachineScaleSetOSProfile{}
	winConfig := compute.WindowsConfiguration{}

	if len(input) > 0 {
		osProfile.AdminUsername = utils.String(input["admin_username"].(string))
		osProfile.AdminPassword = utils.String(input["admin_password"].(string))

		if computerPrefix := input["computer_name_prefix"].(string); computerPrefix != "" {
			osProfile.ComputerNamePrefix = utils.String(computerPrefix)
		}

		if secrets := input["secret"].([]interface{}); len(secrets) > 0 {
			osProfile.Secrets = expandWindowsSecrets(secrets)
		}

		winConfig.AdditionalUnattendContent = expandWindowsConfigurationAdditionalUnattendContent(input["additional_unattend_content"].([]interface{}))
		winConfig.EnableAutomaticUpdates = utils.Bool(input["enable_automatic_updates"].(bool))
		winConfig.ProvisionVMAgent = utils.Bool(input["provision_vm_agent"].(bool))
		winConfig.TimeZone = utils.String(input["timezone"].(string))
		winRmListenersRaw := input["winrm_listener"].(*pluginsdk.Set).List()
		winConfig.WinRM = expandWinRMListener(winRmListenersRaw)
	}

	osProfile.WindowsConfiguration = &winConfig

	return &osProfile
}

func expandOrchestratedVirtualMachineScaleSetOsProfileWithLinuxConfiguration(input map[string]interface{}) *compute.VirtualMachineScaleSetOSProfile {
	osProfile := compute.VirtualMachineScaleSetOSProfile{}
	linConfig := compute.LinuxConfiguration{}

	if len(input) > 0 {
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

		if sshPublicKeys := ExpandSSHKeys(input["admin_ssh_key"].([]interface{})); len(sshPublicKeys) > 0 {
			linConfig.SSH.PublicKeys = &sshPublicKeys
		}

		linConfig.DisablePasswordAuthentication = utils.Bool(input["disable_password_authentication"].(bool))
		linConfig.ProvisionVMAgent = utils.Bool(input["provision_vm_agent"].(bool))
	}

	osProfile.LinuxConfiguration = &linConfig

	return &osProfile
}

func expandWindowsConfigurationAdditionalUnattendContent(input []interface{}) *[]compute.AdditionalUnattendContent {
	output := make([]compute.AdditionalUnattendContent, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		output = append(output, compute.AdditionalUnattendContent{
			SettingName: compute.SettingNames(raw["setting"].(string)),
			Content:     utils.String(raw["content"].(string)),

			// no other possible values
			PassName:      compute.OobeSystem,
			ComponentName: compute.MicrosoftWindowsShellSetup,
		})
	}

	return &output
}

func ExpandOrchestratedVirtualMachineScaleSetIdentity(input []interface{}) (*compute.VirtualMachineScaleSetIdentity, error) {
	if len(input) == 0 {
		// TODO: Does this want to be this, or nil?
		return &compute.VirtualMachineScaleSetIdentity{
			Type: compute.ResourceIdentityTypeNone,
		}, nil
	}

	raw := input[0].(map[string]interface{})

	identity := compute.VirtualMachineScaleSetIdentity{
		Type: compute.ResourceIdentityType(raw["type"].(string)),
	}

	identityIdsRaw := raw["identity_ids"].(*pluginsdk.Set).List()
	identityIds := make(map[string]*compute.VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue)
	for _, v := range identityIdsRaw {
		identityIds[v.(string)] = &compute.VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue{}
	}

	if len(identityIds) > 0 {
		if identity.Type != compute.ResourceIdentityTypeUserAssigned && identity.Type != compute.ResourceIdentityTypeSystemAssignedUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}

		identity.UserAssignedIdentities = identityIds
	}

	return &identity, nil
}

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
	if primary && version == compute.IPv6 {
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

	if primary && version == compute.IPv6 {
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

		disk := compute.VirtualMachineScaleSetDataDisk{
			Caching:    compute.CachingTypes(raw["caching"].(string)),
			DiskSizeGB: utils.Int32(int32(raw["disk_size_gb"].(int))),
			Lun:        utils.Int32(int32(raw["lun"].(int))),
			ManagedDisk: &compute.VirtualMachineScaleSetManagedDiskParameters{
				StorageAccountType: compute.StorageAccountTypes(raw["storage_account_type"].(string)),
			},
			WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),
			CreateOption:            compute.DiskCreateOptionTypes(raw["create_option"].(string)),
		}

		if id := raw["disk_encryption_set_id"].(string); id != "" {
			disk.ManagedDisk.DiskEncryptionSet = &compute.DiskEncryptionSetParameters{
				ID: utils.String(id),
			}
		}

		if iops := raw["disk_iops_read_write"].(int); iops != 0 {
			if !ultraSSDEnabled {
				return nil, fmt.Errorf("`disk_iops_read_write` are only available for UltraSSD disks")
			}
			disk.DiskIOPSReadWrite = utils.Int64(int64(iops))
		}

		if mbps := raw["disk_mbps_read_write"].(int); mbps != 0 {
			if !ultraSSDEnabled {
				return nil, fmt.Errorf("`disk_mbps_read_write` are only available for UltraSSD disks")
			}
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
			Option: compute.DiffDiskOptions(diffDiskRaw["option"].(string)),
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

func ExpandOrchestratedVirtualMachineScaleSetAutomaticRepairsPolicy(input []interface{}) *compute.AutomaticRepairsPolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &compute.AutomaticRepairsPolicy{
		Enabled:     utils.Bool(raw["enabled"].(bool)),
		GracePeriod: utils.String(raw["grace_period"].(string)),
	}
}

func expandOrchestratedVirtualMachineScaleSetExtensions(input []interface{}) (extensionProfile *compute.VirtualMachineScaleSetExtensionProfile, err error) {
	extensionProfile = &compute.VirtualMachineScaleSetExtensionProfile{}
	if len(input) == 0 {
		return nil, nil
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
			AutoUpgradeMinorVersion:  utils.Bool(extensionRaw["auto_upgrade_minor_version"].(bool)),
			ProvisionAfterExtensions: utils.ExpandStringSlice(extensionRaw["provision_after_extensions"].([]interface{})),
		}

		// Leaving this here as it is going to be in the GA API
		// if extensionType == "ApplicationHealthLinux" || extensionType == "ApplicationHealthWindows" {
		// 	hasHealthExtension = true
		// }

		if forceUpdateTag := extensionRaw["force_update_tag"]; forceUpdateTag != nil {
			extensionProps.ForceUpdateTag = utils.String(forceUpdateTag.(string))
		}

		if val, ok := extensionRaw["settings"]; ok && val.(string) != "" {
			settings, err := pluginsdk.ExpandJsonFromString(val.(string))
			if err != nil {
				return nil, fmt.Errorf("failed to parse JSON from `settings`: %+v", err)
			}
			extensionProps.Settings = settings
		}

		if val, ok := extensionRaw["protected_settings"]; ok && val.(string) != "" {
			protectedSettings, err := pluginsdk.ExpandJsonFromString(val.(string))
			if err != nil {
				return nil, fmt.Errorf("failed to parse JSON from `protected_settings`: %+v", err)
			}
			extensionProps.ProtectedSettings = protectedSettings
		}

		extension.VirtualMachineScaleSetExtensionProperties = &extensionProps
		extensions = append(extensions, extension)
	}
	extensionProfile.Extensions = &extensions

	return extensionProfile, nil
}

func flattenOrchestratedVirtualMachineScaleSetExtensions(input *compute.VirtualMachineScaleSetExtensionProfile, d *pluginsdk.ResourceData) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	if input == nil || input.Extensions == nil {
		return result, nil
	}

	for k, v := range *input.Extensions {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		autoUpgradeMinorVersion := false
		forceUpdateTag := ""
		provisionAfterExtension := make([]interface{}, 0)
		protectedSettings := ""
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

			if props.ForceUpdateTag != nil {
				forceUpdateTag = *props.ForceUpdateTag
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
		}
		// protected_settings isn't returned, so we attempt to get it from config otherwise set to empty string
		if protectedSettingsFromConfig, ok := d.GetOk(fmt.Sprintf("extension.%d.protected_settings", k)); ok {
			if protectedSettingsFromConfig.(string) != "" && protectedSettingsFromConfig.(string) != "{}" {
				protectedSettings = protectedSettingsFromConfig.(string)
			}
		}

		result = append(result, map[string]interface{}{
			"name":                       name,
			"auto_upgrade_minor_version": autoUpgradeMinorVersion,
			"force_update_tag":           forceUpdateTag,
			"provision_after_extensions": provisionAfterExtension,
			"protected_settings":         protectedSettings,
			"publisher":                  extPublisher,
			"settings":                   extSettings,
			"type":                       extType,
			"type_handler_version":       extTypeVersion,
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

	return map[string]interface{}{
		"name":                    name,
		"domain_name_label":       domainNameLabel,
		"idle_timeout_in_minutes": idleTimeoutInMinutes,
		"ip_tag":                  ipTags,
		"public_ip_prefix_id":     publicIPPrefixId,
	}
}

func flattenOrchestratedVirtualMachineScaleSetWindowsConfiguration(input *compute.VirtualMachineScaleSetOSProfile, d *pluginsdk.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})
	winConfig := input.WindowsConfiguration

	if v := input.AdminUsername; v != nil {
		output["admin_username"] = *v
	}

	if v := d.Get("admin_password").(string); v != "" {
		output["admin_password"] = v
	}

	if v := input.ComputerNamePrefix; v != nil {
		output["computer_name_prefix"] = *v
	}

	if v := winConfig.AdditionalUnattendContent; v != nil {
		output["additional_unattend_content"] = flattenWindowsConfigurationAdditionalUnattendContent(winConfig)
	}

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

	return []interface{}{output}
}

func flattenOrchestratedVirtualMachineScaleSetLinuxConfiguration(input *compute.VirtualMachineScaleSetOSProfile) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})
	linConfig := input.LinuxConfiguration

	if v := input.AdminUsername; v != nil {
		output["admin_username"] = *v
	}

	if v := input.AdminPassword; v != nil {
		output["admin_password"] = *v
	}

	if v := linConfig.SSH; v != nil {
		if sshKeys, _ := FlattenSSHKeys(v); sshKeys != nil {
			output["admin_ssh_key"] = sshKeys
		}
	}

	if v := input.ComputerNamePrefix; v != nil {
		output["computer_name_prefix"] = *v
	}

	if v := linConfig.DisablePasswordAuthentication; v != nil {
		output["disable_password_authentication"] = *v
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

func FlattenOrchestratedVirtualMachineScaleSetIdentity(input *compute.VirtualMachineScaleSetIdentity) ([]interface{}, error) {
	if input == nil || input.Type == compute.ResourceIdentityTypeNone {
		return []interface{}{}, nil
	}

	identityIds := make([]string, 0)
	if input.UserAssignedIdentities != nil {
		for key := range input.UserAssignedIdentities {
			parsedId, err := msiparse.UserAssignedIdentityIDInsensitively(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
	}

	principalId := ""
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
			"principal_id": principalId,
		},
	}, nil
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
			"caching":                   string(v.Caching),
			"create_option":             string(v.CreateOption),
			"lun":                       lun,
			"disk_encryption_set_id":    diskEncryptionSetId,
			"disk_size_gb":              diskSizeGb,
			"storage_account_type":      storageAccountType,
			"write_accelerator_enabled": writeAcceleratorEnabled,
			"disk_iops_read_write":      iops,
			"disk_mbps_read_write":      mbps,
		})
	}

	return output
}

func flattenWindowsConfigurationAdditionalUnattendContent(input *compute.WindowsConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, v := range *input.AdditionalUnattendContent {
		// content isn't returned by the API since it's sensitive data so we need to look it up later
		// where we can pull it out of the state file.
		output = append(output, map[string]interface{}{
			"content": "",
			"setting": string(v.SettingName),
		})
	}

	return output
}

func FlattenOrchestratedVirtualMachineScaleSetOSDisk(input *compute.VirtualMachineScaleSetOSDisk) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	diffDiskSettings := make([]interface{}, 0)
	if input.DiffDiskSettings != nil {
		diffDiskSettings = append(diffDiskSettings, map[string]interface{}{
			"option": string(input.DiffDiskSettings.Option),
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

func FlattenOrchestratedVirtualMachineScaleSetAutomaticRepairsPolicy(input *compute.AutomaticRepairsPolicy) []interface{} {
	// if enabled is set to false, there will be no AutomaticRepairsPolicy in response, to avoid plan non empty when
	// a user explicitly set enabled to false, we need to assign a default block to this field

	enabled := false
	if input != nil && input.Enabled != nil {
		enabled = *input.Enabled
	}

	gracePeriod := "PT30M"
	if input != nil && input.GracePeriod != nil {
		gracePeriod = *input.GracePeriod
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":      enabled,
			"grace_period": gracePeriod,
		},
	}
}
