package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccSiteRecoveryReplicatedVm_managedDiskMembership(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedDiskMembership(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_disk.#").HasValue("1"),
			),
		},
		{
			Config: r.managedDiskMembershipAdded(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_disk.#").HasValue("2"),
			),
		},
		{
			Config: r.managedDiskMembershipUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_disk.#").HasValue("2"),
			),
		},
		{
			Config: r.managedDiskMembership(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("managed_disk.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (r SiteRecoveryReplicatedVmResource) managedDiskMembership(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replicated_vm" "test" {
  name                                      = "repl-%[2]d"
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  source_vm_id                              = azurerm_virtual_machine.test.id
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.test1.name
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.test1.name

  target_resource_group_id                = azurerm_resource_group.test2.id
  target_recovery_fabric_id               = azurerm_site_recovery_fabric.test2.id
  target_recovery_protection_container_id = azurerm_site_recovery_protection_container.test2.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}
`, r.templateWithDataDisk(data), data.RandomInteger)
}

func (r SiteRecoveryReplicatedVmResource) managedDiskMembershipAdded(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replicated_vm" "test" {
  name                                      = "repl-%[2]d"
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  source_vm_id                              = azurerm_virtual_machine.test.id
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.test1.name
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.test1.name

  target_resource_group_id                = azurerm_resource_group.test2.id
  target_recovery_fabric_id               = azurerm_site_recovery_fabric.test2.id
  target_recovery_protection_container_id = azurerm_site_recovery_protection_container.test2.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  managed_disk {
    disk_id                    = azurerm_managed_disk.test.id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
    azurerm_virtual_machine_data_disk_attachment.test,
  ]
}
`, r.templateWithDataDisk(data), data.RandomInteger)
}


func (r SiteRecoveryReplicatedVmResource) managedDiskMembershipUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replicated_vm" "test" {
  name                                      = "repl-%[2]d"
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  source_vm_id                              = azurerm_virtual_machine.test.id
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.test1.name
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.test1.name

  target_resource_group_id                = azurerm_resource_group.test2.id
  target_recovery_fabric_id               = azurerm_site_recovery_fabric.test2.id
  target_recovery_protection_container_id = azurerm_site_recovery_protection_container.test2.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Standard_LRS"
    target_replica_disk_type   = "Standard_LRS"
  }

  managed_disk {
    disk_id                    = azurerm_managed_disk.test.id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Standard_LRS"
    target_replica_disk_type   = "Standard_LRS"
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
    azurerm_virtual_machine_data_disk_attachment.test,
  ]
}
`, r.templateWithDataDisk(data), data.RandomInteger)
}

func (r SiteRecoveryReplicatedVmResource) templateWithDataDisk(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_managed_disk" "test" {
  name                 = "datadisk-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = 10
  caching            = "ReadWrite"
}
`, r.template(data), data.RandomInteger)
}
