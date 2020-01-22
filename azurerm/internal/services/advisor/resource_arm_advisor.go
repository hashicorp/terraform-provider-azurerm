package advisor

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/Azure/azure-sdk-for-go/services/advisor/mgmt/2017-04-19/advisor"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAdvisor() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAdvisorCreateUpdate,
		Read:   resourceArmAdvisorRead,
		Update: resourceArmAdvisorCreateUpdate,
		Delete: resourceArmAdvisorDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceGroupName,
			},

			"exclude": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"low_cpu_threshold": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"5", "10", "15", "20"}, true),
				ConflictsWith: []string{"resource_group_name"},
			},
		},
	}
}

func resourceArmAdvisorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.ConfigurationsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Advisor creation.")

	p := advisor.ConfigDataProperties{}
	//resource_group
	var resourceGroup string
	if resourceGroupName, ok := d.GetOkExists("resource_group_name"); ok {
		resourceGroup = resourceGroupName.(string)
	}
	//exclude
	if exclude, ok := d.GetOkExists("exclude"); ok {
		exclude := exclude.(bool)
		p.Exclude = &exclude
	}
	//low_cpu_threshold
	if lowCpuThreshold, ok := d.GetOkExists("low_cpu_threshold"); ok {
		lowCpuThreshold := lowCpuThreshold.(string)
		p.LowCPUThreshold = &lowCpuThreshold
	}

	parameters := advisor.ConfigData{
		Properties: &p,
	}

	if resourceGroup != "" {
		_, err := client.CreateInResourceGroup(ctx, parameters, resourceGroup)
		if err != nil {
			return fmt.Errorf("Error creating Advisor (Resource Group %q): %+v", resourceGroup, err)
		}

		readlist, err := client.ListByResourceGroup(ctx, resourceGroup)
		if err != nil {
			return fmt.Errorf("Error retrieving Advisor (Resource Group %q): %+v", resourceGroup, err)
		}
		if readlist.IsEmpty() {
			return fmt.Errorf("Error retrieving Advisor (Resource Group %q)", resourceGroup)
		}
		read := (*readlist.Value)[0]

		if read.ID == nil {
			return fmt.Errorf("Cannot read Advisor (resource group %q) ID", resourceGroup)
		}

		d.SetId(*read.ID)
	} else {
		_, err := client.CreateInSubscription(ctx, parameters)
		if err != nil {
			return fmt.Errorf("Error creating Advisor: %+v", err)
		}

		readlist, err := client.ListBySubscription(ctx)
		if err != nil {
			return fmt.Errorf("Error retrieving Advisor: %+v", err)
		}
		// here is a sdk problem, which NotDone return false when the response is empty
		if !readlist.NotDone() {
			return fmt.Errorf("Error retrieving Advisor, the response page enumeration should be started or is not yet complete")
		}
		read := readlist.Values()[0]

		if read.ID == nil {
			return fmt.Errorf("Cannot read Advisor ")
		}

		d.SetId(*read.ID)
	}

	return resourceArmAdvisorRead(d, meta)
}

func resourceArmAdvisorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	if resourceGroup!="" {
		resplist, err := client.ListByResourceGroup(ctx, resourceGroup)
		if err != nil {
			if utils.ResponseWasNotFound(resplist.Response) || resplist.IsEmpty() {
				d.SetId("")
				log.Printf("[DEBUG] Advisor Configuration of Resource Group %q was not found  - removing from state!", resourceGroup)
				return nil
			}
			return fmt.Errorf("Error reading the state of Advisor Configuration: %+v", err)
		}
		resp := (*resplist.Value)[0]

		if resourceGroup != "" {
			d.Set("resource_group_name", resourceGroup)
		}

		if exclude := resp.Properties.Exclude; exclude != nil {
			d.Set("exclude", exclude)
		}

		if lowCPUThreshold := resp.Properties.LowCPUThreshold; lowCPUThreshold != nil {
			d.Set("low_cpu_threshold", lowCPUThreshold)
		}

		return nil
	} else {
		resplist, err := client.ListBySubscription(ctx)
		if err != nil {
			if !resplist.NotDone() {
				d.SetId("")
				log.Printf("[DEBUG] Advisor Configuration was not found  - removing from state!")
				return nil
			}
			return fmt.Errorf("Error reading the state of Advisor Configuration: %+v", err)
		}

		resp := resplist.Values()[0]

		if exclude := resp.Properties.Exclude; exclude != nil {
			d.Set("exclude", exclude)
		}

		if lowCPUThreshold := resp.Properties.LowCPUThreshold; lowCPUThreshold != nil {
			d.Set("low_cpu_threshold", lowCPUThreshold)
		}

		return nil
	}
}

func resourceArmAdvisorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.ConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	if resourceGroup!="" {
		parameters := advisor.ConfigData{
			Properties: &advisor.ConfigDataProperties{
				Exclude: utils.Bool(true),
			},
		}

		_, err := client.CreateInResourceGroup(ctx, parameters, resourceGroup)
		if err != nil {
			return fmt.Errorf("Error deleting Advisor (Resource Group %q): %+v", resourceGroup, err)
		}

		return nil
	} else {
		parameters := advisor.ConfigData{
			Properties: &advisor.ConfigDataProperties{
				Exclude: utils.Bool(true),
			},
		}

		_, err := client.CreateInSubscription(ctx, parameters)
		if err != nil {
			return fmt.Errorf("Error deleting Advisor (Resource Group %q): %+v", resourceGroup, err)
		}

		return nil
	}
}
