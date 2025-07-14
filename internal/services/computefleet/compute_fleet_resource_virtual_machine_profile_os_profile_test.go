// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccComputeFleet_virtualMachineProfileOsProfile_linuxBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osProfileLinuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOsProfile_linuxComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osProfileLinuxComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"virtual_machine_profile.0.os_profile.0.custom_data_base64"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOsProfile_windowsBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osProfileWindowsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOsProfile_windowsComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osProfileWindowsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password",
			"virtual_machine_profile.0.os_profile.0.custom_data_base64",
			"virtual_machine_profile.0.os_profile.0.windows_configuration.0.additional_unattend_content.0.xml"),
	})
}

func (r ComputeFleetTestResource) osProfileWindowsBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_certificates_on_destroy = false
    }
  }
}

%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  compute_api_version = "2024-03-01"

  spot_priority_profile {
    min_capacity     = 1
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  %[4]s

}
`, r.baseAndAdditionalLocationWindowsTemplateWithOutProvider(data), data.RandomInteger, data.Locations.Primary, r.basicBaseWindowsVirtualMachineProfile(), data.Locations.Secondary)
}

func (r ComputeFleetTestResource) osProfileWindowsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s
resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"

  spot_priority_profile {
    min_capacity     = 1
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "MicrosoftWindowsServer"
      offer     = "WindowsServer"
      sku       = "2022-datacenter-azure-edition-core"
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
        additional_unattend_content {
          setting = "FirstLogonCommands"
          xml     = "<FirstLogonCommands><SynchronousCommand><CommandLine>shutdown /r /t 0 /c \"initial reboot\"</CommandLine><Description>reboot</Description><Order>1</Order></SynchronousCommand></FirstLogonCommands>"
        }
        automatic_updates_enabled         = true
        vm_agent_platform_updates_enabled = true

        patch_mode                            = "AutomaticByPlatform"
        hot_patching_enabled                  = true
        bypass_platform_safety_checks_enabled = false
        patch_rebooting                       = "Always"
        provision_vm_agent_enabled            = true
        secret {
          key_vault_id = azurerm_key_vault.test.id
          certificate {
            url   = azurerm_key_vault_certificate.first.secret_id
            store = "My"
          }
        }

        time_zone = "W. Europe Standard Time"

        winrm_listener {
          certificate_url = azurerm_key_vault_certificate.first.secret_id
          protocol        = "Https"
        }
      }
    }

    extension {
      name                               = "HealthExtension"
      publisher                          = "Microsoft.ManagedServices"
      type                               = "ApplicationHealthWindows"
      type_handler_version               = "1.0"
      auto_upgrade_minor_version_enabled = true

      settings_json = jsonencode({
        "protocol"    = "http"
        "port"        = 80
        "requestPath" = "/healthEndpoint"
      })
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
  }
}
`, r.baseAndAdditionalLocationWindowsTemplateWithOutProvider(data), r.secretWindowsResourceDependencies(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) osProfileLinuxBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_certificates_on_destroy = false
    }
  }
}

%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  compute_api_version = "2024-03-01"

  spot_priority_profile {
    min_capacity     = 1
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  %[4]s
}
`, r.baseAndAdditionalLocationLinuxTemplateWithOutProvider(data), data.RandomInteger, data.Locations.Primary, r.basicBaseLinuxVirtualMachineProfile(), data.Locations.Secondary)
}

func (r ComputeFleetTestResource) osProfileLinuxComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[4]s"

  spot_priority_profile {
    min_capacity     = 1
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      offer     = "0001-com-ubuntu-server-focal"
      publisher = "canonical"
      sku       = "20_04-lts-gen2"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    os_profile {
      custom_data_base64 = "Y3VzdG9tIGRhdGEh"
      linux_configuration {
        computer_name_prefix                  = "testvm"
        admin_username                        = local.admin_username
        admin_password                        = local.admin_password
        password_authentication_enabled       = true
        admin_ssh_keys                        = [local.first_public_key]
        provision_vm_agent_enabled            = true
        vm_agent_platform_updates_enabled     = true
        patch_mode                            = "AutomaticByPlatform"
        bypass_platform_safety_checks_enabled = true
        patch_rebooting                       = "Always"
        secret {
          key_vault_id = azurerm_key_vault.test.id
          certificate {
            url = azurerm_key_vault_certificate.first.secret_id
          }
        }
      }
    }
    extension {
      name                 = "HealthExtension"
      publisher            = "Microsoft.ManagedServices"
      type                 = "ApplicationHealthLinux"
      type_handler_version = "1.0"

      settings_json = jsonencode({
        "protocol"    = "http"
        "port"        = 80
        "requestPath" = "/healthEndpoint"
      })
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
}
`, r.baseAndAdditionalLocationLinuxTemplateWithOutProvider(data), r.secretLinuxResourceDependencies(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) secretLinuxResourceDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy               = false
      purge_soft_deleted_certificates_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%[1]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                        = "standard"
  enabled_for_template_deployment = true
  enabled_for_deployment          = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Update",
    ]

    key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]

    secret_permissions = [
      "Set",
    ]

    storage_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-first"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_certificate" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-second"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault" "linux_test" {
  name                = "acctestkvlinux%[1]s"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                        = "standard"
  enabled_for_template_deployment = true
  enabled_for_deployment          = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Update",
    ]

    key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]

    secret_permissions = [
      "Set",
    ]

    storage_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "linux_test_first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.linux_test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-first"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_certificate" "linux_test_second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.linux_test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-second"
      validity_in_months = 12
    }
  }
}
`, data.RandomString)
}

func (r ComputeFleetTestResource) secretWindowsResourceDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy               = false
      purge_soft_deleted_certificates_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%[1]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                        = "standard"
  enabled_for_template_deployment = true
  enabled_for_deployment          = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Update",
    ]

    key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]

    secret_permissions = [
      "Set",
    ]

    storage_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "first" {
  name         = "first"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-first"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_certificate" "second" {
  name         = "second"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-second"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault" "windows_test" {
  name                = "acctestkvwin%[1]s"
  location            = azurerm_resource_group.windows_test.location
  resource_group_name = azurerm_resource_group.windows_test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                        = "standard"
  enabled_for_template_deployment = true
  enabled_for_deployment          = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Update",
    ]

    key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]

    secret_permissions = [
      "Set",
    ]

    storage_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "windows_test_first" {
  name         = "firstwin"
  key_vault_id = azurerm_key_vault.windows_test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-firstwin"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_certificate" "windows_test_second" {
  name         = "secondwin"
  key_vault_id = azurerm_key_vault.windows_test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=hello-world-secondwin"
      validity_in_months = 12
    }
  }
}
`, data.RandomString)
}
