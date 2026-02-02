// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisfirewallrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name redis_firewall_rule -service-package-name redis -properties "name,resource_group_name,redis_name:redis_cache_name" -known-values "subscription_id:data.Subscriptions.Primary"

func resourceRedisFirewallRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceRedisFirewallRuleCreateUpdate,
		Read:     resourceRedisFirewallRuleRead,
		Update:   resourceRedisFirewallRuleCreateUpdate,
		Delete:   resourceRedisFirewallRuleDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&redisfirewallrules.FirewallRuleId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&redisfirewallrules.FirewallRuleId{}),
		},

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
	client := meta.(*clients.Client).Redis.FirewallRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := redisfirewallrules.NewFirewallRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("redis_cache_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.FirewallRulesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_redis_firewall_rule", id.ID())
		}
	}

	payload := redisfirewallrules.RedisFirewallRule{
		Properties: redisfirewallrules.RedisFirewallRuleProperties{
			StartIP: d.Get("start_ip").(string),
			EndIP:   d.Get("end_ip").(string),
		},
	}

	if _, err := client.FirewallRulesCreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceRedisFirewallRuleRead(d, meta)
}

func resourceRedisFirewallRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.FirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redisfirewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.FirewallRulesGet(ctx, *id)
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

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceRedisFirewallRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.FirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redisfirewallrules.ParseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.FirewallRulesDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
