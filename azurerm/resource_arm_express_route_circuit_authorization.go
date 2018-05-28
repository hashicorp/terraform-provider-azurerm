package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmExpressRouteCircuitAuthorization() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmExpressRouteCircuitAuthorizationCreateUpdate,
		Read:   resourceArmExpressRouteCircuitAuthorizationRead,
		Delete: resourceArmExpressRouteCircuitAuthorizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"express_route_circuit_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"authorization_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"authorization_use_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmExpressRouteCircuitAuthorizationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).expressRouteAuthsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	circuitName := d.Get("express_route_circuit_name").(string)

	properties := network.ExpressRouteCircuitAuthorization{
		AuthorizationPropertiesFormat: &network.AuthorizationPropertiesFormat{},
	}

	azureRMLockByName(circuitName, expressRouteCircuitResourceName)
	defer azureRMUnlockByName(circuitName, expressRouteCircuitResourceName)

	future, err := client.CreateOrUpdate(ctx, resourceGroup, circuitName, name, properties)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", name, circuitName, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for Express Route Circuit Authorization %q (Circuit %q / Resource Group %q) to finish creating/updating: %+v", name, circuitName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, circuitName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", name, circuitName, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmExpressRouteCircuitAuthorizationRead(d, meta)
}

func resourceArmExpressRouteCircuitAuthorizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).expressRouteAuthsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	circuitName := id.Path["expressRouteCircuits"]
	name := id.Path["authorizations"]

	resp, err := client.Get(ctx, resourceGroup, circuitName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", name, circuitName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("express_route_circuit_name", circuitName)

	if props := resp.AuthorizationPropertiesFormat; props != nil {
		d.Set("authorization_key", props.AuthorizationKey)
		d.Set("authorization_use_status", string(props.AuthorizationUseStatus))
	}

	return nil
}

func resourceArmExpressRouteCircuitAuthorizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).expressRouteAuthsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	circuitName := id.Path["expressRouteCircuits"]
	name := id.Path["authorizations"]

	azureRMLockByName(circuitName, expressRouteCircuitResourceName)
	defer azureRMUnlockByName(circuitName, expressRouteCircuitResourceName)

	future, err := client.Delete(ctx, resourceGroup, circuitName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", name, circuitName, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for Express Route Circuit Authorization %q (Circuit %q / Resource Group %q) to be deleted: %+v", name, circuitName, resourceGroup, err)
	}

	return nil
}
