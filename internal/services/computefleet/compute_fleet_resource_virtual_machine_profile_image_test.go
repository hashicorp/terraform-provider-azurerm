// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/ssh"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestAccComputeFleet_virtualMachineProfileImage_imageFromImageSourceReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromSourceImageReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.linux-test-source"),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileImage_imageFromImageId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromSourceImageReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.linux-test-source"),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.imageFromImageId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileImage_imageFromCommunitySharedImageGallery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromSourceImageReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.linux-test-source"),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.imageFromCommunitySharedImageGallery(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileImage_imageFromCommunitySharedImageGalleryVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromSourceImageReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.linux-test-source"),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.imageFromCommunitySharedImageGalleryVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileImage_imageFromSharedImageGallery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromSourceImageReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.linux-test-source"),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.imageFromSharedImageGallery(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileImage_imageFromSharedImageGalleryVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromSourceImageReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.linux-test-source"),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
		{
			Config: r.imageFromSharedImageGalleryVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password",
			"additional_location_profile.0.virtual_machine_profile_override.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetTestResource) imageFromExistingMachinePrep(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    virtual_machine {
      delete_os_disk_on_deletion     = true
      skip_shutdown_and_force_delete = true
    }
  }
}

%[1]s

resource "azurerm_public_ip" "public" {
  name                = "acctpip-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "acctestsourcevm-%[2]d"
  sku                 = "Basic"
}

resource "azurerm_network_interface" "public" {
  name                = "acctestnic-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.public.id
  }
}

resource "azurerm_linux_virtual_machine" "source" {
  name                            = "acctestsourceVM-%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_DS1_v2"
  admin_username                  = local.admin_username
  disable_password_authentication = false
  admin_password                  = local.admin_password

  network_interface_ids = [
    azurerm_network_interface.public.id,
  ]

  admin_ssh_key {
    username   = local.admin_username
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_public_ip" "linux-test-public" {
  name                = "acctpip-%[2]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name
  allocation_method   = "Static"
  domain_name_label   = "acctestsourcevm-%[2]d"
  sku                 = "Basic"
}

resource "azurerm_network_interface" "linux-test-public" {
  name                = "acctestnic-%[2]d"
  location            = azurerm_resource_group.linux_test.location
  resource_group_name = azurerm_resource_group.linux_test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.linux_test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.linux-test-public.id
  }
}

resource "azurerm_linux_virtual_machine" "linux-test-source" {
  name                            = "acctestsourceVM-%[2]d"
  resource_group_name             = azurerm_resource_group.linux_test.name
  location                        = azurerm_resource_group.linux_test.location
  size                            = "Standard_DS1_v2"
  admin_username                  = local.admin_username
  disable_password_authentication = false
  admin_password                  = local.admin_password

  network_interface_ids = [
    azurerm_network_interface.linux-test-public.id,
  ]

  admin_ssh_key {
    username   = local.admin_username
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, r.baseAndAdditionalLocationLinuxTemplateWithOutProvider(data), data.RandomInteger)
}

func (r ComputeFleetTestResource) imageFromSourceImageReference(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-refer-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  compute_api_version = "2024-03-01"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
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

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
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
`, r.imageFromExistingMachinePrep(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) imageFromImageId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

resource "azurerm_image" "linux_test" {
  name                      = "testlinux"
  location                  = azurerm_resource_group.linux_test.location
  resource_group_name       = azurerm_resource_group.linux_test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.linux-test-source.id
}

resource "azurerm_compute_fleet" "image_id" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_api_version = "2024-03-01"

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_id     = azurerm_image.test.id

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
      accelerated_networking_enabled    = false
      ip_forwarding_enabled             = true
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
    }
  }
  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      source_image_id     = azurerm_image.linux_test.id

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
        accelerated_networking_enabled    = false
        ip_forwarding_enabled             = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.linux_test.id]
          primary_ip_configuration_enabled       = true
          subnet_id                              = azurerm_subnet.linux_test.id
        }
      }
    }
  }
}
`, r.imageFromSourceImageReference(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) imageFromSharedImageGallery(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "1.0.1"
  gallery_name        = azurerm_shared_image.test.gallery_name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_shared_image.test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}

resource "azurerm_image" "linux_test" {
  name                      = "test"
  location                  = azurerm_resource_group.linux_test.location
  resource_group_name       = azurerm_resource_group.linux_test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.linux-test-source.id
}

resource "azurerm_shared_image_gallery" "linux_test" {
  name                = "acctestsiglinux%[2]d"
  resource_group_name = "${azurerm_resource_group.linux_test.name}"
  location            = "${azurerm_resource_group.linux_test.location}"
}

resource "azurerm_shared_image" "linux_test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.linux_test.name
  resource_group_name = azurerm_resource_group.linux_test.name
  location            = azurerm_resource_group.linux_test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "linux_test" {
  name                = "1.0.1"
  gallery_name        = azurerm_shared_image.linux_test.gallery_name
  image_name          = azurerm_shared_image.linux_test.name
  resource_group_name = azurerm_shared_image.linux_test.resource_group_name
  location            = azurerm_shared_image.linux_test.location
  managed_image_id    = azurerm_image.linux_test.id

  target_region {
    name                   = azurerm_shared_image.linux_test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}

resource "azurerm_compute_fleet" "image_id" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }
  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_id     = azurerm_shared_image.test.id

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
      accelerated_networking_enabled    = false
      ip_forwarding_enabled             = true
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
    }
  }

  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      source_image_id     = azurerm_shared_image.linux_test.id

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
        accelerated_networking_enabled    = false
        ip_forwarding_enabled             = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.linux_test.id]
          primary_ip_configuration_enabled       = true
          subnet_id                              = azurerm_subnet.linux_test.id
        }
      }
    }
  }
  depends_on = [azurerm_shared_image_version.test, azurerm_shared_image_version.linux_test]
}
`, r.imageFromSourceImageReference(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) imageFromSharedImageGalleryVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.test.gallery_name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_shared_image.test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}

resource "azurerm_image" "linux_test" {
  name                      = "test"
  location                  = azurerm_resource_group.linux_test.location
  resource_group_name       = azurerm_resource_group.linux_test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.linux-test-source.id
}

resource "azurerm_shared_image_gallery" "linux_test" {
  name                = "acctestsiglinux%[2]d"
  resource_group_name = "${azurerm_resource_group.linux_test.name}"
  location            = "${azurerm_resource_group.linux_test.location}"
}

resource "azurerm_shared_image" "linux_test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.linux_test.name
  resource_group_name = azurerm_resource_group.linux_test.name
  location            = azurerm_resource_group.linux_test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "linux_test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.linux_test.gallery_name
  image_name          = azurerm_shared_image.linux_test.name
  resource_group_name = azurerm_shared_image.linux_test.resource_group_name
  location            = azurerm_shared_image.linux_test.location
  managed_image_id    = azurerm_image.linux_test.id

  target_region {
    name                   = azurerm_shared_image.linux_test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}

resource "azurerm_compute_fleet" "image_id" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_id     = azurerm_shared_image_version.test.id

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
      accelerated_networking_enabled    = false
      ip_forwarding_enabled             = true
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
    }
  }
  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      source_image_id     = azurerm_shared_image_version.linux_test.id

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
        accelerated_networking_enabled    = false
        ip_forwarding_enabled             = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.linux_test.id]
          primary_ip_configuration_enabled       = true
          subnet_id                              = azurerm_subnet.linux_test.id
        }
      }
    }
  }
}
`, r.imageFromSourceImageReference(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) imageFromCommunitySharedImageGallery(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sharing {
    permission = "Community"
    community_gallery {
      eula            = "https://eula.net"
      prefix          = "prefix"
      publisher_email = "publisher@test.net"
      publisher_uri   = "https://publisher.net"
    }
  }
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.test.gallery_name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_shared_image.test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}

resource "azurerm_image" "linux_test" {
  name                      = "test"
  location                  = azurerm_resource_group.linux_test.location
  resource_group_name       = azurerm_resource_group.linux_test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.linux-test-source.id
}

resource "azurerm_shared_image_gallery" "linux_test" {
  name                = "acctestsiglinux%[2]d"
  resource_group_name = "${azurerm_resource_group.linux_test.name}"
  location            = "${azurerm_resource_group.linux_test.location}"

  sharing {
    permission = "Community"
    community_gallery {
      eula            = "https://eula.net"
      prefix          = "prefix"
      publisher_email = "publisher@test.net"
      publisher_uri   = "https://publisher.net"
    }
  }
}

resource "azurerm_shared_image" "linux_test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.linux_test.name
  resource_group_name = azurerm_resource_group.linux_test.name
  location            = azurerm_resource_group.linux_test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "linux_test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.linux_test.gallery_name
  image_name          = azurerm_shared_image.linux_test.name
  resource_group_name = azurerm_shared_image.linux_test.resource_group_name
  location            = azurerm_shared_image.linux_test.location
  managed_image_id    = azurerm_image.linux_test.id

  target_region {
    name                   = azurerm_shared_image.linux_test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}

resource "azurerm_compute_fleet" "image_id" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_id     = "/communityGalleries/${azurerm_shared_image_gallery.test.sharing.0.community_gallery.0.name}/images/${azurerm_shared_image_version.test.image_name}"

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
      accelerated_networking_enabled    = false
      ip_forwarding_enabled             = true
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
    }
  }

  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      source_image_id     = "/communityGalleries/${azurerm_shared_image_gallery.linux_test.sharing.0.community_gallery.0.name}/images/${azurerm_shared_image_version.linux_test.image_name}"

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
        accelerated_networking_enabled    = false
        ip_forwarding_enabled             = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.linux_test.id]
          primary_ip_configuration_enabled       = true
          subnet_id                              = azurerm_subnet.linux_test.id
        }
      }
    }
  }
}
`, r.imageFromSourceImageReference(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ComputeFleetTestResource) imageFromCommunitySharedImageGalleryVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sharing {
    permission = "Community"
    community_gallery {
      eula            = "https://eula.net"
      prefix          = "prefix"
      publisher_email = "publisher@test.net"
      publisher_uri   = "https://publisher.net"
    }
  }
}

resource "azurerm_shared_image" "test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.test.gallery_name
  image_name          = azurerm_shared_image.test.name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  managed_image_id    = azurerm_image.test.id

  target_region {
    name                   = azurerm_shared_image.test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}


resource "azurerm_image" "linux_test" {
  name                      = "test"
  location                  = azurerm_resource_group.linux_test.location
  resource_group_name       = azurerm_resource_group.linux_test.name
  source_virtual_machine_id = azurerm_linux_virtual_machine.linux-test-source.id
}

resource "azurerm_shared_image_gallery" "linux_test" {
  name                = "acctestsiglinux%[2]d"
  resource_group_name = "${azurerm_resource_group.linux_test.name}"
  location            = "${azurerm_resource_group.linux_test.location}"

  sharing {
    permission = "Community"
    community_gallery {
      eula            = "https://eula.net"
      prefix          = "prefix"
      publisher_email = "publisher@test.net"
      publisher_uri   = "https://publisher.net"
    }
  }
}

resource "azurerm_shared_image" "linux_test" {
  name                = "acctest-gallery-image"
  gallery_name        = azurerm_shared_image_gallery.linux_test.name
  resource_group_name = azurerm_resource_group.linux_test.name
  location            = azurerm_resource_group.linux_test.location
  os_type             = "Linux"

  identifier {
    publisher = "AcceptanceTest-Publisher"
    offer     = "AcceptanceTest-Offer"
    sku       = "AcceptanceTest-Sku"
  }
}

resource "azurerm_shared_image_version" "linux_test" {
  name                = "0.0.1"
  gallery_name        = azurerm_shared_image.linux_test.gallery_name
  image_name          = azurerm_shared_image.linux_test.name
  resource_group_name = azurerm_shared_image.linux_test.resource_group_name
  location            = azurerm_shared_image.linux_test.location
  managed_image_id    = azurerm_image.linux_test.id

  target_region {
    name                   = azurerm_shared_image.linux_test.location
    regional_replica_count = "5"
    storage_account_type   = "Standard_LRS"
  }
}

resource "azurerm_compute_fleet" "image_id" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  regular_priority_profile {
    capacity     = 1
    min_capacity = 1
  }

  vm_sizes_profile {
    name = "Standard_DS1_v2"
  }

  compute_api_version = "2024-03-01"
  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_id     = "/communityGalleries/${azurerm_shared_image_gallery.test.sharing.0.community_gallery.0.name}/images/${azurerm_shared_image_version.test.image_name}/versions/${azurerm_shared_image_version.test.name}"

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
      accelerated_networking_enabled    = false
      ip_forwarding_enabled             = true
      ip_configuration {
        name                                   = "ipConfigTest"
        load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
        primary_ip_configuration_enabled       = true
        subnet_id                              = azurerm_subnet.test.id
      }
    }
  }

  additional_location_profile {
    location = "%[4]s"
    virtual_machine_profile_override {
      network_api_version = "2020-11-01"
      source_image_id     = "/communityGalleries/${azurerm_shared_image_gallery.linux_test.sharing.0.community_gallery.0.name}/images/${azurerm_shared_image_version.linux_test.image_name}/versions/${azurerm_shared_image_version.linux_test.name}"

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
        accelerated_networking_enabled    = false
        ip_forwarding_enabled             = true
        ip_configuration {
          name                                   = "ipConfigTest"
          load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.linux_test.id]
          primary_ip_configuration_enabled       = true
          subnet_id                              = azurerm_subnet.linux_test.id
        }
      }
    }
  }
}
`, r.imageFromSourceImageReference(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (ComputeFleetTestResource) generalizeVirtualMachine() func(context.Context, *clients.Client, *pluginsdk.InstanceState) error {
	return func(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
		id, err := virtualmachines.ParseVirtualMachineID(state.ID)
		if err != nil {
			return err
		}

		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
			defer cancel()
		}

		// these are nested in a Set in the Legacy VM resource, simpler to compute them
		userName := "testadmin1234"
		password := "Password1234!"

		// first retrieve the Virtual Machine, since we need to find
		nicIdRaw := state.Attributes["network_interface_ids.0"]
		nicId, err := commonids.ParseNetworkInterfaceID(nicIdRaw)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Retrieving Network Interface..")
		nic, err := client.Network.NetworkInterfaces.Get(ctx, *nicId, networkinterfaces.DefaultGetOperationOptions())
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *nicId, err)
		}

		publicIpRaw := ""
		if model := nic.Model; model != nil {
			if props := model.Properties; props != nil {
				if configs := props.IPConfigurations; configs != nil {
					for _, config := range *props.IPConfigurations {
						if configProps := config.Properties; configProps != nil {
							if configProps.PublicIPAddress == nil {
								continue
							}

							if configProps.PublicIPAddress.Id == nil {
								continue
							}

							publicIpRaw = *configProps.PublicIPAddress.Id
							break
						}
					}
				}
			}
		}
		if publicIpRaw == "" {
			return fmt.Errorf("retrieving %s: could not determine Public IP Address ID", *nicId)
		}

		log.Printf("[DEBUG] Retrieving Public IP Address %q..", publicIpRaw)
		publicIpId, err := commonids.ParsePublicIPAddressID(publicIpRaw)
		if err != nil {
			return err
		}

		publicIpAddress, err := client.Network.PublicIPAddresses.Get(ctx, *publicIpId, publicipaddresses.DefaultGetOperationOptions())
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *publicIpId, err)
		}
		fqdn := ""

		if model := publicIpAddress.Model; model != nil {
			if props := model.Properties; props != nil {
				if dns := props.DnsSettings; dns != nil {
					if dns.Fqdn != nil {
						fqdn = *dns.Fqdn
					}
				}
			}
		}
		if fqdn == "" {
			return fmt.Errorf("unable to determine FQDN for %q", *publicIpId)
		}

		log.Printf("[DEBUG] Running Generalization Command..")
		sshGeneralizationCommand := ssh.Runner{
			Hostname: fqdn,
			Port:     22,
			Username: userName,
			Password: password,
			CommandsToRun: []string{
				ssh.LinuxAgentDeprovisionCommand,
			},
		}
		if err := sshGeneralizationCommand.Run(ctx); err != nil {
			return fmt.Errorf("Bad: running generalization command: %+v", err)
		}

		log.Printf("[DEBUG] Deallocating VM..")
		if err := client.Compute.VirtualMachinesClient.DeallocateThenPoll(ctx, *id, virtualmachines.DefaultDeallocateOperationOptions()); err != nil {
			return fmt.Errorf("Bad: deallocating %s: %+v", *id, err)
		}

		log.Printf("[DEBUG] Generalizing VM..")
		if _, err = client.Compute.VirtualMachinesClient.Generalize(ctx, *id); err != nil {
			return fmt.Errorf("Bad: Generalizing %s: %+v", *id, err)
		}

		return nil
	}
}
