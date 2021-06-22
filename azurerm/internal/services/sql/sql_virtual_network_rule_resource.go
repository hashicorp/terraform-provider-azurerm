package sql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSqlVirtualNetworkRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSqlVirtualNetworkRuleCreateUpdate,
		Read:   resourceSqlVirtualNetworkRuleRead,
		Update: resourceSqlVirtualNetworkRuleCreateUpdate,
		Delete: resourceSqlVirtualNetworkRuleDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VirtualNetworkRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlServerName,
			},

			"subnet_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"ignore_missing_vnet_service_endpoint": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false, // When not provided, Azure defaults to false
			},
		},
	}
}

func resourceSqlVirtualNetworkRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	virtualNetworkSubnetId := d.Get("subnet_id").(string)
	ignoreMissingVnetServiceEndpoint := d.Get("ignore_missing_vnet_service_endpoint").(bool)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_virtual_network_rule", *existing.ID)
		}
	}

	parameters := sql.VirtualNetworkRule{
		VirtualNetworkRuleProperties: &sql.VirtualNetworkRuleProperties{
			VirtualNetworkSubnetID:           utils.String(virtualNetworkSubnetId),
			IgnoreMissingVnetServiceEndpoint: utils.Bool(ignoreMissingVnetServiceEndpoint),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, name, parameters); err != nil {
		return fmt.Errorf("Error creating SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	// Wait for the provisioning state to become ready
	log.Printf("[DEBUG] Waiting for SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q) to become ready", name, serverName, resourceGroup)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Initializing", "InProgress", "Unknown", "ResponseNotFound"},
		Target:                    []string{"Ready"},
		Refresh:                   sqlVirtualNetworkStateStatusCodeRefreshFunc(ctx, client, resourceGroup, serverName, name),
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 5,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q) to be created or updated: %+v", name, serverName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceSqlVirtualNetworkRuleRead(d, meta)
}

func resourceSqlVirtualNetworkRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] SQL Virtual Network Rule %q (Server %q / Resource Group %q) was not found - removing from state", id.Name, id.ServerName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Virtual Network Rule %q (Server %q / Resource Group: %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.VirtualNetworkRuleProperties; props != nil {
		d.Set("subnet_id", props.VirtualNetworkSubnetID)
		d.Set("ignore_missing_vnet_service_endpoint", props.IgnoreMissingVnetServiceEndpoint)
	}

	return nil
}

func resourceSqlVirtualNetworkRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("deleting SQL Virtual Network Rule %q (Server %q / Resource Group: %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("waiting for deletion of SQL Virtual Network Rule %q (Server %q / Resource Group: %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	return nil
}

/*
	This function refreshes and checks the state of the SQL Virtual Network Rule.

	Response will contain a VirtualNetworkRuleProperties struct with a State property. The state property contain one of the following states (except ResponseNotFound).
	* Deleting
	* Initializing
	* InProgress
	* Unknown
	* Ready
	* ResponseNotFound (Custom state in case of 404)
*/
func sqlVirtualNetworkStateStatusCodeRefreshFunc(ctx context.Context, client *sql.VirtualNetworkRulesClient, resourceGroup string, serverName string, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving Virtual Network Rule %q (Server %q / Resource Group %q) returned 404.", name, serverName, resourceGroup)
				return nil, "ResponseNotFound", nil
			}

			return nil, "", fmt.Errorf("polling for the state of the Virtual Network Rule %q (Server %q / Resource Group %q): %+v", name, serverName, resourceGroup, err)
		}

		if props := resp.VirtualNetworkRuleProperties; props != nil {
			log.Printf("[DEBUG] Retrieving Virtual Network Rule %q (Server %q / Resource Group %q) returned Status %s", name, serverName, resourceGroup, props.State)
			return resp, string(props.State), nil
		}

		// Valid response was returned but VirtualNetworkRuleProperties was nil. Basically the rule exists, but with no properties for some reason. Assume Unknown instead of returning error.
		log.Printf("[DEBUG] Retrieving Virtual Network Rule %q (SQL Server %q / Resource Group %q) returned empty VirtualNetworkRuleProperties", name, serverName, resourceGroup)
		return resp, "Unknown", nil
	}
}
