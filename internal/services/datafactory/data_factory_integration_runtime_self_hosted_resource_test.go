// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IntegrationRuntimeSelfHostedResource struct{}

func TestAccDataFactoryIntegrationRuntimeSelfHosted_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_self_hosted", "test")
	r := IntegrationRuntimeSelfHostedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeSelfHosted_selfContainedInteractiveAuthoringEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_self_hosted", "test")
	r := IntegrationRuntimeSelfHostedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.selfContainedInteractiveAuthoringEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeSelfHosted_selfContainedInteractiveAuthoringEnabledUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_self_hosted", "test")
	r := IntegrationRuntimeSelfHostedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.selfContainedInteractiveAuthoringEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeSelfHosted_rbac(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_self_hosted", "target")
	r := IntegrationRuntimeSelfHostedResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.rbac(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("rbac_authorization.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (IntegrationRuntimeSelfHostedResource) basic(data acceptance.TestData) string {
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
  name            = "acctestSIR%d"
  data_factory_id = azurerm_data_factory.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (IntegrationRuntimeSelfHostedResource) selfContainedInteractiveAuthoringEnabled(data acceptance.TestData) string {
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
  name                                         = "acctestSIR%d"
  data_factory_id                              = azurerm_data_factory.test.id
  self_contained_interactive_authoring_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (IntegrationRuntimeSelfHostedResource) rbac(data acceptance.TestData) string {
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
  address_prefixes     = ["10.0.2.0/24"]
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
  name                          = "acctvm%s"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  network_interface_ids         = [azurerm_network_interface.test.id]
  vm_size                       = "Standard_F4"
  delete_os_disk_on_termination = true

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
    "fileUris"         = ["https://raw.githubusercontent.com/Azure/azure-quickstart-templates/5661e3290f1d072195d26a5fc9d52bb372a55f48/quickstarts/microsoft.compute/vms-with-selfhost-integration-runtime/gatewayInstall.ps1"],
    "commandToExecute" = "powershell -ExecutionPolicy Unrestricted -File gatewayInstall.ps1 ${azurerm_data_factory_integration_runtime_self_hosted.host.primary_authorization_key} && timeout /t 120"
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
  name            = "acctestirshh%d"
  data_factory_id = azurerm_data_factory.host.id
}

resource "azurerm_resource_group" "target" {
  name     = "acctesttargetRG-%d"
  location = "%s"
}

resource "azurerm_role_assignment" "target" {
  scope                = azurerm_data_factory_integration_runtime_self_hosted.host.id
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
  name            = "acctestirsht%d"
  data_factory_id = azurerm_data_factory.target.id

  rbac_authorization {
    resource_id = azurerm_data_factory_integration_runtime_self_hosted.host.id
  }

  depends_on = [azurerm_role_assignment.target, azurerm_virtual_machine_extension.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (t IntegrationRuntimeSelfHostedResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IntegrationRuntimeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.IntegrationRuntimesClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
