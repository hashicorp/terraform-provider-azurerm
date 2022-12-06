package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SiteRecoveryReplicationRecoveryPlanDataSource struct{}

func TestAccDataSourceSiteRecoveryReplicationRecoverPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlanDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, 24*60, 4*60),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_name").Exists(),
				check.That(data.ResourceName).Key("source_recovery_fabric_id").Exists(),
				check.That(data.ResourceName).Key("target_recovery_fabric_id").Exists(),
				check.That(data.ResourceName).Key("recovery_group.0.group_type").HasValue("Boot"),
				check.That(data.ResourceName).Key("recovery_group.1.group_type").HasValue("Failover"),
				check.That(data.ResourceName).Key("recovery_group.2.group_type").HasValue("Shutdown"),
			),
		},
	})
}

func (SiteRecoveryReplicationRecoveryPlanDataSource) basic(data acceptance.TestData, retentionInMinutes int, snapshotFrequencyInMinutes int) string {
	return fmt.Sprintf(`
%s

data "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                = azurerm_site_recovery_replication_policy.test.name
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
}
`, SiteRecoveryReplicationRecoveryPlan{}.basic(data))
}
