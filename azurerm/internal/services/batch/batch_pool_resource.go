package batch

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2019-08-01/batch"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBatchPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBatchPoolCreate,
		Read:   resourceArmBatchPoolRead,
		Update: resourceArmBatchPoolUpdate,
		Delete: resourceArmBatchPoolDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.BatchPoolID(id)
			return err
		}),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAzureRMBatchPoolName,
			},

			// TODO: make this case sensitive once this API bug has been fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/5574
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAzureRMBatchAccountName,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vm_size": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"max_tasks_per_node": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"fixed_scale": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_dedicated_nodes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(0, 2000),
						},
						"target_low_priority_nodes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"resize_timeout": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "PT15M",
						},
					},
				},
			},
			"auto_scale": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"evaluation_interval": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "PT15M",
						},
						"formula": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
								return strings.TrimSpace(old) == strings.TrimSpace(new)
							},
						},
					},
				},
			},
			"container_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"container_image_names": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"container_registries": {
							Type:       schema.TypeList,
							Optional:   true,
							ForceNew:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"registry_server": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"user_name": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"password": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
			"storage_image_reference": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"publisher": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"offer": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"sku": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.StringIsNotEmpty,
						},

						"version": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"node_agent_sku_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stop_pending_resize_operation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
							// The ID returned for the certificate in the batch account and the certificate applied to the pool
							// are not consistent in their casing which causes issues when referencing IDs across resources
							// (as Terraform still sees differences to apply due to the casing)
							// Handling by ignoring casing for now. Raised as an issue: https://github.com/Azure/azure-rest-api-specs/issues/5574
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"store_location": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CurrentUser",
								"LocalMachine",
							}, false),
						},
						"store_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"visibility": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"StartTask",
									"Task",
									"RemoteUser",
								}, false),
							},
						},
					},
				},
			},
			"start_task": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command_line": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"max_task_retry_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},

						"wait_for_success": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"environment": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"user_identity": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"auto_user": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"elevation_level": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  string(batch.NonAdmin),
													ValidateFunc: validation.StringInSlice([]string{
														string(batch.NonAdmin),
														string(batch.Admin),
													}, false),
												},
												"scope": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  string(batch.AutoUserScopeTask),
													ValidateFunc: validation.StringInSlice([]string{
														string(batch.AutoUserScopeTask),
														string(batch.AutoUserScopePool),
													}, false),
												},
											},
										},
									},
								},
							},
						},

						"resource_file": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_storage_container_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"blob_prefix": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"file_mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"file_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"http_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"storage_container_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"network_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"public_ips": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
						"endpoint_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"protocol": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(batch.TCP),
											string(batch.UDP),
										}, false),
									},
									"backend_port": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.All(
											validation.IsPortNumber,
											validation.IntNotInSlice([]int{29876, 29877}),
										),
									},
									"frontend_port_range": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validateFrontendPortRangeRange,
									},
									"network_security_group_rules": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"priority": {
													Type:         schema.TypeInt,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.IntAtLeast(150),
												},
												"access": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(batch.Allow),
														string(batch.Deny),
													}, false),
												},
												"source_address_prefix": {
													Type:         schema.TypeString,
													Required:     true,
													ForceNew:     true,
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
	}
}

func resourceArmBatchPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch pool creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	vmSize := d.Get("vm_size").(string)
	maxTasksPerNode := int32(d.Get("max_tasks_per_node").(int))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, accountName, poolName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Batch Pool %q (Account %q / Resource Group %q): %+v", poolName, accountName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_batch_pool", *existing.ID)
		}
	}

	parameters := batch.Pool{
		PoolProperties: &batch.PoolProperties{
			VMSize:          &vmSize,
			DisplayName:     &displayName,
			MaxTasksPerNode: &maxTasksPerNode,
		},
	}

	scaleSettings, err := expandBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("Error expanding scale settings: %+v", err)
	}

	parameters.PoolProperties.ScaleSettings = scaleSettings

	nodeAgentSkuID := d.Get("node_agent_sku_id").(string)

	storageImageReferenceSet := d.Get("storage_image_reference").([]interface{})
	imageReference, err := ExpandBatchPoolImageReference(storageImageReferenceSet)
	if err != nil {
		return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
	}

	if imageReference != nil {
		// if an image reference ID is specified, the user wants use a custom image. This property is mutually exclusive with other properties.
		if imageReference.ID != nil && (imageReference.Offer != nil || imageReference.Publisher != nil || imageReference.Sku != nil || imageReference.Version != nil) {
			return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): Properties version, offer, publish cannot be defined when using a custom image id", poolName, resourceGroup)
		} else if imageReference.ID == nil && (imageReference.Offer == nil || imageReference.Publisher == nil || imageReference.Sku == nil || imageReference.Version == nil) {
			return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): Properties version, offer, publish and sku are mandatory when not using a custom image", poolName, resourceGroup)
		}
	} else {
		return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): image reference property can not be empty", poolName, resourceGroup)
	}

	if startTaskValue, startTaskOk := d.GetOk("start_task"); startTaskOk {
		startTaskList := startTaskValue.([]interface{})
		startTask, startTaskErr := ExpandBatchPoolStartTask(startTaskList)

		if startTaskErr != nil {
			return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, startTaskErr)
		}

		// start task should have a user identity defined
		userIdentity := startTask.UserIdentity
		if userIdentityError := validateUserIdentity(userIdentity); userIdentityError != nil {
			return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, userIdentityError)
		}

		parameters.PoolProperties.StartTask = startTask
	}

	containerConfiguration, err := ExpandBatchPoolContainerConfiguration(d.Get("container_configuration").([]interface{}))
	if err != nil {
		return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
	}

	parameters.PoolProperties.DeploymentConfiguration = &batch.DeploymentConfiguration{
		VirtualMachineConfiguration: &batch.VirtualMachineConfiguration{
			NodeAgentSkuID:         &nodeAgentSkuID,
			ImageReference:         imageReference,
			ContainerConfiguration: containerConfiguration,
		},
	}

	certificates := d.Get("certificate").([]interface{})
	certificateReferences, err := ExpandBatchPoolCertificateReferences(certificates)
	if err != nil {
		return fmt.Errorf("Error expanding `certificate`: %+v", err)
	}
	parameters.PoolProperties.Certificates = certificateReferences

	if err := validateBatchPoolCrossFieldRules(&parameters); err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	parameters.PoolProperties.Metadata = ExpandBatchMetaData(metaDataRaw)

	networkConfiguration := d.Get("network_configuration").([]interface{})
	parameters.PoolProperties.NetworkConfiguration, err = ExpandBatchPoolNetworkConfiguration(networkConfiguration)
	if err != nil {
		return fmt.Errorf("Error expanding `network_configuration`: %+v", err)
	}

	future, err := client.Create(ctx, resourceGroup, accountName, poolName, parameters, "", "")
	if err != nil {
		return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, accountName, poolName)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Batch pool %q (resource group %q) ID", poolName, resourceGroup)
	}

	d.SetId(*read.ID)

	// if the pool is not Steady after the create operation, wait for it to be Steady
	if props := read.PoolProperties; props != nil && props.AllocationState != batch.Steady {
		if err = waitForBatchPoolPendingResizeOperation(ctx, client, resourceGroup, accountName, poolName); err != nil {
			return fmt.Errorf("Error waiting for Batch pool %q (resource group %q) being ready", poolName, resourceGroup)
		}
	}

	return resourceArmBatchPoolRead(d, meta)
}

func resourceArmBatchPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BatchPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving the Batch pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if resp.PoolProperties.AllocationState != batch.Steady {
		log.Printf("[INFO] there is a pending resize operation on this pool...")
		stopPendingResizeOperation := d.Get("stop_pending_resize_operation").(bool)
		if !stopPendingResizeOperation {
			return fmt.Errorf("Error updating the Batch pool %q (Resource Group %q) because of pending resize operation. Set flag `stop_pending_resize_operation` to true to force update", id.Name, id.ResourceGroup)
		}

		log.Printf("[INFO] stopping the pending resize operation on this pool...")
		if _, err = client.StopResize(ctx, id.ResourceGroup, id.AccountName, id.Name); err != nil {
			return fmt.Errorf("Error stopping resize operation for Batch pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		// waiting for the pool to be in steady state
		if err = waitForBatchPoolPendingResizeOperation(ctx, client, id.ResourceGroup, id.AccountName, id.Name); err != nil {
			return fmt.Errorf("Error waiting for Batch pool %q (resource group %q) being ready", id.Name, id.ResourceGroup)
		}
	}

	parameters := batch.Pool{
		PoolProperties: &batch.PoolProperties{},
	}

	scaleSettings, err := expandBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("Error expanding scale settings: %+v", err)
	}

	parameters.PoolProperties.ScaleSettings = scaleSettings

	if startTaskValue, startTaskOk := d.GetOk("start_task"); startTaskOk {
		startTaskList := startTaskValue.([]interface{})
		startTask, startTaskErr := ExpandBatchPoolStartTask(startTaskList)

		if startTaskErr != nil {
			return fmt.Errorf("Error updating Batch pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, startTaskErr)
		}

		// start task should have a user identity defined
		userIdentity := startTask.UserIdentity
		if userIdentityError := validateUserIdentity(userIdentity); userIdentityError != nil {
			return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, userIdentityError)
		}

		parameters.PoolProperties.StartTask = startTask
	}
	certificates := d.Get("certificate").([]interface{})
	certificateReferences, err := ExpandBatchPoolCertificateReferences(certificates)
	if err != nil {
		return fmt.Errorf("Error expanding `certificate`: %+v", err)
	}
	parameters.PoolProperties.Certificates = certificateReferences

	if err := validateBatchPoolCrossFieldRules(&parameters); err != nil {
		return err
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for Batch pool %q (Account name %q / Resource Group %q)..", id.Name, id.AccountName, id.ResourceGroup)
		metaDataRaw := d.Get("metadata").(map[string]interface{})

		parameters.PoolProperties.Metadata = ExpandBatchMetaData(metaDataRaw)
	}

	result, err := client.Update(ctx, id.ResourceGroup, id.AccountName, id.Name, parameters, "")
	if err != nil {
		return fmt.Errorf("Error updating Batch pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	// if the pool is not Steady after the update, wait for it to be Steady
	if props := result.PoolProperties; props != nil && props.AllocationState != batch.Steady {
		if err := waitForBatchPoolPendingResizeOperation(ctx, client, id.ResourceGroup, id.AccountName, id.Name); err != nil {
			return fmt.Errorf("Error waiting for Batch pool %q (resource group %q) being ready", id.Name, id.ResourceGroup)
		}
	}

	return resourceArmBatchPoolRead(d, meta)
}

func resourceArmBatchPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BatchPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Batch pool %q in account %q (Resource Group %q) was not found", id.Name, id.AccountName, id.ResourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Batch pool %q: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("account_name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.PoolProperties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("vm_size", props.VMSize)

		if scaleSettings := props.ScaleSettings; scaleSettings != nil {
			if err := d.Set("auto_scale", flattenBatchPoolAutoScaleSettings(scaleSettings.AutoScale)); err != nil {
				return fmt.Errorf("Error flattening `auto_scale`: %+v", err)
			}
			if err := d.Set("fixed_scale", flattenBatchPoolFixedScaleSettings(scaleSettings.FixedScale)); err != nil {
				return fmt.Errorf("Error flattening `fixed_scale `: %+v", err)
			}
		}

		d.Set("max_tasks_per_node", props.MaxTasksPerNode)

		if props.DeploymentConfiguration != nil &&
			props.DeploymentConfiguration.VirtualMachineConfiguration != nil &&
			props.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference != nil {
			imageReference := props.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference

			d.Set("storage_image_reference", flattenBatchPoolImageReference(imageReference))
			d.Set("node_agent_sku_id", props.DeploymentConfiguration.VirtualMachineConfiguration.NodeAgentSkuID)
		}

		if dcfg := props.DeploymentConfiguration; dcfg != nil {
			if vmcfg := dcfg.VirtualMachineConfiguration; vmcfg != nil {
				d.Set("container_configuration", flattenBatchPoolContainerConfiguration(d, vmcfg.ContainerConfiguration))
			}
		}

		if err := d.Set("certificate", flattenBatchPoolCertificateReferences(props.Certificates)); err != nil {
			return fmt.Errorf("Error flattening `certificate`: %+v", err)
		}

		d.Set("start_task", flattenBatchPoolStartTask(props.StartTask))
		d.Set("metadata", FlattenBatchMetaData(props.Metadata))

		if props.NetworkConfiguration != nil {
			if err := d.Set("network_configuration", FlattenBatchPoolNetworkConfiguration(props.NetworkConfiguration)); err != nil {
				return fmt.Errorf("error setting `network_configuration`: %v", err)
			}
		}
	}

	return nil
}

func resourceArmBatchPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BatchPoolID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Batch pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Batch pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}
	return nil
}

func expandBatchPoolScaleSettings(d *schema.ResourceData) (*batch.ScaleSettings, error) {
	scaleSettings := &batch.ScaleSettings{}

	autoScaleValue, autoScaleOk := d.GetOk("auto_scale")
	fixedScaleValue, fixedScaleOk := d.GetOk("fixed_scale")

	if !autoScaleOk && !fixedScaleOk {
		return nil, fmt.Errorf("Error: auto_scale block or fixed_scale block need to be specified")
	}

	if autoScaleOk && fixedScaleOk {
		return nil, fmt.Errorf("Error: auto_scale and fixed_scale blocks cannot be specified at the same time")
	}

	if autoScaleOk {
		autoScale := autoScaleValue.([]interface{})
		if len(autoScale) == 0 {
			return nil, fmt.Errorf("Error: when scale mode is Auto, auto_scale block is required")
		}

		autoScaleSettings := autoScale[0].(map[string]interface{})

		autoScaleEvaluationInterval := autoScaleSettings["evaluation_interval"].(string)
		autoScaleFormula := autoScaleSettings["formula"].(string)

		scaleSettings.AutoScale = &batch.AutoScaleSettings{
			EvaluationInterval: &autoScaleEvaluationInterval,
			Formula:            &autoScaleFormula,
		}
	} else if fixedScaleOk {
		fixedScale := fixedScaleValue.([]interface{})
		if len(fixedScale) == 0 {
			return nil, fmt.Errorf("Error: when scale mode is Fixed, fixed_scale block is required")
		}

		fixedScaleSettings := fixedScale[0].(map[string]interface{})

		targetDedicatedNodes := int32(fixedScaleSettings["target_dedicated_nodes"].(int))
		targetLowPriorityNodes := int32(fixedScaleSettings["target_low_priority_nodes"].(int))
		resizeTimeout := fixedScaleSettings["resize_timeout"].(string)

		scaleSettings.FixedScale = &batch.FixedScaleSettings{
			ResizeTimeout:          &resizeTimeout,
			TargetDedicatedNodes:   &targetDedicatedNodes,
			TargetLowPriorityNodes: &targetLowPriorityNodes,
		}
	}

	return scaleSettings, nil
}

func waitForBatchPoolPendingResizeOperation(ctx context.Context, client *batch.PoolClient, resourceGroup string, accountName string, poolName string) error {
	// waiting for the pool to be in steady state
	log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped...")
	isSteady := false
	for !isSteady {
		resp, err := client.Get(ctx, resourceGroup, accountName, poolName)
		if err != nil {
			return fmt.Errorf("Error retrieving the Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
		}

		isSteady = resp.PoolProperties.AllocationState == batch.Steady
		time.Sleep(time.Second * 30)
		log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped... New try in 30 seconds...")
	}

	return nil
}

// validateUserIdentity validates that the user identity for a start task has been well specified
// it should have a auto_user block or a user_name defined, but not both at the same time.
func validateUserIdentity(userIdentity *batch.UserIdentity) error {
	if userIdentity == nil {
		return errors.New("user_identity block needs to be specified")
	}

	if userIdentity.AutoUser == nil && userIdentity.UserName == nil {
		return errors.New("auto_user or user_name needs to be specified in the user_identity block")
	}

	if userIdentity.AutoUser != nil && userIdentity.UserName != nil {
		return errors.New("auto_user and user_name cannot be specified in the user_identity at the same time")
	}

	return nil
}

func validateBatchPoolCrossFieldRules(pool *batch.Pool) error {
	// Perform validation across multiple fields as per https://docs.microsoft.com/en-us/rest/api/batchmanagement/pool/create#resourcefile

	if pool.StartTask != nil {
		startTask := *pool.StartTask
		if startTask.ResourceFiles != nil {
			for _, referenceFile := range *startTask.ResourceFiles {
				// Must specify exactly one of AutoStorageContainerName, StorageContainerUrl or HttpUrl
				sourceCount := 0
				if referenceFile.AutoStorageContainerName != nil {
					sourceCount++
				}
				if referenceFile.StorageContainerURL != nil {
					sourceCount++
				}
				if referenceFile.HTTPURL != nil {
					sourceCount++
				}
				if sourceCount != 1 {
					return fmt.Errorf("Exactly one of auto_storage_container_name, storage_container_url and http_url must be specified")
				}

				if referenceFile.BlobPrefix != nil {
					if referenceFile.AutoStorageContainerName == nil && referenceFile.StorageContainerURL == nil {
						return fmt.Errorf("auto_storage_container_name or storage_container_url must be specified when using blob_prefix")
					}
				}

				if referenceFile.HTTPURL != nil {
					if referenceFile.FilePath == nil {
						return fmt.Errorf("file_path must be specified when using http_url")
					}
				}
			}
		}
	}

	return nil
}

func validateFrontendPortRangeRange(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	parts := strings.Split(v, "-")
	if len(parts) != 2 {
		errors = append(errors, fmt.Errorf("expected %s to contain a single '-', got %v", k, i))
		return warnings, errors
	}

	startPort, err := strconv.Atoi(parts[0])
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s on the left of - to be an integer, got %v: %v", k, i, err))
		return warnings, errors
	}

	endPort, err := strconv.Atoi(parts[1])
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s on the right of - to be an integer, got %v: %v", k, i, err))
		return warnings, errors
	}

	if !validPortNumber(startPort) || !validPortNumber(endPort) {
		errors = append(errors, fmt.Errorf("expect values range between 1 and 65534 except ports from `50000` to `55000`, got %v: %v", k, i))
		return warnings, errors
	}
	if endPort-startPort < 100 {
		errors = append(errors, fmt.Errorf("values must be a range of at least 100, got %v: %v", k, i))
		return warnings, errors
	}

	return warnings, errors
}

func validPortNumber(port int) bool {
	return 1 <= port && port < 50000 || 55000 < port && port <= 65535
}
