package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	"github.com/hashicorp/terraform/helper/schema"
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
				ValidateFunc: azure.ValidateAzureRMBatchPoolName,
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
			"fixed_scale":             azure.SchemaBatchPoolFixedScale(),
			"auto_scale":              azure.SchemaBatchPoolAutoScale(),
			"storage_image_reference": azure.SchemaBatchPoolImageReference(),
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
			"start_task": azure.SchemaBatchPoolStartTask(),
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

	scaleSettings, err := getBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("Error: scale settings: %+v", err)
	}

	parameters.ScaleSettings = scaleSettings

	nodeAgentSkuID := d.Get("node_agent_sku_id").(string)
	storageImageReferenceSet := d.Get("storage_image_reference").(*schema.Set)
	imageReference, err := azure.ExpandBatchPoolImageReference(storageImageReferenceSet)
	if err != nil {
		return fmt.Errorf("Error creating Batch pool %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	parameters.DeploymentConfiguration = &batch.DeploymentConfiguration{
		VirtualMachineConfiguration: &batch.VirtualMachineConfiguration{
			NodeAgentSkuID: &nodeAgentSkuID,
			ImageReference: imageReference,
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
			return fmt.Errorf("Error updating the Batch pool %q (Resource Group %q) because of pending resize operation. Set flag `stop_pending_resize_operation` to true to force update", poolName, resourceGroup)
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

	scaleSettings, err := getBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("Error: scale settings: %+v", err)
	}

	parameters.ScaleSettings = scaleSettings

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
			d.Set("auto_scale", azure.FlattenBatchPoolAutoScaleSettings(resp.ScaleSettings.AutoScale))
		} else if resp.ScaleSettings.FixedScale != nil {
			d.Set("fixed_scale", azure.FlattenBatchPoolFixedScaleSettings(resp.ScaleSettings.FixedScale))
		}
	}

	if resp.DeploymentConfiguration != nil &&
		resp.DeploymentConfiguration.VirtualMachineConfiguration != nil &&
		resp.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference != nil {

		imageReference := resp.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference

		d.Set("storage_image_reference", azure.FlattenBatchPoolImageReference(imageReference))
		d.Set("node_agent_sku_id", resp.DeploymentConfiguration.VirtualMachineConfiguration.NodeAgentSkuID)
	}

	return nil
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

func getBatchPoolScaleSettings(d *schema.ResourceData) (*batch.ScaleSettings, error) {
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
	} else {
		return nil, fmt.Errorf("Error: scale mode should be either AutoScale of FixedScale")
	}

	return scaleSettings, nil
}
