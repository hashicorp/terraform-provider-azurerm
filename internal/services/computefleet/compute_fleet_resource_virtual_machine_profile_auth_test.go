// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccComputeFleet_virtualMachineProfileAuth_authPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authPassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileAuth_authSSHKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSSHKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccComputeFleet_virtualMachineProfileAuth_authMultipleSSHPublicKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authMultipleSSHPublicKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccComputeFleet_virtualMachineProfileAuth_authSSHKeyAndPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authSSHKeyAndPassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileAuth_authEd25519SSHPublicKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authEd25519SSHPublicKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ComputeFleetResource) authPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_capacity {
    minimum_capacity          = 1
    maintain_capacity_enabled = false
    target_capacity           = 1
  }

  vm_sizes_profile {
    name = "Standard_F1alds_v7"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "MicrosoftWindowsServer"
      offer     = "WindowsServer"
      sku       = "2025-datacenter-core-g2"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    os_profile {
      windows_configuration {
        computer_name_prefix = "testvm"
        admin_username       = local.admin_username
        admin_password       = local.admin_password
      }
    }

    network_interface {
      name = "networkProTest"
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
      primary_network_interface_enabled = true
    }
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) authSSHKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_capacity {
    minimum_capacity          = 1
    maintain_capacity_enabled = false
    target_capacity           = 1
  }

  vm_sizes_profile {
    name = "Standard_F1alds_v7"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "canonical"
      offer     = "ubuntu-24_04-lts"
      sku       = "server"
      version   = "latest"
    }
    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        password_authentication_enabled = false
        admin_ssh_keys                  = [local.first_public_key]
      }
    }

    network_interface {
      name = "networkProTest"
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
      primary_network_interface_enabled = true
    }
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) authMultipleSSHPublicKeys(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_capacity {
    minimum_capacity          = 1
    maintain_capacity_enabled = false
    target_capacity           = 1
  }

  vm_sizes_profile {
    name = "Standard_F1alds_v7"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "canonical"
      offer     = "ubuntu-24_04-lts"
      sku       = "server"
      version   = "latest"
    }
    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        password_authentication_enabled = false
        admin_ssh_keys                  = [local.first_public_key, local.second_public_key]
      }
    }
    network_interface {
      name = "networkProTest"
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
      primary_network_interface_enabled = true
    }
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) authSSHKeyAndPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_capacity {
    minimum_capacity          = 1
    maintain_capacity_enabled = false
    target_capacity           = 1
  }

  vm_sizes_profile {
    name = "Standard_F1alds_v7"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "canonical"
      offer     = "ubuntu-24_04-lts"
      sku       = "server"
      version   = "latest"
    }
    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        admin_ssh_keys                  = [local.first_public_key]
        password_authentication_enabled = true
      }
    }

    network_interface {
      name = "networkProTest"
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
      primary_network_interface_enabled = true
    }
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) authEd25519SSHPublicKeys(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_capacity {
    minimum_capacity          = 1
    maintain_capacity_enabled = false
    target_capacity           = 1
  }

  vm_sizes_profile {
    name = "Standard_F1alds_v7"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "canonical"
      offer     = "ubuntu-24_04-lts"
      sku       = "server"
      version   = "latest"
    }
    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        password_authentication_enabled = false
        admin_ssh_keys                  = [local.first_ed25519_public_key]
      }
    }

    network_interface {
      name = "networkProTest"
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
      primary_network_interface_enabled = true
    }
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}
