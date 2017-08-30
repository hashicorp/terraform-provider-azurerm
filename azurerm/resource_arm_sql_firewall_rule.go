package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlFirewallRuleCreateUpdate,
		Read:   resourceArmSqlFirewallRuleRead,
		Update: resourceArmSqlFirewallRuleCreateUpdate,
		Delete: resourceArmSqlFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"start_ip_address": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: validation?
			},

			"end_ip_address": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: validation?
			},
		},
	}
}

func resourceArmSqlFirewallRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlFirewallRulesClient

	name := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	startIPAddress := d.Get("start_ip_address").(string)
	endIPAddress := d.Get("end_ip_address").(string)

	parameters := sql.FirewallRule{
		FirewallRuleProperties: &sql.FirewallRuleProperties{
			StartIPAddress: utils.String(startIPAddress),
			EndIPAddress:   utils.String(endIPAddress),
		},
	}

	_, err := client.CreateOrUpdate(resourceGroup, serverName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating SQL Firewall Rule: %+v", err)
	}

	resp, err := client.Get(resourceGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving SQL Firewall Rule: %+v", err)
	}

	d.SetId(*resp.ID)

	return resourceArmSqlFirewallRuleRead(d, meta)
}

func resourceArmSqlFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlFirewallRulesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["firewallRules"]

	resp, err := client.Get(resourceGroup, serverName, name)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Firewall Rule %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Firewall Rule: %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)
	d.Set("start_ip_address", resp.StartIPAddress)
	d.Set("end_ip_address", resp.EndIPAddress)

	return nil
}

func resourceArmSqlFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlFirewallRulesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["firewallRules"]

	resp, err := client.Delete(resourceGroup, serverName, name)
	if err != nil {
		if responseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error deleting SQL Firewall Rule: %+v", err)
	}

	return nil
}
