// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BackupProtectionPolicyVMWorkloadResource struct{}

func TestAccBackupProtectionPolicyVMWorkload_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm_workload", "test")
	r := BackupProtectionPolicyVMWorkloadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVMWorkload_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm_workload", "test")
	r := BackupProtectionPolicyVMWorkloadResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t BackupProtectionPolicyVMWorkloadResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := protectionpolicies.ParseBackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ProtectionPoliciesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service Protection Policy (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r BackupProtectionPolicyVMWorkloadResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bpvmw-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-rsv-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm_workload" "test" {
  name                = "acctest-bpvmw-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  workload_type = "SQLDataBase"

  settings {
    time_zone           = "UTC"
    compression_enabled = false
  }

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Daily"
      time      = "15:00"
    }

    retention_daily {
      count = 8
    }

    retention_monthly {
      format_type = "Daily"
      count       = 10
      monthdays   = [27, 28]
    }

    retention_yearly {
      format_type = "Daily"
      count       = 10
      months      = ["February"]
      monthdays   = [27, 28]
    }
  }

  protection_policy {
    policy_type = "Log"

    backup {
      frequency_in_minutes = 15
    }

    simple_retention {
      count = 8
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r BackupProtectionPolicyVMWorkloadResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bpvmw-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-rsv-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm_workload" "test" {
  name                = "acctest-bpvmw-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  workload_type = "SAPHanaDatabase"

  settings {
    time_zone           = "UTC"
    compression_enabled = true
  }

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Weekly"
      time      = "15:00"
      weekdays  = ["Monday", "Tuesday"]
    }

    retention_weekly {
      weekdays = ["Monday", "Tuesday"]
      count    = 4
    }

    retention_monthly {
      format_type = "Weekly"
      weeks       = ["Third"]
      weekdays    = ["Tuesday"]
      count       = 10
    }

    retention_yearly {
      format_type = "Weekly"
      months      = ["May", "February"]
      weeks       = ["Third"]
      weekdays    = ["Tuesday"]
      count       = 8
    }
  }

  protection_policy {
    policy_type = "Incremental"

    backup {
      frequency = "Weekly"
      weekdays  = ["Saturday", "Friday"]
      time      = "23:00"
    }

    simple_retention {
      count = 11
    }
  }

  protection_policy {
    policy_type = "Log"

    backup {
      frequency_in_minutes = 15
    }

    simple_retention {
      count = 8
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r BackupProtectionPolicyVMWorkloadResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-bpvmw-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-rsv-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  soft_delete_enabled = false
}

resource "azurerm_backup_policy_vm_workload" "test" {
  name                = "acctest-bpvmw-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  workload_type = "SAPHanaDatabase"

  settings {
    time_zone           = "Pacific Standard Time"
    compression_enabled = false
  }

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Weekly"
      time      = "16:00"
      weekdays  = ["Tuesday", "Thursday"]
    }

    retention_weekly {
      weekdays = ["Tuesday", "Thursday"]
      count    = 5
    }

    retention_monthly {
      format_type = "Weekly"
      weeks       = ["Third", "First"]
      weekdays    = ["Tuesday", "Thursday"]
      count       = 11
    }

    retention_yearly {
      format_type = "Weekly"
      months      = ["July", "February"]
      weeks       = ["Third", "First"]
      weekdays    = ["Tuesday", "Thursday"]
      count       = 9
    }
  }

  protection_policy {
    policy_type = "Differential"

    backup {
      frequency = "Weekly"
      weekdays  = ["Saturday", "Sunday"]
      time      = "17:00"
    }

    simple_retention {
      count = 12
    }
  }

  protection_policy {
    policy_type = "Log"

    backup {
      frequency_in_minutes = 30
    }

    simple_retention {
      count = 9
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
