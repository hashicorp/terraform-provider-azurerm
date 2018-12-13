package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBatchPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBatchPoolCreate,
		Read:   resourceArmBatchPoolRead,
		Update: resourceArmBatchPoolUpdate,
		Delete: resourceArmBatchPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMBatchPoolName,
			},
			"resource_group_name": resourceGroupNameDiffSuppressSchema(),
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMBatchAccountName,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vm_size": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "Standard_A1",
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ForceNew:         true,
			},
			"scale_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(azure.BatchPoolFixedScale),
				ValidateFunc: validation.StringInSlice([]string{
					string(azure.BatchPoolFixedScale),
					string(azure.BatchPoolAutoScale),
				}, false),
			},
			"target_dedicated_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"target_low_priority_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"resize_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PT15M",
			},
			"autoscale_evaluation_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PT15M",
			},
			"autoscale_formula": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_image_reference": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"offer": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"sku": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
				Set: resourceArmVirtualMachineStorageImageReferenceHash,
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
		},
	}
}

func resourceArmBatchPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchPoolClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Batch pool creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	vmSize := d.Get("vm_size").(string)

	parameters := batch.Pool{
		PoolProperties: &batch.PoolProperties{
			VMSize:      &vmSize,
			DisplayName: &displayName,
		},
	}

	scaleMode := azure.BatchPoolScaleMode(d.Get("scale_mode").(string))

	if scaleMode == azure.BatchPoolAutoScale {
		autoScaleEvaluationInterval := d.Get("autoscale_evaluation_interval").(string)
		autoScaleFormula := d.Get("autoscale_formula").(string)

		if autoScaleFormula == "" {
			return fmt.Errorf("Error: when scale mode is set to Auto, auto_scale formula cannot be empty")
		}

		parameters.PoolProperties.ScaleSettings = &batch.ScaleSettings{
			AutoScale: &batch.AutoScaleSettings{
				EvaluationInterval: &autoScaleEvaluationInterval,
				Formula:            &autoScaleFormula,
			},
		}
	} else if scaleMode == azure.BatchPoolFixedScale {
		targetDedicatedNodes := int32(d.Get("target_dedicated_nodes").(int))
		targetLowPriorityNodes := int32(d.Get("target_low_priority_nodes").(int))
		resizeTimeout := d.Get("resize_timeout").(string)

		parameters.PoolProperties.ScaleSettings = &batch.ScaleSettings{
			FixedScale: &batch.FixedScaleSettings{
				ResizeTimeout:          &resizeTimeout,
				TargetDedicatedNodes:   &targetDedicatedNodes,
				TargetLowPriorityNodes: &targetLowPriorityNodes,
			},
		}
	} else {
		return fmt.Errorf("Error: scale mode should be either AutoScale of FixedScale")
	}

	nodeAgentSkuID := d.Get("node_agent_sku_id").(string)
	storageImageReferenceSet := d.Get("storage_image_reference").(*schema.Set)

	if storageImageReferenceSet == nil || storageImageReferenceSet.Len() == 0 {
		return fmt.Errorf("Error: storage image reference should be defined")
	}

	storageImageRef := storageImageReferenceSet.List()[0].(map[string]interface{})

	storageImageRefOffer, storageImageRefOfferOk := storageImageRef["offer"].(string)
	if !storageImageRefOfferOk {
		return fmt.Errorf("Error: storage image reference offer should be defined")
	}

	storageImageRefPublisher, storageImageRefPublisherOK := storageImageRef["publisher"].(string)
	if !storageImageRefPublisherOK {
		return fmt.Errorf("Error: storage image reference publisher should be defined")
	}

	storageImageRefSku, storageImageRefSkuOK := storageImageRef["sku"].(string)
	if !storageImageRefSkuOK {
		return fmt.Errorf("Error: storage image reference sku should be defined")
	}

	storageImageRefVersion, storageImageRefVersionOK := storageImageRef["version"].(string)
	if !storageImageRefVersionOK {
		return fmt.Errorf("Error: storage image reference version should be defined")
	}

	parameters.DeploymentConfiguration = &batch.DeploymentConfiguration{
		VirtualMachineConfiguration: &batch.VirtualMachineConfiguration{
			NodeAgentSkuID: &nodeAgentSkuID,
			ImageReference: &batch.ImageReference{
				Offer:     &storageImageRefOffer,
				Publisher: &storageImageRefPublisher,
				Sku:       &storageImageRefSku,
				Version:   &storageImageRefVersion,
			},
		},
	}

	future, err := client.Create(ctx, resourceGroupName, accountName, name, parameters, "", "")
	if err != nil {
		return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Batch pool %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch pool %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Batch pool %q (resource group %q) ID", name, resourceGroupName)
	}

	d.SetId(*read.ID)

	return resourceArmBatchPoolRead(d, meta)
}

func resourceArmBatchPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).batchPoolClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	poolName := id.Path["pools"]
	accountName := id.Path["batchAccounts"]

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName)
	if err != nil {
		return fmt.Errorf("Error retrieving the Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
	}

	if resp.AllocationState != batch.Steady {
		log.Printf("[INFO] there is a pending resize operation on this pool...")
		stopPendingResizeOperation := d.Get("stop_pending_resize_operation").(bool)
		if !stopPendingResizeOperation {
			return fmt.Errorf("Error updating the Batch pool %q (Resource Group %q) because of pending resize operation. Set flag `stop_pending_resize_operation` to true to force update.", poolName, resourceGroup)
		}

		log.Printf("[INFO] stopping the pending resize operation on this pool...")
		if _, err = client.StopResize(ctx, resourceGroup, accountName, poolName); err != nil {
			return fmt.Errorf("Error stopping resize operation for Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
		}

		// waiting for the pool to be in steady state
		log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped...")
		isSteady := false
		for !isSteady {
			resp, err = client.Get(ctx, resourceGroup, accountName, poolName)
			if err != nil {
				return fmt.Errorf("Error retrieving the Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
			}

			isSteady = resp.AllocationState == batch.Steady
			time.Sleep(time.Minute * 2)
			log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped... New try in 2 minutes...")
		}
	}

	parameters := batch.Pool{
		PoolProperties: &batch.PoolProperties{},
	}

	scaleMode := azure.BatchPoolScaleMode(d.Get("scale_mode").(string))
	if scaleMode == azure.BatchPoolAutoScale {
		autoScaleEvaluationInterval := d.Get("autoscale_evaluation_interval").(string)
		autoScaleFormula := d.Get("autoscale_formula").(string)

		if autoScaleFormula == "" {
			return fmt.Errorf("Error: when scale mode is set to Auto, auto_scale formula cannot be empty")
		}

		parameters.PoolProperties.ScaleSettings = &batch.ScaleSettings{
			AutoScale: &batch.AutoScaleSettings{
				EvaluationInterval: &autoScaleEvaluationInterval,
				Formula:            &autoScaleFormula,
			},
		}
	} else if scaleMode == azure.BatchPoolFixedScale {
		targetDedicatedNodes := int32(d.Get("target_dedicated_nodes").(int))
		targetLowPriorityNodes := int32(d.Get("target_low_priority_nodes").(int))
		resizeTimeout := d.Get("resize_timeout").(string)

		parameters.PoolProperties.ScaleSettings = &batch.ScaleSettings{
			FixedScale: &batch.FixedScaleSettings{
				ResizeTimeout:          &resizeTimeout,
				TargetDedicatedNodes:   &targetDedicatedNodes,
				TargetLowPriorityNodes: &targetLowPriorityNodes,
			},
		}
	} else {
		return fmt.Errorf("Error: scale mode should be either AutoScale of FixedScale")
	}

	if _, err = client.Update(ctx, resourceGroup, accountName, poolName, parameters, ""); err != nil {
		return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, accountName, poolName)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Batch pool %q (resource group %q) ID", poolName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmBatchPoolRead(d, meta)
}

func resourceArmBatchPoolDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).batchPoolClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	poolName := id.Path["pools"]
	accountName := id.Path["batchAccounts"]

	future, err := client.Delete(ctx, resourceGroup, accountName, poolName)
	if err != nil {
		return fmt.Errorf("Error deleting Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
		}
	}
	return nil
}

func resourceArmBatchPoolRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).batchPoolClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	poolName := id.Path["pools"]
	accountName := id.Path["batchAccounts"]

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Batch pool %q in account %q (Resource Group %q) was not found", poolName, accountName, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Batch pool %q: %+v", poolName, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", poolName)
	d.Set("account_name", accountName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("vm_size", resp.VMSize)

	if resp.ScaleSettings != nil {
		if resp.ScaleSettings.AutoScale != nil {
			d.Set("scale_mode", azure.BatchPoolAutoScale)
			d.Set("autoscale_evaluation_interval", resp.ScaleSettings.AutoScale.EvaluationInterval)
			d.Set("autoscale_formula", resp.ScaleSettings.AutoScale.Formula)
		} else if resp.ScaleSettings.FixedScale != nil {
			d.Set("scale_mode", azure.BatchPoolFixedScale)
			d.Set("target_dedicated_nodes", resp.ScaleSettings.FixedScale.TargetDedicatedNodes)
			d.Set("target_low_priority_nodes", resp.ScaleSettings.FixedScale.TargetLowPriorityNodes)
			d.Set("resize_timeout", resp.ScaleSettings.FixedScale.ResizeTimeout)
		}
	}

	if resp.DeploymentConfiguration != nil &&
		resp.DeploymentConfiguration.VirtualMachineConfiguration != nil &&
		resp.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference != nil {

		imageReference := resp.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference
		storageImageRef := &schema.Set{F: resourceArmVirtualMachineStorageImageReferenceHash}
		storageImageRef.Add(
			map[string]string{
				"publisher": *imageReference.Publisher,
				"offer":     *imageReference.Offer,
				"sku":       *imageReference.Sku,
				"version":   *imageReference.Version,
			})

		d.Set("storage_image_reference", storageImageRef)
		d.Set("node_agent_sku_id", resp.DeploymentConfiguration.VirtualMachineConfiguration.NodeAgentSkuID)
	}

	return nil
}

func validateAzureRMBatchPoolName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"any combination of alphanumeric characters including hyphens and underscores are allowed in %q: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 64 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}

func flattenAzureRmBatchPoolImageReference(image *batch.ImageReference) []interface{} {
	result := make(map[string]interface{})
	if image.Publisher != nil {
		result["publisher"] = *image.Publisher
	}
	if image.Offer != nil {
		result["offer"] = *image.Offer
	}
	if image.Sku != nil {
		result["sku"] = *image.Sku
	}
	if image.Version != nil {
		result["version"] = *image.Version
	}
	if image.ID != nil {
		result["id"] = *image.ID
	}

	return []interface{}{result}
}
