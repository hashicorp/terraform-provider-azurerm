// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccComputeFleet_virtualMachineProfileDataDisk_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataDiskBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileDataDisk_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataDiskComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetTestResource) dataDiskBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      create_option   = "Empty"
      disk_size_in_gb = 10
      lun             = 0
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

  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      source_image_reference {
        publisher = "Canonical"
        offer     = "0001-com-ubuntu-server-jammy"
        sku       = "22_04-lts"
        version   = "latest"
      }

      os_disk {
        caching              = "ReadWrite"
        storage_account_type = "Standard_LRS"
      }

      data_disk {
        create_option   = "Empty"
        disk_size_in_gb = 10
        lun             = 0
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
          subnet_id                        = azurerm_subnet.linux_test.id
          primary_ip_configuration_enabled = true
          public_ip_address {
            name                    = "TestPublicIPConfiguration"
            domain_name_label       = "test-domain-label"
            idle_timeout_in_minutes = 4
          }
        }
      }
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) dataDiskComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"

  spot_priority_profile {
    min_capacity     = 0
    maintain_enabled = false
    capacity         = 0
  }

  vm_sizes_profile {
    name = "Standard_M8ms"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      create_option             = "Empty"
      disk_size_in_gb           = 10
      lun                       = 0
      caching                   = "ReadOnly"
      delete_option             = "Delete"
      disk_encryption_set_id    = azurerm_disk_encryption_set.test.id
      storage_account_type      = "Premium_LRS"
      write_accelerator_enabled = true
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
  additional_location_profile {
    location = "%[5]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      source_image_reference {
        publisher = "Canonical"
        offer     = "0001-com-ubuntu-server-jammy"
        sku       = "22_04-lts"
        version   = "latest"
      }

      os_disk {
        caching              = "ReadWrite"
        storage_account_type = "Standard_LRS"
      }

      data_disk {
        create_option             = "Empty"
        disk_size_in_gb           = 10
        lun                       = 0
        caching                   = "ReadOnly"
        delete_option             = "Delete"
        disk_encryption_set_id    = azurerm_disk_encryption_set.linux_test.id
        storage_account_type      = "Premium_LRS"
        write_accelerator_enabled = true
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
          subnet_id                        = azurerm_subnet.linux_test.id
          primary_ip_configuration_enabled = true
          public_ip_address {
            name                    = "TestPublicIPConfiguration"
            domain_name_label       = "test-domain-label"
            idle_timeout_in_minutes = 4
          }
        }
      }
    }
  }
}
`, r.diskEncryptionSetResourceDependencies(data), r.baseAndAdditionalLocationLinuxTemplateWithOutProvider(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) diskEncryptionSetResourceDependencies(data acceptance.TestData) string {
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
