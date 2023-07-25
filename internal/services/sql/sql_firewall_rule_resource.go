// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSqlFirewallRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSqlFirewallRuleCreateUpdate,
		Read:   resourceSqlFirewallRuleRead,
		Update: resourceSqlFirewallRuleCreateUpdate,
		Delete: resourceSqlFirewallRuleDelete,

		DeprecationMessage: "The `azurerm_sql_firewall_rule` resource is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azurerm_mssql_firewall_rule` resource instead.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FirewallRuleID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlServerName,
			},

			"start_ip_address": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.IsIPAddress,
					validation.StringIsNotEmpty,
				),
			},

			"end_ip_address": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.IsIPAddress,
					validation.StringIsNotEmpty,
				),
			},
		},
	}
}

func resourceSqlFirewallRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.FirewallRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFirewallRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_sql_firewall_rule", id.ID())
		}
	}

	parameters := sql.FirewallRule{
		FirewallRuleProperties: &sql.FirewallRuleProperties{
			StartIPAddress: utils.String(d.Get("start_ip_address").(string)),
			EndIPAddress:   utils.String(d.Get("end_ip_address").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSqlFirewallRuleRead(d, meta)
}

func resourceSqlFirewallRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.FirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.FirewallRuleProperties; props != nil {
		d.Set("start_ip_address", props.StartIPAddress)
		d.Set("end_ip_address", props.EndIPAddress)
	}

	return nil
}

func resourceSqlFirewallRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.FirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
