package mysql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMySQLVirtualNetworkRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMySQLVirtualNetworkRuleCreateUpdate,
		Read:   resourceArmMySQLVirtualNetworkRuleRead,
		Update: resourceArmMySQLVirtualNetworkRuleCreateUpdate,
		Delete: resourceArmMySQLVirtualNetworkRuleDelete,
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.VirtualNetworkRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmMySQLVirtualNetworkRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subnetId := d.Get("subnet_id").(string)

	if d.IsNewResource() {
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

	parameters := mysql.VirtualNetworkRule{
		VirtualNetworkRuleProperties: &mysql.VirtualNetworkRuleProperties{
			VirtualNetworkSubnetID:           utils.String(subnetId),
			IgnoreMissingVnetServiceEndpoint: utils.Bool(false),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, name, parameters); err != nil {
		return fmt.Errorf("Error creating MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	// Wait for the provisioning state to become ready
	log.Printf("[DEBUG] Waiting for MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q) to become ready", name, serverName, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Initializing", "InProgress", "Unknown", "ResponseNotFound"},
		Target:                    []string{"Ready"},
		Refresh:                   mySQLVirtualNetworkStateStatusCodeRefreshFunc(ctx, client, resourceGroup, serverName, name),
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 5,
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q) to be created or updated: %+v", name, serverName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmMySQLVirtualNetworkRuleRead(d, meta)
}

func resourceArmMySQLVirtualNetworkRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

func resourceArmMySQLVirtualNetworkRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

		// Valid response was returned but VirtualNetworkRuleProperties was nil. Basically the rule exists, but with no properties for some reason. Assume Unknown instead of returning error.
		log.Printf("[DEBUG] Retrieving MySQL Virtual Network Rule %q (MySQL Server: %q, Resource Group: %q) returned empty VirtualNetworkRuleProperties", resourceGroup, serverName, name)
		return resp, "Unknown", nil
	}
}
