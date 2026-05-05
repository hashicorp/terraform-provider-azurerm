// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/firewallrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name mysql_flexible_server_firewall_rule -service-package-name mysql -properties "name,resource_group_name,flexible_server_name:server_name" -known-values "subscription_id:data.Subscriptions.Primary"

var mysqlFlexibleServerFirewallResourceName = "azurerm_mysql_flexible_server_firewall_rule"

func resourceMySqlFlexibleServerFirewallRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySqlFlexibleServerFirewallRuleCreateUpdate,
		Read:   resourceMySqlFlexibleServerFirewallRuleRead,
		Update: resourceMySqlFlexibleServerFirewallRuleCreateUpdate,
		Delete: resourceMySqlFlexibleServerFirewallRuleDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&firewallrules.FirewallRuleId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&firewallrules.FirewallRuleId{}),
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
				ValidateFunc: validate.FlexibleServerName,
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

func resourceMySqlFlexibleServerFirewallRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.FirewallRules
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := firewallrules.NewFirewallRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(mysqlFlexibleServerFirewallResourceName, id.ID())
		}
	}

	payload := firewallrules.FirewallRule{
		Properties: firewallrules.FirewallRuleProperties{
			StartIPAddress: d.Get("start_ip_address").(string),
			EndIPAddress:   d.Get("end_ip_address").(string),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceMySqlFlexibleServerFirewallRuleRead(d, meta)
}

func resourceMySqlFlexibleServerFirewallRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.FirewallRules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceMySqlFlexibleServerFirewallRuleFlatten(d, id, resp.Model)
}

func resourceMySqlFlexibleServerFirewallRuleFlatten(d *pluginsdk.ResourceData, id *firewallrules.FirewallRuleId, rule *firewallrules.FirewallRule) error {
	d.Set("name", id.FirewallRuleName)
	d.Set("server_name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if rule != nil {
		d.Set("start_ip_address", rule.Properties.StartIPAddress)
		d.Set("end_ip_address", rule.Properties.EndIPAddress)
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceMySqlFlexibleServerFirewallRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.FirewallRules
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
