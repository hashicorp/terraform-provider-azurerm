package recoveryservices_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryVMWareReplicatedVmResource struct{}

func (r SiteRecoveryVMWareReplicatedVmResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ReplicationProtectedItemsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func TestAccSiteVMWareRecoveryReplicatedVM_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_vmware_replicated_vm", "test")
	r := SiteRecoveryVMWareReplicatedVmResource{}

	vaultId := os.Getenv("ARM_TEST_VMWARE_VAULT_ID")
	sourceVMName := os.Getenv("ARM_TEST_VMWARE_SOURCE_VM_NAME")
	applianceName := os.Getenv("ARM_TEST_VMWARE_APPLIANCE_NAME")
	location := os.Getenv("ARM_TEST_VMWARE_VAULT_LOCATION")

	if vaultId == "" || sourceVMName == "" || applianceName == "" || location == "" {
		t.Skip("Skipping since ARM_TEST_VMWARE_VAULT_ID, ARM_TEST_VMWARE_SOURCE_VM_NAME, ARM_TEST_VMWARE_LOCATION and ARM_TEST_VMWARE_APPLIANCE_NAME are not specified")
		return
	}

	data.Locations.Primary = location

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, vaultId, sourceVMName, applianceName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (SiteRecoveryVMWareReplicatedVmResource) basic(data acceptance.TestData, vaultId, sourceVMName, applianceName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_site_recovery_vmware_replication_policy" "test" {
  recovery_vault_id                                    = "%[1]s"
  name                                                 = "acctest-policy-%[2]d"
  recovery_point_retention_in_minutes                  = 1440
  application_consistent_snapshot_frequency_in_minutes = 240
}

resource "azurerm_site_recovery_vmware_replication_policy_association" "test" {
  name              = "acctest-%[2]d"
  recovery_vault_id = "%[1]s"
  policy_id         = azurerm_site_recovery_vmware_replication_policy.test.id
}

resource "azurerm_resource_group" "target" {
  name     = "acctestRG-%[2]d"
  location = "%[4]s"
}

resource "azurerm_storage_account" "target" {
  name                     = "acct%[3]s"
  resource_group_name      = azurerm_resource_group.target.name
  location                 = azurerm_resource_group.target.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "target" {
  name                = "acctestvn-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.target.location
  resource_group_name = azurerm_resource_group.target.name
}

resource "azurerm_subnet" "target" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.target.name
  virtual_network_name = azurerm_virtual_network.target.name
  address_prefixes     = ["10.0.2.0/24"]
}


resource "azurerm_site_recovery_vmware_replicated_vm" "test" {
  name                                       = "acct%[2]d"
  recovery_vault_id                          = "%[1]s"
  source_vm_name                             = "%[5]s"
  appliance_name                             = "%[6]s"
  recovery_replication_policy_id             = azurerm_site_recovery_vmware_replication_policy.test.id
  credential_type                            = "lincreds"
  license_type                               = "NotSpecified"
  target_boot_diagnostics_storage_account_id = azurerm_storage_account.target.id
  target_vm_name                             = "%[5]s"
  target_resource_group_id                   = azurerm_resource_group.target.id
  default_log_storage_account_id             = azurerm_storage_account.target.id
  default_recovery_disk_type                 = "Standard_LRS"
  target_network_id                          = azurerm_virtual_network.target.id

  network_interface {
    target_subnet_name = azurerm_subnet.target.name
    is_primary         = true
  }

  timeouts {
    create = "600m"
  }

  lifecycle {
    ignore_changes = [
      target_vm_size,
      credential_type,
      default_log_storage_account_id,
      default_recovery_disk_type,
      managed_disk,
      network_interface
    ]
  }
}
`, vaultId, data.RandomInteger, data.RandomString, data.Locations.Primary, sourceVMName, applianceName)
}
