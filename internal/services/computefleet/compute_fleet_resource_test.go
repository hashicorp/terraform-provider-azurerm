// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ComputeFleetTestResource struct{}

func TestAccComputeFleet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccComputeFleet_completeExceptVMSS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeExceptVMSS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeExceptVMSS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.completeExceptVMSSUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.completeExceptVMSS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fleets.ParseFleetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ComputeFleet.ComputeFleetClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ComputeFleetTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

`, r.templateWithOutProvider(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) templateWithOutProvider(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  first_public_key          = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  second_public_key         = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"
  first_ed25519_public_key  = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDqzSi9IHoYnbE3YQ+B2fQEVT8iGFemyPovpEtPziIVB you@me.com"
  second_ed25519_public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDqzSi9IHoYnbE3YQ+B2fQEVT8iGFemyPovpEtPziIVB hello@world.com"
  admin_username            = "testadmin1234"
  admin_password            = "Password1234!"
  admin_password_update     = "Password1234!Update"
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-fleet-%[1]d"
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

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicIP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "test" {
  name                = "acctest-loadbalancer-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) baseAndAdditionalLocationLinuxTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "linux_test" {
  name     = "acctest-rg-fleet-al-%[2]d"
  location = "%[3]s"
}

resource "azurerm_virtual_network" "linux_test" {
  name                = "acctvn-al-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
}

resource "azurerm_subnet" "linux_test" {
  name                 = "acctsub-%[2]d"
  resource_group_name  = azurerm_resource_group.linux_test.name
  virtual_network_name = azurerm_virtual_network.linux_test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "linux_test" {
  name                = "acctestpublicIP%[2]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "linux_test" {
  name                = "acctest-loadbalancer-%[2]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.linux_test.id
  }
}

resource "azurerm_lb_backend_address_pool" "linux_test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.linux_test.id
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) baseAndAdditionalLocationWindowsTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "windows_test" {
  name     = "acctest-rg-fleet-al-win-%[2]d"
  location = "%[3]s"
}

resource "azurerm_virtual_network" "windows_test" {
  name                = "acctvn-al-win-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
}

resource "azurerm_subnet" "windows_test" {
  name                 = "acctsub-al-win-%[2]d"
  resource_group_name  = azurerm_resource_group.windows_test.name
  virtual_network_name = azurerm_virtual_network.windows_test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "windows_test" {
  name                = "acctestpublicIP-al-%[2]d"
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "windows_test" {
  name                = "acctest-loadbalancer-al-%[2]d"
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal-al-%[2]d"
    public_ip_address_id = azurerm_public_ip.windows_test.id
  }
}

resource "azurerm_lb_backend_address_pool" "windows_test" {
  name            = "internal-al"
  loadbalancer_id = azurerm_lb.windows_test.id
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) baseAndAdditionalLocationLinuxTemplateWithOutProvider(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "linux_test" {
  name     = "acctest-rg-fleet-al-%[2]d"
  location = "%[3]s"
}

resource "azurerm_virtual_network" "linux_test" {
  name                = "acctvn-al-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
}

resource "azurerm_subnet" "linux_test" {
  name                 = "acctsub-%[2]d"
  resource_group_name  = azurerm_resource_group.linux_test.name
  virtual_network_name = azurerm_virtual_network.linux_test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "linux_test" {
  name                = "acctestpublicIP%[2]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "linux_test" {
  name                = "acctest-loadbalancer-%[2]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.linux_test.id
  }
}

resource "azurerm_lb_backend_address_pool" "linux_test" {
  name            = "internal"
  loadbalancer_id = azurerm_lb.linux_test.id
}
`, r.templateWithOutProvider(data), data.RandomInteger, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) baseAndAdditionalLocationWindowsTemplateWithOutProvider(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "windows_test" {
  name     = "acctest-rg-fleet-al-win-%[2]d"
  location = "%[3]s"
}

resource "azurerm_virtual_network" "windows_test" {
  name                = "acctvn-al-win-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
}

resource "azurerm_subnet" "windows_test" {
  name                 = "acctsub-al-win-%[2]d"
  resource_group_name  = azurerm_resource_group.windows_test.name
  virtual_network_name = azurerm_virtual_network.windows_test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "windows_test" {
  name                = "acctestpublicIP-al-%[2]d"
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "windows_test" {
  name                = "acctest-loadbalancer-al-%[2]d"
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "internal-al-%[2]d"
    public_ip_address_id = azurerm_public_ip.windows_test.id
  }
}

resource "azurerm_lb_backend_address_pool" "windows_test" {
  name            = "internal-al"
  loadbalancer_id = azurerm_lb.windows_test.id
}
`, r.templateWithOutProvider(data), data.RandomInteger, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    min_capacity     = 1
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      offer     = "0001-com-ubuntu-server-focal"
      publisher = "canonical"
      sku       = "20_04-lts-gen2"
      version   = "latest"
    }

    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "networkProTest"
      primary_network_interface_enabled = true
      ip_configuration {
        name                             = "TestIPConfiguration"
        subnet_id                        = azurerm_subnet.test.id
        primary_ip_configuration_enabled = true
        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }
  }
  # ignore_changes os_disk as os_disk block is not specified the API return default values for caching, delete_option, disk_size_in_gib and storage_account_type
  # ignore_changes compute_api_version as the default value returned by API will be the latest supported computeApiVersion if it is not specified
  lifecycle {
    ignore_changes = [compute_api_version, virtual_machine_profile.0.os_disk]
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_compute_fleet" "import" {
  name                = azurerm_compute_fleet.test.name
  resource_group_name = azurerm_compute_fleet.test.resource_group_name
  location            = azurerm_compute_fleet.test.location

  spot_priority_profile {
    min_capacity     = 1
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      offer     = "0001-com-ubuntu-server-focal"
      publisher = "canonical"
      sku       = "20_04-lts-gen2"
      version   = "latest"
    }

    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "networkProTest"
      primary_network_interface_enabled = true
      ip_configuration {
        name                             = "TestIPConfiguration"
        subnet_id                        = azurerm_subnet.test.id
        primary_ip_configuration_enabled = true
        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }
  }
  # ignore_changes os_disk as os_disk block is not specified the API return default values for caching, delete_option, disk_size_in_gib and storage_account_type
  # ignore_changes compute_api_version as the default value returned by API will be the latest supported computeApiVersion if it is not specified
  lifecycle {
    ignore_changes = [compute_api_version, virtual_machine_profile.0.os_disk]
  }
}
`, r.basic(data))
}

func (r ComputeFleetTestResource) completeExceptVMSS(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%[1]s

resource "azurerm_marketplace_agreement" "barracuda" {
  publisher = "micro-focus"
  offer     = "arcsight-logger"
  plan      = "arcsight_logger_72_byol"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest2%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 1
  compute_api_version         = "2024-03-01"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  plan {
    name           = "arcsight_logger_72_byol"
    product        = "arcsight-logger"
    publisher      = "micro-focus"
    promotion_code = "test"
  }

  spot_priority_profile {
    allocation_strategy     = "PriceCapacityOptimized"
    eviction_policy         = "Delete"
    max_hourly_price_per_vm = -1
    min_capacity            = 0
    maintain_enabled        = false
    capacity                = 0
  }

  regular_priority_profile {
    allocation_strategy = "LowestPrice"
    capacity            = 0
    min_capacity        = 0
  }

  vm_sizes_profile {
    name = "Standard_DS1"
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "micro-focus"
      offer     = "arcsight-logger"
      sku       = "arcsight_logger_72_byol"
      version   = "7.2.0"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      caching              = "ReadWrite"
      create_option        = "FromImage"
      storage_account_type = "Standard_LRS"
    }

    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "networkProTest"
      primary_network_interface_enabled = true
      ip_configuration {
        name                             = "TestIPConfiguration"
        subnet_id                        = azurerm_subnet.test.id
        primary_ip_configuration_enabled = true
        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }
  }

  tags = {
    Hello = "There"
    World = "Example"
  }
  zones = ["1", "2", "3"]

  depends_on = [azurerm_marketplace_agreement.barracuda]
}
	`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) completeExceptVMSSUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%[1]s

resource "azurerm_marketplace_agreement" "barracuda" {
  publisher = "micro-focus"
  offer     = "arcsight-logger"
  plan      = "arcsight_logger_72_byol"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest2%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 1
  compute_api_version         = "2023-09-01"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test2.id]
  }

  plan {
    name           = "arcsight_logger_72_byol"
    product        = "arcsight-logger"
    publisher      = "micro-focus"
    promotion_code = "test"
  }

  spot_priority_profile {
    allocation_strategy     = "PriceCapacityOptimized"
    eviction_policy         = "Delete"
    max_hourly_price_per_vm = -1
    min_capacity            = 0
    maintain_enabled        = false
    capacity                = 1
  }

  regular_priority_profile {
    allocation_strategy = "LowestPrice"
    capacity            = 0
    min_capacity        = 0
  }

  vm_sizes_profile {
    name = "Standard_DS1"
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "micro-focus"
      offer     = "arcsight-logger"
      sku       = "arcsight_logger_72_byol"
      version   = "7.2.0"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      caching              = "ReadWrite"
      create_option        = "FromImage"
      storage_account_type = "Standard_LRS"
    }

    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name                              = "networkProTest"
      primary_network_interface_enabled = true
      ip_configuration {
        name                             = "TestIPConfiguration"
        subnet_id                        = azurerm_subnet.test.id
        primary_ip_configuration_enabled = true
        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }
  }

  tags = {
    Hello = "ThereUpdate"
    World = "ExampleUpdate"
  }
  zones = ["1", "2", "3"]

  depends_on = [azurerm_marketplace_agreement.barracuda]
}
		`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) basicBaseLinuxVirtualMachineProfile() string {
	return `
virtual_machine_profile {
	network_api_version = "2020-11-01"
	source_image_reference {
		offer     = "0001-com-ubuntu-server-focal"
		publisher = "canonical"
		sku       = "20_04-lts-gen2"
		version   = "latest"
	}
	
	os_profile {
		linux_configuration {
			computer_name_prefix = "testvm"
			admin_username       = local.admin_username
			admin_password       = local.admin_password
			password_authentication_enabled = true
		}
	}

	network_interface {
		name                            = "networkProTest"
   	primary_network_interface_enabled 												= true
		ip_configuration {
			name      = "TestIPConfiguration"
        subnet_id = azurerm_subnet.test.id
        primary_ip_configuration_enabled   = true
        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
		}
	}
}
# ignore_changes as when os_disk block is not specified the API return default values for caching, delete_option, disk_size_in_gib and storage_account_type
lifecycle {
	ignore_changes = [virtual_machine_profile.0.os_disk]
}
`
}

func (r ComputeFleetTestResource) basicBaseWindowsVirtualMachineProfile() string {
	return `
virtual_machine_profile {
	network_api_version = "2020-11-01"
	source_image_reference {
		publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-Server-Core"
    version   = "latest"
	}

	os_profile {
		windows_configuration {
			computer_name_prefix = "testvm"
      admin_username       = local.admin_username
      admin_password       = local.admin_password
		}
	}

	network_interface {
		name                            = "networkProTest"
   	primary_network_interface_enabled 												= true
		ip_configuration {
			name      = "TestIPConfiguration"
      primary_ip_configuration_enabled   = true
      subnet_id = azurerm_subnet.test.id
      public_ip_address {
        name                    = "TestPublicIPConfiguration"
        domain_name_label       = "test-domain-label"
        idle_timeout_in_minutes = 4
      }
		}
	}
}
# ignore_changes as when os_disk block is not specified the API return default values for caching, delete_option, disk_size_in_gib and storage_account_type
lifecycle {
	ignore_changes = [virtual_machine_profile.0.os_disk]
}
`
}
