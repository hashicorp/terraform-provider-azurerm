// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func virtualMachineProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"network_api_version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"network_interface": networkInterfaceSchema(),

				"os_profile": osProfileSchema(),

				"boot_diagnostic_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"boot_diagnostic_storage_account_endpoint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"capacity_reservation_group_id": commonschema.ResourceIDReferenceOptionalForceNew(&capacityreservationgroups.CapacityReservationGroupId{}),

				"data_disk": storageProfileDataDiskSchema(),

				"encryption_at_host_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"extension": extensionSchema(),

				"extension_operations_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  true,
				},

				"extensions_time_budget_duration": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
				},

				"gallery_application": galleryApplicationSchema(),

				"license_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						"RHEL_BYOS",
						"SLES_BYOS",
						"Windows_Client",
						"Windows_Server",
					}, false),
				},

				"os_disk": storageProfileOsDiskSchema(),

				"scheduled_event_os_image_timeout_duration": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						"PT15M",
					}, false),
				},

				"scheduled_event_termination_timeout_duration": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
				},

				"secure_boot_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},

				"source_image_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.Any(
						images.ValidateImageID,
						computeValidate.SharedImageID,
						computeValidate.SharedImageVersionID,
						computeValidate.CommunityGalleryImageID,
						computeValidate.CommunityGalleryImageVersionID,
						computeValidate.SharedGalleryImageID,
						computeValidate.SharedGalleryImageVersionID,
					),
				},

				"source_image_reference": storageProfileSourceImageReferenceSchema(),

				"user_data_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsBase64,
				},

				"vtpm_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},
			},
		},
	}
}

func galleryApplicationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 100,
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

func extensionSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
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

				"extensions_to_provision_after_vm_creation": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
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
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
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
								ValidateFunc: keyVaultValidate.NestedItemId,
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

func networkInterfaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
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
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						string(fleets.NetworkInterfaceAuxiliaryModeAcceleratedConnections),
						string(fleets.NetworkInterfaceAuxiliaryModeFloating),
					}, false),
				},

				"auxiliary_sku": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
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
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
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
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
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
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					ForceNew: true,
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
					ForceNew: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
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
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDeleteOptions(), false),
				},

				"domain_name_label": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
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
					ValidateFunc: validation.IntBetween(4, 32),
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

func osProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"custom_data_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsBase64,
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
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: computeValidate.LinuxAdminUsername,
							},

							"computer_name_prefix": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: computeValidate.LinuxComputerNamePrefix,
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
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},

							"bypass_platform_safety_checks_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  false,
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

							"provision_vm_agent_enabled": {
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
											MinItems: 1,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													// whilst we /could/ flatten this to `certificate_urls` we're intentionally not to keep this
													// closer to the Windows VMSS resource, which will also take a `store` param
													"url": {
														Type:         pluginsdk.TypeString,
														Required:     true,
														ForceNew:     true,
														ValidateFunc: keyVaultValidate.NestedItemId,
													},
												},
											},
										},
									},
								},
							},

							"vm_agent_platform_updates_enabled": {
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
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: computeValidate.WindowsAdminUsername,
							},

							"admin_password": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								Sensitive:    true,
								ValidateFunc: computeValidate.WindowsAdminPassword,
							},

							"computer_name_prefix": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: computeValidate.WindowsComputerNamePrefix,
							},

							"additional_unattend_content": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								ForceNew: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"xml": {
											Type:      pluginsdk.TypeString,
											Required:  true,
											ForceNew:  true,
											Sensitive: true,
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

							"provision_vm_agent_enabled": {
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
														ValidateFunc: keyVaultValidate.NestedItemId,
													},

													"store": {
														Type:     pluginsdk.TypeString,
														Optional: true,
														ForceNew: true,
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
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"vm_agent_platform_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  false,
							},

							"winrm_listener": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								ForceNew: true,
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
											ValidateFunc: keyVaultValidate.NestedItemId,
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

func storageProfileDataDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
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
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
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
					ValidateFunc: computeValidate.DiskEncryptionSetID,
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
					ValidateFunc: validation.IntBetween(0, 2000),
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

func storageProfileOsDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"caching": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						string(fleets.CachingTypesReadOnly),
						string(fleets.CachingTypesReadWrite),
					}, false),
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      string(fleets.DiskDeleteOptionTypesDelete),
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
				},

				"diff_disk_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiffDiskOptions(), false),
				},

				"diff_disk_placement": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiffDiskPlacement(), false),
				},

				"disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: computeValidate.DiskEncryptionSetID,
				},

				"disk_size_in_gib": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"security_encryption_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSecurityEncryptionTypes(), false),
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

func storageProfileSourceImageReferenceSchema() *pluginsdk.Schema {
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
	}
}

func expandVirtualMachineProfileModel(inputList []VirtualMachineProfileModel, d *schema.ResourceData) (*fleets.BaseVirtualMachineProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fleets.BaseVirtualMachineProfile{
		OsProfile:          expandOSProfileModel(inputList),
		ApplicationProfile: expandApplicationProfileModel(input.GalleryApplicationProfile),
		DiagnosticsProfile: &fleets.DiagnosticsProfile{
			BootDiagnostics: &fleets.BootDiagnostics{
				Enabled:    pointer.To(input.BootDiagnosticEnabled),
				StorageUri: pointer.To(input.BootDiagnosticStorageAccountEndpoint),
			},
		},
		NetworkProfile: &fleets.VirtualMachineScaleSetNetworkProfile{
			NetworkApiVersion:              pointer.To(fleets.NetworkApiVersion(input.NetworkApiVersion)),
			NetworkInterfaceConfigurations: expandNetworkInterfaceModel(input.NetworkInterface),
		},
		StorageProfile: &fleets.VirtualMachineScaleSetStorageProfile{
			ImageReference: expandImageReference(input.SourceImageReference, input.SourceImageId),
			OsDisk:         expandOSDiskModel(input),
			DataDisks:      expandDataDiskModel(input.DataDisks),
		},
	}

	if input.ScheduledEventTerminationTimeoutDuration != "" {
		output.ScheduledEventsProfile = &fleets.ScheduledEventsProfile{
			TerminateNotificationProfile: &fleets.TerminateNotificationProfile{
				Enable:           pointer.To(true),
				NotBeforeTimeout: pointer.To(input.ScheduledEventTerminationTimeoutDuration),
			},
		}
	}

	if input.ScheduledEventOsImageTimeoutDuration != "" {
		if output.ScheduledEventsProfile == nil {
			output.ScheduledEventsProfile = &fleets.ScheduledEventsProfile{}
		}
		output.ScheduledEventsProfile.OsImageNotificationProfile = &fleets.OSImageNotificationProfile{
			Enable:           pointer.To(true),
			NotBeforeTimeout: pointer.To(input.ScheduledEventOsImageTimeoutDuration),
		}
	}

	if input.CapacityReservationGroupId != "" {
		output.CapacityReservation = &fleets.CapacityReservationProfile{
			CapacityReservationGroup: expandSubResource(input.CapacityReservationGroupId),
		}
	}

	encryptionAtHostEnabledExist := d.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["encryption_at_host_enabled"]
	if !encryptionAtHostEnabledExist.IsNull() {
		output.SecurityProfile = &fleets.SecurityProfile{
			EncryptionAtHost: pointer.To(input.EncryptionAtHostEnabled),
		}
	}

	if v := input.OsDisk; len(v) > 0 && v[0].SecurityEncryptionType != "" {
		if output.SecurityProfile == nil {
			output.SecurityProfile = &fleets.SecurityProfile{}
		}
		output.SecurityProfile.SecurityType = pointer.To(fleets.SecurityTypesConfidentialVM)
		output.SecurityProfile.UefiSettings = &fleets.UefiSettings{
			SecureBootEnabled: pointer.To(input.SecureBootEnabled),
			VTpmEnabled:       pointer.To(input.VTpmEnabled),
		}
	} else {
		if input.SecureBootEnabled {
			if output.SecurityProfile == nil {
				output.SecurityProfile = &fleets.SecurityProfile{}
			}
			output.SecurityProfile.UefiSettings = &fleets.UefiSettings{
				SecureBootEnabled: pointer.To(input.SecureBootEnabled),
			}
			output.SecurityProfile.SecurityType = pointer.To(fleets.SecurityTypesTrustedLaunch)
		}

		if input.VTpmEnabled {
			if output.SecurityProfile == nil {
				output.SecurityProfile = &fleets.SecurityProfile{}
			}
			if output.SecurityProfile.UefiSettings == nil {
				output.SecurityProfile.UefiSettings = &fleets.UefiSettings{}
			}
			output.SecurityProfile.UefiSettings.VTpmEnabled = pointer.To(input.VTpmEnabled)

			output.SecurityProfile.SecurityType = pointer.To(fleets.SecurityTypesTrustedLaunch)
		}
	}

	extensionProfileValue, err := expandExtensionModel(input.Extension, input.ExtensionsTimeBudgetDuration)
	if err != nil {
		return nil, err
	}
	output.ExtensionProfile = extensionProfileValue

	output.LicenseType = pointer.To("None")
	if input.LicenseType != "" {
		output.LicenseType = pointer.To(input.LicenseType)
	}

	if input.UserDataBase64 != "" {
		output.UserData = pointer.To(input.UserDataBase64)
	}

	return &output, nil
}

func expandApplicationProfileModel(inputList []GalleryApplicationModel) *fleets.ApplicationProfile {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VMGalleryApplication, 0)
	for _, input := range inputList {
		output := fleets.VMGalleryApplication{
			EnableAutomaticUpgrade:          pointer.To(input.AutomaticUpgradeEnabled),
			Order:                           pointer.To(input.Order),
			PackageReferenceId:              input.VersionId,
			TreatFailureAsDeploymentFailure: pointer.To(input.TreatFailureAsDeploymentFailureEnabled),
		}

		if input.ConfigurationBlobUri != "" {
			output.ConfigurationReference = pointer.To(input.ConfigurationBlobUri)
		}

		if input.Tag != "" {
			output.Tags = pointer.To(input.Tag)
		}
		outputList = append(outputList, output)
	}

	return &fleets.ApplicationProfile{
		GalleryApplications: &outputList,
	}
}

func expandSubResource(input string) *fleets.SubResource {
	if input == "" {
		return nil
	}

	return &fleets.SubResource{
		Id: pointer.To(input),
	}
}

func expandSubResources(inputList []string) *[]fleets.SubResource {
	if len(inputList) == 0 {
		return nil
	}
	outputList := make([]fleets.SubResource, 0)
	for _, input := range inputList {
		output := expandSubResource(input)
		if output != nil {
			outputList = append(outputList, pointer.From(output))
		}
	}

	return &outputList
}

func expandExtensionModel(inputList []ExtensionModel, timeBudget string) (*fleets.VirtualMachineScaleSetExtensionProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	output := fleets.VirtualMachineScaleSetExtensionProfile{}
	extensions := make([]fleets.VirtualMachineScaleSetExtension, 0)
	for _, v := range inputList {
		extension := fleets.VirtualMachineScaleSetExtension{
			Name: pointer.To(v.Name),
			Properties: &fleets.VirtualMachineScaleSetExtensionProperties{
				AutoUpgradeMinorVersion: pointer.To(v.AutoUpgradeMinorVersionEnabled),
				SuppressFailures:        pointer.To(v.FailureSuppressionEnabled),
				EnableAutomaticUpgrade:  pointer.To(v.AutomaticUpgradeEnabled),
			},
		}

		if len(v.ProtectedSettingsFromKeyVault) > 0 {
			extension.Properties.ProtectedSettingsFromKeyVault = &fleets.KeyVaultSecretReference{
				SecretURL:   v.ProtectedSettingsFromKeyVault[0].SecretUrl,
				SourceVault: pointer.From(expandSubResource(v.ProtectedSettingsFromKeyVault[0].SourceVaultId)),
			}
		}

		if v.ForceExtensionExecutionOnChange != "" {
			extension.Properties.ForceUpdateTag = pointer.To(v.ForceExtensionExecutionOnChange)
		}

		if v.ProtectedSettingsJson != "" {
			protectedSettingsValue := make(map[string]interface{})
			err := json.Unmarshal([]byte(v.ProtectedSettingsJson), &protectedSettingsValue)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling `protected_settings_json`: %+v", err)
			}
			extension.Properties.ProtectedSettings = pointer.To(protectedSettingsValue)
		}

		if len(v.ExtensionsToProvisionAfterVmCreation) > 0 {
			extension.Properties.ProvisionAfterExtensions = pointer.To(v.ExtensionsToProvisionAfterVmCreation)
		}

		if v.Publisher != "" {
			extension.Properties.Publisher = pointer.To(v.Publisher)
		}

		if v.SettingsJson != "" {
			result := make(map[string]interface{})
			err := json.Unmarshal([]byte(v.SettingsJson), &result)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling `settings_json`: %+v", err)
			}
			extension.Properties.Settings = pointer.To(result)
		}

		if v.Type != "" {
			extension.Properties.Type = pointer.To(v.Type)
		}

		if v.TypeHandlerVersion != "" {
			extension.Properties.TypeHandlerVersion = pointer.To(v.TypeHandlerVersion)
		}

		extensions = append(extensions, extension)
	}

	output.Extensions = &extensions

	if timeBudget != "" {
		output.ExtensionsTimeBudget = pointer.To(timeBudget)
	}
	return &output, nil
}

func expandNetworkInterfaceModel(inputList []NetworkInterfaceModel) *[]fleets.VirtualMachineScaleSetNetworkConfiguration {
	outputList := make([]fleets.VirtualMachineScaleSetNetworkConfiguration, 0)
	for _, v := range inputList {
		output := fleets.VirtualMachineScaleSetNetworkConfiguration{
			Name: v.Name,
			Properties: &fleets.VirtualMachineScaleSetNetworkConfigurationProperties{
				EnableAcceleratedNetworking: pointer.To(v.AcceleratedNetworkingEnabled),
				EnableIPForwarding:          pointer.To(v.IPForwardingEnabled),
				NetworkSecurityGroup:        expandSubResource(v.NetworkSecurityGroupId),
				Primary:                     pointer.To(v.PrimaryNetworkInterfaceEnabled),
				IPConfigurations:            pointer.From(expandIPConfigurationModel(v.IPConfiguration)),
			},
		}

		if len(v.DnsServers) > 0 {
			output.Properties.DnsSettings = &fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
				DnsServers: pointer.To(v.DnsServers),
			}
		}

		auxiliaryMode := fleets.NetworkInterfaceAuxiliaryModeNone
		if v.AuxiliaryMode != "" {
			auxiliaryMode = fleets.NetworkInterfaceAuxiliaryMode(v.AuxiliaryMode)
		}
		output.Properties.AuxiliaryMode = pointer.To(auxiliaryMode)

		if v.DeleteOption != "" {
			output.Properties.DeleteOption = pointer.To(fleets.DeleteOptions(v.DeleteOption))
		}

		auxiliarySku := fleets.NetworkInterfaceAuxiliarySkuNone
		if v.AuxiliarySku != "" {
			auxiliarySku = fleets.NetworkInterfaceAuxiliarySku(v.AuxiliarySku)
		}
		output.Properties.AuxiliarySku = pointer.To(auxiliarySku)

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandIPConfigurationModel(inputList []IPConfigurationModel) *[]fleets.VirtualMachineScaleSetIPConfiguration {
	outputList := make([]fleets.VirtualMachineScaleSetIPConfiguration, 0)
	for _, input := range inputList {
		output := fleets.VirtualMachineScaleSetIPConfiguration{
			Name: input.Name,
			Properties: &fleets.VirtualMachineScaleSetIPConfigurationProperties{
				ApplicationGatewayBackendAddressPools: expandSubResources(input.ApplicationGatewayBackendAddressPoolIds),
				ApplicationSecurityGroups:             expandSubResources(input.ApplicationSecurityGroupIds),
				LoadBalancerBackendAddressPools:       expandSubResources(input.LoadBalancerBackendAddressPoolIds),
				Primary:                               pointer.To(input.PrimaryIpConfigurationEnabled),
				PublicIPAddressConfiguration:          expandPublicIPAddressModel(input.PublicIPAddress),
			},
		}

		if input.SubnetId != "" {
			output.Properties.Subnet = &fleets.ApiEntityReference{
				Id: pointer.To(input.SubnetId),
			}
		}

		if input.Version != "" {
			output.Properties.PrivateIPAddressVersion = pointer.To(fleets.IPVersion(input.Version))
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandPublicIPAddressModel(inputList []PublicIPAddressModel) *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetPublicIPAddressConfiguration{
		Name: input.Name,
		Properties: &fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
			PublicIPPrefix: expandSubResource(input.PublicIPPrefix),
		},
	}

	if input.DomainNameLabel != "" {
		output.Properties.DnsSettings = &fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
			DomainNameLabel: input.DomainNameLabel,
		}
	}

	if input.DomainNameLabelScope != "" {
		if output.Properties.DnsSettings == nil {
			output.Properties.DnsSettings = &fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{}
		}
		output.Properties.DnsSettings.DomainNameLabelScope = pointer.To(fleets.DomainNameLabelScopeTypes(input.DomainNameLabelScope))
	}

	if v := input.SkuName; v != "" {
		skuParts := strings.Split(v, "_")
		output.Sku = &fleets.PublicIPAddressSku{
			Name: pointer.To(fleets.PublicIPAddressSkuName(skuParts[0])),
			Tier: pointer.To(fleets.PublicIPAddressSkuTier(skuParts[1])),
		}
	}

	if input.IdleTimeoutInMinutes > 0 {
		output.Properties.IdleTimeoutInMinutes = pointer.To(input.IdleTimeoutInMinutes)
	}
	if input.DeleteOption != "" {
		output.Properties.DeleteOption = pointer.To(fleets.DeleteOptions(input.DeleteOption))
	}
	if input.Version != "" {
		output.Properties.PublicIPAddressVersion = pointer.To(fleets.IPVersion(input.Version))
	}

	return &output
}

func expandOSProfileModel(inputList []VirtualMachineProfileModel) *fleets.VirtualMachineScaleSetOSProfile {
	osProfile := &inputList[0].OsProfile[0]
	output := fleets.VirtualMachineScaleSetOSProfile{
		AllowExtensionOperations: pointer.To(inputList[0].ExtensionOperationsEnabled),
	}
	if osProfile.CustomDataBase64 != "" {
		output.CustomData = pointer.To(osProfile.CustomDataBase64)
	}

	if lConfig := osProfile.LinuxConfiguration; len(lConfig) > 0 {
		linuxConfig := fleets.LinuxConfiguration{
			DisablePasswordAuthentication: pointer.To(!lConfig[0].PasswordAuthenticationEnabled),
			ProvisionVMAgent:              pointer.To(lConfig[0].ProvisionVMAgentEnabled),
			EnableVMAgentPlatformUpdates:  pointer.To(lConfig[0].VMAgentPlatformUpdatesEnabled),
			PatchSettings: &fleets.LinuxPatchSettings{
				// 'AssessmentMode' can only be set `ImageDefault` on Virtual Machine Scale Sets.
				AssessmentMode: pointer.To(fleets.LinuxPatchAssessmentModeImageDefault),
			},
		}

		// AutomaticByPlatformSettings cannot be set if the PatchMode is not `AutomaticByPlatform`
		if lConfig[0].PatchMode == string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform) {
			linuxConfig.PatchSettings.AutomaticByPlatformSettings = &fleets.LinuxVMGuestPatchAutomaticByPlatformSettings{
				BypassPlatformSafetyChecksOnUserSchedule: pointer.To(lConfig[0].BypassPlatformSafetyChecksEnabled),
			}

			if lConfig[0].PatchRebooting != "" {
				linuxConfig.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.To(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSetting(lConfig[0].PatchRebooting))
			}
		}
		if lConfig[0].PatchMode != "" {
			linuxConfig.PatchSettings.PatchMode = pointer.To(fleets.LinuxVMGuestPatchMode(lConfig[0].PatchMode))
		}
		if lConfig[0].AdminUsername != "" {
			output.AdminUsername = pointer.To(lConfig[0].AdminUsername)
		}
		if lConfig[0].AdminPassword != "" {
			output.AdminPassword = pointer.To(lConfig[0].AdminPassword)
		}
		if lConfig[0].ComputerNamePrefix != "" {
			output.ComputerNamePrefix = pointer.To(lConfig[0].ComputerNamePrefix)
		}
		output.Secrets = expandOsProfileLinuxSecretsModel(lConfig[0].Secret)

		if lConfig[0].AdminUsername != "" || len(lConfig[0].AdminSSHKeys) > 0 {
			publicKeys := make([]fleets.SshPublicKey, 0)
			for _, v := range lConfig[0].AdminSSHKeys {
				output := fleets.SshPublicKey{
					Path: pointer.To(fmt.Sprintf("/home/%s/.ssh/authorized_keys", lConfig[0].AdminUsername)),
				}
				if v != "" {
					output.KeyData = pointer.To(v)
				}
				publicKeys = append(publicKeys, output)
			}

			linuxConfig.Ssh = &fleets.SshConfiguration{
				PublicKeys: pointer.To(publicKeys),
			}
		}

		output.LinuxConfiguration = &linuxConfig
	}

	if winConfig := osProfile.WindowsConfiguration; len(winConfig) > 0 {
		windowsConfig := fleets.WindowsConfiguration{
			AdditionalUnattendContent:    expandAdditionalUnAttendContentModel(winConfig[0].AdditionalUnattendContent),
			EnableAutomaticUpdates:       pointer.To(winConfig[0].AutomaticUpdatesEnabled),
			EnableVMAgentPlatformUpdates: pointer.To(winConfig[0].VMAgentPlatformUpdatesEnabled),
			ProvisionVMAgent:             pointer.To(winConfig[0].ProvisionVMAgentEnabled),
			PatchSettings: &fleets.PatchSettings{
				EnableHotpatching: pointer.To(winConfig[0].HotPatchingEnabled),
				// 'AssessmentMode' can only be set `ImageDefault` on Virtual Machine Scale Sets.
				AssessmentMode: pointer.To(fleets.WindowsPatchAssessmentModeImageDefault),
			},
		}

		if winRm := winConfig[0].WinRM; len(winRm) > 0 {
			listenerList := make([]fleets.WinRMListener, 0)
			for _, v := range winRm {
				output := fleets.WinRMListener{
					Protocol: pointer.To(fleets.ProtocolTypes(v.Protocol)),
				}

				if v.CertificateUrl != "" {
					output.CertificateURL = pointer.To(v.CertificateUrl)
				}
				listenerList = append(listenerList, output)
			}
			windowsConfig.WinRM = &fleets.WinRMConfiguration{
				Listeners: pointer.To(listenerList),
			}
		}

		if winConfig[0].AdminUsername != "" {
			output.AdminUsername = pointer.To(winConfig[0].AdminUsername)
		}
		if winConfig[0].AdminPassword != "" {
			output.AdminPassword = pointer.To(winConfig[0].AdminPassword)
		}
		if winConfig[0].ComputerNamePrefix != "" {
			output.ComputerNamePrefix = pointer.To(winConfig[0].ComputerNamePrefix)
		}

		// AutomaticByPlatformSettings cannot be set if the PatchMode is not `AutomaticByPlatform`
		if winConfig[0].PatchMode == string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
			windowsConfig.PatchSettings.AutomaticByPlatformSettings = &fleets.WindowsVMGuestPatchAutomaticByPlatformSettings{
				BypassPlatformSafetyChecksOnUserSchedule: pointer.To(winConfig[0].BypassPlatformSafetyChecksEnabled),
			}

			if winConfig[0].PatchRebooting != "" {
				windowsConfig.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.To(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSetting(winConfig[0].PatchRebooting))
			}
		}
		if winConfig[0].PatchMode != "" {
			windowsConfig.PatchSettings.PatchMode = pointer.To(fleets.WindowsVMGuestPatchMode(winConfig[0].PatchMode))
		}
		if winConfig[0].TimeZone != "" {
			windowsConfig.TimeZone = pointer.To(winConfig[0].TimeZone)
		}
		output.WindowsConfiguration = &windowsConfig

		output.Secrets = expandOsProfileWindowsSecretsModel(winConfig[0].Secret)
	}

	return &output
}

func validateWindowsSetting(inputList []VirtualMachineProfileModel, d *schema.ResourceDiff) error {
	if len(inputList) == 0 || len(inputList[0].OsProfile) == 0 {
		return nil
	}

	input := &inputList[0]
	if v := input.OsProfile[0].WindowsConfiguration; len(v) > 0 {
		patchMode := v[0].PatchMode
		hotPatchingEnabled := v[0].HotPatchingEnabled
		provisionVMAgentEnabled := v[0].ProvisionVMAgentEnabled

		rebootSetting := v[0].PatchRebooting
		bypassPlatformSafetyChecksEnabledExist := d.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["os_profile"].AsValueSlice()[0].AsValueMap()["windows_configuration"].AsValueSlice()[0].AsValueMap()["bypass_platform_safety_checks_enabled"]
		if !bypassPlatformSafetyChecksEnabledExist.IsNull() || rebootSetting != "" {
			if patchMode != string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
				return fmt.Errorf("`bypass_platform_safety_checks_enabled` and `patch_rebooting` cannot be set if the `PatchMode` is not `AutomaticByPlatform`")
			}
		}

		if input.ExtensionOperationsEnabled && !provisionVMAgentEnabled {
			return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent_enabled` is set to `false`")
		}

		isHotPatchEnabledImage := isValidHotPatchSourceImageReference(input.SourceImageReference)
		hasHealthExtension := false
		if v := input.Extension; len(v) > 0 && (v[0].Type == "ApplicationHealthLinux" || v[0].Type == "ApplicationHealthWindows") {
			hasHealthExtension = true
		}

		if isHotPatchEnabledImage {
			// it is a hot patching enabled image, validate hot patching enabled settings
			if patchMode != string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
				return fmt.Errorf("when referencing a hot patching enabled image the `patch_mode` field must always be set to %q", fleets.WindowsVMGuestPatchModeAutomaticByPlatform)
			}
			if !provisionVMAgentEnabled {
				return fmt.Errorf("when referencing a hot patching enabled image the `provision_vm_agent_enabled` field must always be set to `true`")
			}
			if !hasHealthExtension {
				return fmt.Errorf("when referencing a hot patching enabled image the `extension` field must always contain a `application health extension`")
			}
			if !hotPatchingEnabled {
				return fmt.Errorf("when referencing a hot patching enabled image the `hot_patching_enabled` field must always be set to `true`")
			}
		} else {
			// not a hot patching enabled image verify Automatic VM Guest Patching settings
			if patchMode == string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
				if !provisionVMAgentEnabled {
					return fmt.Errorf("when `patch_mode` is set to %q then `provision_vm_agent_enabled` must be set to `true`", patchMode)
				}
				if !hasHealthExtension {
					return fmt.Errorf("when `patch_mode` is set to %q then the `extension` field must always contain a `application health extension`", patchMode)
				}
			}

			if hotPatchingEnabled {
				return fmt.Errorf("`hot_patching_enabled` field is not supported unless you are using one of the following hot patching enable images, `2022-datacenter-azure-edition-core`, `2022-datacenter-azure-edition-core-smalldisk`, `2022-datacenter-azure-edition-hotpatch` or `2022-datacenter-azure-edition-hotpatch-smalldisk`")
			}
		}
	}
	return nil
}

func validateSecuritySetting(inputList []VirtualMachineProfileModel) error {
	if len(inputList) == 0 || len(inputList[0].OsProfile) == 0 {
		return nil
	}

	input := &inputList[0]
	if v := input.OsDisk; len(v) > 0 {
		secureBootEnabled := input.SecureBootEnabled
		vTpmEnabled := input.VTpmEnabled
		if v[0].SecurityEncryptionType != "" {
			if fleets.SecurityEncryptionTypesDiskWithVMGuestState == fleets.SecurityEncryptionTypes(v[0].SecurityEncryptionType) && (!secureBootEnabled || !vTpmEnabled) {
				return fmt.Errorf("`secure_boot_enabled` and `vtpm_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
			}
			if !vTpmEnabled {
				return fmt.Errorf("`vtpm_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set")
			}
		}
	}
	return nil
}

func validateLinuxSetting(inputList []VirtualMachineProfileModel, d *schema.ResourceDiff) error {
	if len(inputList) == 0 || len(inputList[0].OsProfile) == 0 {
		return nil
	}

	input := &inputList[0]
	if v := input.OsProfile[0].LinuxConfiguration; len(v) > 0 {
		patchMode := v[0].PatchMode
		provisionVMAgentEnabled := v[0].ProvisionVMAgentEnabled

		rebootSetting := v[0].PatchRebooting
		bypassPlatformSafetyChecksEnabledExist := d.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["os_profile"].AsValueSlice()[0].AsValueMap()["linux_configuration"].AsValueSlice()[0].AsValueMap()["bypass_platform_safety_checks_enabled"]
		if !bypassPlatformSafetyChecksEnabledExist.IsNull() || rebootSetting != "" {
			if patchMode != string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform) {
				return fmt.Errorf("`bypass_platform_safety_checks_enabled` and `patch_rebooting` cannot be set if the `PatchMode` is not `AutomaticByPlatform`")
			}
		}

		if input.ExtensionOperationsEnabled && !provisionVMAgentEnabled {
			return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent_enabled` is set to `false`")
		}

		hasHealthExtension := false
		if v := input.Extension; len(v) > 0 && (v[0].Type == "ApplicationHealthLinux" || v[0].Type == "ApplicationHealthWindows") {
			hasHealthExtension = true
		}

		if patchMode == string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform) {
			if !provisionVMAgentEnabled {
				return fmt.Errorf("when the `patch_mode` field is set to %q the `provision_vm_agent_enabled` field must always be set to `true`, got %q", patchMode, strconv.FormatBool(provisionVMAgentEnabled))
			}

			if !hasHealthExtension {
				return fmt.Errorf("when the `patch_mode` field is set to %q the `extension` field must contain at least one `application health extension`, got 0", patchMode)
			}
		}

		if v[0].AdminPassword == "" && v[0].PasswordAuthenticationEnabled {
			return fmt.Errorf("`admin_password` is required when `password_authentication_enabled` is enabled")
		}
	}
	return nil
}

func isValidHotPatchSourceImageReference(referenceInput []SourceImageReferenceModel) bool {
	if len(referenceInput) == 0 {
		return false
	}
	raw := referenceInput[0]
	pub := raw.Publisher
	offer := raw.Offer
	sku := raw.Sku

	if pub == "MicrosoftWindowsServer" && offer == "WindowsServer" && (sku == "2022-datacenter-azure-edition-core" || sku == "2022-datacenter-azure-edition-core-smalldisk" || sku == "2022-datacenter-azure-edition-hotpatch" || sku == "2022-datacenter-azure-edition-hotpatch-smalldisk") {
		return true
	}

	return false
}

func expandOsProfileLinuxSecretsModel(inputList []LinuxSecretModel) *[]fleets.VaultSecretGroup {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VaultSecretGroup, 0)
	for _, v := range inputList {
		output := fleets.VaultSecretGroup{
			SourceVault: expandSubResource(v.KeyVaultId),
		}

		if len(v.Certificate) > 0 {
			vcs := make([]fleets.VaultCertificate, 0)
			for _, v := range v.Certificate {
				vc := fleets.VaultCertificate{}
				if v.Url != "" {
					vc.CertificateURL = pointer.To(v.Url)
					vcs = append(vcs, vc)
				}
			}
			output.VaultCertificates = &vcs
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandOsProfileWindowsSecretsModel(inputList []WindowsSecretModel) *[]fleets.VaultSecretGroup {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VaultSecretGroup, 0)
	for _, v := range inputList {
		output := fleets.VaultSecretGroup{
			SourceVault: expandSubResource(v.KeyVaultId),
		}

		if len(v.Certificate) > 0 {
			vcs := make([]fleets.VaultCertificate, 0)
			for _, v := range v.Certificate {
				vc := fleets.VaultCertificate{}
				if v.Store != "" {
					vc.CertificateStore = pointer.To(v.Store)
				}
				if v.Url != "" {
					vc.CertificateURL = pointer.To(v.Url)
				}
				vcs = append(vcs, vc)
			}

			output.VaultCertificates = &vcs
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandAdditionalUnAttendContentModel(inputList []AdditionalUnattendContentModel) *[]fleets.AdditionalUnattendContent {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.AdditionalUnattendContent, 0)
	for _, input := range inputList {
		output := fleets.AdditionalUnattendContent{
			SettingName: pointer.To(fleets.SettingNames(input.Setting)),
			Content:     pointer.To(input.Xml),
			// no other possible values
			ComponentName: pointer.To(fleets.ComponentNameMicrosoftNegativeWindowsNegativeShellNegativeSetup),
			PassName:      pointer.To(fleets.PassNameOobeSystem),
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandDataDiskModel(inputList []DataDiskModel) *[]fleets.VirtualMachineScaleSetDataDisk {
	outputList := make([]fleets.VirtualMachineScaleSetDataDisk, 0)
	for _, input := range inputList {
		output := fleets.VirtualMachineScaleSetDataDisk{
			CreateOption:            fleets.DiskCreateOptionTypes(input.CreateOption),
			Lun:                     input.Lun,
			WriteAcceleratorEnabled: pointer.To(input.WriteAcceleratorEnabled),
		}

		if input.DeleteOption != "" {
			output.DeleteOption = pointer.To(fleets.DiskDeleteOptionTypes(input.DeleteOption))
		}

		if input.DiskSizeInGiB > 0 {
			output.DiskSizeGB = pointer.To(input.DiskSizeInGiB)
		}

		caching := string(fleets.CachingTypesNone)
		if input.Caching != "" {
			caching = input.Caching
		}
		output.Caching = pointer.To(fleets.CachingTypes(caching))

		if input.StorageAccountType != "" {
			output.ManagedDisk = &fleets.VirtualMachineScaleSetManagedDiskParameters{
				StorageAccountType: pointer.To(fleets.StorageAccountTypes(input.StorageAccountType)),
			}
		}

		if input.DiskEncryptionSetId != "" {
			if output.ManagedDisk == nil {
				output.ManagedDisk = &fleets.VirtualMachineScaleSetManagedDiskParameters{}
			}
			output.ManagedDisk.DiskEncryptionSet = &fleets.DiskEncryptionSetParameters{
				Id: pointer.To(input.DiskEncryptionSetId),
			}
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandImageReference(inputList []SourceImageReferenceModel, imageId string) *fleets.ImageReference {
	if imageId != "" {
		// With Version            : "/communityGalleries/publicGalleryName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/communityGalleries/publicGalleryName/images/myGalleryImageName"
		if _, errors := validation.Any(computeValidate.CommunityGalleryImageID, computeValidate.CommunityGalleryImageVersionID)(imageId, "source_image_id"); len(errors) == 0 {
			return &fleets.ImageReference{
				CommunityGalleryImageId: pointer.To(imageId),
			}
		}

		// Shared Image Gallery with Cross-Tenant Sharing
		// With Version            : "/sharedGalleries/galleryUniqueName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/sharedGalleries/galleryUniqueName/images/myGalleryImageName"
		if _, errors := validation.Any(computeValidate.SharedGalleryImageID, computeValidate.SharedGalleryImageVersionID)(imageId, "source_image_id"); len(errors) == 0 {
			return &fleets.ImageReference{
				SharedGalleryImageId: pointer.To(imageId),
			}
		}

		return &fleets.ImageReference{
			// Standard Shared Image ID: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{imageDefinitionName}/versions/{imageVersion}
			// Standard Image ID: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.Compute/images/{imageName}
			Id: pointer.To(imageId),
		}
	}

	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	return &fleets.ImageReference{
		Publisher: pointer.To(input.Publisher),
		Offer:     pointer.To(input.Offer),
		Sku:       pointer.To(input.Sku),
		Version:   pointer.To(input.Version),
	}
}

func expandOSDiskModel(input *VirtualMachineProfileModel) *fleets.VirtualMachineScaleSetOSDisk {
	osType := fleets.OperatingSystemTypesLinux
	if len(input.OsProfile) > 0 && len(input.OsProfile[0].WindowsConfiguration) > 0 {
		osType = fleets.OperatingSystemTypesWindows
	}

	if input == nil || len(input.OsDisk) == 0 {
		return nil
	}

	inputOsDisk := &input.OsDisk[0]
	output := fleets.VirtualMachineScaleSetOSDisk{
		DeleteOption:            pointer.To(fleets.DiskDeleteOptionTypes(inputOsDisk.DeleteOption)),
		OsType:                  pointer.To(osType),
		WriteAcceleratorEnabled: pointer.To(inputOsDisk.WriteAcceleratorEnabled),
		ManagedDisk: &fleets.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: pointer.To(fleets.StorageAccountTypes(inputOsDisk.StorageAccountType)),
		},
		// these have to be hard-coded so there's no point exposing them
		CreateOption: fleets.DiskCreateOptionTypesFromImage,
	}

	if inputOsDisk != nil {
		if inputOsDisk.DiffDiskOption != "" {
			output.DiffDiskSettings = &fleets.DiffDiskSettings{
				Option: pointer.To(fleets.DiffDiskOptions(inputOsDisk.DiffDiskOption)),
			}
		}

		if inputOsDisk.DiffDiskPlacement != "" {
			if output.DiffDiskSettings == nil {
				output.DiffDiskSettings = &fleets.DiffDiskSettings{}
			}
			output.DiffDiskSettings.Placement = pointer.To(fleets.DiffDiskPlacement(inputOsDisk.DiffDiskPlacement))
		}
	}

	if inputOsDisk.DiskSizeInGiB > 0 {
		output.DiskSizeGB = pointer.To(inputOsDisk.DiskSizeInGiB)
	}

	caching := fleets.CachingTypesNone
	if v := inputOsDisk.Caching; v != "" {
		caching = fleets.CachingTypes(v)
	}
	output.Caching = pointer.To(caching)

	if inputOsDisk.DiskEncryptionSetId != "" {
		output.ManagedDisk.DiskEncryptionSet = &fleets.DiskEncryptionSetParameters{
			Id: pointer.To(inputOsDisk.DiskEncryptionSetId),
		}
	}

	if inputOsDisk.SecurityEncryptionType != "" {
		output.ManagedDisk.SecurityProfile = &fleets.VMDiskSecurityProfile{
			SecurityEncryptionType: pointer.To(fleets.SecurityEncryptionTypes(inputOsDisk.SecurityEncryptionType)),
		}
	}

	return &output
}

func flattenVirtualMachineProfileModel(input *fleets.BaseVirtualMachineProfile, metadata sdk.ResourceMetaData) ([]VirtualMachineProfileModel, error) {
	outputList := make([]VirtualMachineProfileModel, 0)
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineProfileModel{
		GalleryApplicationProfile: flattenApplicationProfileModel(input.ApplicationProfile),
		NetworkInterface:          flattenNetworkInterfaceModel(input.NetworkProfile),
		UserDataBase64:            pointer.From(input.UserData),
	}

	if v := input.NetworkProfile; v != nil {
		output.NetworkApiVersion = string(pointer.From(v.NetworkApiVersion))
	}
	if v := input.SecurityProfile; v != nil {
		output.EncryptionAtHostEnabled = pointer.From(v.EncryptionAtHost)
		if v.UefiSettings != nil {
			output.SecureBootEnabled = pointer.From(v.UefiSettings.SecureBootEnabled)
			output.VTpmEnabled = pointer.From(v.UefiSettings.VTpmEnabled)
		}
	}

	if v := input.OsProfile; v != nil {
		osProfile, err := flattenOSProfileModel(v, metadata.ResourceData)
		if err != nil {
			return outputList, err
		}
		output.OsProfile = osProfile
		output.ExtensionOperationsEnabled = pointer.From(v.AllowExtensionOperations)
	}

	if v := input.StorageProfile; v != nil {
		output.DataDisks = flattenDataDiskModel(v.DataDisks)
		storageImageId := ""
		if v.ImageReference != nil {
			if v.ImageReference.Id != nil {
				storageImageId = pointer.From(v.ImageReference.Id)
			}
			if v.ImageReference.CommunityGalleryImageId != nil {
				storageImageId = *v.ImageReference.CommunityGalleryImageId
			}
			if v.ImageReference.SharedGalleryImageId != nil {
				storageImageId = *v.ImageReference.SharedGalleryImageId
			}
		}
		output.SourceImageId = storageImageId
		output.SourceImageReference = flattenImageReference(v.ImageReference, storageImageId != "")
		output.OsDisk = flattenOSDiskModel(v.OsDisk)
	}

	if se := input.ScheduledEventsProfile; se != nil {
		if v := se.TerminateNotificationProfile; v != nil {
			output.ScheduledEventTerminationTimeoutDuration = pointer.From(v.NotBeforeTimeout)
		}
		if v := se.OsImageNotificationProfile; v != nil {
			output.ScheduledEventOsImageTimeoutDuration = pointer.From(v.NotBeforeTimeout)
		}
	}

	if cr := input.CapacityReservation; cr != nil {
		if v := cr.CapacityReservationGroup; v != nil {
			output.CapacityReservationGroupId = pointer.From(v.Id)
		}
	}

	if dp := input.DiagnosticsProfile; dp != nil {
		if v := dp.BootDiagnostics; v != nil {
			output.BootDiagnosticEnabled = pointer.From(v.Enabled)
			output.BootDiagnosticStorageAccountEndpoint = pointer.From(v.StorageUri)
		}
	}

	extensionProfileValue, err := flattenExtensionModel(input.ExtensionProfile, metadata)
	if err != nil {
		return nil, err
	}
	output.Extension = extensionProfileValue

	if input.ExtensionProfile != nil {
		output.ExtensionsTimeBudgetDuration = pointer.From(input.ExtensionProfile.ExtensionsTimeBudget)
	}

	licenseType := ""
	if v := pointer.From(input.LicenseType); v != "None" {
		licenseType = v
	}
	output.LicenseType = licenseType

	return append(outputList, output), nil
}

func flattenAdminSshKeyModel(input *fleets.SshConfiguration) ([]string, error) {
	outputList := make([]string, 0)
	if input == nil || input.PublicKeys == nil {
		return outputList, nil
	}

	for _, input := range *input.PublicKeys {
		username := parseUsernameFromAuthorizedKeysPath(*input.Path)
		if username == nil {
			return nil, fmt.Errorf("parsing username from %q", pointer.From(input.Path))
		}
		outputList = append(outputList, pointer.From(input.KeyData))
	}

	return outputList, nil
}

func parseUsernameFromAuthorizedKeysPath(input string) *string {
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
			value := values[i]
			return &value
		}
	}

	return nil
}

func flattenApplicationProfileModel(input *fleets.ApplicationProfile) []GalleryApplicationModel {
	outputList := make([]GalleryApplicationModel, 0)
	if input == nil {
		return outputList
	}

	for _, input := range *input.GalleryApplications {
		output := GalleryApplicationModel{
			VersionId:                              input.PackageReferenceId,
			ConfigurationBlobUri:                   pointer.From(input.ConfigurationReference),
			AutomaticUpgradeEnabled:                pointer.From(input.EnableAutomaticUpgrade),
			Order:                                  pointer.From(input.Order),
			Tag:                                    pointer.From(input.Tags),
			TreatFailureAsDeploymentFailureEnabled: pointer.From(input.TreatFailureAsDeploymentFailure),
		}
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenNetworkInterfaceModel(input *fleets.VirtualMachineScaleSetNetworkProfile) []NetworkInterfaceModel {
	outputList := make([]NetworkInterfaceModel, 0)
	if input == nil {
		return outputList
	}

	for _, input := range *input.NetworkInterfaceConfigurations {
		output := NetworkInterfaceModel{
			Name: input.Name,
		}

		if props := input.Properties; props != nil {
			if v := props.AuxiliaryMode; v != nil && *v != fleets.NetworkInterfaceAuxiliaryModeNone {
				output.AuxiliaryMode = string(*v)
			}

			if v := props.AuxiliarySku; v != nil && *v != fleets.NetworkInterfaceAuxiliarySkuNone {
				output.AuxiliarySku = string(*v)
			}

			output.DeleteOption = string(pointer.From(props.DeleteOption))

			if v := props.DnsSettings; v != nil {
				output.DnsServers = pointer.From(v.DnsServers)
			}

			output.AcceleratedNetworkingEnabled = pointer.From(props.EnableAcceleratedNetworking)

			output.IPForwardingEnabled = pointer.From(props.EnableIPForwarding)

			output.IPConfiguration = flattenIPConfigurationModel(props.IPConfigurations)

			if v := props.NetworkSecurityGroup; v != nil {
				output.NetworkSecurityGroupId = pointer.From(v.Id)
			}

			output.PrimaryNetworkInterfaceEnabled = pointer.From(props.Primary)
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenOsProfileLinuxSecretsModel(inputList *[]fleets.VaultSecretGroup) []LinuxSecretModel {
	outputList := make([]LinuxSecretModel, 0)
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := LinuxSecretModel{
			Certificate: flattenLinuxVaultCertificateModel(input.VaultCertificates),
		}
		if v := input.SourceVault; v != nil {
			output.KeyVaultId = pointer.From(v.Id)
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenOsProfileWindowsSecretsModel(inputList *[]fleets.VaultSecretGroup) []WindowsSecretModel {
	outputList := make([]WindowsSecretModel, 0)
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := WindowsSecretModel{
			Certificate: flattenWindowsVaultCertificateModel(input.VaultCertificates),
		}
		if v := input.SourceVault; v != nil {
			output.KeyVaultId = pointer.From(v.Id)
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenOSProfileModel(input *fleets.VirtualMachineScaleSetOSProfile, d *schema.ResourceData) ([]OSProfileModel, error) {
	outputList := make([]OSProfileModel, 0)
	if input == nil {
		return outputList, nil
	}

	output := OSProfileModel{}
	if input.CustomData != nil {
		output.CustomDataBase64 = pointer.From(input.CustomData)
	} else {
		output.CustomDataBase64 = utils.Base64EncodeIfNot(d.Get("virtual_machine_profile.0.os_profile.0.custom_data_base64").(string))
	}

	windowsConfigs := make([]WindowsConfigurationModel, 0)
	if v := input.WindowsConfiguration; v != nil {
		windowsConfig := WindowsConfigurationModel{
			AdditionalUnattendContent:     flattenAdditionalUnAttendContentModel(v.AdditionalUnattendContent, d),
			WinRM:                         flattenWinRMModel(v.WinRM),
			AdminUsername:                 pointer.From(input.AdminUsername),
			AdminPassword:                 d.Get("virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password").(string),
			AutomaticUpdatesEnabled:       pointer.From(v.EnableAutomaticUpdates),
			ComputerNamePrefix:            pointer.From(input.ComputerNamePrefix),
			VMAgentPlatformUpdatesEnabled: pointer.From(v.EnableVMAgentPlatformUpdates),
			ProvisionVMAgentEnabled:       pointer.From(v.ProvisionVMAgent),
			TimeZone:                      pointer.From(v.TimeZone),
			Secret:                        flattenOsProfileWindowsSecretsModel(input.Secrets),
		}

		if p := v.PatchSettings; p != nil {
			windowsConfig.PatchMode = string(pointer.From(p.PatchMode))
			if a := p.AutomaticByPlatformSettings; a != nil {
				windowsConfig.BypassPlatformSafetyChecksEnabled = pointer.From(a.BypassPlatformSafetyChecksOnUserSchedule)
				windowsConfig.PatchRebooting = string(pointer.From(a.RebootSetting))
			}
			windowsConfig.HotPatchingEnabled = pointer.From(p.EnableHotpatching)
		}
		windowsConfigs = append(windowsConfigs, windowsConfig)
	}
	output.WindowsConfiguration = windowsConfigs

	linuxConfigs := make([]LinuxConfigurationModel, 0)
	if v := input.LinuxConfiguration; v != nil {
		linuxConfig := LinuxConfigurationModel{
			AdminUsername:                 pointer.From(input.AdminUsername),
			AdminPassword:                 d.Get("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password").(string),
			ComputerNamePrefix:            pointer.From(input.ComputerNamePrefix),
			PasswordAuthenticationEnabled: !pointer.From(v.DisablePasswordAuthentication),
			VMAgentPlatformUpdatesEnabled: pointer.From(v.EnableVMAgentPlatformUpdates),
			ProvisionVMAgentEnabled:       pointer.From(v.ProvisionVMAgent),
			Secret:                        flattenOsProfileLinuxSecretsModel(input.Secrets),
		}

		if p := v.PatchSettings; p != nil {
			linuxConfig.PatchMode = string(pointer.From(p.PatchMode))
			if a := p.AutomaticByPlatformSettings; a != nil {
				linuxConfig.BypassPlatformSafetyChecksEnabled = pointer.From(a.BypassPlatformSafetyChecksOnUserSchedule)
				linuxConfig.PatchRebooting = string(pointer.From(a.RebootSetting))
			}
		}

		flattenedSSHPublicKeys, err := flattenAdminSshKeyModel(v.Ssh)
		if err != nil {
			return nil, fmt.Errorf("flattening `linux_configuration.0.admin_ssh_keys`: %+v", err)
		}
		linuxConfig.AdminSSHKeys = flattenedSSHPublicKeys

		linuxConfigs = append(linuxConfigs, linuxConfig)
	}

	output.LinuxConfiguration = linuxConfigs

	return append(outputList, output), nil
}

func flattenLinuxVaultCertificateModel(inputList *[]fleets.VaultCertificate) []LinuxCertificateModel {
	outputList := make([]LinuxCertificateModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := LinuxCertificateModel{
			Url: pointer.From(input.CertificateURL),
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenWindowsVaultCertificateModel(inputList *[]fleets.VaultCertificate) []WindowsCertificateModel {
	outputList := make([]WindowsCertificateModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := WindowsCertificateModel{
			Store: pointer.From(input.CertificateStore),
			Url:   pointer.From(input.CertificateURL),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenAdditionalUnAttendContentModel(inputList *[]fleets.AdditionalUnattendContent, d *schema.ResourceData) []AdditionalUnattendContentModel {
	outputList := make([]AdditionalUnattendContentModel, 0)
	if inputList == nil {
		return outputList
	}
	for i, input := range *inputList {
		output := AdditionalUnattendContentModel{
			Setting: string(pointer.From(input.SettingName)),
		}
		existing := make([]interface{}, 0)
		if v, ok := d.GetOk("virtual_machine_profile.0.os_profile.0.windows_configuration.0.additional_unattend_content"); ok {
			existing = v.([]interface{})
		}

		// content isn't returned by the API since it's sensitive data so we need to pull from the state file.
		content := ""
		if len(existing) > i {
			existingVal := existing[i]
			existingRaw, ok := existingVal.(map[string]interface{})
			if ok {
				contentRaw, ok := existingRaw["xml"]
				if ok {
					content = contentRaw.(string)
				}
			}
		}
		output.Xml = content

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenWinRMModel(input *fleets.WinRMConfiguration) []WinRMModel {
	outputList := make([]WinRMModel, 0)
	if input == nil || input.Listeners == nil {
		return outputList
	}

	for _, input := range *input.Listeners {
		output := WinRMModel{
			CertificateUrl: pointer.From(input.CertificateURL),
			Protocol:       string(pointer.From(input.Protocol)),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenDataDiskModel(inputList *[]fleets.VirtualMachineScaleSetDataDisk) []DataDiskModel {
	outputList := make([]DataDiskModel, 0)
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := DataDiskModel{
			CreateOption:            string(input.CreateOption),
			Lun:                     input.Lun,
			DeleteOption:            string(pointer.From(input.DeleteOption)),
			DiskSizeInGiB:           pointer.From(input.DiskSizeGB),
			WriteAcceleratorEnabled: pointer.From(input.WriteAcceleratorEnabled),
		}

		caching := ""
		if v := input.Caching; v != nil && *v != fleets.CachingTypesNone {
			caching = string(*v)
		}
		output.Caching = caching

		if md := input.ManagedDisk; md != nil {
			if v := md.DiskEncryptionSet; v != nil {
				output.DiskEncryptionSetId = pointer.From(v.Id)
			}
			output.StorageAccountType = string(pointer.From(md.StorageAccountType))
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenImageReference(input *fleets.ImageReference, hasImageId bool) []SourceImageReferenceModel {
	outputList := make([]SourceImageReferenceModel, 0)
	// since the image id is pulled out as a separate field, if that's set we should return an empty block here
	if input == nil || hasImageId {
		return outputList
	}
	output := SourceImageReferenceModel{
		Publisher: pointer.From(input.Publisher),
		Offer:     pointer.From(input.Offer),
		Sku:       pointer.From(input.Sku),
		Version:   pointer.From(input.Version),
	}

	return append(outputList, output)
}

func flattenOSDiskModel(input *fleets.VirtualMachineScaleSetOSDisk) []OSDiskModel {
	outputList := make([]OSDiskModel, 0)
	if input == nil {
		return outputList
	}

	output := OSDiskModel{
		DeleteOption:            string(pointer.From(input.DeleteOption)),
		WriteAcceleratorEnabled: pointer.From(input.WriteAcceleratorEnabled),
	}

	if v := input.DiffDiskSettings; v != nil {
		output.DiffDiskOption = string(pointer.From(v.Option))
		output.DiffDiskPlacement = string(pointer.From(v.Placement))
	}

	caching := ""
	if v := input.Caching; v != nil && *v != fleets.CachingTypesNone {
		caching = string(*v)
	}
	output.Caching = caching

	if input.DiskSizeGB != nil {
		output.DiskSizeInGiB = pointer.From(input.DiskSizeGB)
	}

	if md := input.ManagedDisk; md != nil {
		if v := md.DiskEncryptionSet; v != nil {
			output.DiskEncryptionSetId = pointer.From(v.Id)
		}
		if sp := md.SecurityProfile; sp != nil {
			output.SecurityEncryptionType = string(pointer.From(sp.SecurityEncryptionType))
		}
		output.StorageAccountType = string(pointer.From(md.StorageAccountType))
	}

	return append(outputList, output)
}

func flattenExtensionModel(input *fleets.VirtualMachineScaleSetExtensionProfile, metadata sdk.ResourceMetaData) ([]ExtensionModel, error) {
	outputList := make([]ExtensionModel, 0)
	if input == nil || input.Extensions == nil {
		return outputList, nil
	}

	for i, input := range *input.Extensions {
		output := ExtensionModel{
			Name: pointer.From(input.Name),
		}

		if props := input.Properties; props != nil {
			output.Publisher = pointer.From(props.Publisher)
			output.Type = pointer.From(props.Type)
			output.TypeHandlerVersion = pointer.From(props.TypeHandlerVersion)
			output.AutoUpgradeMinorVersionEnabled = pointer.From(props.AutoUpgradeMinorVersion)
			output.FailureSuppressionEnabled = pointer.From(props.SuppressFailures)
			output.AutomaticUpgradeEnabled = pointer.From(props.EnableAutomaticUpgrade)
			output.ForceExtensionExecutionOnChange = pointer.From(props.ForceUpdateTag)
			// Sensitive data isn't returned, so we get it from config
			output.ProtectedSettingsJson = metadata.ResourceData.Get("virtual_machine_profile.0.extension." + strconv.Itoa(i) + ".protected_settings_json").(string)
			output.ProtectedSettingsFromKeyVault = flattenProtectedSettingsFromKeyVaultModel(props.ProtectedSettingsFromKeyVault)
			output.ExtensionsToProvisionAfterVmCreation = pointer.From(props.ProvisionAfterExtensions)
			extSettings := ""
			if props.Settings != nil {
				settings, err := json.Marshal(props.Settings)
				if err != nil {
					return nil, fmt.Errorf("unmarshaling `settings`: %+v", err)
				}

				extSettings = string(settings)
			}
			output.SettingsJson = extSettings
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func flattenProtectedSettingsFromKeyVaultModel(input *fleets.KeyVaultSecretReference) []ProtectedSettingsFromKeyVaultModel {
	outputList := make([]ProtectedSettingsFromKeyVaultModel, 0)
	if input == nil {
		return outputList
	}

	output := ProtectedSettingsFromKeyVaultModel{
		SecretUrl:     input.SecretURL,
		SourceVaultId: pointer.From(input.SourceVault.Id),
	}

	return append(outputList, output)
}

func flattenIPConfigurationModel(inputList []fleets.VirtualMachineScaleSetIPConfiguration) []IPConfigurationModel {
	outputList := make([]IPConfigurationModel, 0)
	if len(inputList) == 0 {
		return outputList
	}
	for _, input := range inputList {
		output := IPConfigurationModel{
			Name: input.Name,
		}
		if props := input.Properties; props != nil {
			output.PrimaryIpConfigurationEnabled = pointer.From(props.Primary)
			output.Version = string(pointer.From(props.PrivateIPAddressVersion))

			addressPools := make([]string, 0)
			if v := props.ApplicationGatewayBackendAddressPools; v != nil {
				addressPools = flattenSubResourceId(*v)
			}
			output.ApplicationGatewayBackendAddressPoolIds = addressPools

			lbAddressPools := make([]string, 0)
			if v := props.LoadBalancerBackendAddressPools; v != nil {
				lbAddressPools = flattenSubResourceId(*v)
			}
			output.LoadBalancerBackendAddressPoolIds = lbAddressPools

			groupIds := make([]string, 0)
			if v := props.ApplicationSecurityGroups; v != nil {
				groupIds = flattenSubResourceId(*v)
			}
			output.ApplicationSecurityGroupIds = groupIds

			if v := props.PublicIPAddressConfiguration; v != nil {
				output.PublicIPAddress = flattenPublicIPAddressModel(v)
			}

			if v := props.Subnet; v != nil {
				output.SubnetId = pointer.From(v.Id)
			}
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenPublicIPAddressModel(input *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration) []PublicIPAddressModel {
	outputList := make([]PublicIPAddressModel, 0)
	if input == nil {
		return outputList
	}
	output := PublicIPAddressModel{
		Name: input.Name,
	}

	if v := input.Sku; v != nil {
		if v.Name != nil && v.Tier != nil {
			output.SkuName = fmt.Sprintf("%s_%s", pointer.From(v.Name), pointer.From(v.Tier))
		}
	}

	if props := input.Properties; props != nil {
		output.DeleteOption = string(pointer.From(props.DeleteOption))
		if v := props.DnsSettings; v != nil {
			output.DomainNameLabel = v.DomainNameLabel
			output.DomainNameLabelScope = string(pointer.From(v.DomainNameLabelScope))
		}
		output.IdleTimeoutInMinutes = pointer.From(props.IdleTimeoutInMinutes)
		output.Version = string(pointer.From(props.PublicIPAddressVersion))
	}
	return append(outputList, output)
}

func flattenSubResourceId(inputList []fleets.SubResource) []string {
	outputList := make([]string, 0)
	if len(inputList) == 0 {
		return outputList
	}
	for _, input := range inputList {
		outputList = append(outputList, pointer.From(input.Id))
	}
	return outputList
}
