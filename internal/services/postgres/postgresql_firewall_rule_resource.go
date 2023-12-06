// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePostgreSQLFirewallRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgreSQLFirewallRuleCreate,
		Read:   resourcePostgreSQLFirewallRuleRead,
		Delete: resourcePostgreSQLFirewallRuleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := firewallrules.ParseFirewallRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Rule name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"start_ip_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"end_ip_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
		},
	}
}

func resourcePostgreSQLFirewallRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FirewallRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := firewallrules.NewFirewallRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_postgresql_firewall_rule", id.ID())
	}

	properties := firewallrules.FirewallRule{
		Properties: firewallrules.FirewallRuleProperties{
			StartIPAddress: d.Get("start_ip_address").(string),
			EndIPAddress:   d.Get("end_ip_address").(string),
		},
	}

	if err = client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePostgreSQLFirewallRuleRead(d, meta)
}

func resourcePostgreSQLFirewallRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.FirewallRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("server_name", id.ServerName)

	if model := resp.Model; model != nil {
		d.Set("start_ip_address", resp.Model.Properties.StartIPAddress)
		d.Set("end_ip_address", resp.Model.Properties.EndIPAddress)
	}

	return nil
}

func resourcePostgreSQLFirewallRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
