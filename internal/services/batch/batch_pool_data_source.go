package batch

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceBatchPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBatchPoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.PoolName,
			},
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},
			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"vm_size": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"max_tasks_per_node": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
			"fixed_scale": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"node_deallocation_option": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"target_dedicated_nodes": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"target_low_priority_nodes": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"resize_timeout": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
			"auto_scale": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"evaluation_interval": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"formula": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
			"application_licenses": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
			"application_packages": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
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
			"deployment_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cloud_service_configuration": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"os_family": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"os_version": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
						"virtual_machine_configuration": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"container_configuration": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"type": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"container_image_names": {
													Type:     pluginsdk.TypeSet,
													Computed: true,
												},
												"container_registries": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &pluginsdk.Resource{
														Schema: batchPoolDataContainerRegistry(),
													},
												},
											},
										},
									},
									"data_disks": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"lun": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"caching": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"disk_size_gb": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"storage_account_type": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
											},
										},
									},
									"disk_encryption_configuration": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"disk_encryption_target": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
											},
										},
									},
									"extensions": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"publisher": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"type": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"type_handler_version": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"auto_upgrade_minor_version": {
													Type:     pluginsdk.TypeBool,
													Computed: true,
												},
												"settings": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"protected_settings": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"provision_after_extensions": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type:     pluginsdk.TypeString,
														Computed: true,
													},
												},
											},
										},
									},
									"image_reference": {
										Type:     pluginsdk.TypeList,
										Computed: true,
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
									"license_type": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"node_agent_sku_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"node_placement_configuration": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"policy": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
											},
										},
									},
									"os_disk": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"ephemeral_os_disk_settings": {
													Type:     pluginsdk.TypeList,
													Computed: true,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"placement": {
																Type:     pluginsdk.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"windows_configuration": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"enable_automatic_updates": {
													Type:     pluginsdk.TypeBool,
													Computed: true,
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
			"inter_node_communication": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"mount_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"azure_blob_file_system_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"account_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"container_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"relative_mount_path": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"account_key": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"sas_key": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"identity_reference": {
										Type:     pluginsdk.TypeList,
										Computed: true,
									},
									"blobfuse_options": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
						"azure_file_share_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"account_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"azure_file_url": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"account_key": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"relative_mount_path": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"mount_options": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cifs_mount_configuration": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"user_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"source": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"relative_mount_path": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"mount_options": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"password": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
						"nfs_mount_configuration": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"source": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"relative_mount_path": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"mount_options": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"metadata": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"network_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subnet_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						//TODO test whether this is forcenew
						"dynamic_vnet_assignment_scope": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"public_ips": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
						},
						"public_address_provisioning_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"endpoint_configuration": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"backend_port": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
									"frontend_port_range": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"network_security_group_rules": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"priority": {
													Type:     pluginsdk.TypeInt,
													Computed: true,
												},
												"access": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"source_address_prefix": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"source_port_ranges": {
													Type:     pluginsdk.TypeList,
													Computed: true,
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
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"storage_image_reference": {
				Type:     pluginsdk.TypeList,
				Computed: true,
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
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"node_fill_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
			//"container_configuration": {
			//	Type:     pluginsdk.TypeList,
			//	Computed: true,
			//	Elem: &pluginsdk.Resource{
			//		Schema: map[string]*pluginsdk.Schema{
			//			"type": {
			//				Type:     pluginsdk.TypeString,
			//				Computed: true,
			//			},
			//
			//			"container_image_names": {
			//				Type:     pluginsdk.TypeSet,
			//				Computed: true,
			//				Elem: &pluginsdk.Schema{
			//					Type: pluginsdk.TypeString,
			//				},
			//			},
			//
			//			"container_registries": {
			//				Type:     pluginsdk.TypeList,
			//				Computed: true,
			//				Elem: &pluginsdk.Resource{
			//					Schema: map[string]*pluginsdk.Schema{
			//						"registry_server": {
			//							Type:     pluginsdk.TypeString,
			//							Computed: true,
			//						},
			//						"user_name": {
			//							Type:     pluginsdk.TypeString,
			//							Computed: true,
			//						},
			//						"password": {
			//							Type:      pluginsdk.TypeString,
			//							Computed:  true,
			//							Sensitive: true,
			//						},
			//					},
			//				},
			//			},
			//		},
			//	},
			//},
			"certificate": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"store_location": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"store_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"visibility": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},
			"start_task": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: batchPoolDataStartTaskDSSchema(),
				},
			},
			"user_accounts": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"password": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"elevation_level": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"linux_user_configuration": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"uid": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
									"gid": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
									"ssh_private_key": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
						"windows_user_configuration": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"login_mode": {
										Type:     pluginsdk.TypeString,
										Computed: true,
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

func batchPoolDataIdentityReference() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"identity_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func batchPoolDataContainerRegistry() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"password": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"registry_server": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"identity_reference": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     batchPoolDataIdentityReference(),
		},
	}
}

func batchPoolDataStartTaskDSSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"command_line": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"container_settings": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"container_run_options": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"image_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"registry": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Resource{
							Schema: batchPoolDataContainerRegistry(),
						},
					},
					"working_directory": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"task_retry_maximum": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"wait_for_success": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"common_environment_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			// Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"user_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"user_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"auto_user": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"elevation_level": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"scope": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"resource_file": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"auto_storage_container_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"blob_prefix": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"file_mode": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"file_path": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"http_url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"storage_container_url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
	return s
}

func dataSourceBatchPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("account_name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.PoolProperties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("vm_size", props.VMSize)
		d.Set("max_tasks_per_node", props.TaskSlotsPerNode)
		d.Set("inter_node_communication", string(props.InterNodeCommunication))

		if props.ApplicationLicenses != nil {
			applicationLicenses := make([]interface{}, len(*props.ApplicationLicenses))
			for _, license := range *props.ApplicationLicenses {
				applicationLicenses = append(applicationLicenses, license)
			}
			d.Set("application_licenses", applicationLicenses)
		}

		if props.ApplicationPackages != nil {
			applicationPackages := make([]interface{}, len(*props.ApplicationPackages))
			for _, pkg := range *props.ApplicationPackages {
				appPkg := make(map[string]interface{}, 1)
				appPkg["id"] = *pkg.ID
				if pkg.Version != nil {
					appPkg["version"] = *pkg.Version
				}
				applicationPackages = append(applicationPackages, appPkg)
			}
			d.Set("application_packages", applicationPackages)
		}

		if props.TaskSchedulingPolicy != nil {
			taskSchedulingPolicy := make([]interface{}, 1)
			nodeFillType := make(map[string]interface{}, 1)
			nodeFillType["node_fill_type"] = string(props.TaskSchedulingPolicy.NodeFillType)
			taskSchedulingPolicy = append(taskSchedulingPolicy, nodeFillType)
			d.Set("task_scheduling_policy", taskSchedulingPolicy)
		}

		if props.UserAccounts != nil {
			userAccounts := make([]interface{}, len(*props.UserAccounts))
			for _, userAccount := range *props.UserAccounts {
				userAccounts = append(userAccounts, flattenBatchPoolUserAccount(&userAccount))
			}
			d.Set("user_accounts", userAccounts)
		}

		if props.MountConfiguration != nil {
			mountConfigs := make([]interface{}, len(*props.MountConfiguration))
			for _, mountConfig := range *props.MountConfiguration {
				mountConfigs = append(mountConfigs, flattenBatchPoolMountConfig(&mountConfig))
			}
			d.Set("mount_configuration", mountConfigs)
		}

		if props.DeploymentConfiguration != nil {
			deploymentConfigMap := make(map[string]interface{}, 0)
			if props.DeploymentConfiguration.CloudServiceConfiguration != nil {
				deploymentConfigMap["cloud_service_configuration"] = flattenBatchPoolCloudServiceConfiguration(props.DeploymentConfiguration.CloudServiceConfiguration)

			}
			if props.DeploymentConfiguration.VirtualMachineConfiguration != nil {
				deploymentConfigMap["cloud_service_configuration"] = flattenBatchPoolVirtualMachineConfiguration(props.DeploymentConfiguration.VirtualMachineConfiguration)
			}

			d.Set("deployment_configuration", deploymentConfigMap)
		}

		if scaleSettings := props.ScaleSettings; scaleSettings != nil {
			if err := d.Set("auto_scale", flattenBatchPoolAutoScaleSettings(scaleSettings.AutoScale)); err != nil {
				return fmt.Errorf("flattening `auto_scale`: %+v", err)
			}
			if err := d.Set("fixed_scale", flattenBatchPoolFixedScaleSettings(scaleSettings.FixedScale)); err != nil {
				return fmt.Errorf("flattening `fixed_scale `: %+v", err)
			}
		}

		//if dcfg := props.DeploymentConfiguration; dcfg != nil {
		//	if vmcfg := dcfg.VirtualMachineConfiguration; vmcfg != nil {
		//		if err := d.Set("container_configuration", flattenBatchPoolContainerConfiguration(vmcfg.ContainerConfiguration)); err != nil {
		//			return fmt.Errorf("setting `container_configuration`: %v", err)
		//		}
		//
		//		if err := d.Set("storage_image_reference", flattenBatchPoolImageReference(vmcfg.ImageReference)); err != nil {
		//			return fmt.Errorf("setting `storage_image_reference`: %v", err)
		//		}
		//
		//		if err := d.Set("node_agent_sku_id", vmcfg.NodeAgentSkuID); err != nil {
		//			return fmt.Errorf("setting `node_agent_sku_id`: %v", err)
		//		}
		//	}
		//}

		if err := d.Set("certificate", flattenBatchPoolCertificateReferences(props.Certificates)); err != nil {
			return fmt.Errorf("setting `certificate`: %v", err)
		}

		d.Set("start_task", flattenBatchPoolStartTask(props.StartTask))
		d.Set("metadata", FlattenBatchMetaData(props.Metadata))

		if err := d.Set("network_configuration", flattenBatchPoolNetworkConfiguration(props.NetworkConfiguration)); err != nil {
			return fmt.Errorf("setting `network_configuration`: %v", err)
		}
	}

	return nil
}
