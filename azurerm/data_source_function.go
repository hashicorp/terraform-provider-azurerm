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

	webmgmt "github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
)

func dataSourceArmFunction() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmFunctionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"function_name": {
				Type:     schema.TypeString,
				Required: true,
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
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      functionAvailabilityChecker(ctx, name, functionName, resourceGroup, webClient),
		Timeout:      10 * time.Minute,
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}

	resp, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Function %q in Function app %q (Resource Group %q) to become available", functionName, name, resourceGroup)
	}

	functionSecret := resp.(webmgmt.FunctionSecrets)

	if functionSecret.TriggerURL == nil {
		return fmt.Errorf("Error retrieving key for Function %q in Function app %q (Resource Group %q). TriggerURL returned nil from API", functionName, name, resourceGroup)
	}
	if functionSecret.Key == nil {
		return fmt.Errorf("Error retrieving key for Function %q in Function app %q (Resource Group %q). Key returned nil from API", functionName, name, resourceGroup)
	}

	d.SetId(*functionSecret.TriggerURL)

	d.Set("trigger_url", functionSecret.TriggerURL)
	d.Set("key", functionSecret.Key)

	return nil
}

func functionAvailabilityChecker(ctx context.Context, name, functionName, resourceGroup string, client *web.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if Function %q is available..", functionName)

		function, err := client.AppServicesClient.ListFunctionSecrets(ctx, resourceGroup, name, functionName)

		if err != nil || function.StatusCode != 200 {
			log.Printf("[DEBUG] Didn't find Function at %q", name)
			return nil, "pending", err
		}

		log.Printf("[DEBUG] Found function at %q", functionName)
		return function, "available", nil
	}
}
