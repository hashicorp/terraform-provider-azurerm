// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_pool

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceBatchPoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PoolName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"vm_size": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
		},
		"max_tasks_per_node": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"fixed_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// Property `node_deallocation_method` is set to be a writeOnly property by service team
					// It can only perform on PUT operation and is not able to perform GET operation
					// Here we treat `node_deallocation_method` the same as a secret value.
					// Issue link: https://github.com/Azure/azure-rest-api-specs/issues/20948
					"node_deallocation_method": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForComputeNodeDeallocationOption(), false),
					},
					"target_dedicated_nodes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(0, 2000),
					},
					"target_low_priority_nodes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"resize_timeout": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "PT15M",
					},
				},
			},
		},
		"auto_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"evaluation_interval": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "PT15M",
					},
					"formula": {
						Type:     pluginsdk.TypeString,
						Required: true,
						DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
							return strings.TrimSpace(old) == strings.TrimSpace(new)
						},
					},
				},
			},
		},
		"container_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"container_configuration.0.type", "container_configuration.0.container_image_names", "container_configuration.0.container_registries"},
					},
					"container_image_names": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						AtLeastOneOf: []string{"container_configuration.0.type", "container_configuration.0.container_image_names", "container_configuration.0.container_registries"},
					},
					"container_registries": {
						Type:       pluginsdk.TypeList,
						Optional:   true,
						ForceNew:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: containerRegistry(),
						},
						AtLeastOneOf: []string{"container_configuration.0.type", "container_configuration.0.container_image_names", "container_configuration.0.container_registries"},
					},
				},
			},
		},
		"storage_image_reference": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
						AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
					},

					"publisher": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
					},

					"offer": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
					},

					"sku": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     validation.StringIsNotEmpty,
						AtLeastOneOf:     []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
					},

					"version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
					},
				},
			},
		},
		"node_agent_sku_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"stop_pending_resize_operation": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"start_task": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: startTaskSchema(),
			},
		},
		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"mount": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"azure_blob_file_system": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"account_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"container_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"relative_mount_path": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"account_key": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"sas_key": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"identity_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: commonids.ValidateUserAssignedIdentityID,
								},
								"blobfuse_options": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"azure_file_share": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"account_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"azure_file_url": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsURLWithHTTPS,
								},
								"account_key": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"relative_mount_path": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"mount_options": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"cifs_mount": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"user_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"source": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"relative_mount_path": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"mount_options": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"password": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"nfs_mount": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"source": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"relative_mount_path": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"mount_options": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},
		"network_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dynamic_vnet_assignment_scope": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ForceNew:         true,
						Default:          string(pool.DynamicVNetAssignmentScopeNone),
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     validation.StringInSlice(pool.PossibleValuesForDynamicVNetAssignmentScope(), false),
					},
					"accelerated_networking_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
						ForceNew: true,
					},
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"public_ips": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},
					"public_address_provisioning_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForIPAddressProvisioningType(), false),
					},
					"endpoint_configuration": {
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
								"protocol": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(pool.PossibleValuesForInboundEndpointProtocol(), false),
								},
								"backend_port": {
									Type:     pluginsdk.TypeInt,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.All(
										validation.IsPortNumber,
										validation.IntNotInSlice([]int{29876, 29877}),
									),
								},
								"frontend_port_range": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.FrontendPortRange,
								},
								"network_security_group_rules": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"priority": {
												Type:         pluginsdk.TypeInt,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.IntAtLeast(150),
											},
											"access": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringInSlice(pool.PossibleValuesForNetworkSecurityGroupRuleAccess(), false),
											},
											"source_address_prefix": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"source_port_ranges": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Computed: true, // azignore:AZS007 - pre-existing violation
												ForceNew: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													Default:      "*",
													ValidateFunc: validation.StringIsNotEmpty,
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
		},
		"data_disks": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"lun": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 63),
					},
					"caching": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(pool.CachingTypeReadOnly),
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForCachingType(), false),
					},
					"disk_size_gb": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"storage_account_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  pool.StorageAccountTypeStandardLRS,
						ValidateFunc: validation.StringInSlice([]string{
							string(pool.StorageAccountTypeStandardLRS),
							string(pool.StorageAccountTypePremiumLRS),
						}, false),
					},
				},
			},
		},
		"disk_encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disk_encryption_target": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForDiskEncryptionTarget(), false),
					},
				},
			},
		},
		"extensions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
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
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"auto_upgrade_minor_version": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"automatic_upgrade_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"settings_json": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsJSON,
					},
					"protected_settings": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},
					"provision_after_extensions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
		"node_placement": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"policy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(pool.NodePlacementPolicyTypeRegional),
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForNodePlacementPolicyType(), false),
					},
				},
			},
		},
		"license_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"os_disk_placement": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice(
				pool.PossibleValuesForDiffDiskPlacement(), false,
			),
		},
		"inter_node_communication": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(pool.InterNodeCommunicationStateEnabled),
			ValidateFunc: validation.StringInSlice(pool.PossibleValuesForInterNodeCommunicationState(), false),
		},

		"security_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host_encryption_enabled": {
						Type:     pluginsdk.TypeBool,
						ForceNew: true,
						Optional: true,
					},
					"security_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForSecurityTypes(), false),
					},
					"secure_boot_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						ForceNew:     true,
						RequiredWith: []string{"security_profile.0.security_type"},
					},
					"vtpm_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						ForceNew:     true,
						RequiredWith: []string{"security_profile.0.security_type"},
					},
				},
			},
		},

		"target_node_communication_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(pool.PossibleValuesForNodeCommunicationMode(), false),
		},

		"task_scheduling_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"node_fill_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForComputeNodeFillType(), false),
					},
				},
			},
		},
		"user_accounts": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"elevation_level": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForElevationLevel(), false),
					},
					"linux_user_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"uid": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"gid": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"ssh_private_key": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},
							},
						},
					},
					"windows_user_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"login_mode": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(pool.PossibleValuesForLoginMode(), false),
								},
							},
						},
					},
				},
			},
		},
		"windows": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enable_automatic_updates": {
						Type:     pluginsdk.TypeBool,
						Default:  true,
						Optional: true,
					},
				},
			},
		},
	}
}

func startTaskSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"command_line": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"container": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"run_options": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"image_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"registry": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: containerRegistry(),
						},
					},
					"working_directory": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForContainerWorkingDirectory(), false),
					},
				},
			},
		},

		"task_retry_maximum": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(-1),
		},

		"wait_for_success": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"common_environment_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"user_identity": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"user_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						AtLeastOneOf: []string{"start_task.0.user_identity.0.user_name", "start_task.0.user_identity.0.auto_user"},
					},
					"auto_user": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"elevation_level": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(pool.ElevationLevelNonAdmin),
									ValidateFunc: validation.StringInSlice(pool.PossibleValuesForElevationLevel(), false),
								},
								"scope": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(pool.AutoUserScopeTask),
									ValidateFunc: validation.StringInSlice(pool.PossibleValuesForAutoUserScope(), false),
								},
							},
						},
						AtLeastOneOf: []string{"start_task.0.user_identity.0.user_name", "start_task.0.user_identity.0.auto_user"},
					},
				},
			},
		},
		// lintignore:XS003
		"resource_file": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"auto_storage_container_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"blob_prefix": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"file_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"file_path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"http_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"storage_container_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func containerRegistry() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"registry_server": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			Description:  "The User Assigned Identity to use for Container Registry access.",
		},
		"user_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
