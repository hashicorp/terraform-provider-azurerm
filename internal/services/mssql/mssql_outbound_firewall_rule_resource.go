// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlOutboundFirewallRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlOutboundFirewallRuleCreate,
		Read:   resourceMsSqlOutboundFirewallRuleRead,
		// Update: resourceMsSqlOutboundFirewallRuleCreateUpdate,
		Delete: resourceMsSqlOutboundFirewallRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.OutboundFirewallRuleID(id)
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

	id := parse.NewOutboundFirewallRuleID(serverId.SubscriptionId, serverId.ResourceGroup, serverId.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing MSSQL %s: %+v", id.String(), err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_outbound_firewall_rule", id.ID())
		}
	}

	parameters := sql.OutboundFirewallRule{
		OutboundFirewallRuleProperties: &sql.OutboundFirewallRuleProperties{},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating MSSQL %s: %+v", id.String(), err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	return resourceMsSqlOutboundFirewallRuleRead(d, meta)
}

func resourceMsSqlOutboundFirewallRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.OutboundFirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.OutboundFirewallRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] MSSQL %s was not found - removing from state", id.String())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving MSSQL %s: %+v", id.String(), err)
	}

	d.Set("name", id.Name)

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	d.Set("server_id", serverId.ID())

	return nil
}

func resourceMsSqlOutboundFirewallRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.OutboundFirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.OutboundFirewallRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting MSSQL %s: %+v", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for delete of %s: %+v", id.String(), err)
	}

	return nil
}
