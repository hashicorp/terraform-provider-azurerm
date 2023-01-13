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
				check.That(data.ResourceName).Key("recovery_group.0.type").HasValue("Boot"),
				check.That(data.ResourceName).Key("recovery_group.1.type").HasValue("Failover"),
				check.That(data.ResourceName).Key("recovery_group.2.type").HasValue("Shutdown"),
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
