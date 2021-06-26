package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SiteRecoveryReplicatedVmResource struct {
}

func TestAccSiteRecoveryReplicatedVm_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (SiteRecoveryReplicatedVmResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d-1"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-recovery-%[1]d-2"
  location = "%[3]s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%[1]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_site_recovery_fabric" "test1" {
  resource_group_name = azurerm_resource_group.test2.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  name                = "acctest-fabric1-%[1]d"
  location            = azurerm_resource_group.test.location
}

resource "azurerm_site_recovery_fabric" "test2" {
  resource_group_name = azurerm_resource_group.test2.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  name                = "acctest-fabric2-%[1]d"
  location            = azurerm_resource_group.test2.location
  depends_on          = [azurerm_site_recovery_fabric.test1]
}

resource "azurerm_site_recovery_protection_container" "test1" {
  resource_group_name  = azurerm_resource_group.test2.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  name                 = "acctest-protection-cont1-%[1]d"
}

resource "azurerm_site_recovery_protection_container" "test2" {
  resource_group_name  = azurerm_resource_group.test2.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test2.name
  name                 = "acctest-protection-cont2-%[1]d"
}

resource "azurerm_site_recovery_replication_policy" "test" {
  resource_group_name                                  = azurerm_resource_group.test2.name
  recovery_vault_name                                  = azurerm_recovery_services_vault.test.name
  name                                                 = "acctest-policy-%[1]d"
  recovery_point_retention_in_minutes                  = 24 * 60
  application_consistent_snapshot_frequency_in_minutes = 4 * 60
}

resource "azurerm_site_recovery_protection_container_mapping" "test" {
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  recovery_fabric_name                      = azurerm_site_recovery_fabric.test1.name
  recovery_source_protection_container_name = azurerm_site_recovery_protection_container.test1.name
  recovery_target_protection_container_id   = azurerm_site_recovery_protection_container.test2.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  name                                      = "mapping-%[1]d"
}

resource "azurerm_virtual_network" "test1" {
  name                = "net-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_site_recovery_fabric.test1.location
}

resource "azurerm_subnet" "test1" {
  name                 = "snet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test1.name
  address_prefixes     = ["192.168.1.0/24"]
}

resource "azurerm_virtual_network" "test2" {
  name                = "net-%[1]d"
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_site_recovery_fabric.test2.location
}

resource "azurerm_subnet" "test2_1" {
  name                 = "acctest-snet-%[1]d_1"
  resource_group_name  = "${azurerm_resource_group.test2.name}"
  virtual_network_name = "${azurerm_virtual_network.test2.name}"
  address_prefixes     = ["192.168.2.0/27"]
}

resource "azurerm_subnet" "test2_2" {
  name                 = "snet-%[1]d_2"
  resource_group_name  = "${azurerm_resource_group.test2.name}"
  virtual_network_name = "${azurerm_virtual_network.test2.name}"
  address_prefixes     = ["192.168.2.32/27"]
}

resource "azurerm_site_recovery_network_mapping" "test" {
  resource_group_name         = azurerm_resource_group.test2.name
  recovery_vault_name         = azurerm_recovery_services_vault.test.name
  name                        = "mapping-%[1]d"
  source_recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  target_recovery_fabric_name = azurerm_site_recovery_fabric.test2.name
  source_network_id           = azurerm_virtual_network.test1.id
  target_network_id           = azurerm_virtual_network.test2.id
}

resource "azurerm_network_interface" "test" {
  name                = "vm-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "vm-%[1]d"
    subnet_id                     = azurerm_subnet.test1.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test-source.id
  }
}

resource "azurerm_virtual_machine" "test" {
  name                = "vm-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  vm_size = "Standard_B1s"

  storage_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7.5"
    version   = "latest"
  }

  storage_os_disk {
    name              = "disk-%[1]d"
    os_type           = "Linux"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    admin_username = "testadmin"
    admin_password = "Password1234!"
    computer_name  = "vm-%[1]d"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
  network_interface_ids = [azurerm_network_interface.test.id]
}

resource "azurerm_public_ip" "test-source" {
  name                = "pubip%[1]d-source"
  allocation_method   = "Static"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_public_ip" "test-recovery" {
  name                = "pubip%[1]d-recovery"
  allocation_method   = "Static"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  sku                 = "Basic"
}

resource "azurerm_storage_account" "test" {
  name                     = "acct%[1]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_site_recovery_replicated_vm" "test" {
  name                                      = "repl-%[1]d"
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

  network_interface {
    source_network_interface_id   = azurerm_network_interface.test.id
    target_subnet_name            = "snet-%[1]d_2"
    recovery_public_ip_address_id = azurerm_public_ip.test-recovery.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (t SiteRecoveryReplicatedVmResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	protectionContainerName := id.Path["replicationProtectionContainers"]
	name := id.Path["replicationProtectedItems"]

	resp, err := clients.RecoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName).Get(ctx, fabricName, protectionContainerName, name)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery replicated vm (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
