package advisor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/advisor/mgmt/2017-04-19/advisor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/advisor/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAdvisor() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAdvisorCreateUpdate,
		Read:   resourceArmAdvisorRead,
		Update: resourceArmAdvisorCreateUpdate,
		Delete: resourceArmAdvisorDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			err := parse.AdvisorSubscriptionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// API will return a default value for low_cpu_threshold if not assigned
			"low_cpu_threshold": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"5", "10", "15", "20"}, true),
			},

			"exclude_resource_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceGroupName,
				},
			},
		},
	}
}

func resourceArmAdvisorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.ConfigurationsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Advisor creation.")

	props := advisor.ConfigDataProperties{
		Exclude: utils.Bool(false),
	}
	//low_cpu_threshold
	if lowCpuThreshold, ok := d.GetOk("low_cpu_threshold"); ok {
		props.LowCPUThreshold = utils.String(lowCpuThreshold.(string))
	}
	//exclude_resource_groups
	var excludeResourceGroup []string
	if v, ok := d.GetOk("exclude_resource_groups"); ok {
		excludeResourceGroup = *utils.ExpandStringSlice(v.(*schema.Set).List())
	}

	parameters := advisor.ConfigData{
		Properties: &props,
	}
	_, err := client.CreateInSubscription(ctx, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Advisor: %+v", err)
	}

	respIter, err := client.ListBySubscriptionComplete(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving Advisor: %+v", err)
	}
	if !respIter.NotDone() {
		return fmt.Errorf("Error retrieving Advisor, the response is empty")
	}
	resp := respIter.Value()

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Advisor ")
	}

	// exclude resource groups in this subscription
	for _, resourceGroup := range excludeResourceGroup {
		parameters := advisor.ConfigData{
			Properties: &advisor.ConfigDataProperties{
				Exclude: utils.Bool(true),
			},
		}
		_, err := client.CreateInResourceGroup(ctx, parameters, resourceGroup)
		if err != nil {
			return fmt.Errorf("Error excluding Advisor (Resource Group %q): %+v", resourceGroup, err)
		}
	}

	d.SetId(*resp.ID)
	return resourceArmAdvisorRead(d, meta)
}

func resourceArmAdvisorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err := parse.AdvisorSubscriptionID(d.Id())
	if err != nil {
		return err
	}
	respIter, err := client.ListBySubscriptionComplete(ctx)
	if err != nil {
		if !respIter.NotDone() {
			d.SetId("")
			log.Printf("[DEBUG] Advisor Configuration was not found  - removing from state!")
			return nil
		}
		return fmt.Errorf("Error reading the state of Advisor Configuration: %+v", err)
	}

	subscriptionResp := respIter.Value()
	var lowCPUThreshold string
	if v := subscriptionResp.Properties.LowCPUThreshold; v != nil {
		lowCPUThreshold = *v
	}
	d.Set("low_cpu_threshold", lowCPUThreshold)
	var excludeResourceGroup []string
	for err := respIter.NextWithContext(ctx); respIter.NotDone() && err == nil; err = respIter.NextWithContext(ctx) {
		resourceGroupResp := respIter.Value()
		var resGroupId *parse.AdvisorResGroupId
		if resGroupId, err = parse.AdvisorResGroupID(*resourceGroupResp.ID); err != nil {
			return err
		}
		if *resourceGroupResp.Properties.Exclude {
			excludeResourceGroup = append(excludeResourceGroup, (*resGroupId).ResourceGroup)
		}
	}
	flattenExcludeResGroups := utils.FlattenStringSlice(&excludeResourceGroup)
	if err := d.Set("exclude_resource_groups", flattenExcludeResGroups); err != nil {
		return fmt.Errorf("Error setting `exclude_resource_groups`: %+v", err)
	}

	return nil
}

func resourceArmAdvisorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.ConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	//we exclude the whole subscription as deleting
	parameters := advisor.ConfigData{
		Properties: &advisor.ConfigDataProperties{
			Exclude: utils.Bool(true),
		},
	}

	_, err := client.CreateInSubscription(ctx, parameters)
	if err != nil {
		return fmt.Errorf("Error deleting Advisor: %+v", err)
	}

	return nil
}
