// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccComputeFleet_virtualMachineProfileOthers_bootDiagnosticEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bootDiagnostic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOthers_capacityReservationGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capacityReservationGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOthers_galleryApplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.galleryApplication(data, "test"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOthers_licenseType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.licenseType(data, "Windows_Client"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.windows_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOthers_scheduledEvent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scheduledEvent(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOthers_UserData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userData(data, "Hello World"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOthers_additionalCapabilitiesUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.additionalCapabilitiesUltraSSD(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileOthers_additionalCapabilitiesHibernation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.additionalCapabilitiesHibernation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetTestResource) bootDiagnostic(data acceptance.TestData, enabled bool) string {
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

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version                      = "2020-11-01"
    boot_diagnostic_enabled                  = %[4]t
    boot_diagnostic_storage_account_endpoint = azurerm_storage_account.test.primary_blob_endpoint

    source_image_reference {
      publisher = "Canonical"
      offer     = "0001-com-ubuntu-server-jammy"
      sku       = "22_04-lts-gen2"
      version   = "latest"
    }

    os_disk {
      storage_account_type = "Standard_LRS"
      caching              = "ReadWrite"
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
    location = "%[6]s"
    virtual_machine_profile_override {
      network_api_version                      = "2020-11-01"
      boot_diagnostic_enabled                  = %[4]t
      boot_diagnostic_storage_account_endpoint = azurerm_storage_account.linux_test.primary_blob_endpoint

      source_image_reference {
        publisher = "Canonical"
        offer     = "0001-com-ubuntu-server-jammy"
        sku       = "22_04-lts-gen2"
        version   = "latest"
      }

      os_disk {
        storage_account_type = "Standard_LRS"
        caching              = "ReadWrite"
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

resource "azurerm_storage_account" "test" {
  name                            = "accteststr%[5]s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false
}

resource "azurerm_storage_account" "linux_test" {
  name                            = "accteststrlinux%[5]s"
  resource_group_name             = azurerm_resource_group.linux_test.name
  location                        = azurerm_resource_group.linux_test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, enabled, data.RandomString, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) capacityReservationGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zones               = ["1", "2", "3"]
}

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%[2]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  sku {
    name     = "Standard_F2"
    capacity = 2
  }
}

resource "azurerm_capacity_reservation_group" "linux_test" {
  name                = "acctest-ccrg-%[2]d"
  resource_group_name = azurerm_resource_group.linux_test.name
  location            = azurerm_resource_group.linux_test.location
  zones               = ["1", "2", "3"]
}

resource "azurerm_capacity_reservation" "linux_test" {
  name                          = "acctest-ccr-%[2]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.linux_test.id
  sku {
    name     = "Standard_F2"
    capacity = 2
  }
}

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
    name = "Standard_F2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version           = "2020-11-01"
    capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id

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
      capacity_reservation_group_id = azurerm_capacity_reservation_group.linux_test.id

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
  depends_on = [azurerm_capacity_reservation.test, azurerm_capacity_reservation.linux_test]
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) galleryApplication(data acceptance.TestData, tag string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "accteststr%[4]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test" {
  name                   = "script"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_storage_blob" "test2" {
  name                   = "script2"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%[2]d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_shared_image_gallery.test.location
  supported_os_type = "Linux"
}

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  source {
    media_link                 = azurerm_storage_blob.test.id
    default_configuration_link = azurerm_storage_blob.test.id
  }

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
    storage_account_type   = "Premium_LRS"
  }
}


resource "azurerm_storage_account" "linux_test" {
  name                     = "accteststrlinux%[4]s"
  resource_group_name      = azurerm_resource_group.linux_test.name
  location                 = azurerm_resource_group.linux_test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "linux_test" {
  name                  = "testlinux"
  storage_account_name  = azurerm_storage_account.linux_test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "linux_test" {
  name                   = "scriptlinux"
  storage_account_name   = azurerm_storage_account.linux_test.name
  storage_container_name = azurerm_storage_container.linux_test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_storage_blob" "linux_test2" {
  name                   = "script2linux"
  storage_account_name   = azurerm_storage_account.linux_test.name
  storage_container_name = azurerm_storage_container.linux_test.name
  type                   = "Page"
  size                   = 512
}

resource "azurerm_shared_image_gallery" "linux_test" {
  name                = "acctestsiglinux%[2]d"
  resource_group_name = azurerm_resource_group.linux_test.name
  location            = azurerm_resource_group.linux_test.location
}

resource "azurerm_gallery_application" "linux_test" {
  name              = "acctest-applinux-%[2]d"
  gallery_id        = azurerm_shared_image_gallery.linux_test.id
  location          = azurerm_shared_image_gallery.linux_test.location
  supported_os_type = "Linux"
}

resource "azurerm_gallery_application_version" "linux_test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.linux_test.id
  location               = azurerm_gallery_application.linux_test.location

  source {
    media_link                 = azurerm_storage_blob.linux_test.id
    default_configuration_link = azurerm_storage_blob.linux_test.id
  }

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  target_region {
    name                   = azurerm_gallery_application.linux_test.location
    regional_replica_count = 1
    storage_account_type   = "Premium_LRS"
  }
}

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

    gallery_application {
      version_id                                  = azurerm_gallery_application_version.test.id
      configuration_blob_uri                      = azurerm_storage_blob.test2.id
      order                                       = 1
      tag                                         = "%[5]s"
      automatic_upgrade_enabled                   = false
      treat_failure_as_deployment_failure_enabled = false
    }
  }

  additional_location_profile {
    location = "%[6]s"
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

      gallery_application {
        version_id                                  = azurerm_gallery_application_version.linux_test.id
        configuration_blob_uri                      = azurerm_storage_blob.linux_test2.id
        order                                       = 1
        tag                                         = "%[5]s"
        automatic_upgrade_enabled                   = false
        treat_failure_as_deployment_failure_enabled = false
      }
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, data.RandomString, tag, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) licenseType(data acceptance.TestData, lType string) string {
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

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "MicrosoftWindowsServer"
      offer     = "WindowsServer"
      sku       = "2016-Datacenter-Server-Core"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    os_profile {
      windows_configuration {
        computer_name_prefix       = "testvm"
        admin_username             = local.admin_username
        admin_password             = local.admin_password
        automatic_updates_enabled  = true
        provision_vm_agent_enabled = true
        time_zone                  = "W. Europe Standard Time"

        winrm_listener {
          protocol = "Http"
        }
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
    license_type = "%[4]s"
  }

  additional_location_profile {
    location = "%[5]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      source_image_reference {
        publisher = "MicrosoftWindowsServer"
        offer     = "WindowsServer"
        sku       = "2016-Datacenter-Server-Core"
        version   = "latest"
      }

      os_disk {
        caching              = "ReadWrite"
        storage_account_type = "Standard_LRS"
      }

      os_profile {
        windows_configuration {
          computer_name_prefix       = "testvm"
          admin_username             = local.admin_username
          admin_password             = local.admin_password
          automatic_updates_enabled  = true
          provision_vm_agent_enabled = true
          time_zone                  = "W. Europe Standard Time"

          winrm_listener {
            protocol = "Http"
          }
        }
      }

      network_interface {
        name                              = "networkProTest"
        primary_network_interface_enabled = true
        ip_configuration {
          name                             = "TestIPConfiguration"
          subnet_id                        = azurerm_subnet.windows_test.id
          primary_ip_configuration_enabled = true
          public_ip_address {
            name                    = "TestPublicIPConfiguration"
            domain_name_label       = "test-domain-label"
            idle_timeout_in_minutes = 4
          }
        }
      }
      license_type = "%[4]s"
    }
  }
}
`, r.baseAndAdditionalLocationWindowsTemplate(data), data.RandomInteger, data.Locations.Primary, lType, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) scheduledEvent(data acceptance.TestData) string {
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

    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm-%[2]d"
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

    scheduled_event_termination_timeout_duration = "PT5M"
    scheduled_event_os_image_timeout_duration    = "PT15M"
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

      os_profile {
        linux_configuration {
          computer_name_prefix            = "testvm-%[2]d"
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

      scheduled_event_termination_timeout_duration = "PT5M"
      scheduled_event_os_image_timeout_duration    = "PT15M"
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) userData(data acceptance.TestData, userDta string) string {
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

    os_profile {
      linux_configuration {
        computer_name_prefix            = "testvm-%[2]d"
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
    user_data_base64 = base64encode("%[4]s")
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

      os_profile {
        linux_configuration {
          computer_name_prefix            = "testvm-%[2]d"
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
      user_data_base64 = base64encode("%[4]s")
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, userDta, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) additionalCapabilitiesUltraSSD(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  zones               = ["1", "2", "3"]

  additional_capabilities {
    ultra_ssd_enabled = true
  }

  spot_priority_profile {
    min_capacity     = 1
    maintain_enabled = false
    capacity         = 1
  }

  vm_sizes_profile {
    name = "Standard_D2s_v3"
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
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetTestResource) additionalCapabilitiesHibernation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  additional_capabilities {
    hibernation_enabled = true
  }

  regular_priority_profile {
    allocation_strategy = "LowestPrice"
    capacity            = 1
    min_capacity        = 0
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
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
    }
  }
}
`, r.baseAndAdditionalLocationLinuxTemplate(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
