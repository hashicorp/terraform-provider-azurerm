package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/sql"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmSqlFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlFirewallRuleCreateOrUpdate,
		Read:   resourceArmSqlFirewallRuleRead,
		Update: resourceArmSqlFirewallRuleCreateOrUpdate,
		Delete: resourceArmSqlFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"server_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"start_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"end_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmSqlFirewallRuleCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	sqlFirewallRulesClient := meta.(*ArmClient).sqlFirewallRulesClient

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)
	startIpAddress := d.Get("start_ip_address").(string)
	endIpAddress := d.Get("end_ip_address").(string)

	props := sql.FirewallRuleProperties{
		StartIPAddress: &startIpAddress,
		EndIPAddress:   &endIpAddress,
	}

	parameters := sql.FirewallRule{
		Name: &name,
		FirewallRuleProperties: &props,
	}

	//last parameter is set to empty to allow updates to records after creation
	// (per SDK, set it to '*' to prevent updates, all other values are ignored)
	result, err := sqlFirewallRulesClient.CreateOrUpdate(resGroup, serverName, name, parameters)
	if err != nil {
		return err
	}

	if result.ID == nil {
		return fmt.Errorf("Cannot create sql firewall rule %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*result.ID)

	return resourceArmSqlFirewallRuleRead(d, meta)
}

func resourceArmSqlFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	sqlFirewallRulesClient := meta.(*ArmClient).sqlFirewallRulesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["firewallRules"]
	serverName := id.Path["servers"]

	result, err := sqlFirewallRulesClient.Get(resGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error reading DNS PTR record %s: %v", name, err)
	}
	if result.Response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	props := *result.FirewallRuleProperties

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("server_name", serverName)
	d.Set("start_ip_address", *props.StartIPAddress)
	d.Set("end_ip_address", *props.EndIPAddress)

	return nil
}

func resourceArmSqlFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	sqlFirewallRulesClient := meta.(*ArmClient).sqlFirewallRulesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["firewallRules"]
	serverName := id.Path["servers"]

	result, error := sqlFirewallRulesClient.Delete(resGroup, serverName, name)
	if result.Response.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting sql firewall rule %s: %s", name, error)
	}

	return nil
}
