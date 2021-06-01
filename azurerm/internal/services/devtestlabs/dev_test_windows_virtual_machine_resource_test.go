package devtestlabs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DevTestVirtualMachineResource struct {
}

func TestAccDevTestVirtualMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_windows_virtual_machine", "test")
	r := DevTestVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gallery_image_reference.0.publisher").HasValue("MicrosoftWindowsServer"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(
			// not returned from the API
			"lab_subnet_name",
			"lab_virtual_network_id",
			"password",
		),
	})
}

func TestAccDevTestVirtualMachine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_windows_virtual_machine", "test")
	r := DevTestVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_dev_test_windows_virtual_machine"),
		},
	})
}

func TestAccDevTestWindowsVirtualMachine_inboundNatRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_windows_virtual_machine", "test")
	r := DevTestVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.inboundNatRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disallow_public_ip_address").HasValue("true"),
				check.That(data.ResourceName).Key("gallery_image_reference.0.publisher").HasValue("MicrosoftWindowsServer"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Acceptance").HasValue("Test"),
			),
		},
		data.ImportStep(
			// not returned from the API
			"inbound_nat_rule",
			"lab_subnet_name",
			"lab_virtual_network_id",
			"password",
		),
	})
}

func TestAccDevTestWindowsVirtualMachine_updateStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_windows_virtual_machine", "test")
	r := DevTestVirtualMachineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storage(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gallery_image_reference.0.publisher").HasValue("MicrosoftWindowsServer"),
				check.That(data.ResourceName).Key("storage_type").HasValue("Standard"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.storage(data, "Premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gallery_image_reference.0.publisher").HasValue("MicrosoftWindowsServer"),
				check.That(data.ResourceName).Key("storage_type").HasValue("Premium"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (DevTestVirtualMachineResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	labName := id.Path["labs"]
	name := id.Path["virtualmachines"]

	resp, err := clients.DevTestLabs.VirtualMachinesClient.Get(ctx, id.ResourceGroup, labName, name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving DevTest Windows Virtual Machine %q (Lab %q / Resource Group: %q): %v", name, labName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.LabVirtualMachineProperties != nil), nil
}

func (r DevTestVirtualMachineResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "test" {
  name                   = "acctestvm%d"
  lab_name               = azurerm_dev_test_lab.test.name
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  size                   = "Standard_F2"
  username               = "acct5stU5er"
  password               = "Pa$w0rd1234!"
  lab_virtual_network_id = azurerm_dev_test_virtual_network.test.id
  lab_subnet_name        = azurerm_dev_test_virtual_network.test.subnet[0].name
  storage_type           = "Standard"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger%1000000)
}

func (r DevTestVirtualMachineResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "import" {
  name                   = azurerm_dev_test_windows_virtual_machine.test.name
  lab_name               = azurerm_dev_test_windows_virtual_machine.test.lab_name
  resource_group_name    = azurerm_dev_test_windows_virtual_machine.test.resource_group_name
  location               = azurerm_dev_test_windows_virtual_machine.test.location
  size                   = azurerm_dev_test_windows_virtual_machine.test.size
  username               = "acct5stU5er"
  password               = "Pa$w0rd1234!"
  lab_virtual_network_id = azurerm_dev_test_virtual_network.test.id
  lab_subnet_name        = azurerm_dev_test_virtual_network.test.subnet[0].name
  storage_type           = "Standard"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }
}
`, r.template(data))
}

func (r DevTestVirtualMachineResource) inboundNatRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "test" {
  name                       = "acctestvm%d"
  lab_name                   = azurerm_dev_test_lab.test.name
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  size                       = "Standard_F2"
  username                   = "acct5stU5er"
  password                   = "Pa$w0rd1234!"
  disallow_public_ip_address = true
  lab_virtual_network_id     = azurerm_dev_test_virtual_network.test.id
  lab_subnet_name            = azurerm_dev_test_virtual_network.test.subnet[0].name
  storage_type               = "Standard"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }

  inbound_nat_rule {
    protocol     = "Tcp"
    backend_port = 22
  }

  inbound_nat_rule {
    protocol     = "Tcp"
    backend_port = 3389
  }

  tags = {
    "Acceptance" = "Test"
  }
}
`, r.template(data), data.RandomInteger%1000000)
}

func (r DevTestVirtualMachineResource) storage(data acceptance.TestData, storageType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "test" {
  name                   = "acctestvm%d"
  lab_name               = azurerm_dev_test_lab.test.name
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  size                   = "Standard_B1ms"
  username               = "acct5stU5er"
  password               = "Pa$w0rd1234!"
  lab_virtual_network_id = azurerm_dev_test_virtual_network.test.id
  lab_subnet_name        = azurerm_dev_test_virtual_network.test.subnet[0].name
  storage_type           = "%s"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }
}
`, r.template(data), data.RandomInteger%1000000, storageType)
}

func (DevTestVirtualMachineResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_virtual_network" "test" {
  name                = "acctestdtvn%d"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name

  subnet {
    use_public_ip_address           = "Allow"
    use_in_virtual_machine_creation = "Allow"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
