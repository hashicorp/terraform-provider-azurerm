package mysql

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMySqlFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceMySqlFirewallRuleCreateUpdate,
		Read:   resourceMySqlFirewallRuleRead,
		Update: resourceMySqlFirewallRuleCreateUpdate,
		Delete: resourceMySqlFirewallRuleDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"start_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azValidate.IPv4Address,
			},

			"end_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azValidate.IPv4Address,
			},
		},
	}
}

func resourceMySqlFirewallRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FirewallRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Firewall Rule creation.")
	startIPAddress := d.Get("start_ip_address").(string)
	endIPAddress := d.Get("end_ip_address").(string)

	id := parse.NewFirewallRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_mysql_firewall_rule", id.ID())
		}
	}

	properties := mysql.FirewallRule{
		FirewallRuleProperties: &mysql.FirewallRuleProperties{
			StartIPAddress: utils.String(startIPAddress),
			EndIPAddress:   utils.String(endIPAddress),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, properties)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %v", id, err)
	}

	d.SetId(id.ID())
	return resourceMySqlFirewallRuleRead(d, meta)
}

func resourceMySqlFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("start_ip_address", resp.StartIPAddress)
	d.Set("end_ip_address", resp.EndIPAddress)

	return nil
}

func resourceMySqlFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
