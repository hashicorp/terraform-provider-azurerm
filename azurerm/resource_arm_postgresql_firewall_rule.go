package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPostgreSQLFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLFirewallRuleCreate,
		Read:   resourceArmPostgreSQLFirewallRuleRead,
		Delete: resourceArmPostgreSQLFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
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
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Firewall Rule creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	// first check if there's one in this subscription requiring import
	resp, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error checking for the existence of Firewall Rule %q (PostgreSQL Server %q / Resource Group %q): %+v", name, serverName, resGroup, err)
		}
	}

	if resp.ID != nil {
		return tf.ImportAsExistsError("azurerm_postgresql_firewall_rule", *resp.ID)
	}

	startIPAddress := d.Get("start_ip_address").(string)
	endIPAddress := d.Get("end_ip_address").(string)

	properties := postgresql.FirewallRule{
		FirewallRuleProperties: &postgresql.FirewallRuleProperties{
			StartIPAddress: utils.String(startIPAddress),
			EndIPAddress:   utils.String(endIPAddress),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, name, properties)
	if err != nil {
		return err
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, serverName, name)
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
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["firewallRules"]

	resp, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL Firewall Rule '%s' was not found (resource group '%s')", name, resGroup)
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
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["firewallRules"]

	future, err := client.Delete(ctx, resGroup, serverName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return nil
}
