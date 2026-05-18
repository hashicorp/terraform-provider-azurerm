// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ComputeFleetResource struct{}

func TestAccComputeFleet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}
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
	r := ComputeFleetResource{}
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
	r := ComputeFleetResource{}
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
	r := ComputeFleetResource{}
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

// Generated resource identity test cannot be used due to non-empty plan caused by `virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password`
func TestAccComputeFleet_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")
	r := ComputeFleetResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"name":                {},
		"resource_group_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_compute_fleet.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_compute_fleet.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_compute_fleet.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_compute_fleet.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		{
			ResourceName: data.ResourceName,
			// `ImportStateVerify` and `ImportStateVerifyIgnore` cannot be used on `virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password` in plannable import tests, hence `ExpectNonEmptyPlan` is used instead
			ExpectNonEmptyPlan: true,
			ImportState:        true,
			ImportStateKind:    resource.ImportBlockWithResourceIdentity,
		},
		{
			ResourceName:       data.ResourceName,
			ExpectNonEmptyPlan: true,
			ImportState:        true,
			ImportStateKind:    resource.ImportBlockWithID,
		},
	}, false)
}

func (r ComputeFleetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r ComputeFleetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

`, r.templateWithOutProvider(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) templateWithOutProvider(data acceptance.TestData) string {
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

func (r ComputeFleetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_compute_fleet" "test" {
  name                = "acctest-fleet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  spot_capacity {
    maintain_capacity_enabled = false
    target_capacity           = 1
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1alds_v7"
  }

  virtual_machine_profile {
    source_image_reference {
      offer     = "ubuntu-24_04-lts"
      publisher = "canonical"
      sku       = "server"
      version   = "latest"
    }

    os_profile {
      linux_configuration {
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name = "networkProTest"

      ip_configuration {
        name      = "TestIPConfiguration"
        subnet_id = azurerm_subnet.test.id
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_compute_fleet" "import" {
  name                = azurerm_compute_fleet.test.name
  resource_group_name = azurerm_compute_fleet.test.resource_group_name
  location            = azurerm_compute_fleet.test.location

  spot_capacity {
    maintain_capacity_enabled = false
    target_capacity           = 1
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1alds_v7"
  }

  virtual_machine_profile {
    source_image_reference {
      offer     = "ubuntu-24_04-lts"
      publisher = "canonical"
      sku       = "server"
      version   = "latest"
    }

    os_profile {
      linux_configuration {
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name = "networkProTest"

      ip_configuration {
        name      = "TestIPConfiguration"
        subnet_id = azurerm_subnet.test.id
      }
    }
  }
}
`, r.basic(data))
}

func (r ComputeFleetResource) completeExceptVMSS(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%[1]s

resource "azurerm_marketplace_agreement" "barracuda" {
  publisher = "nvidia"
  offer     = "ngc_azure_17_11"
  plan      = "ngc-base-version-25_9_1_gen2"
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
    name           = "ngc-base-version-25_9_1_gen2"
    product        = "ngc_azure_17_11"
    publisher      = "nvidia"
    promotion_code = "test"
  }

  spot_capacity {
    allocation_strategy                          = "CapacityOptimized"
    eviction_policy                              = "Delete"
    maximum_hourly_price_per_virtual_machine_usd = -1
    minimum_capacity                             = 0
    maintain_capacity_enabled                    = true
    target_capacity                              = 1
  }

  on_demand_capacity {
    allocation_strategy       = "Prioritized"
    target_capacity           = 1
    minimum_starting_capacity = 0
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1alds_v7"
    rank = 0
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1as_v7"
    rank = 1
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1als_v7"
    rank = 2
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "nvidia"
      offer     = "ngc_azure_17_11"
      sku       = "ngc-base-version-25_9_1_gen2"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      caching              = "ReadWrite"
      create_option        = "Empty"
      disk_size_in_gib     = 10
      lun                  = 42
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
}
	`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) completeExceptVMSSUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%[1]s

resource "azurerm_marketplace_agreement" "barracuda" {
  publisher = "nvidia"
  offer     = "ngc_azure_17_11"
  plan      = "ngc-base-version-25_9_1_gen2"
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
    identity_ids = [azurerm_user_assigned_identity.test2.id]
  }

  plan {
    name           = "ngc-base-version-25_9_1_gen2"
    product        = "ngc_azure_17_11"
    publisher      = "nvidia"
    promotion_code = "test"
  }

  spot_capacity {
    allocation_strategy                          = "CapacityOptimized"
    eviction_policy                              = "Delete"
    maximum_hourly_price_per_virtual_machine_usd = -1
    minimum_capacity                             = 0
    maintain_capacity_enabled                    = true
    target_capacity                              = 0
  }

  on_demand_capacity {
    allocation_strategy       = "Prioritized"
    target_capacity           = 2
    minimum_starting_capacity = 0
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1als_v7"
    rank = 0
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1as_v7"
    rank = 1
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1alds_v7"
    rank = 2
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1ads_v7"
    rank = 3
  }

  virtual_machine_profile {
    network_api_version = "2020-11-01"
    source_image_reference {
      publisher = "nvidia"
      offer     = "ngc_azure_17_11"
      sku       = "ngc-base-version-25_9_1_gen2"
      version   = "latest"
    }

    os_disk {
      caching              = "ReadWrite"
      storage_account_type = "Standard_LRS"
    }

    data_disk {
      caching              = "ReadWrite"
      create_option        = "Empty"
      disk_size_in_gib     = 10
      lun                  = 42
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
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r ComputeFleetResource) basicBaseLinuxVirtualMachineProfile() string {
	return `
virtual_machine_profile {
	network_api_version = "2020-11-01"
	source_image_reference {
		offer     = "ubuntu-24_04-lts"
		publisher = "canonical"
		sku       = "server"
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
		name                               = "networkProTest"
   	    primary_network_interface_enabled  = true
		ip_configuration {
		name                               = "TestIPConfiguration"
        subnet_id                          = azurerm_subnet.test.id
        primary_ip_configuration_enabled   = true
        public_ip_address {
            name                    = "TestPublicIPConfiguration"
            domain_name_label       = "test-domain-label"
            idle_timeout_in_minutes = 4
          }
		}
	}
}
`
}

func (r ComputeFleetResource) basicBaseWindowsVirtualMachineProfile() string {
	return `
virtual_machine_profile {
	network_api_version = "2020-11-01"
	source_image_reference {
      publisher = "MicrosoftWindowsServer"
      offer     = "WindowsServer"
      sku       = "2025-datacenter-core-g2"
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
   	    primary_network_interface_enabled 	 = true
		ip_configuration {
          name                               = "TestIPConfiguration"
          primary_ip_configuration_enabled   = true
          subnet_id                          = azurerm_subnet.test.id
          public_ip_address {
            name                    = "TestPublicIPConfiguration"
            domain_name_label       = "test-domain-label"
            idle_timeout_in_minutes = 4
          }
		}
	}
}
`
}
