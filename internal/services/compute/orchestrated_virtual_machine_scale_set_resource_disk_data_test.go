// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccOrchestratedVirtualMachineScaleSet_basicLinux_managedDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux_managedDisk(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_importBasic_managedDisk_withZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicLinux_managedDisk_withZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_loadBalancerManagedDataDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.loadBalancerTemplateManagedDataDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("data_disk.#").HasValue("1"),
			),
		},
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_disksDataDiskStorageAccountTypePremiumLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountType(data, "Premium_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_disksDataDiskStorageAccountTypePremiumV2LRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypePremiumV2LRS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_disksDataDiskStorageAccountTypePremiumZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypeWithRestrictedLocation(data, "Premium_ZRS", "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_disksDataDiskStorageAccountTypeStandardLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountType(data, "Standard_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_disksDataDiskStorageAccountTypeStandardSSDLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountType(data, "StandardSSD_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_disksDataDiskStorageAccountTypeStandardSSDZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypeWithRestrictedLocation(data, "StandardSSD_ZRS", "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_disksDataDiskStorageAccountTypeUltraSSDLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypeUltraSSDLRS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccOrchestratedVirtualMachineScaleSet_disksDataDiskSizeFromMarketPlaceImage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_orchestrated_virtual_machine_scale_set", "test")
	r := OrchestratedVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataDiskMarketPlaceImage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("data_disk.0.disk_size_gb").HasValue("900"),
				check.That(data.ResourceName).Key("data_disk.0.lun").HasValue("0"),
			),
		},
		data.ImportStep("os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (OrchestratedVirtualMachineScaleSetResource) basicLinux_managedDisk(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name  = "Standard_D1_v2"
  instances = 2

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) basicLinux_managedDisk_withZones(data acceptance.TestData) string {
	r := OrchestratedVirtualMachineScaleSetResource{}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  zones                       = ["1"]
  platform_fault_domain_count = 1

  sku_name  = "Standard_D1_v2"
  instances = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile-%[1]d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, r.natgateway_template(data))
}

func (OrchestratedVirtualMachineScaleSetResource) loadBalancerTemplateManagedDataDisks(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
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

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "test"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name  = "Standard_F2"
  instances = 1

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    lun                  = 0
    caching              = "ReadWrite"
    create_option        = "Empty"
    disk_size_gb         = 10
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (OrchestratedVirtualMachineScaleSetResource) disksDataDiskStorageAccountType(data acceptance.TestData, storageAccountType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
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

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "test"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name  = "Standard_F2s_v2"
  instances = 1

  platform_fault_domain_count = 2

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    lun                  = 0
    caching              = "ReadWrite"
    create_option        = "Empty"
    disk_size_gb         = 10
    storage_account_type = "%[3]s"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, storageAccountType)
}

func (r OrchestratedVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypeWithRestrictedLocation(data acceptance.TestData, storageAccountType string, location string) string {
	// Limited regional availability for some storage account type
	data.Locations.Primary = location
	return r.disksDataDiskStorageAccountType(data, storageAccountType)
}

func (OrchestratedVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypePremiumV2LRS(data acceptance.TestData) string {
	// Limited regional availability for `PremiumV2_LRS`
	// `PremiumV2_LRS` disks can only be can only be attached to zonal VMs currently
	data.Locations.Primary = "westeurope"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
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

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "test"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  zones               = ["1"]

  sku_name  = "Standard_F2s_v2"
  instances = 1

  platform_fault_domain_count = 1

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    lun                  = 0
    caching              = "None"
    create_option        = "Empty"
    disk_size_gb         = 10
    storage_account_type = "PremiumV2_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (OrchestratedVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypeUltraSSDLRS(data acceptance.TestData) string {
	// Limited regional availability for `UltraSSD_LRS`
	// `UltraSSD_LRS` disks needs to be used with `ultra_ssd_enabled` set to `true`
	data.Locations.Primary = "eastus2"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
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

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "test"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  zones               = ["1"]

  sku_name  = "Standard_F2s_v2"
  instances = 1

  platform_fault_domain_count = 1

  additional_capabilities {
    ultra_ssd_enabled = true
  }

  os_profile {
    linux_configuration {
      computer_name_prefix = "testvm-%[1]d"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    lun                  = 0
    caching              = "None"
    create_option        = "Empty"
    disk_size_gb         = 10
    storage_account_type = "UltraSSD_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (OrchestratedVirtualMachineScaleSetResource) dataDiskMarketPlaceImage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-OVMSS-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/8"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsn-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku = "Standard"

  frontend_ip_configuration {
    name                          = "default"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "acctestbap-%[1]d"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_marketplace_agreement" "barracuda" {
  publisher = "micro-focus"
  offer     = "arcsight-logger"
  plan      = "arcsight_logger_72_byol"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "test" {
  name                = "acctestOVMSS-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name  = "Standard_F2"
  instances = 1

  platform_fault_domain_count = 2

  os_profile {

    linux_configuration {
      computer_name_prefix = "testvm-test"
      admin_username       = "myadmin"
      admin_password       = "Passwword1234"

      disable_password_authentication = false
    }
  }

  network_interface {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id

      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    caching              = "ReadWrite"
    disk_size_gb         = 900
    create_option        = "FromImage"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "micro-focus"
    offer     = "arcsight-logger"
    sku       = "arcsight_logger_72_byol"
    version   = "7.2.0"
  }

  plan {
    name      = "arcsight_logger_72_byol"
    product   = "arcsight-logger"
    publisher = "micro-focus"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
