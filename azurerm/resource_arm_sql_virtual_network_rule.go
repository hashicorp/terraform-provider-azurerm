package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2015-05-01-preview/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlVnetRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlVnetRuleCreateUpdate,
		Read:   resourceArmSqlVnetRuleRead,
		Update: resourceArmSqlVnetRuleCreateUpdate,
		Delete: resourceArmSqlVnetRuleDelete,
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

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"virtual_network_subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: validation?
			},

			"ignore_missing_vnet_service_endpoint": {
				Type:     schema.TypeBool,
				Optional: true,
				// TODO: validation?
			},
		},
	}
}

func resourceArmSqlVnetRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlVnetRulesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	virtualNetworkSubnetId := d.Get("virtual_network_subnet_id").(string)
	ignoreMissingVnetServiceEndpoint := d.Get("ignore_missing_vnet_service_endpoint").(bool)

	parameters := sql.VnetRule{
		VnetRuleProperties: &sql.VnetRuleProperties{
			virtualNetworkSubnetId: utils.String(virtualNetworkSubnetId),
			ignoreMissingVnetServiceEndpoint: utils.String(ignoreMissingVnetServiceEndpoint),
		},
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating SQL Virtual Network Rule: %+v", err)
	}

	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving SQL Virtual Network Rule: %+v", err)
	}

	d.SetId(*resp.ID)

	return resourceArmSqlVnetRuleRead(d, meta)
}

func resourceArmSqlVnetRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlVnetRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["virtualNetworkRules"]

	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Virtual Network Rule %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Virtual Network Rule: %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)
	d.Set("virtual_network_subnet_id", resp.virtualNetworkSubnetId)
	d.Set("ignore_missing_vnet_service_endpoint", resp.ignoreMissingVnetServiceEndpoint)

	return nil
}

func resourceArmSqlVnetRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlVnetRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["virtualNetworkRules"]

	resp, err := client.Delete(ctx, resourceGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error deleting SQL Virtual Network Rule: %+v", err)
	}

	return nil
}
