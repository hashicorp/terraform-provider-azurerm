package azurerm

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/postgresql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jen20/riviera/azure"
)

func resourceArmPostgreSQLFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLFirewallRuleCreate,
		Read:   resourceArmPostgreSQLFirewallRuleRead,
		Delete: resourceArmPostgreSQLFirewallRuleDelete,
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
				ForceNew: true,
			},

			"end_ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmPostgreSQLFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlFirewallRulesClient

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Firewall Rule creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)
	startIPAddress := d.Get("start_ip_address").(string)
	endIPAddress := d.Get("end_ip_address").(string)

	properties := postgresql.FirewallRule{
		FirewallRuleProperties: &postgresql.FirewallRuleProperties{
			StartIPAddress: azure.String(startIPAddress),
			EndIPAddress:   azure.String(endIPAddress),
		},
	}

	_, error := client.Create(resGroup, serverName, name, properties, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, serverName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Firewall Rule %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPostgreSQLFirewallRuleRead(d, meta)
}

func resourceArmPostgreSQLFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlFirewallRulesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["firewallRules"]

	resp, err := client.Get(resGroup, serverName, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure PostgreSQL Firewall Rule %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("server_name", serverName)
	d.Set("start_ip_address", resp.StartIPAddress)
	d.Set("end_ip_address", resp.EndIPAddress)

	return nil
}

func resourceArmPostgreSQLFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlFirewallRulesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["firewallRules"]

	_, error := client.Delete(resGroup, serverName, name, make(chan struct{}))
	err = <-error

	return err
}
