// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/outboundfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlOutboundFirewallRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlOutboundFirewallRuleCreate,
		Read:   resourceMsSqlOutboundFirewallRuleRead,
		Delete: resourceMsSqlOutboundFirewallRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.OutboundFirewallRuleID(id)
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
			},

			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerID,
			},
		},
	}
}

func resourceMsSqlOutboundFirewallRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.OutboundFirewallRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := parse.ServerID(d.Get("server_id").(string))
	if err != nil {
		return fmt.Errorf("parsing server ID %q: %+v", d.Get("server_id"), err)
	}

	id := outboundfirewallrules.NewOutboundFirewallRuleID(serverId.SubscriptionId, serverId.ResourceGroup, serverId.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_mssql_outbound_firewall_rule", id.ID())
		}
	}

	err = client.CreateOrUpdateThenPoll(ctx, id)
	if err != nil {
		return fmt.Errorf("creating MSSQL %s: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	return resourceMsSqlOutboundFirewallRuleRead(d, meta)
}

func resourceMsSqlOutboundFirewallRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.OutboundFirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outboundfirewallrules.ParseOutboundFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] MSSQL %s was not found - removing from state", id.String())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving MSSQL %s: %+v", id.String(), err)
	}

	d.Set("name", id.OutboundFirewallRuleName)

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
	d.Set("server_id", serverId.ID())

	return nil
}

func resourceMsSqlOutboundFirewallRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.OutboundFirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outboundfirewallrules.ParseOutboundFirewallRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting MSSQL %s: %+v", id.String(), err)
	}

	return nil
}
