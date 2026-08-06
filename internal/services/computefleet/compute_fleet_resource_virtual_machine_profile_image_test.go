// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestAccComputeFleet_virtualMachineProfileImage_imageFromImageId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromExistingMachinePrep(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
			),
		},
		{
			Config: r.imageFromImageId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileImage_imageFromCommunitySharedImageGallery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromExistingMachinePrep(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
			),
		},
		{
			Config: r.imageFromCommunitySharedImageGallery(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileImage_imageFromCommunitySharedImageGalleryVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromExistingMachinePrep(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
			),
		},
		{
			Config: r.imageFromCommunitySharedImageGalleryVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileImage_imageFromSharedImageGallery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromExistingMachinePrep(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
			),
		},
		{
			Config: r.imageFromSharedImageGallery(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func TestAccComputeFleet_virtualMachineProfileImage_imageFromSharedImageGalleryVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.imageFromExistingMachinePrep(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.generalizeVirtualMachine(), "azurerm_linux_virtual_machine.source"),
			),
		},
		{
			Config: r.imageFromSharedImageGalleryVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password"),
	})
}

func (r ComputeFleetResource) imageFromExistingMachinePrep(data acceptance.TestData) string {
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
  sku                 = "Standard"
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
  size                            = "Standard_F1alds_v7"
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

  custom_data = base64encode(<<-EOT
    #!/bin/bash
    sudo waagent -verbose -deprovision+user -force
  EOT
  )

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "canonical"
    offer     = "ubuntu-24_04-lts"
    sku       = "server"
    version   = "latest"
  }
}
`, r.templateWithOutProvider(data), data.RandomInteger)
}

func (r ComputeFleetResource) imageFromImageId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  hyper_v_generation        = "V2"
  source_virtual_machine_id = azurerm_linux_virtual_machine.source.id
}

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  on_demand_capacity {
    target_capacity = 1
  }

  virtual_machine_sizes_profile {
    name = "Standard_D2as_v5"
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
}
`, r.imageFromExistingMachinePrep(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) imageFromSharedImageGallery(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  hyper_v_generation        = "V2"
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
  hyper_v_generation  = "V2"
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

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  on_demand_capacity {
    target_capacity = 1
  }
  virtual_machine_sizes_profile {
    name = "Standard_D2as_v5"
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
  depends_on = [azurerm_shared_image_version.test]
}
`, r.imageFromExistingMachinePrep(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) imageFromSharedImageGalleryVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`




%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  hyper_v_generation        = "V2"
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
  hyper_v_generation  = "V2"
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

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  on_demand_capacity {
    target_capacity = 1
  }

  virtual_machine_sizes_profile {
    name = "Standard_D2as_v5"
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
}
`, r.imageFromExistingMachinePrep(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) imageFromCommunitySharedImageGallery(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  hyper_v_generation        = "V2"
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
  hyper_v_generation  = "V2"
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

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  on_demand_capacity {
    target_capacity = 1
  }

  virtual_machine_sizes_profile {
    name = "Standard_D2as_v5"
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
}
`, r.imageFromExistingMachinePrep(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) imageFromCommunitySharedImageGalleryVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_image" "test" {
  name                      = "test"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  hyper_v_generation        = "V2"
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
  hyper_v_generation  = "V2"
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

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-id-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  on_demand_capacity {
    target_capacity = 1
  }

  virtual_machine_sizes_profile {
    name = "Standard_D2as_v5"
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
}
`, r.imageFromExistingMachinePrep(data), data.RandomInteger, data.Locations.Primary)
}

func (ComputeFleetResource) generalizeVirtualMachine() func(context.Context, *clients.Client, *pluginsdk.InstanceState) error {
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
