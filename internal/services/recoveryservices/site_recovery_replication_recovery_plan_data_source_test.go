// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SiteRecoveryReplicationRecoveryPlanDataSource struct{}

func TestAccDataSourceSiteRecoveryReplicationRecoveryPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlanDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_id").Exists(),
				check.That(data.ResourceName).Key("source_recovery_fabric_id").Exists(),
				check.That(data.ResourceName).Key("target_recovery_fabric_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceSiteRecoveryReplicationRecoveryPlan_withZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlanDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_id").Exists(),
				check.That(data.ResourceName).Key("source_recovery_fabric_id").Exists(),
				check.That(data.ResourceName).Key("target_recovery_fabric_id").Exists(),
				check.That(data.ResourceName).Key("azure_to_azure_settings.0.primary_zone").HasValue("1"),
				check.That(data.ResourceName).Key("azure_to_azure_settings.0.recovery_zone").HasValue("2"),
			),
		},
	})
}

func (SiteRecoveryReplicationRecoveryPlanDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_site_recovery_replication_recovery_plan" "test" {
  name              = azurerm_site_recovery_replication_recovery_plan.test.name
  recovery_vault_id = azurerm_site_recovery_replication_recovery_plan.test.recovery_vault_id
}
`, SiteRecoveryReplicationRecoveryPlan{}.basic(data))
}

func (SiteRecoveryReplicationRecoveryPlanDataSource) withZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_site_recovery_replication_recovery_plan" "test" {
  name              = azurerm_site_recovery_replication_recovery_plan.test.name
  recovery_vault_id = azurerm_site_recovery_replication_recovery_plan.test.recovery_vault_id
}
`, SiteRecoveryReplicationRecoveryPlan{}.withZones(data))
}
