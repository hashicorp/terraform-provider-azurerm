package batch

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
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
			"node_agent_sku_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
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
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"container_registries": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"registry_server": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"user_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"password": {
										Type:      pluginsdk.TypeString,
										Computed:  true,
										Sensitive: true,
									},
								},
							},
						},
					},
				},
			},
			"certificate": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"store_location": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CurrentUser",
								"LocalMachine",
							}, false),
						},
						"store_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"visibility": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},
			"start_task": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"command_line": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"max_task_retry_count": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  1,
						},

						"wait_for_success": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"environment": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
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
	}
}

func dataSourceBatchPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountName := d.Get("account_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Batch pool %q in account %q (Resource Group %q) was not found", name, accountName, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Batch pool %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("account_name", accountName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.PoolProperties; props != nil {
		d.Set("vm_size", props.VMSize)
		d.Set("max_tasks_per_node", props.MaxTasksPerNode)

		if scaleSettings := props.ScaleSettings; scaleSettings != nil {
			if err := d.Set("auto_scale", flattenBatchPoolAutoScaleSettings(scaleSettings.AutoScale)); err != nil {
				return fmt.Errorf("Error flattening `auto_scale`: %+v", err)
			}
			if err := d.Set("fixed_scale", flattenBatchPoolFixedScaleSettings(scaleSettings.FixedScale)); err != nil {
				return fmt.Errorf("Error flattening `fixed_scale `: %+v", err)
			}
		}

		if dcfg := props.DeploymentConfiguration; dcfg != nil {
			if vmcfg := dcfg.VirtualMachineConfiguration; vmcfg != nil {
				if err := d.Set("container_configuration", flattenBatchPoolContainerConfiguration(d, vmcfg.ContainerConfiguration)); err != nil {
					return fmt.Errorf("error setting `container_configuration`: %v", err)
				}

				if err := d.Set("storage_image_reference", flattenBatchPoolImageReference(vmcfg.ImageReference)); err != nil {
					return fmt.Errorf("error setting `storage_image_reference`: %v", err)
				}

				if err := d.Set("node_agent_sku_id", vmcfg.NodeAgentSkuID); err != nil {
					return fmt.Errorf("error setting `node_agent_sku_id`: %v", err)
				}
			}
		}

		if err := d.Set("certificate", flattenBatchPoolCertificateReferences(props.Certificates)); err != nil {
			return fmt.Errorf("error setting `certificate`: %v", err)
		}

		d.Set("start_task", flattenBatchPoolStartTask(props.StartTask))
		d.Set("metadata", FlattenBatchMetaData(props.Metadata))

		if err := d.Set("network_configuration", flattenBatchPoolNetworkConfiguration(props.NetworkConfiguration)); err != nil {
			return fmt.Errorf("error setting `network_configuration`: %v", err)
		}
	}

	return nil
}
