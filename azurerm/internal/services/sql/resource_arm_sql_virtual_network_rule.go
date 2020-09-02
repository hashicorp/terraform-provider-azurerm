package sql

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlVirtualNetworkRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlVirtualNetworkRuleCreateUpdate,
		Read:   resourceArmSqlVirtualNetworkRuleRead,
		Update: resourceArmSqlVirtualNetworkRuleCreateUpdate,
		Delete: resourceArmSqlVirtualNetworkRuleDelete,

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
				ValidateFunc: ValidateSqlVirtualNetworkRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ignore_missing_vnet_service_endpoint": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false, // When not provided, Azure defaults to false
			},
		},
	}
}

func resourceArmSqlVirtualNetworkRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Initializing", "InProgress", "Unknown", "ResponseNotFound"},
		Target:                    []string{"Ready"},
		Refresh:                   sqlVirtualNetworkStateStatusCodeRefreshFunc(ctx, client, resourceGroup, serverName, name),
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 5,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q) to be created or updated: %+v", name, serverName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmSqlVirtualNetworkRuleRead(d, meta)
}

func resourceArmSqlVirtualNetworkRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.VirtualNetworkRulesClient
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
			log.Printf("[INFO] Error reading SQL Virtual Network Rule %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Virtual Network Rule: %q (SQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)

	if props := resp.VirtualNetworkRuleProperties; props != nil {
		d.Set("subnet_id", props.VirtualNetworkSubnetID)
		d.Set("ignore_missing_vnet_service_endpoint", props.IgnoreMissingVnetServiceEndpoint)
	}

	return nil
}

func resourceArmSqlVirtualNetworkRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.VirtualNetworkRulesClient
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

		return fmt.Errorf("Error deleting SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
	}

	return nil
}

/*
	This function checks the format of the SQL Virtual Network Rule Name to make sure that
	it does not contain any potentially invalid values.
*/
func ValidateSqlVirtualNetworkRuleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Cannot be empty
	if len(value) == 0 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be an empty string: %q", k, value))
	}

	// Cannot be shorter than 2 characters
	if len(value) == 1 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be shorter than 2 characters: %q", k, value))
	}

	// Cannot be more than 64 characters
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}

	// Must only contain alphanumeric characters, underscores, periods or hyphens
	if !regexp.MustCompile(`^[A-Za-z0-9-\._]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q can only contain alphanumeric characters, underscores, periods and hyphens: %q",
			k, value))
	}

	// Cannot end in a hyphen
	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a hyphen: %q", k, value))
	}

	// Cannot end in a period
	if regexp.MustCompile(`\.$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot end with a period: %q", k, value))
	}

	// Cannot start with a period, underscore or hyphen
	if regexp.MustCompile(`^[\._-]`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot start with a period, underscore or hyphen: %q", k, value))
	}

	// There are multiple returns in the case that there is more than one invalid
	// case applied to the name.
	return warnings, errors
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
func sqlVirtualNetworkStateStatusCodeRefreshFunc(ctx context.Context, client *sql.VirtualNetworkRulesClient, resourceGroup string, serverName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceGroup, serverName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q) returned 404.", resourceGroup, serverName, name)
				return nil, "ResponseNotFound", nil
			}

			return nil, "", fmt.Errorf("Error polling for the state of the SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q): %+v", name, serverName, resourceGroup, err)
		}

		if props := resp.VirtualNetworkRuleProperties; props != nil {
			log.Printf("[DEBUG] Retrieving SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q) returned Status %s", resourceGroup, serverName, name, props.State)
			return resp, string(props.State), nil
		}

		// Valid response was returned but VirtualNetworkRuleProperties was nil. Basically the rule exists, but with no properties for some reason. Assume Unknown instead of returning error.
		log.Printf("[DEBUG] Retrieving SQL Virtual Network Rule %q (SQL Server: %q, Resource Group: %q) returned empty VirtualNetworkRuleProperties", resourceGroup, serverName, name)
		return resp, "Unknown", nil
	}
}
