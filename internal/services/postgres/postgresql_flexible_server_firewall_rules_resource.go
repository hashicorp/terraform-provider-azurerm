// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

type FirewallRuleWithId struct {
	Id           firewallrules.FirewallRuleId
	FirewallRule firewallrules.FirewallRule
}

func resourcePostgresqlFlexibleServerFirewallRules() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgresqlFlexibleServerFirewallRulesCreateUpdate,
		Read:   resourcePostgresqlFlexibleServerFirewallRulesRead,
		Update: resourcePostgresqlFlexibleServerFirewallRulesCreateUpdate,
		Delete: resourcePostgresqlFlexibleServerFirewallRulesDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := firewallrules.ParseFirewallRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: firewallrules.ValidateFlexibleServerID,
			},

			"firewall_rule": {
				Type:       pluginsdk.TypeSet,
				Required:   true,
				ConfigMode: pluginsdk.SchemaConfigModeBlock,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.FlexibleServerFirewallRuleName,
						},

						"end_ip_address": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},

						"start_ip_address": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
					},
				},
				Set: hashPostgresqlFlexibleServerFirewallRule,
			},
		},
	}
}

func resourcePostgresqlFlexibleServerFirewallRulesCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	firewall_rules_client := meta.(*clients.Client).Postgres.FlexibleServerFirewallRuleClient
	flexible_servers_client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

	resp, err := flexible_servers_client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Postgresql Flexibleserver %q does not exist", d.Id())
			return err
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	flexibleServerId := firewallrules.NewFlexibleServerID(subscriptionId, id.ResourceGroupName, id.FlexibleServerName)

	listFirewallRulesResult, err := firewall_rules_client.ListByServerComplete(ctx, flexibleServerId)
	if err != nil {
		return err
	}

	currentFirewallRules := listFirewallRulesResult.Items

	// Build a list of what the firewall rules should look like
	correctFirewallRules := make([]FirewallRuleWithId, 0)

	for _, address := range d.Get("firewall_rule").(*pluginsdk.Set).List() {
		addressMap := address.(map[string]interface{})
		fwRule := firewallrules.FirewallRule{
			Properties: firewallrules.FirewallRuleProperties{
				EndIPAddress:   addressMap["end_ip_address"].(string),
				StartIPAddress: addressMap["start_ip_address"].(string),
			},
		}
		fwRuleId := firewallrules.NewFirewallRuleID(subscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName, addressMap["name"].(string))
		correctFirewallRules = append(correctFirewallRules, FirewallRuleWithId{Id: fwRuleId, FirewallRule: fwRule})
	}

	rulesToDelete := make([]firewallrules.FirewallRuleId, 0)
	rulesToCreateUpdate := make([]FirewallRuleWithId, 0)

	// Iterate through the current firewall rules and compare them to the desired state
	// Any firewall rules that do not appear in the desired state should be deleted
	for _, currentRule := range currentFirewallRules {
		found := false
		for _, correctRule := range correctFirewallRules {
			if *currentRule.Name == *&correctRule.Id.FirewallRuleName {
				found = true
				break
			}
		}
		if !found {
			rulesToDelete = append(rulesToDelete, firewallrules.NewFirewallRuleID(subscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName, *currentRule.Name))
		}
	}

	// Iterate through the desired firewall rules and compare them to the current state
	// Any firewall rules that do not appear in the current state should be created
	// Any firewall rules that do appear in the current state but have different properties should be updated
	for _, correctRule := range correctFirewallRules {
		found := false
		for _, currentRule := range currentFirewallRules {
			if *currentRule.Name == *&correctRule.Id.FirewallRuleName {
				found = true
				if currentRule.Properties.StartIPAddress != correctRule.FirewallRule.Properties.StartIPAddress || currentRule.Properties.EndIPAddress != correctRule.FirewallRule.Properties.EndIPAddress {
					rulesToCreateUpdate = append(rulesToCreateUpdate, correctRule)
				}
				break
			}
		}
		if !found {
			rulesToCreateUpdate = append(rulesToCreateUpdate, correctRule)
		}
	}

	// The lists of rules to Create/Update/Delete have been built
	pollers := make([]pollers.Poller, 0)

	for _, rule := range rulesToDelete {
		poller, err := firewall_rules_client.Delete(ctx, rule)
		if err != nil {
			return fmt.Errorf("deleting %q: %+v", rule, err)
		}
		pollers = append(pollers, poller.Poller)
	}

	for _, rule := range rulesToCreateUpdate {
		poller, err := firewall_rules_client.CreateOrUpdate(ctx, rule.Id, rule.FirewallRule)
		if err != nil {
			return fmt.Errorf("creating/updating %q: %+v", rule.Id, err)
		}
		pollers = append(pollers, poller.Poller)
	}

	wg := sync.WaitGroup{}

	for _, poller := range pollers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := poller.PollUntilDone(ctx); err != nil {
				fmt.Errorf("polling after CreateOrUpdate: %+v", err)
			}
		}()

	}
	wg.Wait()

	d.SetId(id.ID())
	return resourcePostgresqlFlexibleServerFirewallRulesRead(d, meta)
}

func resourcePostgresqlFlexibleServerFirewallRulesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	firewall_rules_client := meta.(*clients.Client).Postgres.FlexibleServerFirewallRuleClient
	flexible_servers_client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := flexible_servers_client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Postgresql Flexibleserver %q does not exist - removing firewall rules from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	flexibleServerId := firewallrules.NewFlexibleServerID(subscriptionId, id.ResourceGroupName, id.FlexibleServerName)

	d.Set("server_id", id.ID())
	addresses := make([]map[string]interface{}, 0)
	fwRules, err := firewall_rules_client.ListByServerComplete(ctx, flexibleServerId)
	for _, rule := range fwRules.Items {
		addresses = append(addresses, map[string]interface{}{
			"name":             rule.Name,
			"end_ip_address":   rule.Properties.EndIPAddress,
			"start_ip_address": rule.Properties.StartIPAddress,
		})
	}
	d.Set("firewall_rule", addresses)
	return nil
}

func resourcePostgresqlFlexibleServerFirewallRulesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	firewall_rules_client := meta.(*clients.Client).Postgres.FlexibleServerFirewallRuleClient
	flexible_servers_client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

	resp, err := flexible_servers_client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Postgresql Flexibleserver %q does not exist", d.Id())
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	flexibleServerId := firewallrules.NewFlexibleServerID(subscriptionId, id.ResourceGroupName, id.FlexibleServerName)

	listFirewallRulesResult, err := firewall_rules_client.ListByServerComplete(ctx, flexibleServerId)
	if err != nil {
		return err
	}

	pollers := make([]pollers.Poller, 0)

	for _, rule := range listFirewallRulesResult.Items {
		poller, err := firewall_rules_client.Delete(ctx, firewallrules.NewFirewallRuleID(subscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName, *rule.Name))
		if err != nil {
			return fmt.Errorf("deleting %q: %+v", rule.Name, err)
		}
		pollers = append(pollers, poller.Poller)
	}

	wg := sync.WaitGroup{}

	for _, poller := range pollers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := poller.PollUntilDone(ctx); err != nil {
				fmt.Errorf("polling after Delete: %+v", err)
			}
		}()
	}

	wg.Wait()
	return nil
}

func hashPostgresqlFlexibleServerFirewallRule(v interface{}) int {
	var buf bytes.Buffer
	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["name"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(v.(string))))
		}
		if v, ok := m["start_ip_address"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["end_ip_address"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}
	return pluginsdk.HashString(buf.String())
}
