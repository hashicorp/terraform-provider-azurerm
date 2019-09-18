package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmExpressRouteCircuitAuthorization() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmExpressRouteCircuitAuthorizationCreate,
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

			"resource_group_name": azure.SchemaResourceGroupName(),

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

func resourceArmExpressRouteCircuitAuthorizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ExpressRouteAuthsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	circuitName := d.Get("express_route_circuit_name").(string)

	locks.ByName(circuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(circuitName, expressRouteCircuitResourceName)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, circuitName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %s", name, circuitName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_express_route_circuit_authorization", *existing.ID)
		}
	}

	properties := network.ExpressRouteCircuitAuthorization{
		AuthorizationPropertiesFormat: &network.AuthorizationPropertiesFormat{},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, circuitName, name, properties)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", name, circuitName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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
	client := meta.(*ArmClient).network.ExpressRouteAuthsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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
	client := meta.(*ArmClient).network.ExpressRouteAuthsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	circuitName := id.Path["expressRouteCircuits"]
	name := id.Path["authorizations"]

	locks.ByName(circuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(circuitName, expressRouteCircuitResourceName)

	future, err := client.Delete(ctx, resourceGroup, circuitName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", name, circuitName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for Express Route Circuit Authorization %q (Circuit %q / Resource Group %q) to be deleted: %+v", name, circuitName, resourceGroup, err)
	}

	return nil
}
