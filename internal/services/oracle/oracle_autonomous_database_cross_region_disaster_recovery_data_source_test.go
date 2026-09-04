// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type AutonomousDatabaseCrossRegionDisasterRecoveryDataSource struct{}

func TestAccAutonomousDatabaseCrossRegionDisasterRecoveryDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, fmt.Sprintf("data.%s", oracle.AutonomousDatabaseCrossRegionDisasterRecoveryDataSource{}.ResourceType()), "adbs_secondary_crdr")
	r := AutonomousDatabaseCrossRegionDisasterRecoveryDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("auto_scaling_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("auto_scaling_for_storage_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("backup_retention_period_in_days").HasValue("12"),
				check.That(data.ResourceName).Key("character_set").HasValue("AL32UTF8"),
				check.That(data.ResourceName).Key("compute_count").HasValue("2"),
				check.That(data.ResourceName).Key("customer_contacts.#").HasValue("1"),
				check.That(data.ResourceName).Key("customer_contacts.0").HasValue("test@test.com"),
				check.That(data.ResourceName).Key("data_storage_size_in_tb").HasValue("1"),
				check.That(data.ResourceName).Key("database_version").HasValue("19c"),
				check.That(data.ResourceName).Key("database_workload").HasValue("DW"),
				check.That(data.ResourceName).Key("license_model").HasValue("LicenseIncluded"),
				check.That(data.ResourceName).Key("local_adg_auto_failover_max_data_loss_limit_in_seconds").Exists(),
				check.That(data.ResourceName).Key("national_character_set").HasValue("AL16UTF16"),
				check.That(data.ResourceName).Key("remote_disaster_recovery_type").HasValue("Adg"),
				check.That(data.ResourceName).Key("source_autonomous_database_id").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("virtual_network_id").Exists(),
			),
		},
	})
}

func (d AutonomousDatabaseCrossRegionDisasterRecoveryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "adbs_secondary_crdr" {
  name                = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.name
  resource_group_name = azurerm_oracle_autonomous_database_cross_region_disaster_recovery.adbs_secondary_crdr.resource_group_name
}
`, AdbsCrossRegionDisasterRecoveryResource{}.complete(data))
}
