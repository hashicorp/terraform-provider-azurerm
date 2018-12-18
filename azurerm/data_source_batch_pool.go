package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmBatchPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBatchPoolRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateAzureRMBatchPoolName,
			},
			"resource_group_name": resourceGroupNameDiffSuppressSchema(),
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAzureRMBatchAccountName,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_scale":             azure.SchemaBatchPoolFixedScaleForDataSource(),
			"auto_scale":              azure.SchemaBatchPoolAutoScaleForDataSource(),
			"storage_image_reference": azure.SchemaBatchPoolImageReferenceForDataSource(),
			"node_agent_sku_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_task": azure.SchemaBatchPoolStartTaskForDataSource(),
		},
	}
}

func dataSourceArmBatchPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchPoolClient

	name := d.Get("name").(string)
	accountName := d.Get("account_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	ctx := meta.(*ArmClient).StopContext
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

		vmCfg := resp.DeploymentConfiguration.VirtualMachineConfiguration

		if err := d.Set("storage_image_reference", azure.FlattenBatchPoolImageReference(vmCfg.ImageReference)); err != nil {
			return fmt.Errorf("Error setting AzureRM Batch Pool Image Reference: %#v", err)
		}

		d.Set("node_agent_sku_id", vmCfg.NodeAgentSkuID)
	}

	return nil
}
