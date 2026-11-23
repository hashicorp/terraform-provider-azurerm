// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package batch_pool

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func batchPoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Batch pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PoolName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Description:  "The name of the Batch Account in which the pool will be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},
		"display_name": {
			Description: "The display name of the Batch pool.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"vm_size": {
			Description:      "The size of the VM created in the Batch pool.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
		},
		"max_tasks_per_node": {
			Description:  "The maximum number of tasks that can run concurrently on a single compute node in the pool.",
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
						Description:  "It determines what to do with a node and its running task(s) if the pool size is decreasing. Values are `Requeue`, `RetainedData`, `TaskCompletion` and `Terminate`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForComputeNodeDeallocationOption(), false),
					},
					"target_dedicated_nodes": {
						Description:  "The number of nodes in the Batch pool.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(0, 2000),
					},
					"target_low_priority_nodes": {
						Description:  "The number of low priority nodes in the Batch pool.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"resize_timeout": {
						Description: "The timeout for resize operations.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     "PT15M",
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
						Description: "The interval to wait before evaluating if the pool needs to be scaled.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     "PT15M",
					},
					"formula": {
						Description: "The autoscale formula that needs to be used for scaling the Batch pool.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
							return strings.TrimSpace(old) == strings.TrimSpace(new)
						},
					},
				},
			},
		},
		"container_configuration": {
			Description: "The container configuration used in the pool's VMs.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Description:  "The type of container configuration.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"container_configuration.0.type", "container_configuration.0.container_image_names", "container_configuration.0.container_registries"},
					},
					"container_image_names": {
						Description: "A list of container image names to use, as would be specified by `docker pull`.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						AtLeastOneOf: []string{"container_configuration.0.type", "container_configuration.0.container_image_names", "container_configuration.0.container_registries"},
					},
					"container_registries": {
						Description: "Additional container registries from which container images can be pulled by the pool's VMs.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						ForceNew:    true,
						ConfigMode:  pluginsdk.SchemaConfigModeAttr,
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
						Description:  "The ID of the Custom Image which the Virtual Machines should be created from. See [official documentation](https://docs.microsoft.com/azure/batch/batch-custom-images) for more details.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
						AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
					},

					"publisher": {
						Description:  "The publisher of the image used to create the Virtual Machines.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
					},

					"offer": {
						Description:  "The offer of the image used to create the Virtual Machines.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
					},

					"sku": {
						Description:      "The SKU of the image used to create the Virtual Machines.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     validation.StringIsNotEmpty,
						AtLeastOneOf:     []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
					},

					"version": {
						Description:  "The version of the image used to create the Virtual Machines.",
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
			Description: "The SKU of the node agents that will be created in the Batch pool.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"stop_pending_resize_operation": {
			Description: "Whether to stop if there is a pending resize operation on this pool.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
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
			Description: "A map of custom batch pool metadata.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
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
									Description:  "The Azure Storage Account name.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"container_name": {
									Description:  "The Azure Blob Storage Container name.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"relative_mount_path": {
									Description:  "The relative path on compute node where the file system will be mounted. All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"account_key": {
									Description:  "The Azure Storage Account key.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"sas_key": {
									Description:  "The Azure Storage SAS token.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"identity_id": {
									Description:  "The ARM resource ID of the User Assigned Identity.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: commonids.ValidateUserAssignedIdentityID,
								},
								"blobfuse_options": {
									Description:  "Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.",
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
									Description:  "The Azure Storage Account name.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"azure_file_url": {
									Description:  "The Azure Files URL. This is of the form 'https://{account}.file.core.windows.net/'.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsURLWithHTTPS,
								},
								"account_key": {
									Description:  "The Azure Storage Account key.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"relative_mount_path": {
									Description:  "The relative path on compute node where the file system will be mounted. All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"mount_options": {
									Description:  "Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.",
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
									Description:  "The user to use for authentication against the CIFS file system.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"source": {
									Description:  "The URI of the file system to mount.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"relative_mount_path": {
									Description:  "The relative path on compute node where the file system will be mounted. All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"mount_options": {
									Description:  "Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"password": {
									Description:  "The password to use for authentication against the CIFS file system.",
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
									Description:  "The URI of the file system to mount.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"relative_mount_path": {
									Description:  "The relative path on compute node where the file system will be mounted. All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"mount_options": {
									Description:  "Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.",
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
						Description:      "The scope of dynamic VNet assignment.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ForceNew:         true,
						Default:          string(pool.DynamicVNetAssignmentScopeNone),
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     validation.StringInSlice(pool.PossibleValuesForDynamicVNetAssignmentScope(), false),
					},
					"accelerated_networking_enabled": {
						Description: "Whether to enable accelerated networking.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
						ForceNew:    true,
					},
					"subnet_id": {
						Description:  "The ARM resource identifier of the Virtual Network subnet which the compute nodes of the pool will join.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"public_ips": {
						Description: "A list of public IP IDs that will be allocated to nodes.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},
					"public_address_provisioning_type": {
						Description:  "Type of public IP address provisioning.",
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
									Description:  "The name of the endpoint. The name must be unique within a Batch pool, can contain letters, numbers, underscores, periods, and hyphens. Names must start with a letter or number, must end with a letter, number, or underscore, and cannot exceed 77 characters.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"protocol": {
									Description:  "The protocol of the endpoint. Acceptable values are `TCP` and `UDP`.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(pool.PossibleValuesForInboundEndpointProtocol(), false),
								},
								"backend_port": {
									Description: "The port number on the compute node. Acceptable values are between `1` and `65535` except for `29876`, `29877` as these are reserved.",
									Type:        pluginsdk.TypeInt,
									Required:    true,
									ForceNew:    true,
									ValidateFunc: validation.All(
										validation.IsPortNumber,
										validation.IntNotInSlice([]int{29876, 29877}),
									),
								},
								"frontend_port_range": {
									Description:  "The range of external ports that will be used to provide inbound access to the backendPort on individual compute nodes in the format of `1000-1100`. Acceptable values range between `1` and `65534` except ports from `50000` to `55000` which are reserved by the Batch service. All ranges within a pool must be distinct and cannot overlap.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.FrontendPortRange,
								},
								"network_security_group_rules": {
									Description: "The maximum number of rules that can be specified across all the endpoints on a Batch pool is `25`. If no Network Security Group rules are specified, a default rule will be created to allow inbound access to the specified backendPort.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									ForceNew:    true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"priority": {
												Description:  "The priority for this rule.",
												Type:         pluginsdk.TypeInt,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.IntAtLeast(150),
											},
											"access": {
												Description:  "The action that should be taken for a specified IP address, subnet range or tag. Acceptable values are `Allow` and `Deny`.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringInSlice(pool.PossibleValuesForNetworkSecurityGroupRuleAccess(), false),
											},
											"source_address_prefix": {
												Description:  "The source address prefix or tag to match for the rule.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"source_port_ranges": {
												Description: "The source port ranges to match for the rule. `100-200`). The ports should in the range of 0 to 65535 and the port ranges or ports can't overlap. If any other values are provided the request fails with HTTP status code 400.",
												Type:        pluginsdk.TypeList,
												Optional:    true,
												Computed:    true, // azignore:AZS007 - pre-existing violation
												ForceNew:    true,
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
						Description:  "The LUN is used to uniquely identify each data disk. If attaching multiple disks, each should have a distinct LUN.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 63),
					},
					"caching": {
						Description:  "Values are: \"none\" - The caching mode for the disk is not enabled. \"readOnly\" - The caching mode for the disk is read only. \"readWrite\" - The caching mode for the disk is read and write. For information about the caching options see: <https://blogs.msdn.microsoft.com/windowsazurestorage/2012/06/27/exploring-windows-azure-drives-disks-and-images/>.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(pool.CachingTypeReadOnly),
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForCachingType(), false),
					},
					"disk_size_gb": {
						Description:  "The initial disk size in GB when creating new data disk.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"storage_account_type": {
						Description: "The Storage Account type to be used for the data disk.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     pool.StorageAccountTypeStandardLRS,
						ValidateFunc: validation.StringInSlice([]string{
							string(pool.StorageAccountTypeStandardLRS),
							string(pool.StorageAccountTypePremiumLRS),
						}, false),
					},
				},
			},
		},
		"disk_encryption": {
			Description: "Disk encryption configuration is not supported on Linux pool created with Virtual Machine Image or Shared Image Gallery Image.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disk_encryption_target": {
						Description:  "On Linux pool, only \"TemporaryDisk\" is supported; on Windows pool, \"OsDisk\" and \"TemporaryDisk\" must be specified.",
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
						Description:  "The name of the Virtual Machine Extension.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"publisher": {
						Description:  "The name of the extension handler publisher.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Description:  "The type of the extensions.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type_handler_version": {
						Description:  "The version of script handler.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"auto_upgrade_minor_version": {
						Description: "Whether the extension should use a newer minor version if one is available at deployment time. Once deployed, however, the extension will not upgrade minor versions unless redeployed, even with this property set to true.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"automatic_upgrade_enabled": {
						Description: "Whether the extension should be automatically upgraded by the platform if there is a newer version available.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"settings_json": {
						Description:  "JSON formatted public settings for the extension, the value should be encoded with [`jsonencode`](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsJSON,
					},
					"protected_settings": {
						Description: "JSON formatted protected settings for the extension, the value should be encoded with [`jsonencode`](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function. The extension can contain either `protected_settings` or `provision_after_extensions` or no protected settings at all.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Sensitive:   true,
					},
					"provision_after_extensions": {
						Description: "The collection of extension names. Collection of extension names after which this extension needs to be provisioned.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
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
						Description:  "The placement policy for allocating nodes in the pool. Values are: \"Regional\": All nodes in the pool will be allocated in the same region; \"Zonal\": Nodes in the pool will be spread across different zones with the best effort balancing.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(pool.NodePlacementPolicyTypeRegional),
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForNodePlacementPolicyType(), false),
					},
				},
			},
		},
		"license_type": {
			Description:  "The type of on-premises license to be used when deploying the operating system. This only applies to images that contain the Windows operating system, and should only be used when you hold valid on-premises licenses for the nodes which will be deployed. If omitted, no on-premises licensing discount is applied. Values are: \"Windows_Server\" - The on-premises license is for Windows Server. \"Windows_Client\" - The on-premises license is for Windows Client.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"os_disk_placement": {
			Description: "The ephemeral disk placement for operating system disk for all VMs in the pool. This property can be used by user in the request to choose which location the operating system should be in. e.g., cache disk space for Ephemeral OS disk provisioning. For more information on Ephemeral OS disk size requirements, please refer to Ephemeral OS disk size requirements for Windows VMs at <https://docs.microsoft.com/azure/virtual-machines/windows/ephemeral-os-disks#size-requirements> and Linux VMs at <https://docs.microsoft.com/azure/virtual-machines/linux/ephemeral-os-disks#size-requirements>.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice(
				pool.PossibleValuesForDiffDiskPlacement(), false,
			),
		},
		"inter_node_communication": {
			Description:  "Whether the pool permits direct communication between nodes. This imposes restrictions on which nodes can be assigned to the pool. Enabling this value can reduce the chance of the requested number of nodes to be allocated in the pool. Values allowed are `Disabled` and `Enabled`.",
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
						Description: "Whether to enable host encryption for the Virtual Machine or Virtual Machine Scale Set. This will enable the encryption for all the disks including Resource/Temp disk at host itself.",
						Type:        pluginsdk.TypeBool,
						ForceNew:    true,
						Optional:    true,
					},
					"security_type": {
						Description:  "The security type of the Virtual Machine.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForSecurityTypes(), false),
					},
					"secure_boot_enabled": {
						Description:  "Whether to enable secure boot for the Virtual Machine or Virtual Machine Scale Set.",
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						ForceNew:     true,
						RequiredWith: []string{"security_profile.0.security_type"},
					},
					"vtpm_enabled": {
						Description:  "Whether to enable virtual trusted platform module (vTPM) for the Virtual Machine or Virtual Machine Scale Set.",
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						ForceNew:     true,
						RequiredWith: []string{"security_profile.0.security_type"},
					},
				},
			},
		},

		"target_node_communication_mode": {
			Description:  "The desired node communication mode for the pool.",
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
						Description:  "\"Pack\" means as many tasks as possible (taskSlotsPerNode) should be assigned to each node in the pool before any tasks are assigned to the next node in the pool. \"Spread\" means that tasks should be assigned evenly across all nodes in the pool.",
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
						Description:  "The name of the user account.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password": {
						Description:  "The password for the user account.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"elevation_level": {
						Description:  "The elevation level of the user account. \"NonAdmin\" - The auto user is a standard user without elevated access. \"Admin\" - The auto user is a user with elevated access and operates with full Administrator permissions.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForElevationLevel(), false),
					},
					"linux_user_configuration": {
						Description: "This property is ignored if specified on a Windows pool. If not specified, the user is created with the default options.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"uid": {
									Description: "The group ID for the user account. The `uid` and `gid` properties must be specified together or not at all. If not specified the underlying operating system picks the gid.",
									Type:        pluginsdk.TypeInt,
									Optional:    true,
								},
								"gid": {
									Description: "The user ID of the user account. The `uid` and `gid` properties must be specified together or not at all. If not specified the underlying operating system picks the uid.",
									Type:        pluginsdk.TypeInt,
									Optional:    true,
								},
								"ssh_private_key": {
									Description: "The SSH private key for the user account. The private key must not be password protected. The private key is used to automatically configure asymmetric-key based authentication for SSH between nodes in a Linux pool when the pool's enableInterNodeCommunication property is true (it is ignored if enableInterNodeCommunication is false). It does this by placing the key pair into the user's .ssh directory. If not specified, password-less SSH is not configured between nodes (no modification of the user's .ssh directory is done).",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									Sensitive:   true,
								},
							},
						},
					},
					"windows_user_configuration": {
						Description: "This property can only be specified if the user is on a Windows pool. If not specified and on a Windows pool, the user is created with the default options.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"login_mode": {
									Description:  "Specifies login mode for the user. Values supported are \"Batch\" and \"Interactive\".",
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
						Description: "Whether automatic updates are enabled on the Virtual Machine.",
						Type:        pluginsdk.TypeBool,
						Default:     true,
						Optional:    true,
					},
				},
			},
		},
	}
}

func startTaskSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"command_line": {
			Description:  "The command line executed by the start task.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"container": {
			Description: "When this is specified, all directories recursively below the `AZ_BATCH_NODE_ROOT_DIR` (the root of Azure Batch directories on the node) are mapped into the container, all task environment variables are mapped into the container, and the task command line is executed in the container.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"run_options": {
						Description:  "Additional options to the container create command. These additional options are supplied as arguments to the \"docker create\" command, in addition to those controlled by the Batch Service.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"image_name": {
						Description:  "The image to use to create the container in which the task will run. This is the full image reference, as would be specified to \"docker pull\". If no tag is provided as part of the image name, the tag \":latest\" is used as a default.",
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
						Description:  "A flag to indicate where the container task working directory is.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(pool.PossibleValuesForContainerWorkingDirectory(), false),
					},
				},
			},
		},

		"task_retry_maximum": {
			Description:  "The number of retries. If this is set to `0`, the Batch service does not retry Tasks. If this is set to `-1`, the Batch service retries Batch Tasks without limit.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(-1),
		},

		"wait_for_success": {
			Description: "A flag that indicates if the Batch pool should wait for the start task to be completed.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"common_environment_properties": {
			Description: "A map of strings (key, value) that represents the environment variables to set in the start task.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
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
						Description:  "The username to be used by the Batch pool start task.",
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
									Description:  "The elevation level of the user identity under which the start task runs.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(pool.ElevationLevelNonAdmin),
									ValidateFunc: validation.StringInSlice(pool.PossibleValuesForElevationLevel(), false),
								},
								"scope": {
									Description:  "The scope of the user identity under which the start task runs.",
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
						Description: "The Storage Container name in the auto Storage Account.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"blob_prefix": {
						Description: "The blob prefix to use when downloading blobs from an Azure Storage Container. Only the blobs whose names begin with the specified prefix will be downloaded. The property is valid only when `auto_storage_container_name` or `storage_container_url` is used. This prefix can be a partial filename or a subdirectory. If a prefix is not specified, all the files in the container will be downloaded.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"file_mode": {
						Description: "The file permission mode represented as a string in octal format (e.g. `\"0644\"`). This property applies only to files being downloaded to Linux compute nodes. It will be ignored if it is specified for a `resource_file` which will be downloaded to a Windows node. If this property is not specified for a Linux node, then a default value of 0770 is applied to the file.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"file_path": {
						Description: "The location on the compute node to which to download the file, relative to the task's working directory. If the `http_url` property is specified, the `file_path` is required and describes the path which the file will be downloaded to, including the filename. Otherwise, if the `auto_storage_container_name` or `storage_container_url` property is specified, `file_path` is optional and is the directory to download the files to. In the case where `file_path` is used as a directory, any directory structure already associated with the input data will be retained in full and appended to the specified filePath directory. The specified relative path cannot break out of the task's working directory (for example by using '.').",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"http_url": {
						Description: "The URL of the file to download. If the URL is Azure Blob Storage, it must be readable using anonymous access; that is, the Batch service does not present any credentials when downloading the blob. There are two ways to get such a URL for a blob in Azure storage: include a Shared Access Signature (SAS) granting read permissions on the blob, or set the ACL for the blob or its container to allow public access.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"storage_container_url": {
						Description: "The URL of the blob container within Azure Blob Storage. This URL must be readable and listable using anonymous access; that is, the Batch service does not present any credentials when downloading the blob. There are two ways to get such a URL for a blob in Azure storage: include a Shared Access Signature (SAS) granting read and list permissions on the blob, or set the ACL for the blob or its container to allow public access.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"user_assigned_identity_id": {
						Description: "An identity reference from pool's User Assigned Managed Identity list.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
				},
			},
		},
	}
}
