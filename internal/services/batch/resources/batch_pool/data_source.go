// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package batch_pool

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Batch pool.",
		Read:        dataSourceBatchPoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Batch pool.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.PoolName,
			},
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
			"account_name": {
				Description:  "The Batch Account name associated with the Batch pool.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},

			"auto_scale": {
				Description: "An `auto_scale` block that describes the scale settings when using auto scale.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"evaluation_interval": {
							Description: "The interval to wait before evaluating if the pool needs to be scaled.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"formula": {
							Description: "The autoscale formula that needs to be used for scaling the Batch pool.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"container_configuration": {
				Description: "The container configuration used in the pool's VMs.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Description: "The type of container configuration.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"container_image_names": {
							Description: "A list of container image names to use, as would be specified by `docker pull`.",
							Type:        pluginsdk.TypeSet,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"container_registries": {
							Description: "Additional container registries from which container images can be pulled by the pool's VMs.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Resource{
								Schema: batchPoolDataContainerRegistry(),
							},
						},
					},
				},
			},

			"data_disks": {
				Description: "A `data_disks` block describes the data disk settings.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"lun": {
							Description: "The LUN is used to uniquely identify each data disk.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
						"caching": {
							Description: "The caching mode of data disks.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"disk_size_gb": {
							Description: "The initial disk size in GB when creating new data disk.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
						"storage_account_type": {
							Description: "The Storage Account type to be used for the data disk.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"disk_encryption": {
				Description: "A `disk_encryption` block describes the disk encryption configuration applied on compute nodes in the pool.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"disk_encryption_target": {
							Description: "On Linux pool, only `TemporaryDisk` is supported; on Windows pool, `OsDisk` and `TemporaryDisk` must be specified.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"extensions": {
				Description: "An `extensions` block describes the extension settings.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the Virtual Machine Extension.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"publisher": {
							Description: "The name of the extension handler publisher.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "The type of the extensions.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"type_handler_version": {
							Description: "The version of script handler.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"auto_upgrade_minor_version": {
							Description: "Whether the extension should use a newer minor version if one is available at deployment time. Once deployed, however, the extension will not upgrade minor versions unless redeployed, even with this property set to true.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
						"settings_json": {
							Description: "JSON formatted public settings for the extension.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"protected_settings": {
							Description: "The extension can contain either `protected_settings` or `provision_after_extensions` or no protected settings at all.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"provision_after_extensions": {
							Description: "The collection of extension names. Collection of extension names after which this extension needs to be provisioned.",
							Type:        pluginsdk.TypeSet,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"fixed_scale": {
				Description: "A `fixed_scale` block that describes the scale settings when using fixed scale.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"target_dedicated_nodes": {
							Description: "The number of nodes in the Batch pool.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
						"target_low_priority_nodes": {
							Description: "The number of low priority nodes in the Batch pool.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
						"resize_timeout": {
							Description: "The timeout for resize operations.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"inter_node_communication": {
				Description: "Whether the pool permits direct communication between nodes. This imposes restrictions on which nodes can be assigned to the pool. Enabling this value can reduce the chance of the requested number of nodes to be allocated in the pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"license_type": {
				Description: "The type of on-premises license to be used when deploying the operating system.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"max_tasks_per_node": {
				Description: "The maximum number of tasks that can run concurrently on a single compute node in the pool.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"metadata": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"mount": {
				Description: "A `mount` block that describes mount configuration.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"azure_blob_file_system": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"account_name": {
										Description: "The Azure Storage Account name.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"container_name": {
										Description: "The Azure Blob Storage Container name.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"relative_mount_path": {
										Description: "The relative path on compute node where the file system will be mounted. All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"account_key": {
										Description: "The Azure Storage Account key.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"sas_key": {
										Description: "The Azure Storage SAS token.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"identity_id": {
										Description: "The ARM resource ID of the User Assigned Identity.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"blobfuse_options": {
										Description: "Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"azure_file_share": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"account_name": {
										Description: "The Azure Storage Account name.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"azure_file_url": {
										Description: "The Azure Files URL. This is of the form 'https://{account}.file.core.windows.net/'.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"account_key": {
										Description: "The Azure Storage Account key.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"relative_mount_path": {
										Description: "The relative path on compute node where the file system will be mounted. All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"mount_options": {
										Description: "Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"cifs_mount": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"user_name": {
										Description: "The user to use for authentication against the CIFS file system.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"source": {
										Description: "The URI of the file system to mount.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"relative_mount_path": {
										Description: "The relative path on compute node where the file system will be mounted. All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"mount_options": {
										Description: "Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"password": {
										Description: "The password to use for authentication against the CIFS file system.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
										Sensitive:   true,
									},
								},
							},
						},
						"nfs_mount": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"source": {
										Description: "The URI of the file system to mount.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"relative_mount_path": {
										Description: "The relative path on compute node where the file system will be mounted. All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"mount_options": {
										Description: "Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},

			"network_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dynamic_vnet_assignment_scope": {
							Description: "The scope of dynamic VNet assignment.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"accelerated_networking_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"subnet_id": {
							Description: "The ARM resource identifier of the Virtual Network subnet which the compute nodes of the pool are joined to.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"public_ips": {
							Description: "A list of public IP IDs that will be allocated to nodes.",
							Type:        pluginsdk.TypeSet,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"public_address_provisioning_type": {
							Description: "Type of public IP address provisioning.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"endpoint_configuration": {
							Description: "The inbound NAT pools that are used to address specific ports on the individual compute node externally.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Description: "The name of the endpoint.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"protocol": {
										Description: "The protocol of the endpoint.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"backend_port": {
										Description: "The port number on the compute node.",
										Type:        pluginsdk.TypeInt,
										Computed:    true,
									},
									"frontend_port_range": {
										Description: "The range of external ports that are used to provide inbound access to the backendPort on the individual compute nodes in the format of `1000-1100`.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
									"network_security_group_rules": {
										Description: "The list of Network Security Group rules that are applied to the endpoint.",
										Type:        pluginsdk.TypeList,
										Computed:    true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"priority": {
													Description: "The priority for this rule.",
													Type:        pluginsdk.TypeInt,
													Computed:    true,
												},
												"access": {
													Description: "The action that should be taken for a specified IP address, subnet range or tag.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
												"source_address_prefix": {
													Description: "The source address prefix or tag to match for the rule.",
													Type:        pluginsdk.TypeString,
													Computed:    true,
												},
												"source_port_ranges": {
													Description: "The source port ranges to match for the rule.",
													Type:        pluginsdk.TypeList,
													Computed:    true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
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

			"node_agent_sku_id": {
				Description: "The SKU of the node agents in the Batch pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"node_placement": {
				Description: "A `node_placement` block that describes the placement policy for allocating nodes in the pool.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"policy": {
							Description: "The placement policy for allocating nodes in the pool.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"os_disk_placement": {
				Description: "The ephemeral disk placement for operating system disk for all VMs in the pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"start_task": {
				Description: "A `start_task` block that describes the start task settings for the Batch pool.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: startTaskDSSchema(),
				},
			},

			"storage_image_reference": {
				Description: "The reference of the storage image used by the nodes in the Batch pool.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"publisher": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"offer": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"sku": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"task_scheduling_policy": {
				Description: "A `task_scheduling_policy` block that describes how tasks are distributed across compute nodes in a pool.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"node_fill_type": {
							Description: "`Pack` means as many tasks as possible (taskSlotsPerNode) should be assigned to each node in the pool before any tasks are assigned to the next node in the pool. `Spread` means that tasks should be assigned evenly across all nodes in the pool.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"user_accounts": {
				Description: "A `user_accounts` block that describes the list of user accounts to be created on each node in the pool.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the user account.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"password": {
							Description: "The password for the user account.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"elevation_level": {
							Description: "The elevation level of the user account. \"NonAdmin\" - The auto user is a standard user without elevated access. \"Admin\" - The auto user is a user with elevated access and operates with full Administrator permissions.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"linux_user_configuration": {
							Description: "This property is ignored if specified on a Windows pool. If not specified, the user is created with the default options.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"uid": {
										Description: "The group ID for the user account.",
										Type:        pluginsdk.TypeInt,
										Computed:    true,
									},
									"gid": {
										Description: "The user ID of the user account.",
										Type:        pluginsdk.TypeInt,
										Computed:    true,
									},
									"ssh_private_key": {
										Description: "The SSH private key for the user account.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
										Sensitive:   true,
									},
								},
							},
						},
						"windows_user_configuration": {
							Description: "This property can only be specified if the user is on a Windows pool. If not specified and on a Windows pool, the user is created with the default options.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"login_mode": {
										Description: "The login mode for the user.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},

			"vm_size": {
				Description: "The size of the VM created in the Batch pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"windows": {
				Description: "A `windows` block that describes the Windows configuration in the pool.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enable_automatic_updates": {
							Description: "Whether automatic updates are enabled on the Virtual Machine.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func startTaskDSSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"command_line": {
			Description: "The command line executed by the start task.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"container": {
			Description: "The settings for the container under which the start task runs.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"run_options": {
						Description: "Additional options to the container create command.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"image_name": {
						Description: "The image to use to create the container in which the task will run.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"registry": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: batchPoolDataContainerRegistry(),
						},
					},
					"working_directory": {
						Description: "A flag to indicate where the container task working directory is.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"task_retry_maximum": {
			Description: "The number of retries.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"wait_for_success": {
			Description: "A flag that indicates if the Batch pool should wait for the start task to be completed.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"common_environment_properties": {
			Description: "A map of strings (key, value) that represents the environment variables to set in the start task.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			// Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"user_identity": {
			Description: "A `user_identity` block that describes the user identity under which the start task runs.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"user_name": {
						Description: "The username to be used by the Batch pool start task.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"auto_user": {
						Description: "An `auto_user` block that describes the user identity under which the start task runs.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"elevation_level": {
									Description: "The elevation level of the user identity under which the start task runs.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"scope": {
									Description: "The scope of the user identity under which the start task runs.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},

		"resource_file": {
			Description: "One or more `resource_file` blocks that describe the files to be downloaded to a compute node.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"auto_storage_container_name": {
						Description: "The Storage Container name in the auto Storage Account.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"blob_prefix": {
						Description: "The blob prefix used when downloading blobs from an Azure Storage Container.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"file_mode": {
						Description: "The file permission mode attribute represented as a string in octal format (e.g. `\"0644\"`).",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"file_path": {
						Description: "The location on the compute node to which to download the file, relative to the task's working directory. If the `http_url` property is specified, the `file_path` is required and describes the path which the file will be downloaded to, including the filename. Otherwise, if the `auto_storage_container_name` or `storage_container_url` property is specified.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"http_url": {
						Description: "The URL of the file to download. If the URL is Azure Blob Storage, it must be readable using anonymous access.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"storage_container_url": {
						Description: "The URL of the blob container within Azure Blob Storage.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"user_assigned_identity_id": {
						Description: "An identity reference from pool's User Assigned Managed Identity list.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

func batchPoolDataContainerRegistry() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"registry_server": {
			Description: "The container registry URL.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"user_assigned_identity_id": {
			Description: "The reference to the User Assigned Identity to use to access an Azure Container Registry instead of username and password.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"user_name": {
			Description: "The user name to log into the registry server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"password": {
			Description: "The password to log into the registry server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}

func dataSourceBatchPoolRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Batch.PoolClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := pool.NewPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.PoolName)
	d.Set("account_name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("display_name", props.DisplayName)
			d.Set("vm_size", props.VMSize)
			d.Set("inter_node_communication", pointer.FromEnum(props.InterNodeCommunication))
			d.Set("max_tasks_per_node", props.TaskSlotsPerNode)

			if scaleSettings := props.ScaleSettings; scaleSettings != nil {
				if err := d.Set("auto_scale", flattenBatchPoolAutoScaleSettings(scaleSettings.AutoScale)); err != nil {
					return fmt.Errorf("flattening `auto_scale`: %+v", err)
				}
				if err := d.Set("fixed_scale", flattenBatchPoolFixedScaleSettings(d, scaleSettings.FixedScale)); err != nil {
					return fmt.Errorf("flattening `fixed_scale `: %+v", err)
				}
			}

			if props.UserAccounts != nil {
				userAccounts := make([]any, 0)
				for _, userAccount := range *props.UserAccounts {
					userAccounts = append(userAccounts, flattenBatchPoolUserAccount(d, &userAccount))
				}
				d.Set("user_accounts", userAccounts)
			}

			if props.MountConfiguration != nil {
				mountConfigs := make([]any, 0)
				for _, mountConfig := range *props.MountConfiguration {
					mountConfigs = append(mountConfigs, flattenBatchPoolMountConfig(d, &mountConfig))
				}
				d.Set("mount_configuration", mountConfigs)
			}

			if props.DeploymentConfiguration != nil {
				if props.DeploymentConfiguration.VirtualMachineConfiguration != nil {
					config := props.DeploymentConfiguration.VirtualMachineConfiguration
					if config.ContainerConfiguration != nil {
						d.Set("container_configuration", flattenBatchPoolContainerConfiguration(d, config.ContainerConfiguration))
					}
					if config.DataDisks != nil {
						dataDisks := make([]any, 0)
						for _, item := range *config.DataDisks {
							dataDisk := make(map[string]any)
							dataDisk["lun"] = item.Lun
							dataDisk["disk_size_gb"] = item.DiskSizeGB

							caching := ""
							if item.Caching != nil {
								caching = string(*item.Caching)
							}
							dataDisk["caching"] = caching

							storageAccountType := ""
							if item.StorageAccountType != nil {
								storageAccountType = string(*item.StorageAccountType)
							}
							dataDisk["storage_account_type"] = storageAccountType

							dataDisks = append(dataDisks, dataDisk)
						}
						d.Set("data_disks", dataDisks)
					}
					if config.DiskEncryptionConfiguration != nil {
						diskEncryptionConfiguration := make([]any, 0)
						if config.DiskEncryptionConfiguration.Targets != nil {
							for _, item := range *config.DiskEncryptionConfiguration.Targets {
								target := make(map[string]any)
								target["disk_encryption_target"] = string(item)
								diskEncryptionConfiguration = append(diskEncryptionConfiguration, target)
							}
						}
						d.Set("disk_encryption", diskEncryptionConfiguration)
					}
					if config.Extensions != nil {
						extensions := make([]any, 0)
						n := len(*config.Extensions)
						for _, item := range *config.Extensions {
							extension := make(map[string]any)
							extension["name"] = item.Name
							extension["publisher"] = item.Publisher
							extension["type"] = item.Type
							if item.TypeHandlerVersion != nil {
								extension["type_handler_version"] = *item.TypeHandlerVersion
							}
							if item.AutoUpgradeMinorVersion != nil {
								extension["auto_upgrade_minor_version"] = *item.AutoUpgradeMinorVersion
							}
							if item.Settings != nil {
								extension["settings_json"] = item.Settings
							}

							for i := range n {
								if v, ok := d.GetOk(fmt.Sprintf("extensions.%d.name", i)); ok && v == item.Name {
									extension["protected_settings"] = d.Get(fmt.Sprintf("extensions.%d.protected_settings", i))
									break
								}
							}

							if item.ProvisionAfterExtensions != nil {
								extension["provision_after_extensions"] = *item.ProvisionAfterExtensions
							}
							extensions = append(extensions, extension)
						}
						d.Set("extensions", extensions)
					}

					d.Set("storage_image_reference", flattenBatchPoolImageReference(&config.ImageReference))

					if config.LicenseType != nil {
						d.Set("license_type", config.LicenseType)
					}

					d.Set("node_agent_sku_id", config.NodeAgentSkuId)

					if config.NodePlacementConfiguration != nil {
						nodePlacementConfiguration := make([]any, 0)
						nodePlacementConfig := make(map[string]any)
						nodePlacementConfig["policy"] = string(*config.NodePlacementConfiguration.Policy)
						nodePlacementConfiguration = append(nodePlacementConfiguration, nodePlacementConfig)
						d.Set("node_placement", nodePlacementConfiguration)
					}
					osDiskPlacement := ""
					if config.OsDisk != nil && config.OsDisk.EphemeralOSDiskSettings != nil && config.OsDisk.EphemeralOSDiskSettings.Placement != nil {
						osDiskPlacement = string(*config.OsDisk.EphemeralOSDiskSettings.Placement)
					}
					d.Set("os_disk_placement", osDiskPlacement)
					if config.WindowsConfiguration != nil {
						windowsConfig := []any{
							map[string]any{
								"enable_automatic_updates": *config.WindowsConfiguration.EnableAutomaticUpdates,
							},
						}
						d.Set("windows", windowsConfig)
					}
				}
			}

			d.Set("start_task", flattenBatchPoolStartTask(d, props.StartTask))
			d.Set("metadata", FlattenBatchMetaData(props.Metadata))

			if err := d.Set("network_configuration", flattenBatchPoolNetworkConfiguration(props.NetworkConfiguration)); err != nil {
				return fmt.Errorf("setting `network_configuration`: %v", err)
			}
		}
	}
	return nil
}
