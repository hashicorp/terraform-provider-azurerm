package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDevTestVirtualMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestWindowsVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestWindowsVirtualMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gallery_image_reference.0.publisher", "MicrosoftWindowsServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(
				// not returned from the API
				"lab_subnet_name",
				"lab_virtual_network_id",
				"password",
			),
		},
	})
}

func TestAccAzureRMDevTestVirtualMachine_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_dev_test_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestWindowsVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestWindowsVirtualMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDevTestWindowsVirtualMachine_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dev_test_windows_virtual_machine"),
			},
		},
	})
}

func TestAccAzureRMDevTestWindowsVirtualMachine_inboundNatRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestWindowsVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestWindowsVirtualMachine_inboundNatRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "disallow_public_ip_address", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "gallery_image_reference.0.publisher", "MicrosoftWindowsServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Acceptance", "Test"),
				),
			},
			data.ImportStep(
				// not returned from the API
				"inbound_nat_rule",
				"lab_subnet_name",
				"lab_virtual_network_id",
				"password",
			),
		},
	})
}

func TestAccAzureRMDevTestWindowsVirtualMachine_updateStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestWindowsVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestWindowsVirtualMachine_storage(data, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gallery_image_reference.0.publisher", "MicrosoftWindowsServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_type", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMDevTestWindowsVirtualMachine_storage(data, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestWindowsVirtualMachineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "gallery_image_reference.0.publisher", "MicrosoftWindowsServer"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_type", "Premium"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testCheckAzureRMDevTestWindowsVirtualMachineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.VirtualMachinesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualMachineName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualMachineName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get devTestVirtualMachinesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevTest Windows Virtual Machine %q (Lab %q / Resource Group: %q) does not exist", virtualMachineName, labName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDevTestWindowsVirtualMachineDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.VirtualMachinesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_windows_virtual_machine" {
			continue
		}

		virtualMachineName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualMachineName, "")

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DevTest Windows Virtual Machine still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDevTestWindowsVirtualMachine_basic(data acceptance.TestData) string {
	template := testAccAzureRMDevTestWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "test" {
  name                   = "acctestvm%d"
  lab_name               = "${azurerm_dev_test_lab.test.name}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  location               = "${azurerm_resource_group.test.location}"
  size                   = "Standard_F2"
  username               = "acct5stU5er"
  password               = "Pa$$w0rd1234!"
  lab_virtual_network_id = "${azurerm_dev_test_virtual_network.test.id}"
  lab_subnet_name        = "${azurerm_dev_test_virtual_network.test.subnet.0.name}"
  storage_type           = "Standard"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }
}
`, template, data.RandomInteger%1000000)
}

func testAccAzureRMDevTestWindowsVirtualMachine_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDevTestWindowsVirtualMachine_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "import" {
  name                   = "${azurerm_dev_test_windows_virtual_machine.test.name}"
  lab_name               = "${azurerm_dev_test_windows_virtual_machine.test.lab_name}"
  resource_group_name    = "${azurerm_dev_test_windows_virtual_machine.test.resource_group_name}"
  location               = "${azurerm_dev_test_windows_virtual_machine.test.location}"
  size                   = "${azurerm_dev_test_windows_virtual_machine.test.size}"
  username               = "acct5stU5er"
  password               = "Pa$$w0rd1234!"
  lab_virtual_network_id = "${azurerm_dev_test_virtual_network.test.id}"
  lab_subnet_name        = "${azurerm_dev_test_virtual_network.test.subnet.0.name}"
  storage_type           = "Standard"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testAccAzureRMDevTestWindowsVirtualMachine_inboundNatRules(data acceptance.TestData) string {
	template := testAccAzureRMDevTestWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "test" {
  name                       = "acctestvm%d"
  lab_name                   = "${azurerm_dev_test_lab.test.name}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
  location                   = "${azurerm_resource_group.test.location}"
  size                       = "Standard_F2"
  username                   = "acct5stU5er"
  password                   = "Pa$$w0rd1234!"
  disallow_public_ip_address = true
  lab_virtual_network_id     = "${azurerm_dev_test_virtual_network.test.id}"
  lab_subnet_name            = "${azurerm_dev_test_virtual_network.test.subnet.0.name}"
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
`, template, data.RandomInteger%1000000)
}

func testAccAzureRMDevTestWindowsVirtualMachine_storage(data acceptance.TestData, storageType string) string {
	template := testAccAzureRMDevTestWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_windows_virtual_machine" "test" {
  name                   = "acctestvm%d"
  lab_name               = "${azurerm_dev_test_lab.test.name}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  location               = "${azurerm_resource_group.test.location}"
  size                   = "Standard_B1ms"
  username               = "acct5stU5er"
  password               = "Pa$$w0rd1234!"
  lab_virtual_network_id = "${azurerm_dev_test_virtual_network.test.id}"
  lab_subnet_name        = "${azurerm_dev_test_virtual_network.test.subnet.0.name}"
  storage_type           = "%s"

  gallery_image_reference {
    offer     = "WindowsServer"
    publisher = "MicrosoftWindowsServer"
    sku       = "2012-Datacenter"
    version   = "latest"
  }
}
`, template, data.RandomInteger%1000000, storageType)
}

func testAccAzureRMDevTestWindowsVirtualMachine_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dev_test_virtual_network" "test" {
  name                = "acctestdtvn%d"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  subnet {
    use_public_ip_address           = "Allow"
    use_in_virtual_machine_creation = "Allow"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
