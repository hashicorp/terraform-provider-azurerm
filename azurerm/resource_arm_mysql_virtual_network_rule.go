package azurerm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMySqlVirtualNetworkRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMySqlVirtualNetworkRuleCreateUpdate,
		Read:   resourceArmMySqlVirtualNetworkRuleRead,
		Update: resourceArmMySqlVirtualNetworkRuleCreateUpdate,
		Delete: resourceArmMySqlVirtualNetworkRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VirtualNetworkRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmMySqlVirtualNetworkRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mysql.VirtualNetworkRulesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subnetId := d.Get("subnet_id").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mysql_virtual_network_rule", *existing.ID)
		}
	}

	// due to a bug in the API we have to ensure the Subnet's configured correctly or the API call will timeout
	// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3719
	subnetsClient := meta.(*ArmClient).network.SubnetsClient
	subnetParsedId, err := azure.ParseAzureResourceID(subnetId)
	if err != nil {
		return err
	}

	subnetResourceGroup := subnetParsedId.ResourceGroup
	virtualNetwork := subnetParsedId.Path["virtualNetworks"]
	subnetName := subnetParsedId.Path["subnets"]
	subnet, err := subnetsClient.Get(ctx, subnetResourceGroup, virtualNetwork, subnetName, "")
	if err != nil {
		if utils.ResponseWasNotFound(subnet.Response) {
			return fmt.Errorf("Subnet with ID %q was not found: %+v", subnetId, err)
		}

		return fmt.Errorf("Error obtaining Subnet %q (Virtual Network %q / Resource Group %q: %+v", subnetName, virtualNetwork, subnetResourceGroup, err)
	}

	containsEndpoint := false
	if props := subnet.SubnetPropertiesFormat; props != nil {
		if endpoints := props.ServiceEndpoints; endpoints != nil {
			for _, e := range *endpoints {
				if e.Service == nil {
					continue
				}

				if strings.EqualFold(*e.Service, "Microsoft.Sql") {
					containsEndpoint = true
					break
				}
			}
		}
	}

	if !containsEndpoint {
		return fmt.Errorf("Error creating MySQL Virtual Network Rule: Subnet %q (Virtual Network %q / Resource Group %q) must contain a Service Endpoint for `Microsoft.Sql`", subnetName, virtualNetwork, subnetResourceGroup)
	}

	parameters := mysql.VirtualNetworkRule{
		VirtualNetworkRuleProperties: &mysql.VirtualNetworkRuleProperties{
			VirtualNetworkSubnetID:           utils.String(subnetId),
			IgnoreMissingVnetServiceEndpoint: utils.Bool(false),
		},
	}

	if _, err = client.CreateOrUpdate(ctx, resourceGroup, serverName, name, parameters); err != nil {
		return fmt.Errorf("Error creating MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	//Wait for the provisioning state to become ready
	log.Printf("[DEBUG] Waiting for MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q) to become ready: %+v", name, serverName, resourceGroup, err)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Initializing", "InProgress", "Unknown", "ResponseNotFound"},
		Target:                    []string{"Ready"},
		Refresh:                   mySQLVirtualNetworkStateStatusCodeRefreshFunc(ctx, client, resourceGroup, serverName, name),
		Timeout:                   30 * time.Minute,
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 5,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q) to be created or updated: %+v", name, serverName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmMySqlVirtualNetworkRuleRead(d, meta)
}

func resourceArmMySqlVirtualNetworkRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mysql.VirtualNetworkRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["virtualNetworkRules"]

	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading MySQL Virtual Network Rule %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading MySQL Virtual Network Rule: %q (MySQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)

	if props := resp.VirtualNetworkRuleProperties; props != nil {
		d.Set("subnet_id", props.VirtualNetworkSubnetID)
	}

	return nil
}

func resourceArmMySqlVirtualNetworkRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mysql.VirtualNetworkRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["virtualNetworkRules"]

	future, err := client.Delete(ctx, resourceGroup, serverName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
		}
	}

	return nil
}

func mySQLVirtualNetworkStateStatusCodeRefreshFunc(ctx context.Context, client *mysql.VirtualNetworkRulesClient, resourceGroup string, serverName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceGroup, serverName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q) returned 404.", resourceGroup, serverName, name)
				return nil, "ResponseNotFound", nil
			}

			return nil, "", fmt.Errorf("Error polling for the state of the MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
		}

		if props := resp.VirtualNetworkRuleProperties; props != nil {
			log.Printf("[DEBUG] Retrieving MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q) returned Status %s", resourceGroup, serverName, name, props.State)
			return resp, string(props.State), nil
		}

		//Valid response was returned but VirtualNetworkRuleProperties was nil. Basically the rule exists, but with no properties for some reason. Assume Unknown instead of returning error.
		log.Printf("[DEBUG] Retrieving MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q) returned empty VirtualNetworkRuleProperties", resourceGroup, serverName, name)
		return resp, "Unknown", nil
	}
}
