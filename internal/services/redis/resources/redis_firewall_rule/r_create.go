// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_firewall_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	redisClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisFirewallRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := redisClient.From(meta).FirewallRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := redisfirewallrules.NewFirewallRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("redis_cache_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
	}

	payload := redisfirewallrules.RedisFirewallRule{
		Properties: redisfirewallrules.RedisFirewallRuleProperties{
			StartIP: d.Get("start_ip").(string),
			EndIP:   d.Get("end_ip").(string),
		},
	}

	if _, err := client.FirewallRulesCreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	}

	return resourceRedisFirewallRuleRead(d, meta)
}
