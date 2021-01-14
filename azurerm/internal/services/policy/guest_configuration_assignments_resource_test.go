package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type GuestConfigurationAssignmentResource struct{}

func TestAccGuestConfigurationAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
	r := GuestConfigurationAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGuestConfigurationAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
	r := GuestConfigurationAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccGuestConfigurationAssignment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
	r := GuestConfigurationAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

//func TestAccAzureRMguestConfigurationAssignment_update(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:     func() { acceptance.PreCheck(t) },
//		Providers:    acceptance.SupportedProviders,
//		CheckDestroy: testCheckAzureRMguestConfigurationAssignmentDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: basic(data),
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
//				),
//			},
//			data.ImportStep(),
//			{
//				Config: complete(data),
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
//				),
//			},
//			data.ImportStep(),
//			{
//				Config: basic(data),
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
//				),
//			},
//			data.ImportStep(),
//		},
//	})
//}

//func TestAccAzureRMguestConfigurationAssignment_updateGuestConfiguration(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_guest_configuration_assignment", "test")
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:     func() { acceptance.PreCheck(t) },
//		Providers:    acceptance.SupportedProviders,
//		CheckDestroy: testCheckAzureRMguestConfigurationAssignmentDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: complete(data),
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
//				),
//			},
//			data.ImportStep(),
//			{
//				Config: testAccAzureRMguestConfigurationAssignment_updateGuestConfiguration(data),
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAzureRMguestConfigurationAssignmentExists(data.ResourceName),
//				),
//			},
//			data.ImportStep(),
//		},
//	})
//}

func (r GuestConfigurationAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.GuestConfigurationAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Policy.GuestConfigurationAssignmentsClient.Get(ctx, id.ResourceGroup, id.Name, id.VMName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Guest Configuration %q (Resource Group %q / Virtual Machine Name %q): %+v", id.Name, id.ResourceGroup, id.VMName, err)
	}
	return utils.Bool(true), nil
}

func (r GuestConfigurationAssignmentResource) templateBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
# note: whilst these aren't used in all tests, it saves us redefining these everywhere
locals {
  first_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  second_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r GuestConfigurationAssignmentResource) template(data acceptance.TestData) string {
	template := r.templateBase(data)
	return fmt.Sprintf(`
%s

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.first_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, template, data.RandomInteger)
}

func (r GuestConfigurationAssignmentResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_guest_configuration_assignment" "test" {
  name = "acctest-gca-%d"
location = azurerm_linux_virtual_machine.test.location
virtual_machine_id = azurerm_linux_virtual_machine.test.id
guest_configuration {
      name = "something"
      version = "1.0"
      configuration_parameter {
          name = "[InstalledApplication]bwhitelistedapp;Name"
          value = "NotePad,sql"
      }
      configuration_setting {
          reboot_if_needed = false
          action_after_reboot = "ContinueConfiguration"
      }
  }
}
`, template, data.RandomInteger)
}

func (r GuestConfigurationAssignmentResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_guest_configuration_assignment" "import" {
  name = azurerm_guest_configuration_assignment.test.name
  resource_group_name = azurerm_guest_configuration_assignment.test.resource_group_name
  location = azurerm_guest_configuration_assignment.test.location
  vm_name = azurerm_guest_configuration_assignment.test.vm_name
}
`, config)
}

func (r GuestConfigurationAssignmentResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_guest_configuration_assignment" "test" {
  name = "acctest-gca-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  vm_name = azurerm_storage_container.test.name
  context = "Azure policy"
  guest_configuration {
    name = "WhitelistedApplication"
    configuration_parameter {
      name = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }

    configuration_setting {
      action_after_reboot = "ContinueConfiguration"
      allow_module_overwrite = ""
      configuration_mode = "MonitorOnly"
      configuration_mode_frequency_mins = 15
      reboot_if_needed = "False"
      refresh_frequency_mins = 30
    }
    kind = ""
    version = "1.*"
  }
}
`, template, data.RandomInteger)
}
