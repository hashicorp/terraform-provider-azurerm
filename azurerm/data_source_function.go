package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmFunction() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmServiceBusNamespaceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"function_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trigger_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmFunctionRead(d *schema.ResourceData, meta interface{}) error {
	webClient := meta.(*ArmClient).web
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	functionName := d.Get("function_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	log.Printf("[DEBUG] Waiting for Function %q in Function app %q (Resource Group %q) to become available", functionName, name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"pending"},
		Target:                    []string{"available"},
		Refresh:                   functionAvailabilityChecker(ctx, name, functionName, resourceGroup, webClient),
		Timeout:                   10 * time.Minute,
		Delay:                     30 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 3,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Function %q in Function app %q (Resource Group %q) to become available", functionName, name, resourceGroup)
	}

	resp, err := webClient.AppServicesClient.ListFunctionSecrets(ctx, resourceGroup, name, functionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Function %q was not found in Function App %q in Resource Group %q", functionName, name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Function %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if triggerUrl := resp.TriggerURL; triggerUrl != nil {
		d.Set("trigger_url", triggerUrl)
	}

	if key := resp.Key; key != nil {
		d.Set("key", key)
	}

	return nil
}

func functionAvailabilityChecker(ctx context.Context, name, functionName, resourceGroup string, client *web.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if Function %q is available..", functionName)

		_, err := client.AppServicesClient.GetFunction(ctx, name, functionName, resourceGroup)
		if err != nil {
			log.Printf("[DEBUG] Didn't find Function at %q", name)
			return nil, "pending", err
		}

		log.Printf("[DEBUG] Found function at %q", functionName)
		return "available", "available", nil
	}
}
