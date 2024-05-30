// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FirewallRuleWithId struct {
	Id           firewallrules.FirewallRuleId `tfschema:"id"`
	FirewallRule firewallrules.FirewallRule   `tfschema:"firewall_rule"`
}

type NamedFirewallRule struct {
	Name           string `tfschema:"name"`
	StartIPAddress string `tfschema:"start_ip_address"`
	EndIPAddress   string `tfschema:"end_ip_address"`
}

type FlexibleServerFirewallRulesModel struct {
	ServerID     string              `tfschema:"server_id"`
	FirewallRule []NamedFirewallRule `tfschema:"firewall_rule"`
}

var (
	_ sdk.Resource           = FlexibleServerFirewallRulesResource{}
	_ sdk.ResourceWithUpdate = FlexibleServerFirewallRulesResource{}
)

type FlexibleServerFirewallRulesResource struct{}

func (r FlexibleServerFirewallRulesResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
	}
}

func (r FlexibleServerFirewallRulesResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r FlexibleServerFirewallRulesResource) ResourceType() string {
	return "azurerm_postgresql_flexible_server_firewall_rules"
}

func (r FlexibleServerFirewallRulesResource) ModelObject() interface{} {
	return &FlexibleServerFirewallRulesModel{}
}

func (r FlexibleServerFirewallRulesResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewallrules.ValidateFirewallRuleID
}

func (r FlexibleServerFirewallRulesResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			firewall_rules_client := metadata.Client.Postgres.FlexibleServerFirewallRuleClient
			flexible_servers_client := metadata.Client.Postgres.FlexibleServersClient

			m := FlexibleServerFirewallRulesModel{}
			if err := metadata.Decode(&m); err != nil {
				return err
			}

			id, err := servers.ParseFlexibleServerID(m.ServerID)
			if err != nil {
				return err
			}

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			resp, err := flexible_servers_client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[INFO] Postgresql Flexibleserver %q does not exist", m.ServerID)
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

			for _, address := range m.FirewallRule {
				fwRule := firewallrules.FirewallRule{
					Properties: firewallrules.FirewallRuleProperties{
						EndIPAddress:   address.EndIPAddress,
						StartIPAddress: address.StartIPAddress,
					},
				}
				fwRuleId := firewallrules.NewFirewallRuleID(subscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName, address.Name)
				correctFirewallRules = append(correctFirewallRules, FirewallRuleWithId{Id: fwRuleId, FirewallRule: fwRule})
			}

			// Iterate through the current firewall rules and compare them to the desired state
			// Any firewall rules that do not appear in the desired state should be deleted
			rulesToDelete := make([]firewallrules.FirewallRuleId, 0)
			for _, currentRule := range currentFirewallRules {
				found := false
				for _, correctRule := range correctFirewallRules {
					// Initially match on name, then check the IPs
					if *currentRule.Name == correctRule.Id.FirewallRuleName {
						if currentRule.Properties.StartIPAddress == correctRule.FirewallRule.Properties.StartIPAddress && currentRule.Properties.EndIPAddress == correctRule.FirewallRule.Properties.EndIPAddress {
							found = true
							break
						}
					}
				}
				if !found {
					rulesToDelete = append(rulesToDelete, firewallrules.NewFirewallRuleID(subscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName, *currentRule.Name))
				}
			}

			// Iterate through the desired firewall rules and compare them to the current state
			// Any firewall rules that do not appear in the current state should be created
			// Any firewall rules that do appear in the current state but have different properties should be updated
			rulesToCreateUpdate := make([]FirewallRuleWithId, 0)
			for _, correctRule := range correctFirewallRules {
				found := false
				for _, currentRule := range currentFirewallRules {
					if *currentRule.Name == correctRule.Id.FirewallRuleName {
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
						fmt.Errorf("polling after change: %+v", err)
					}
				}()

			}
			wg.Wait()

			metadata.SetID(id)
			return nil
		},
	}
}

func (r FlexibleServerFirewallRulesResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			firewall_rules_client := metadata.Client.Postgres.FlexibleServerFirewallRuleClient
			flexible_servers_client := metadata.Client.Postgres.FlexibleServersClient

			id, err := servers.ParseFlexibleServerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := flexible_servers_client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			flexibleServerId := firewallrules.NewFlexibleServerID(subscriptionId, id.ResourceGroupName, id.FlexibleServerName)
			fwRules, err := firewall_rules_client.ListByServerComplete(ctx, flexibleServerId)
			currentRulesData := make([]NamedFirewallRule, 0)
			for _, rule := range fwRules.Items {
				currentRulesData = append(currentRulesData, NamedFirewallRule{
					Name:           *rule.Name,
					StartIPAddress: rule.Properties.StartIPAddress,
					EndIPAddress:   rule.Properties.EndIPAddress,
				})
			}
			m := FlexibleServerFirewallRulesModel{
				ServerID:     id.ID(),
				FirewallRule: currentRulesData,
			}
			return metadata.Encode(&m)
		},
	}
}

func (r FlexibleServerFirewallRulesResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			firewall_rules_client := metadata.Client.Postgres.FlexibleServerFirewallRuleClient
			flexible_servers_client := metadata.Client.Postgres.FlexibleServersClient

			m := FlexibleServerFirewallRulesModel{}
			if err := metadata.Decode(&m); err != nil {
				return err
			}

			id, err := servers.ParseFlexibleServerID(m.ServerID)
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("firewall_rule") {

				locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
				defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

				resp, err := flexible_servers_client.Get(ctx, *id)
				if err != nil {
					if response.WasNotFound(resp.HttpResponse) {
						log.Printf("[INFO] Postgresql Flexibleserver %q does not exist", m.ServerID)
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

				for _, address := range m.FirewallRule {
					fwRule := firewallrules.FirewallRule{
						Properties: firewallrules.FirewallRuleProperties{
							EndIPAddress:   address.EndIPAddress,
							StartIPAddress: address.StartIPAddress,
						},
					}
					fwRuleId := firewallrules.NewFirewallRuleID(subscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName, address.Name)
					correctFirewallRules = append(correctFirewallRules, FirewallRuleWithId{Id: fwRuleId, FirewallRule: fwRule})
				}

				// Iterate through the current firewall rules and compare them to the desired state
				// Any firewall rules that do not appear in the desired state should be deleted
				rulesToDelete := make([]firewallrules.FirewallRuleId, 0)
				for _, currentRule := range currentFirewallRules {
					found := false
					for _, correctRule := range correctFirewallRules {
						// Initially match on name, then check the IPs
						if *currentRule.Name == correctRule.Id.FirewallRuleName {
							if currentRule.Properties.StartIPAddress == correctRule.FirewallRule.Properties.StartIPAddress && currentRule.Properties.EndIPAddress == correctRule.FirewallRule.Properties.EndIPAddress {
								found = true
								break
							}
						}
					}
					if !found {
						rulesToDelete = append(rulesToDelete, firewallrules.NewFirewallRuleID(subscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName, *currentRule.Name))
					}
				}

				// Iterate through the desired firewall rules and compare them to the current state
				// Any firewall rules that do not appear in the current state should be created
				// Any firewall rules that do appear in the current state but have different properties should be updated
				rulesToCreateUpdate := make([]FirewallRuleWithId, 0)
				for _, correctRule := range correctFirewallRules {
					found := false
					for _, currentRule := range currentFirewallRules {
						if *currentRule.Name == correctRule.Id.FirewallRuleName {
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
							fmt.Errorf("polling after change: %+v", err)
						}
					}()

				}
				wg.Wait()
			}
			return nil
		},
	}
}

func (r FlexibleServerFirewallRulesResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			firewall_rules_client := metadata.Client.Postgres.FlexibleServerFirewallRuleClient
			flexible_servers_client := metadata.Client.Postgres.FlexibleServersClient

			id, err := servers.ParseFlexibleServerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			resp, err := flexible_servers_client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[INFO] Postgresql Flexibleserver %q does not exist", metadata.ResourceData.Id())
					return err
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
					return fmt.Errorf("deleting %q: %+v", *rule.Name, err)
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
		},
	}
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
