package datafactory_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataFactoryIntegrationRuntimeSelfHosted_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_self_hosted", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryIntegrationRuntimeSelfHostedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryIntegrationRuntimeSelfHosted_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryIntegrationRuntimeSelfHostedExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryIntegrationRuntimeSelfHosted_rbac(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_self_hosted", "target")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryIntegrationRuntimeSelfHostedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryIntegrationRuntimeSelfHosted_rbac(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryIntegrationRuntimeSelfHostedExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "rbac_authorization.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMDataFactoryIntegrationRuntimeSelfHosted_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirsh%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_self_hosted" "test" {
  name                = "acctestSIR%d"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryIntegrationRuntimeSelfHosted_rbac(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
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

resource "azurerm_public_ip" "test" {
  name                = "acctpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm%s"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F4"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "acctvm%s"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
    timezone           = "Pacific Standard Time"
    provision_vm_agent = true
  }
}

resource "azurerm_virtual_machine_extension" "test" {
  name                 = "acctestExt-%d"
  virtual_machine_id   = azurerm_virtual_machine.test.id
  publisher            = "Microsoft.Compute"
  type                 = "CustomScriptExtension"
  type_handler_version = "1.10"
  settings = jsonencode({
    "fileUris"         = ["https://raw.githubusercontent.com/Azure/azure-quickstart-templates/00b79d2102c88b56502a63041936ef4dd62cf725/101-vms-with-selfhost-integration-runtime/gatewayInstall.ps1"],
    "commandToExecute" = "powershell -ExecutionPolicy Unrestricted -File gatewayInstall.ps1 ${azurerm_data_factory_integration_runtime_self_hosted.host.auth_key_1} && timeout /t 120"
  })
}

resource "azurerm_resource_group" "host" {
  name     = "acctesthostRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "host" {
  name                = "acctestdfirshh%d"
  location            = azurerm_resource_group.host.location
  resource_group_name = azurerm_resource_group.host.name
}

resource "azurerm_data_factory_integration_runtime_self_hosted" "host" {
  name                = "acctestirshh%d"
  data_factory_name   = azurerm_data_factory.host.name
  resource_group_name = azurerm_resource_group.host.name
}

resource "azurerm_resource_group" "target" {
  name     = "acctesttargetRG-%d"
  location = "%s"
}

resource "azurerm_role_assignment" "target" {
  scope                = azurerm_data_factory.host.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_data_factory.target.identity[0].principal_id
}

resource "azurerm_data_factory" "target" {
  name                = "acctestdfirsht%d"
  location            = azurerm_resource_group.target.location
  resource_group_name = azurerm_resource_group.target.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_factory_integration_runtime_self_hosted" "target" {
  name                = "acctestirsht%d"
  data_factory_name   = azurerm_data_factory.target.name
  resource_group_name = azurerm_resource_group.target.name

  rbac_authorization {
    resource_id = azurerm_data_factory_integration_runtime_self_hosted.host.id
  }

  depends_on = [azurerm_role_assignment.target, azurerm_virtual_machine_extension.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testCheckAzureRMDataFactoryIntegrationRuntimeSelfHostedExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.IntegrationRuntimesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		id, err := parse.IntegrationRuntimeID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on IntegrationRuntimesClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Data Factory Self-hosted Integration Runtime %q (Resource Group: %q, Data Factory %q) does not exist", id.Name, id.FactoryName, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryIntegrationRuntimeSelfHostedDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_integration_self_hosted" {
			continue
		}

		id, err := parse.IntegrationRuntimeID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on IntegrationRuntimesClient: %+v", err)
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Data Factory Self-hosted Integration Runtime still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}
