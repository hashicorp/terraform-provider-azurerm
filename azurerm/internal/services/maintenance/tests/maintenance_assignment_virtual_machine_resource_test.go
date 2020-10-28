package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMaintenanceAssignmentVirtualMachine_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMaintenanceAssignmentVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMaintenanceAssignmentVirtualMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMaintenanceAssignmentVirtualMachineExists(data.ResourceName),
				),
			},
			// location not returned by list rest api
			data.ImportStep("location"),
		},
	})
}

func TestAccAzureRMMaintenanceAssignmentVirtualMachine_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMaintenanceAssignmentVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMaintenanceAssignmentVirtualMachine_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMaintenanceAssignmentVirtualMachineExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMaintenanceAssignmentVirtualMachine_requiresImport),
		},
	})
}

func testCheckAzureRMMaintenanceAssignmentVirtualMachineDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_maintenance_assignment_virtual_machine" {
			continue
		}

		id, err := parse.MaintenanceAssignmentVirtualMachineID(rs.Primary.ID)
		if err != nil {
			return err
		}

		listResp, err := conn.List(ctx, id.VirtualMachineId.ResourceGroup, "Microsoft.Compute", "virtualMachines", id.VirtualMachineId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(listResp.Response) {
				return err
			}
			return nil
		}
		if listResp.Value != nil && len(*listResp.Value) > 0 {
			return fmt.Errorf("maintenance assignment (Virtual Machine id: %q) still exists", id.VirtualMachineIdRaw)
		}

		return nil
	}

	return nil
}

func testCheckAzureRMMaintenanceAssignmentVirtualMachineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Maintenance.ConfigurationAssignmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		id, err := parse.MaintenanceAssignmentVirtualMachineID(rs.Primary.ID)
		if err != nil {
			return err
		}

		listResp, err := conn.List(ctx, id.VirtualMachineId.ResourceGroup, "Microsoft.Compute", "virtualMachines", id.VirtualMachineId.Name)
		if err != nil {
			return fmt.Errorf("bad: list on ConfigurationAssignmentsClient: %+v", err)
		}
		if listResp.Value == nil || len(*listResp.Value) == 0 {
			return fmt.Errorf("could not find Maintenance Assignment (Virtual Machine id: %q)", id.VirtualMachineIdRaw)
		}

		return nil
	}
}

func testAccAzureRMMaintenanceAssignmentVirtualMachine_basic(data acceptance.TestData) string {
	template := testAccAzureRMMaintenanceAssignmentVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_virtual_machine" "test" {
  location                     = azurerm_resource_group.test.location
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  virtual_machine_id           = azurerm_linux_virtual_machine.test.id
}
`, template)
}

func testAccAzureRMMaintenanceAssignmentVirtualMachine_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMaintenanceAssignmentVirtualMachine_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment_virtual_machine" "import" {
  location                     = azurerm_maintenance_assignment_virtual_machine.test.location
  maintenance_configuration_id = azurerm_maintenance_assignment_virtual_machine.test.maintenance_configuration_id
  virtual_machine_id           = azurerm_maintenance_assignment_virtual_machine.test.virtual_machine_id
}
`, template)
}

func testAccAzureRMMaintenanceAssignmentVirtualMachine_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maint-%[1]d"
  location = "%[2]s"
}

resource "azurerm_maintenance_configuration" "test" {
  name                = "acctest-MC%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope               = "All"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "internal"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_D15_v2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"

  disable_password_authentication = false

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

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
`, data.RandomInteger, data.Locations.Primary)
}
