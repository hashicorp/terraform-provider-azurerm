// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccPostgresqlFlexibleServerFirewallRule_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_firewall_rule", "test")
	r := PostgresqlFlexibleServerFirewallRuleResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_postgresql_flexible_server_firewall_rule.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_postgresql_flexible_server_firewall_rule.test", tfjsonpath.New("flexible_server_name"), tfjsonpath.New("server_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_postgresql_flexible_server_firewall_rule.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("server_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_postgresql_flexible_server_firewall_rule.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("server_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
