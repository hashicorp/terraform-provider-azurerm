package automanage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AutomanageConfigurationProfileAssignmentResource struct{}

func TestAccAutomanageConfigurationProfileAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile_assignment", "test")
	r := AutomanageConfigurationProfileAssignmentResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAutomanageConfigurationProfileAssignment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile_assignment", "test")
	r := AutomanageConfigurationProfileAssignmentResource{}
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

func (r AutomanageConfigurationProfileAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AutomanageConfigurationProfileAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automanage.ConfigurationProfileAssignmentClient.Get(ctx, id.ResourceGroup, id.ConfigurationProfileAssignmentName, id.VirtualMachineName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q): %+v", id.ConfigurationProfileAssignmentName, id.ResourceGroup, id.VirtualMachineName, err)
	}
	return utils.Bool(true), nil
}

func (r AutomanageConfigurationProfileAssignmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-automanage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VN-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-sub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctnicsource-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsads%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "accttest-sc-%d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctest-vm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.testsource.id]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "winhost01"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
    timezone = "Pacific Standard Time"
  }

  tags = {
    "azsecpack"                                                                = "nonprod"
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r AutomanageConfigurationProfileAssignmentResource) requiresImport(data acceptance.TestData) string {
	config := r.complete(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profile_assignment" "import" {
  name                  = azurerm_automanage_configuration_profile_assignment.test.name
  resource_group_name   = azurerm_automanage_configuration_profile_assignment.test.resource_group_name
  vm_name               = azurerm_automanage_configuration_profile_assignment.test.vm_name
  configuration_profile = azurerm_automanage_configuration_profile_assignment.test.configuration_profile
}
`, config)
}

func (r AutomanageConfigurationProfileAssignmentResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profile_assignment" "test" {
  name                  = "default"
  resource_group_name   = azurerm_resource_group.test.name
  vm_name               = azurerm_virtual_machine.test.name
  configuration_profile = "/providers/Microsoft.Automanage/bestPractices/AzureBestPracticesProduction"
}
`, template)
}
