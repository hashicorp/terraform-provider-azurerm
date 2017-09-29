package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/mysql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMySqlFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMySqlFirewallRuleCreateUpdate,
		Read:   resourceArmMySqlFirewallRuleRead,
		Update: resourceArmMySqlFirewallRuleCreateUpdate,
		Delete: resourceArmMySqlFirewallRuleDelete,
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

			"start_ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},

			"end_ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmMySqlFirewallRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mysqlFirewallRulesClient

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Firewall Rule creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)
	startIPAddress := d.Get("start_ip_address").(string)
	endIPAddress := d.Get("end_ip_address").(string)

	properties := mysql.FirewallRule{
		FirewallRuleProperties: &mysql.FirewallRuleProperties{
			StartIPAddress: utils.String(startIPAddress),
			EndIPAddress:   utils.String(endIPAddress),
		},
	}

	_, createErr := client.CreateOrUpdate(resGroup, serverName, name, properties, make(chan struct{}))
	err := <-createErr
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, serverName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read MySQL Firewall Rule %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMySqlFirewallRuleRead(d, meta)
}

func resourceArmMySqlFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mysqlFirewallRulesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["firewallRules"]

	resp, err := client.Get(resGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure MySQL Firewall Rule %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("server_name", serverName)
	d.Set("start_ip_address", resp.StartIPAddress)
	d.Set("end_ip_address", resp.EndIPAddress)

	return nil
}

func resourceArmMySqlFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mysqlFirewallRulesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["firewallRules"]

	_, deleteErr := client.Delete(resGroup, serverName, name, make(chan struct{}))
	err = <-deleteErr

	return err
}
