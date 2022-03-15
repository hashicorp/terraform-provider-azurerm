package mariadb

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mariadb/mgmt/2018-06-01/mariadb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmMariaDBFirewallRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmMariaDBFirewallRuleCreateUpdate,
		Read:   resourceArmMariaDBFirewallRuleRead,
		Update: resourceArmMariaDBFirewallRuleCreateUpdate,
		Delete: resourceArmMariaDBFirewallRuleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.MariaDBFirewallRuleID(id)
			return err
		}),

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
				ValidateFunc: validate.FirewallRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"start_ip_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azValidate.IPv4Address,
			},

			"end_ip_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azValidate.IPv4Address,
			},
		},
	}
}

func resourceArmMariaDBFirewallRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.FirewallRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MariaDB Firewall Rule creation.")

	id := parse.NewMariaDBFirewallRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	startIPAddress := d.Get("start_ip_address").(string)
	endIPAddress := d.Get("end_ip_address").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.FirewallRuleName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_mariadb_firewall_rule", id.ID())
		}
	}

	properties := mariadb.FirewallRule{
		FirewallRuleProperties: &mariadb.FirewallRuleProperties{
			StartIPAddress: utils.String(startIPAddress),
			EndIPAddress:   utils.String(endIPAddress),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.FirewallRuleName, properties)
	if err != nil {
		return fmt.Errorf("issuing create/update request for %s: %v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmMariaDBFirewallRuleRead(d, meta)
}

func resourceArmMariaDBFirewallRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.FirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MariaDBFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.FirewallRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.FirewallRuleName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)
	d.Set("start_ip_address", resp.StartIPAddress)
	d.Set("end_ip_address", resp.EndIPAddress)

	return nil
}

func resourceArmMariaDBFirewallRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.FirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MariaDBFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.FirewallRuleName)
	if err != nil {
		return err
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
