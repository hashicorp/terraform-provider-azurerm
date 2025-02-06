// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisFirewallRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRedisFirewallRuleCreateUpdate,
		Read:   resourceRedisFirewallRuleRead,
		Update: resourceRedisFirewallRuleCreateUpdate,
		Delete: resourceRedisFirewallRuleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := firewallrules.ParseFirewallRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.FirewallRuleV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallRuleName,
			},

			"redis_cache_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"start_ip": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.IsIPAddress,
					validation.StringIsNotEmpty,
				),
			},

			"end_ip": {
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

func resourceRedisFirewallRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Redis.FirewallRules
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := firewallrules.NewFirewallRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("redis_cache_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_redis_firewall_rule", id.ID())
		}
	}

	payload := firewallrules.RedisFirewallRule{
		Properties: firewallrules.RedisFirewallRuleProperties{
			StartIP: d.Get("start_ip").(string),
			EndIP:   d.Get("end_ip").(string),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceRedisFirewallRuleRead(d, meta)
}

func resourceRedisFirewallRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.FirewallRules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.FirewallRuleName)
	d.Set("redis_cache_name", id.RedisName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("end_ip", model.Properties.EndIP)
		d.Set("start_ip", model.Properties.StartIP)
	}

	return nil
}

func resourceRedisFirewallRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.FirewallRules
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
