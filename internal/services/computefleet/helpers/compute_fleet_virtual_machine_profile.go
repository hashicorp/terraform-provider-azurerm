// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/computefleet/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func GalleryApplicationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 25,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"version_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: galleryapplicationversions.ValidateApplicationVersionID,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"configuration_blob_uri": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},

				"order": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(0, 2147483647),
				},

				"tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"treat_failure_as_deployment_failure_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},
			},
		},
	}
}

func ExtensionSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[^_\W][\w.-]{0,78}[\w]$|^[^_\W]$`),
						"`name` can be up to 80 characters long. It must begin with a word character, and it must end with a word character or with `_`. `name` may contain word characters or `.`, `-`, `_`."),
				},

				"publisher": {
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

				"type_handler_version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"auto_upgrade_minor_version_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"extensions_to_provision_after_virtual_machine_creation": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"failure_suppression_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"force_extension_execution_on_change": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"protected_settings_from_key_vault": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_url": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeAny),
							},

							"source_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),
						},
					},
				},

				"protected_settings_json": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsJSON,
				},

				"settings_json": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					ForceNew:         true,
					ValidateFunc:     validation.StringIsJSON,
					DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				},
			},
		},
	}
}

func NetworkInterfaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 8,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[^_\W][\w.-]{0,69}[\w]$|^[^_\W]$`),
						"`name` can be up to 71 characters long. It must begin with a word character, and it must end with a word character or with `_`. `name` may contain word characters or `.`, `-`, `_`."),
				},

				"ip_configuration": ipConfigurationSchema(),

				"accelerated_networking_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"auxiliary_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.NetworkInterfaceAuxiliaryModeAcceleratedConnections),
						string(fleets.NetworkInterfaceAuxiliaryModeFloating),
					}, false),
				},

				"auxiliary_sku": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.NetworkInterfaceAuxiliarySkuATwo),
						string(fleets.NetworkInterfaceAuxiliarySkuAFour),
						string(fleets.NetworkInterfaceAuxiliarySkuAEight),
						string(fleets.NetworkInterfaceAuxiliarySkuAOne),
					}, false),
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
				},

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 9,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.IsIPAddress,
					},
				},

				"ip_forwarding_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"network_security_group_id": commonschema.ResourceIDReferenceOptionalForceNew(&networksecuritygroups.NetworkSecurityGroupId{}),

				"primary_network_interface_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},
			},
		},
	}
}

func ipConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 256,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[^_\W][\w.-]{0,78}[\w]$|^[^_\W]$`),
						"`name` can be up to 80 characters long. It must begin with a word character, and it must end with a word character or with `_`. `name` may contain word characters or `.`, `-`, `_`."),
				},

				"subnet_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},

				"application_gateway_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					ForceNew: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					ForceNew: true,
					MaxItems: 20,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
					},
				},

				"load_balancer_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					ForceNew: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
					},
				},

				"primary_ip_configuration_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"public_ip_address": publicIPAddressSchema(),

				"version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      string(fleets.IPVersionIPvFour),
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForIPVersion(), false),
				},
			},
		},
	}
}

func publicIPAddressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[^_\W][\w.-]{0,70}$`),
						"`name` can be up to 71 characters long. It must begin with a word character. `name` may contain word characters or `.`, `-`, `_`."),
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDeleteOptions(), false),
				},

				"domain_name_label": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[a-z][a-z0-9-]{1,52}[a-z0-9]$`),
						"`domain_name_label` must be between 1 - 54 characters long (inclusive), contain only lower case letters, numbers, and hyphens, start with a letter, and end with a letter"),
				},

				"domain_name_label_scope": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDomainNameLabelScopeTypes(), false),
				},

				"idle_timeout_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(4, 30),
				},

				"public_ip_prefix_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
				},

				// since "BasicSkuPublicIpIsNotAllowedForVmssFlex", the possible values are `Standard_Regional` and `Standard_Global`
				"sku_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string("Standard_Regional"),
						string("Standard_Global"),
					}, false),
				},

				"version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      string(fleets.IPVersionIPvFour),
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForIPVersion(), false),
				},
			},
		},
	}
}

func OsProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"custom_data_base64": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					ForceNew:  true,
					Sensitive: true,
					ValidateFunc: validation.All(
						validation.StringIsBase64,
						validation.StringLenBetween(1, 87380),
					),
				},

				"linux_configuration": {
					Type:         pluginsdk.TypeList,
					Optional:     true,
					ForceNew:     true,
					MaxItems:     1,
					ExactlyOneOf: []string{"virtual_machine_profile.0.os_profile.0.linux_configuration", "virtual_machine_profile.0.os_profile.0.windows_configuration"},
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"admin_username": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.All(
									computeValidate.LinuxAdminUsername,
									// Validation according to Azure portal. Behavior of Azure REST API is different.
									validation.StringMatch(regexp.MustCompile(`^[a-zA-Z_][-a-zA-Z0-9_]*$`), "`admin_username` must only contain letters, numbers, hyphens, and underscores and may not start with a hyphen or number."),
								),
							},

							"admin_password": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								Sensitive:    true,
								ValidateFunc: computeValidate.LinuxAdminPassword,
							},

							"admin_ssh_keys": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								ForceNew: true,
								MinItems: 1,
								Elem: &pluginsdk.Schema{
									Type:             pluginsdk.TypeString,
									ValidateFunc:     computeValidate.SSHKey,
									DiffSuppressFunc: suppress.SSHKey,
								},
							},

							"bypass_platform_safety_checks_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  false,
							},

							"computer_name_prefix": {
								Type:             pluginsdk.TypeString,
								Optional:         true,
								ForceNew:         true,
								ValidateFunc:     computeValidate.LinuxComputerNamePrefix,
								DiffSuppressFunc: computeFleetComputerNamePrefixDiffSuppress,
							},

							"password_authentication_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  false,
							},

							"patch_mode": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchMode(), false),
							},

							"provision_virtual_machine_agent_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  true,
							},

							"patch_rebooting": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchAutomaticByPlatformRebootSetting(), false),
							},

							"secret": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

										"certificate": {
											Type:     pluginsdk.TypeSet,
											Required: true,
											ForceNew: true,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													// whilst we /could/ flatten this to `certificate_urls` we're intentionally not to keep this
													// closer to the Windows VMSS resource, which will also take a `store` param
													"url": {
														Type:         pluginsdk.TypeString,
														Required:     true,
														ForceNew:     true,
														ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeAny),
													},
												},
											},
										},
									},
								},
							},

							"virtual_machine_agent_platform_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  false,
							},
						},
					},
				},

				"windows_configuration": {
					Type:         pluginsdk.TypeList,
					Optional:     true,
					ForceNew:     true,
					MaxItems:     1,
					ExactlyOneOf: []string{"virtual_machine_profile.0.os_profile.0.linux_configuration", "virtual_machine_profile.0.os_profile.0.windows_configuration"},
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"admin_username": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.All(
									computeValidate.WindowsAdminUsername,
									validation.StringDoesNotContainAny("\\/\"[]:|<>+=;,?*@&"),
								),
							},

							"admin_password": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								Sensitive:    true,
								ValidateFunc: computeValidate.WindowsAdminPassword,
							},

							"additional_unattend_content": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"xml": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ForceNew:     true,
											Sensitive:    true,
											ValidateFunc: validate.Xml,
										},
										"setting": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSettingNames(), false),
										},
									},
								},
							},

							"automatic_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  true,
							},

							"bypass_platform_safety_checks_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  false,
							},

							"computer_name_prefix": {
								Type:             pluginsdk.TypeString,
								Optional:         true,
								ForceNew:         true,
								ValidateFunc:     computeValidate.WindowsComputerNamePrefix,
								DiffSuppressFunc: computeFleetComputerNamePrefixDiffSuppress,
							},

							"hot_patching_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  false,
							},

							"patch_mode": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForWindowsVMGuestPatchMode(), false),
							},

							"provision_virtual_machine_agent_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  true,
							},

							"patch_rebooting": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForWindowsVMGuestPatchAutomaticByPlatformRebootSetting(), false),
							},

							"secret": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

										"certificate": {
											Type:     pluginsdk.TypeSet,
											Required: true,
											ForceNew: true,
											MinItems: 1,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"url": {
														Type:         pluginsdk.TypeString,
														Required:     true,
														ForceNew:     true,
														ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeAny),
													},

													"store": {
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ForceNew:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},
												},
											},
										},
									},
								},
							},

							"time_zone": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: computeValidate.VirtualMachineTimeZone(),
							},

							"virtual_machine_agent_platform_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  false,
							},

							"winrm_listener": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								ForceNew: true,
								MaxItems: 2,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"protocol": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForProtocolTypes(), false),
										},
										"certificate_url": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ForceNew:     true,
											ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeAny),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func StorageProfileDataDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 10,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"create_option": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.DiskCreateOptionTypesEmpty),
						string(fleets.DiskCreateOptionTypesFromImage),
					}, false),
				},

				"caching": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.CachingTypesReadOnly),
						string(fleets.CachingTypesReadWrite),
					}, false),
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
				},

				"disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: commonids.ValidateDiskEncryptionSetID,
				},

				"disk_size_in_gib": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"lun": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(0, 63),
				},

				"storage_account_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForStorageAccountTypes(), false),
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},
			},
		},
	}
}

func StorageProfileOsDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		// NOTE: O+C - `os_disk` block is returned by the API even if not specified
		Computed: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"caching": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.CachingTypesReadOnly),
						string(fleets.CachingTypesReadWrite),
					}, false),
					AtLeastOneOf: computeFleetVirtualMachineProfileOsDiskConstraint(),
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      string(fleets.DiskDeleteOptionTypesDelete),
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
					AtLeastOneOf: computeFleetVirtualMachineProfileOsDiskConstraint(),
				},

				"diff_disk_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiffDiskOptions(), false),
					AtLeastOneOf: computeFleetVirtualMachineProfileOsDiskConstraint(),
				},

				"diff_disk_placement": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiffDiskPlacement(), false),
					AtLeastOneOf: computeFleetVirtualMachineProfileOsDiskConstraint(),
				},

				"disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: commonids.ValidateDiskEncryptionSetID,
					AtLeastOneOf: computeFleetVirtualMachineProfileOsDiskConstraint(),
				},

				"disk_size_in_gib": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(1, 4095),
					AtLeastOneOf: computeFleetVirtualMachineProfileOsDiskConstraint(),
				},

				"security_encryption_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSecurityEncryptionTypes(), false),
					AtLeastOneOf: computeFleetVirtualMachineProfileOsDiskConstraint(),
				},

				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					// NOTE: OS Disks don't support Ultra SSDs or PremiumV2_LRS
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.StorageAccountTypesPremiumLRS),
						string(fleets.StorageAccountTypesPremiumZRS),
						string(fleets.StorageAccountTypesStandardLRS),
						string(fleets.StorageAccountTypesStandardSSDLRS),
						string(fleets.StorageAccountTypesStandardSSDZRS),
					}, false),
					AtLeastOneOf: computeFleetVirtualMachineProfileOsDiskConstraint(),
				},

				"write_accelerator_enabled": {
					Type:         pluginsdk.TypeBool,
					Optional:     true,
					ForceNew:     true,
					Default:      false,
					AtLeastOneOf: computeFleetVirtualMachineProfileOsDiskConstraint(),
				},
			},
		},
	}
}

func StorageProfileSourceImageReferenceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"offer": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"sku": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
		ExactlyOneOf: []string{
			"virtual_machine_profile.0.source_image_reference",
			"virtual_machine_profile.0.source_image_id",
		},
	}
}

func ParseUsernameFromAuthorizedKeysPath(input string) *string {
	// the Azure VM agent hard-codes this to `/home/username/.ssh/authorized_keys`
	// as such we can hard-code this for a better UX
	r := regexp.MustCompile("(/home/)+(?P<username>.*?)(/.ssh/authorized_keys)+")

	keys := r.SubexpNames()
	values := r.FindStringSubmatch(input)

	if values == nil {
		return nil
	}

	for i, k := range keys {
		if k == "username" {
			return &values[i]
		}
	}

	return nil
}

// Suppress difference from default when `computer_name_prefix` is not set. The default of `computer_name_prefix` is the first 9 characters of the compute fleet name according to Azure portal.
func computeFleetComputerNamePrefixDiffSuppress(_, old, new string, d *pluginsdk.ResourceData) bool {
	defaultComputerNamePrefix := d.Get("name").(string)
	if len(defaultComputerNamePrefix) > 9 {
		defaultComputerNamePrefix = defaultComputerNamePrefix[:9]
	}

	if new == "" && old == defaultComputerNamePrefix {
		return true
	}

	return false
}

func computeFleetVirtualMachineProfileOsDiskConstraint() []string {
	return []string{
		"virtual_machine_profile.0.os_disk.0.caching",
		"virtual_machine_profile.0.os_disk.0.delete_option",
		"virtual_machine_profile.0.os_disk.0.diff_disk_option",
		"virtual_machine_profile.0.os_disk.0.diff_disk_placement",
		"virtual_machine_profile.0.os_disk.0.disk_encryption_set_id",
		"virtual_machine_profile.0.os_disk.0.disk_size_in_gib",
		"virtual_machine_profile.0.os_disk.0.security_encryption_type",
		"virtual_machine_profile.0.os_disk.0.storage_account_type",
		"virtual_machine_profile.0.os_disk.0.write_accelerator_enabled",
	}
}
