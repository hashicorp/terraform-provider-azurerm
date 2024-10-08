// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryReplicatedVmResource struct{}

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

func TestAccSiteRecoveryReplicatedVm_withTFOSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTFOSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_interface.0.failover_test_subnet_name").HasValue("snet3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_withProximityPlacementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withProximityPlacementGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_withBootDiagStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBootDiagStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_withCapacityReservationGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCapacityReservationGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_withVMSS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withVMSS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_withMultiVmGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMultiVmGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_withEdgeZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTargetEdgeZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_withUnManagedDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withUnManagedVmDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_des(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.des(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_zone2zone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zone2zone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_targetDiskEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.targetDiskEncryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_withAvailabilitySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAvailabilitySet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicatedVm_targetVirtualMachineSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replicated_vm", "test")
	r := SiteRecoveryReplicatedVmResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.targetVirtualMachineSize(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.targetVirtualMachineSizeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (SiteRecoveryReplicatedVmResource) template(data acceptance.TestData) string {
	tags := ""
	if strings.HasPrefix(strings.ToLower(data.Client().SubscriptionID), "85b3dbca") {
		tags = `
  tags = {
    "azsecpack"                                                                = "nonprod"
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true"
  }
`
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
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
  name                = "net2-%[1]d"
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_site_recovery_fabric.test2.location
}

resource "azurerm_subnet" "test2" {
  name                 = "snet-%[1]d_2"
  resource_group_name  = azurerm_resource_group.test2.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["192.168.2.0/24"]
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

  delete_os_disk_on_termination = true

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

 %[4]s
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
  name                     = "accsa%[1]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, tags)
}

func (SiteRecoveryReplicatedVmResource) vmSizeTemplate(data acceptance.TestData, vmSize string) string {
	tags := ""
	if strings.HasPrefix(strings.ToLower(data.Client().SubscriptionID), "85b3dbca") {
		tags = `
  tags = {
    "azsecpack"                                                                = "nonprod"
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true"
  }
`
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
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
  name                = "net2-%[1]d"
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_site_recovery_fabric.test2.location
}

resource "azurerm_subnet" "test2" {
  name                 = "snet-%[1]d_2"
  resource_group_name  = azurerm_resource_group.test2.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["192.168.2.0/24"]
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

  vm_size = "%[5]s"

  delete_os_disk_on_termination = true

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

 %[4]s
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
  name                     = "accsa%[1]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, tags, vmSize)
}

func (r SiteRecoveryReplicatedVmResource) basic(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicatedVmResource) withTFOSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "tfo" {
  name                = "net3-%[2]d"
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_site_recovery_fabric.test2.location
}

resource "azurerm_subnet" "tfo" {
  name                 = "snet3"
  resource_group_name  = azurerm_resource_group.test2.name
  virtual_network_name = azurerm_virtual_network.tfo.name
  address_prefixes     = ["192.168.2.0/24"]
}

resource "azurerm_public_ip" "tfo" {
  name                = "pubip%[2]d-tfo"
  allocation_method   = "Static"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  sku                 = "Basic"
}

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
  test_network_id                         = azurerm_virtual_network.tfo.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  network_interface {
    source_network_interface_id        = azurerm_network_interface.test.id
    target_subnet_name                 = azurerm_subnet.test2.name
    recovery_public_ip_address_id      = azurerm_public_ip.test-recovery.id
    failover_test_subnet_name          = azurerm_subnet.tfo.name
    failover_test_public_ip_address_id = azurerm_public_ip.tfo.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (SiteRecoveryReplicatedVmResource) des(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d-1"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-recovery-%[1]d-2"
  location = "%[3]s"
}


resource "azurerm_key_vault" "test" {
  name                        = "kv%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}


resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestdes-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_managed_disk" "test" {
  name                   = "acctestd-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  storage_account_type   = "Standard_LRS"
  create_option          = "Empty"
  disk_size_gb           = 1
  disk_encryption_set_id = azurerm_disk_encryption_set.test.id

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}


resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%[1]d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-focal"
    sku       = "20_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osd-%[1]d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    disk_size_gb      = "50"
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name              = "acctmd-%[1]d"
    create_option     = "Empty"
    disk_size_gb      = "1"
    lun               = 0
    managed_disk_type = "Standard_LRS"
  }

  storage_data_disk {
    name            = azurerm_managed_disk.test.name
    create_option   = "Attach"
    disk_size_gb    = "1"
    lun             = 1
    managed_disk_id = azurerm_managed_disk.test.id
  }

  os_profile {
    computer_name  = "hn%[1]d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
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

resource "azurerm_virtual_network" "test2" {
  name                = "net-%[1]d-2"
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_site_recovery_fabric.test2.location
}

resource "azurerm_site_recovery_network_mapping" "test" {
  resource_group_name         = azurerm_resource_group.test2.name
  recovery_vault_name         = azurerm_recovery_services_vault.test.name
  name                        = "mapping-%[1]d"
  source_recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  target_recovery_fabric_name = azurerm_site_recovery_fabric.test2.name
  source_network_id           = azurerm_virtual_network.test.id
  target_network_id           = azurerm_virtual_network.test2.id
}

resource "azurerm_public_ip" "test-recovery" {
  name                = "pubip%[1]d-recovery"
  allocation_method   = "Static"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  sku                 = "Basic"
}

resource "azurerm_key_vault" "test2" {
  name                        = "kv%[1]d2"
  location                    = azurerm_resource_group.test2.location
  resource_group_name         = azurerm_resource_group.test2.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_access_policy" "service-principal2" {
  key_vault_id = azurerm_key_vault.test2.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
  ]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test2.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = ["azurerm_key_vault_access_policy.service-principal2"]
}

resource "azurerm_disk_encryption_set" "test2" {
  name                = "acctestdes-%[1]d2"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location
  key_vault_key_id    = azurerm_key_vault_key.test2.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption2" {
  key_vault_id = azurerm_key_vault.test2.id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_disk_encryption_set.test2.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test2.identity.0.principal_id
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
    disk_id                       = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id    = azurerm_storage_account.test.id
    target_resource_group_id      = azurerm_resource_group.test2.id
    target_disk_type              = "Premium_LRS"
    target_replica_disk_type      = "Premium_LRS"
    target_disk_encryption_set_id = azurerm_disk_encryption_set.test2.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (SiteRecoveryReplicatedVmResource) zone2zone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
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

resource "azurerm_site_recovery_protection_container" "test1" {
  resource_group_name  = azurerm_resource_group.test2.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  name                 = "acctest-protection-cont1-%[1]d"
}

resource "azurerm_site_recovery_protection_container" "test2" {
  resource_group_name  = azurerm_resource_group.test2.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  name                 = "acctest-protection-cont2-t-%[1]d"
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
  location            = azurerm_site_recovery_fabric.test1.location
}

resource "azurerm_network_interface" "test" {
  name                = "vm-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "vm-%[1]d"
    subnet_id                     = azurerm_subnet.test1.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                = "vm-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  vm_size = "Standard_B1s"

  delete_os_disk_on_termination = true

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
  target_zone                               = "2"

  target_resource_group_id                = azurerm_resource_group.test2.id
  target_recovery_fabric_id               = azurerm_site_recovery_fabric.test1.id
  target_recovery_protection_container_id = azurerm_site_recovery_protection_container.test2.id
  target_network_id                       = azurerm_virtual_network.test1.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  network_interface {
    source_network_interface_id = azurerm_network_interface.test.id
    target_subnet_name          = "snet-%[1]d"
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (SiteRecoveryReplicatedVmResource) targetDiskEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy          = false
      purge_soft_deleted_keys_on_destroy    = false
      purge_soft_deleted_secrets_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d-1"
  location = "%[2]s"
}
resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-recovery-%[1]d-2"
  location = "%[3]s"
}

resource "azurerm_key_vault" "test1" {
  name                        = "acckv-%[1]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "premium"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test1.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
  ]
}

resource "azurerm_key_vault_key" "test1" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test1.id
  key_type     = "RSA"
  key_size     = 3072

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault_access_policy.service-principal
  ]
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
resource "azurerm_site_recovery_protection_container" "test1" {
  resource_group_name  = azurerm_resource_group.test2.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  name                 = "acctest-protection-cont1-%[1]d"
}
resource "azurerm_site_recovery_protection_container" "test2" {
  resource_group_name  = azurerm_resource_group.test2.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  name                 = "acctest-protection-cont2-t-%[1]d"
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
resource "azurerm_network_interface" "test" {
  name                = "vm-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_configuration {
    name                          = "vm-%[1]d"
    subnet_id                     = azurerm_subnet.test1.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_windows_virtual_machine" "vm" {
  name                = "acctvm%[4]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D2s_v3"
  admin_username      = "adminuser"
  admin_password      = "P@ssw0rd1234!"
  zone                = "1"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }
}

resource "azurerm_virtual_machine_extension" "test" {
  name                       = "AzureDiskEncryption"
  publisher                  = "Microsoft.Azure.Security"
  type                       = "AzureDiskEncryption"
  type_handler_version       = "2.2"
  auto_upgrade_minor_version = false
  virtual_machine_id         = azurerm_windows_virtual_machine.vm.id

  settings = <<SETTINGS
{
  "EncryptionOperation": "EnableEncryption",
  "KeyEncryptionAlgorithm": "RSA-OAEP",
  "KeyVaultURL": "${azurerm_key_vault.test1.vault_uri}",
  "KeyVaultResourceId": "${azurerm_key_vault.test1.id}",
  "KeyEncryptionKeyURL": "${azurerm_key_vault_key.test1.id}",
  "KekVaultResourceId": "${azurerm_key_vault.test1.id}",
  "VolumeType": "All"
}
SETTINGS
}

resource "azurerm_storage_account" "test" {
  name                     = "acct%[1]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_managed_disk" "test" {
  name                = azurerm_windows_virtual_machine.vm.os_disk[0].name
  resource_group_name = azurerm_windows_virtual_machine.vm.resource_group_name

  depends_on = [
    azurerm_virtual_machine_extension.test
  ]
}

// Use snapshot as a workaround of encryption_settings not yet supported on managed_disk
resource "azurerm_snapshot" "test" {
  name                = "snapshot-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  create_option       = "Copy"
  source_resource_id  = data.azurerm_managed_disk.test.id
  lifecycle {
    ignore_changes = [
      encryption_settings
    ]
  }
}

data "azurerm_snapshot" "test" {
  name                = azurerm_snapshot.test.name
  resource_group_name = azurerm_snapshot.test.resource_group_name
}

resource "azurerm_site_recovery_replicated_vm" "test" {
  name                                      = "repl-%[1]d"
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  source_vm_id                              = azurerm_windows_virtual_machine.vm.id
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.test1.name
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.test1.name
  target_zone                               = "2"
  target_resource_group_id                  = azurerm_resource_group.test2.id
  target_recovery_fabric_id                 = azurerm_site_recovery_fabric.test1.id
  target_recovery_protection_container_id   = azurerm_site_recovery_protection_container.test2.id
  target_network_id                         = azurerm_virtual_network.test1.id

  managed_disk {
    disk_id                    = data.azurerm_managed_disk.test.id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
    target_disk_encryption {
      disk_encryption_key {
        secret_url = data.azurerm_snapshot.test.encryption_settings[0].disk_encryption_key[0].secret_url
        vault_id   = azurerm_key_vault.test1.id
      }
      key_encryption_key {
        key_url  = azurerm_key_vault_key.test1.id
        vault_id = azurerm_key_vault.test1.id
      }
    }
  }
  network_interface {
    source_network_interface_id = azurerm_network_interface.test.id
    target_subnet_name          = "snet-%[1]d"
  }
  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, data.RandomString)
}

func (r SiteRecoveryReplicatedVmResource) withProximityPlacementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctest-replication-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
}

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
  target_proximity_placement_group_id     = azurerm_proximity_placement_group.test.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  network_interface {
    source_network_interface_id   = azurerm_network_interface.test.id
    target_subnet_name            = azurerm_subnet.test2.name
    recovery_public_ip_address_id = azurerm_public_ip.test-recovery.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}


`, r.template(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r SiteRecoveryReplicatedVmResource) withBootDiagStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test2" {
  name                     = "acctre%[2]d"
  location                 = azurerm_resource_group.test2.location
  resource_group_name      = azurerm_resource_group.test2.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_site_recovery_replicated_vm" "test" {
  name                                      = "repl-%[2]d"
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  source_vm_id                              = azurerm_virtual_machine.test.id
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.test1.name
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.test1.name

  target_resource_group_id                  = azurerm_resource_group.test2.id
  target_recovery_fabric_id                 = azurerm_site_recovery_fabric.test2.id
  target_recovery_protection_container_id   = azurerm_site_recovery_protection_container.test2.id
  target_boot_diagnostic_storage_account_id = azurerm_storage_account.test2.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  network_interface {
    source_network_interface_id   = azurerm_network_interface.test.id
    target_subnet_name            = azurerm_subnet.test2.name
    recovery_public_ip_address_id = azurerm_public_ip.test-recovery.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}




`, r.template(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r SiteRecoveryReplicatedVmResource) withUnManagedVmDisks(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d-1"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-recovery-%[1]d-2"
  location = "%[3]s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%[1]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 1
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}


resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%[1]d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_D1_v2"

  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-focal"
    sku       = "20_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "hn%[1]d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Production"
    cost-center = "Ops"
  }
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

resource "azurerm_virtual_network" "test2" {
  name                = "net-%[1]d-2"
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_site_recovery_fabric.test2.location
}

resource "azurerm_site_recovery_network_mapping" "test" {
  resource_group_name         = azurerm_resource_group.test2.name
  recovery_vault_name         = azurerm_recovery_services_vault.test.name
  name                        = "mapping-%[1]d"
  source_recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  target_recovery_fabric_name = azurerm_site_recovery_fabric.test2.name
  source_network_id           = azurerm_virtual_network.test.id
  target_network_id           = azurerm_virtual_network.test2.id
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

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctre%[1]d"
  location                 = azurerm_resource_group.test2.location
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

  unmanaged_disk {
    disk_uri                   = azurerm_virtual_machine.test.storage_os_disk[0].vhd_uri
    staging_storage_account_id = azurerm_storage_account.test.id
    target_storage_account_id  = azurerm_storage_account.test2.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r SiteRecoveryReplicatedVmResource) withCapacityReservationGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-crg-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
}

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%[2]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id

  sku {
    name     = "Standard_B1s"
    capacity = 1
  }
}

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
  target_capacity_reservation_group_id    = azurerm_capacity_reservation_group.test.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  network_interface {
    source_network_interface_id   = azurerm_network_interface.test.id
    target_subnet_name            = azurerm_subnet.test2.name
    recovery_public_ip_address_id = azurerm_public_ip.test-recovery.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
    azurerm_capacity_reservation.test,
  ]
}

`, r.template(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r SiteRecoveryReplicatedVmResource) withVMSS(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name

  sku_name                    = "Standard_B1s"
  instances                   = 2
  platform_fault_domain_count = 2

  network_interface {
    name    = "TestNetworkProfile-%[2]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test2.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-focal"
    sku       = "20_04-lts"
    version   = "latest"
  }

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[2]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

}

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
  target_virtual_machine_scale_set_id     = azurerm_orchestrated_virtual_machine_scale_set.test.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  network_interface {
    source_network_interface_id   = azurerm_network_interface.test.id
    target_subnet_name            = azurerm_subnet.test2.name
    recovery_public_ip_address_id = azurerm_public_ip.test-recovery.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}


`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r SiteRecoveryReplicatedVmResource) withMultiVmGroup(data acceptance.TestData) string {
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
  multi_vm_group_name                       = "accmultivmgroup"

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
    target_subnet_name            = azurerm_subnet.test2.name
    recovery_public_ip_address_id = azurerm_public_ip.test-recovery.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicatedVmResource) withTargetEdgeZone(data acceptance.TestData) string {
	// WestUS has an edge zone available - so hard-code to that for now
	data.Locations.Secondary = "westus"

	return fmt.Sprintf(`
%s

data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}

resource "azurerm_site_recovery_replicated_vm" "test" {
  name                                      = "repl-%[2]d"
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  source_vm_id                              = azurerm_virtual_machine.test.id
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.test1.name
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.test1.name
  target_edge_zone                          = data.azurerm_extended_locations.test.extended_locations[0]

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
    target_subnet_name            = azurerm_subnet.test2.name
    recovery_public_ip_address_id = azurerm_public_ip.test-recovery.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}
`, r.template(data), data.RandomInteger)

}

func (r SiteRecoveryReplicatedVmResource) withAvailabilitySet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_availability_set" "test" {
  name                = "acctestavset-%d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  managed             = true
}

resource "azurerm_site_recovery_replicated_vm" "test" {
  name                                      = "repl-%[2]d"
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  source_vm_id                              = azurerm_virtual_machine.test.id
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.test1.name
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.test1.name

  target_availability_set_id              = azurerm_availability_set.test.id
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
    target_subnet_name            = azurerm_subnet.test2.name
    recovery_public_ip_address_id = azurerm_public_ip.test-recovery.id
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicatedVmResource) targetVirtualMachineSize(data acceptance.TestData) string {
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
  target_virtual_machine_size             = "Standard_B1s"

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
`, r.vmSizeTemplate(data, "Standard_B1s"), data.RandomInteger)
}

func (r SiteRecoveryReplicatedVmResource) targetVirtualMachineSizeUpdated(data acceptance.TestData) string {
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
  target_virtual_machine_size             = "Standard_B2s"

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
`, r.vmSizeTemplate(data, "Standard_B2s"), data.RandomInteger)
}

func (r SiteRecoveryReplicatedVmResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ReplicationProtectedItemsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery replicated vm (%s): %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return nil, fmt.Errorf("reading site recovery replicated vm (%s): model is nil", id.String())
	}

	return utils.Bool(model.Id != nil), nil
}
