package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
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

func (t BackupProtectionPolicyVMWorkloadResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ProtectionPoliciesClient.Get(ctx, id.VaultName, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service Protection Policy (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r BackupProtectionPolicyVMWorkloadResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_backup_policy_vm_workload" "test" {
  name                  = "acctest-bpvmw-test01"
  resource_group_name   = "acctestRG-bpvmw-test01"
  recovery_vault_name   = "acctest-rsv-test01"

  workload_type = "SQLDataBase"

  settings {
    time_zone               = "UTC"
    compression_enabled     = true
    sql_compression_enabled = true
  }

  protection_policy {
    policy_type = "Full"

    backup {
      frequency = "Weekly"
      time      = "14:00"
      weekdays  = ["Monday", "Tuesday"]
    }

    retention_weekly {
      weekdays = ["Monday", "Tuesday"]
      count    = 3
    }

    retention_monthly {
      format_type = "Weekly"
      count       = 6
      weeks       = ["Third"]
      weekdays    = ["Monday"]
    }
  }

  protection_policy {
    policy_type = "Log"

    backup {
      frequency_in_minutes = 30
    }

    simple_retention {
      count = 8
    }
  }
}
`)
}
