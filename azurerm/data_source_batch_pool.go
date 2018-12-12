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
				ValidateFunc: validateAzureRMBatchPoolName,
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
			"scale_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_dedicated_nodes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target_low_priority_nodes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resize_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"autoscale_evaluation_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"autoscale_formula": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_image_reference": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
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
				Computed: true,
			},
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

		if err := d.Set("storage_image_reference", schema.NewSet(resourceArmVirtualMachineStorageImageReferenceHash, flattenAzureRmBatchPoolImageReference(imageReference))); err != nil {
			return fmt.Errorf("[DEBUG] Error setting AzureRM Batch Pool Image Reference: %#v", err)
		}

		d.Set("node_agent_sku_id", resp.DeploymentConfiguration.VirtualMachineConfiguration.NodeAgentSkuID)
	}

	return nil
}
