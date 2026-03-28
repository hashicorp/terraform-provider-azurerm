// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccComputeFleet_virtualMachineProfileExtensions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileExtensions_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"virtual_machine_profile.0.extension.0.protected_settings_json"),
	})
}

func TestAccComputeFleet_virtualMachineProfileExtensions_protectedSettingsFromKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.extensionsProtectedSettingsFromKeyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetTestResource) extensionsBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  spot_priority_profile {
    min_capacity              = 0
    maintain_capacity_enabled = false
    capacity                  = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        password_authentication_enabled = false
        admin_ssh_keys                  = [local.first_public_key]
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
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    extension {
      name                 = "CustomScript"
      publisher            = "Microsoft.Azure.Extensions"
      type                 = "CustomScript"
      type_handler_version = "2.0"
    }
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) extensionsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 1

  spot_priority_profile {
    min_capacity              = 0
    maintain_capacity_enabled = false
    capacity                  = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        password_authentication_enabled = false
        admin_ssh_keys                  = [local.first_public_key]
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
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    extension_operations_enabled = true
    extension {
      name                               = "testOmsAgentForLinux"
      publisher                          = "Microsoft.EnterpriseCloud.Monitoring"
      type                               = "OmsAgentForLinux"
      type_handler_version               = "1.12"
      auto_upgrade_minor_version_enabled = true
      automatic_upgrade_enabled          = true
      failure_suppression_enabled        = true

      settings_json = jsonencode({
        "commandToExecute" = "echo $HOSTNAME",
        "fileUris"         = []
      })
      protected_settings_json = jsonencode({
        "commandToExecute" = "echo 'Hello World!'",
        "fileUris"         = []
      })
    }

    extension {
      name                               = "CustomScript"
      publisher                          = "Microsoft.Azure.Extensions"
      type                               = "CustomScript"
      type_handler_version               = "2.0"
      auto_upgrade_minor_version_enabled = true
    }

    extension {
      name                                      = "Docker"
      publisher                                 = "Microsoft.Azure.Extensions"
      type                                      = "DockerExtension"
      type_handler_version                      = "1.0"
      auto_upgrade_minor_version_enabled        = true
      extensions_to_provision_after_vm_creation = ["CustomScript"]
      force_extension_execution_on_change       = "test"
    }
    extensions_time_budget_duration = "PT30M"
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) extensionsProtectedSettingsFromKeyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy          = false
      purge_soft_deleted_keys_on_destroy    = false
      purge_soft_deleted_secrets_on_destroy = false
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                   = "acctestkv1%[4]s"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  tenant_id              = data.azurerm_client_config.current.tenant_id
  sku_name               = "standard"
  enabled_for_deployment = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret"
  value        = "{\"commandToExecute\":\"echo $HOSTNAME\"}"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_compute_fleet" "test" {
  name                        = "acctest-fleet-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = "%[3]s"
  platform_fault_domain_count = 2

  spot_priority_profile {
    min_capacity              = 0
    maintain_capacity_enabled = false
    capacity                  = 1
  }

  vm_sizes_profile {
    name = "Standard_D1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm"
        admin_username                  = local.admin_username
        password_authentication_enabled = false
        admin_ssh_keys                  = [local.first_public_key]
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
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
    }

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts"
      version   = "latest"
    }

    extension {
      name                 = "CustomScript"
      publisher            = "Microsoft.Azure.Extensions"
      type                 = "CustomScript"
      type_handler_version = "2.1"

      protected_settings_from_key_vault {
        secret_url      = azurerm_key_vault_secret.test.id
        source_vault_id = azurerm_key_vault.test.id
      }
    }
  }
}
`, r.templateWithOutProvider(data), data.RandomInteger, data.Locations.Primary, data.RandomString)
}
