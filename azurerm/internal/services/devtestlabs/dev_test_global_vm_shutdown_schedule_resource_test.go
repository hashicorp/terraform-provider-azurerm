package devtestlabs_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	computeParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccDevTestGlobalVMShutdownSchedule_autoShutdownBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_global_vm_shutdown_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestGlobalVMShutdownScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestGlobalVMShutdownSchedule_autoShutdownBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestGlobalVMShutdownScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "timezone", "Pacific Standard Time"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.time_in_minutes", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.webhook_url", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence_time", "0100"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDevTestGlobalVMShutdownSchedule_autoShutdownComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_global_vm_shutdown_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestGlobalVMShutdownScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestGlobalVMShutdownSchedule_autoShutdownComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestGlobalVMShutdownScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "timezone", "Central Standard Time"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.time_in_minutes", "15"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.webhook_url", "https://www.bing.com/2/4"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence_time", "1100"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Environment", "Production"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDevTestGlobalVMShutdownSchedule_autoShutdownUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_global_vm_shutdown_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevTestGlobalVMShutdownScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevTestGlobalVMShutdownSchedule_autoShutdownBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestGlobalVMShutdownScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "timezone", "Pacific Standard Time"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.time_in_minutes", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.webhook_url", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence_time", "0100"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDevTestGlobalVMShutdownSchedule_autoShutdownComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevTestGlobalVMShutdownScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "timezone", "Central Standard Time"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.time_in_minutes", "15"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.webhook_url", "https://www.bing.com/2/4"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence_time", "1100"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Environment", "Production"),
				),
			},
		},
	})
}

func testCheckDevTestGlobalVMShutdownScheduleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		targetResourceID := rs.Primary.Attributes["virtual_machine_id"]
		exists, err := testCheckDevTestGlobalVMShutdownScheduleExistsInternal(targetResourceID)
		if err != nil {
			return fmt.Errorf("Error checking if item has been created: %s", err)
		}
		if !exists {
			return fmt.Errorf("Bad: Dev Test Lab Global Schedule %q does not exist", targetResourceID)
		}

		return nil
	}
}

func testCheckDevTestGlobalVMShutdownScheduleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_global_vm_shutdown_schedule" {
			continue
		}

		targetResourceID := rs.Primary.Attributes["virtual_machine_id"]
		exists, err := testCheckDevTestGlobalVMShutdownScheduleExistsInternal(targetResourceID)
		if err != nil {
			return fmt.Errorf("Error checking if item has been destroyed: %s", err)
		}
		if exists {
			return fmt.Errorf("Bad: Dev Test Lab Global Schedule %q still exists", targetResourceID)
		}
	}

	return nil
}

func testCheckDevTestGlobalVMShutdownScheduleExistsInternal(vmID string) (bool, error) {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.GlobalLabSchedulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	vm, err := computeParse.VirtualMachineID(vmID)
	if err != nil {
		return false, fmt.Errorf("Bad: Failed to parse ID (id: %s): %+v", vmID, err)
	}

	vmName := vm.Name
	name := "shutdown-computevm-" + vmName // Auto-shutdown schedule must use this naming format for Compute VMs
	resourceGroup := vm.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return false, nil
		}
		return false, fmt.Errorf("Bad: Get on devTestLabsGlobalSchedules client (id: %s): %+v", vmID, err)
	}

	return true, nil
}

func testAccDevTestGlobalVMShutdownSchedule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dtl-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVN-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSN-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestNIC-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                  = "acctestVM-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  size                  = "Standard_B2s"

  admin_username                  = "testadmin"
  admin_password                  = "Password1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  os_disk {
    name                 = "myosdisk-%d"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccDevTestGlobalVMShutdownSchedule_autoShutdownBasic(data acceptance.TestData) string {
	template := testAccDevTestGlobalVMShutdownSchedule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_global_vm_shutdown_schedule" "test" {
  location              = azurerm_resource_group.test.location
  virtual_machine_id    = azurerm_linux_virtual_machine.test.id
  daily_recurrence_time = "0100"
  timezone              = "Pacific Standard Time"

  notification_settings {
    enabled = false
  }

  tags = {
    environment = "Production"
  }
}
`, template)
}

func testAccDevTestGlobalVMShutdownSchedule_autoShutdownComplete(data acceptance.TestData) string {
	template := testAccDevTestGlobalVMShutdownSchedule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_global_vm_shutdown_schedule" "test" {
  location           = azurerm_resource_group.test.location
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  enabled            = false

  daily_recurrence_time = "1100"
  timezone              = "Central Standard Time"

  notification_settings {
    time_in_minutes = 15
    webhook_url     = "https://www.bing.com/2/4"
    enabled         = true
  }

  tags = {
    Environment = "Production"
  }
}

`, template)
}
