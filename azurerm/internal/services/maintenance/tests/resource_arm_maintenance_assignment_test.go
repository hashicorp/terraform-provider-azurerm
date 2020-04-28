package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/maintenance/mgmt/2018-06-01-preview/maintenance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMaintenanceAssignment_vmBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMaintenanceAssignmentDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccAzureRMMaintenanceAssignment_vmTemplate(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				// It may take a few minutes after starting a VM for it to become available to assign to a configuration
				// for newly created machine, wait several minutes
				PreConfig: func() { time.Sleep(3 * time.Minute) },
				Config:    testAccAzureRMMaintenanceAssignment_vmBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMaintenanceAssignmentExists(data.ResourceName),
				),
			},
			// location not returned by list rest api
			data.ImportStep("location"),
		},
	})
}

func TestAccAzureRMMaintenanceAssignment_vmRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMaintenanceAssignmentDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccAzureRMMaintenanceAssignment_vmTemplate(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				// It may take a few minutes after starting a VM for it to become available to assign to a configuration
				// for newly created machine, wait several minutes
				PreConfig: func() { time.Sleep(3 * time.Minute) },
				Config:    testAccAzureRMMaintenanceAssignment_vmBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMaintenanceAssignmentExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMaintenanceAssignment_vmRequiresImport),
		},
	})
}

func TestAccAzureRMMaintenanceAssignment_dedicatedHostBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMaintenanceAssignmentDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccAzureRMMaintenanceAssignment_dedicatedHostTemplate(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				// It may take a few minutes after starting a VM for it to become available to assign to a configuration
				// for newly created machine, wait several minutes
				PreConfig: func() { time.Sleep(3 * time.Minute) },
				Config:    testAccAzureRMMaintenanceAssignment_dedicatedHostBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMaintenanceAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep("location"),
		},
	})
}

func TestAccAzureRMMaintenanceAssignment_dedicatedHostRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_maintenance_assignment", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMaintenanceAssignmentDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccAzureRMMaintenanceAssignment_dedicatedHostTemplate(data),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				// It may take a few minutes after starting a VM for it to become available to assign to a configuration
				// for newly created machine, wait several minutes
				PreConfig: func() { time.Sleep(3 * time.Minute) },
				Config:    testAccAzureRMMaintenanceAssignment_dedicatedHostBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMaintenanceAssignmentExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMMaintenanceAssignment_dedicatedHostRequiresImport),
		},
	})
}

func testCheckAzureRMMaintenanceAssignmentDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_maintenance_assignment" {
			continue
		}

		id, err := parse.MaintenanceAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}

		var listResp maintenance.ListConfigurationAssignmentsResult
		if id.HasParentResource {
			listResp, err = conn.ListParent(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName)
		} else {
			listResp, err = conn.List(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceType, id.ResourceName)
		}

		if err != nil {
			if !utils.ResponseWasNotFound(listResp.Response) {
				return err
			}
			return nil
		}
		if listResp.Value != nil && len(*listResp.Value) > 0 {
			return fmt.Errorf("maintenance assignment to resource %q still exists", id.ResourceId)
		}

		return nil
	}

	return nil
}

func testCheckAzureRMMaintenanceAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Maintenance.ConfigurationAssignmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		id, err := parse.MaintenanceAssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}

		var listResp maintenance.ListConfigurationAssignmentsResult
		if id.HasParentResource {
			listResp, err = conn.ListParent(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName)
		} else {
			listResp, err = conn.List(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceType, id.ResourceName)
		}

		if err != nil {
			return fmt.Errorf("bad: list on ConfigurationAssignmentsClient: %+v", err)
		}
		if listResp.Value == nil || len(*listResp.Value) == 0 {
			return fmt.Errorf("could not find Maintenance Assignment to resource %q", id.ResourceId)
		}

		return nil
	}
}

func testAccAzureRMMaintenanceAssignment_vmBasic(data acceptance.TestData) string {
	template := testAccAzureRMMaintenanceAssignment_vmTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment" "test" {
  location                     = azurerm_resource_group.test.location
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  target_resource_id           = azurerm_linux_virtual_machine.test.id
}
`, template)
}

func testAccAzureRMMaintenanceAssignment_vmRequiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMaintenanceAssignment_vmBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment" "import" {
  location                     = azurerm_maintenance_assignment.test.location
  maintenance_configuration_id = azurerm_maintenance_assignment.test.maintenance_configuration_id
  target_resource_id           = azurerm_maintenance_assignment.test.target_resource_id
}
`, template)
}

func testAccAzureRMMaintenanceAssignment_dedicatedHostBasic(data acceptance.TestData) string {
	template := testAccAzureRMMaintenanceAssignment_dedicatedHostTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment" "test" {
  location                     = azurerm_resource_group.test.location
  maintenance_configuration_id = azurerm_maintenance_configuration.test.id
  target_resource_id           = azurerm_dedicated_host.test.id
}
`, template)
}

func testAccAzureRMMaintenanceAssignment_dedicatedHostRequiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMaintenanceAssignment_dedicatedHostBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_maintenance_assignment" "import" {
  location                     = azurerm_maintenance_assignment.test.location
  maintenance_configuration_id = azurerm_maintenance_assignment.test.maintenance_configuration_id
  target_resource_id           = azurerm_maintenance_assignment.test.target_resource_id
}
`, template)
}

func testAccAzureRMMaintenanceAssignment_vmTemplate(data acceptance.TestData) string {
	template := testAccAzureRMMaintenanceAssignment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[2]d"
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
  name                = "acctni-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-%[2]d"
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
`, template, data.RandomInteger)
}

func testAccAzureRMMaintenanceAssignment_dedicatedHostTemplate(data acceptance.TestData) string {
	template := testAccAzureRMMaintenanceAssignment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctest-DHG-%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%[2]d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
}
`, template, data.RandomInteger)
}

func testAccAzureRMMaintenanceAssignment_template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}
