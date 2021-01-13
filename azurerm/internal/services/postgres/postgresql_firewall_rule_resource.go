package postgres

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2020-01-01/postgresql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePostgreSQLFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgreSQLFirewallRuleCreate,
		Read:   resourcePostgreSQLFirewallRuleRead,
		Delete: resourcePostgreSQLFirewallRuleDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.FirewallRuleID(id)
			return err
		}),

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
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Rule name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
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
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"end_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
		},
	}
}

func resourcePostgreSQLFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FirewallRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Firewall Rule creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)
	startIPAddress := d.Get("start_ip_address").(string)
	endIPAddress := d.Get("end_ip_address").(string)

	existing, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing PostgreSQL Firewall Rule %s (resource group %s) ID", name, resGroup)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_postgresql_firewall_rule", *existing.ID)
	}

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

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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

	return resourcePostgreSQLFirewallRuleRead(d, meta)
}

func resourcePostgreSQLFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL Firewall Rule %q was not found (Resource Group %q)", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving PostgreSQL Firewall Rule %q: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)

	if props := resp.FirewallRuleProperties; props != nil {
		d.Set("start_ip_address", props.StartIPAddress)
		d.Set("end_ip_address", props.EndIPAddress)
	}

	return nil
}

func resourcePostgreSQLFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return nil
}
