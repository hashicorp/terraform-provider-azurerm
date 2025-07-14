// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccComputeFleet_virtualMachineProfileOsDisk_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osDiskBasic(data, data.Locations.Primary, data.Locations.Secondary, "Standard_D1_v2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOsDisk_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// HardCode location is due to the limitation that VM size could be supported in two regions at the same time
			Config: r.osDiskComplete(data, "westeurope", "centralus"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetTestResource) osDiskBasic(data acceptance.TestData, primaryLocation string, secondaryLocation string, vmSize string) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 1

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "%[5]s"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
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
        primary_ip_configuration_enabled = true
        subnet_id                        = azurerm_subnet.test.id

        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }

    os_disk {}

    source_image_reference {
      publisher = "canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }
  }
}
`, r.osDiskTemplate(data, primaryLocation, secondaryLocation), data.RandomInteger, primaryLocation, secondaryLocation, vmSize)
}

func (r ComputeFleetTestResource) osDiskComplete(data acceptance.TestData, primaryLocation string, secondaryLocation string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[3]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[4]s"
  platform_fault_domain_count = 1

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DC8eds_v5"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version        = "2020-11-01"
    encryption_at_host_enabled = true
    secure_boot_enabled        = true
    vtpm_enabled               = true

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
        primary_ip_configuration_enabled = true
        subnet_id                        = azurerm_subnet.test.id

        public_ip_address {
          name                    = "TestPublicIPConfiguration"
          domain_name_label       = "test-domain-label"
          idle_timeout_in_minutes = 4
        }
      }
    }

    os_disk {
      caching                   = "ReadOnly"
      delete_option             = "Delete"
      diff_disk_option          = "Local"
      diff_disk_placement       = "ResourceDisk"
      disk_size_in_gib          = 30
      storage_account_type      = "Premium_LRS"
      security_encryption_type  = "DiskWithVMGuestState"
      disk_encryption_set_id    = azurerm_disk_encryption_set.test.id
      write_accelerator_enabled = false
    }

    source_image_reference {
      publisher = "canonical"
      offer     = "0001-com-ubuntu-confidential-vm-jammy"
      sku       = "22_04-lts-cvm"
      version   = "latest"
    }
  }

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
    "azurerm_role_assignment.linux-test-disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.linux-test-disk-encryption"
  ]
}
`, r.osDiskDiskEncryptionSetResourceDependencies(data), r.osDiskTemplateWithOutProvider(data, primaryLocation, secondaryLocation), data.RandomInteger, primaryLocation, secondaryLocation)
}

func (r ComputeFleetTestResource) osDiskTemplate(data acceptance.TestData, primaryLocation string, secondaryLocation string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

`, r.osDiskTemplateWithOutProvider(data, primaryLocation, secondaryLocation), data.RandomInteger, primaryLocation, secondaryLocation)
}

func (r ComputeFleetTestResource) osDiskTemplateWithOutProvider(data acceptance.TestData, primaryLocation string, secondaryLocation string) string {
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


resource "azurerm_resource_group" "linux_test" {
  name     = "acctest-rg-fleet-al-%[1]d"
  location = "%[3]s"
}

resource "azurerm_virtual_network" "linux_test" {
  name                = "acctvn-al-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
}

resource "azurerm_subnet" "linux_test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.linux_test.name
  virtual_network_name = azurerm_virtual_network.linux_test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "linux_test" {
  name                = "acctestpublicIP%[1]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_lb" "linux_test" {
  name                = "acctest-loadbalancer-%[1]d"
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

`, data.RandomInteger, primaryLocation, secondaryLocation)
}

func (r ComputeFleetTestResource) osDiskDiskEncryptionSetResourceDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv%[1]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  purge_protection_enabled    = true
  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Recover",
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
  name                = "acctestdes-%[2]d"
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

resource "azurerm_key_vault" "linux_test" {
  name                        = "acctestkvlinux%[1]s"
  location                    = azurerm_resource_group.linux_test.location
  resource_group_name         = azurerm_resource_group.linux_test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  purge_protection_enabled    = true
  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_access_policy" "linux-test-service-principal" {
  key_vault_id = azurerm_key_vault.linux_test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Recover",
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

resource "azurerm_key_vault_key" "linux_test" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.linux_test.id
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

  depends_on = ["azurerm_key_vault_access_policy.linux-test-service-principal"]
}

resource "azurerm_disk_encryption_set" "linux_test" {
  name                = "acctestdes-%[2]d"
  resource_group_name = azurerm_resource_group.linux_test.name
  location            = azurerm_resource_group.linux_test.location
  key_vault_key_id    = azurerm_key_vault_key.linux_test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "linux-test-disk-encryption" {
  key_vault_id = azurerm_key_vault.linux_test.id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_disk_encryption_set.linux_test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.linux_test.identity.0.principal_id
}

resource "azurerm_role_assignment" "linux-test-disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.linux_test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.linux_test.identity.0.principal_id
}
`, data.RandomString, data.RandomInteger)
}
