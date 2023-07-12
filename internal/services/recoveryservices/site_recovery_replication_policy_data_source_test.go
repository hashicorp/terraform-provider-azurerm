// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SiteRecoveryReplicationPolicyDataSource struct{}

func TestAccDataSourceSiteRecoveryReplicationPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_policy", "test")
	r := SiteRecoveryReplicationPolicyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, 24*60, 4*60),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_name").Exists(),
				check.That(data.ResourceName).Key("recovery_point_retention_in_minutes").HasValue("1440"),
				check.That(data.ResourceName).Key("application_consistent_snapshot_frequency_in_minutes").HasValue("240"),
			),
		},
	})
}

func (SiteRecoveryReplicationPolicyDataSource) basic(data acceptance.TestData, retentionInMinutes int, snapshotFrequencyInMinutes int) string {
	return fmt.Sprintf(`
%s

data "azurerm_site_recovery_replication_policy" "test" {
  name                = azurerm_site_recovery_replication_policy.test.name
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
}
`, SiteRecoveryReplicationPolicyResource{}.basic(data, retentionInMinutes, snapshotFrequencyInMinutes))
}
