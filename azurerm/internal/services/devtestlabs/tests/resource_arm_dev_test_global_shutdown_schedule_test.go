package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMDevTestLabGlobalShutdownSchedule_autoShutdownBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_global_shutdown_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDevTestLabGlobalShutdownScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestLabGlobalShutdownSchedule_autoShutdownBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabGlobalShutdownScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "time_zone_id", "Pacific Standard Time"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.status", "Disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence.0.time", "0100"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDevTestLabGlobalShutdownSchedule_autoShutdownBasicUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestLabGlobalShutdownScheduleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "Disabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "time_zone_id", "Central Standard Time"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.status", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.time_in_minutes", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "notification_settings.0.webhook_url", "https://www.bing.com/2/4"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_recurrence.0.time", "1100"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
				),
			},
		},
	})
}

func testCheckAzureRMDevTestLabGlobalShutdownScheduleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		targetResourceID := rs.Primary.Attributes["target_resource_id"]
		exists, err := testCheckAzureRMDevTestLabGlobalShutdownScheduleExistsInternal(targetResourceID)

		if err != nil {
			return fmt.Errorf("Error checking if item has been created: %s", err)
		}
		if !exists {
			return fmt.Errorf("Bad: Dev Test Lab Global Schedule %q does not exist", targetResourceID)
		}

		return nil
	}
}

func testCheckAzureRMDevTestLabGlobalShutdownScheduleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_global_shutdown_schedule" {
			continue
		}

		targetResourceID := rs.Primary.Attributes["target_resource_id"]
		exists, err := testCheckAzureRMDevTestLabGlobalShutdownScheduleExistsInternal(targetResourceID)

		if err != nil {
			return fmt.Errorf("Error checking if item has been destroyed: %s", err)
		}
		if exists {
			return fmt.Errorf("Bad: Dev Test Lab Global Schedule %q still exists", targetResourceID)
		}
	}

	return nil
}

func testCheckAzureRMDevTestLabGlobalShutdownScheduleExistsInternal(vmID string) (bool, error) {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DevTestLabs.GlobalLabSchedulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	rid, err := azure.ParseAzureResourceID(vmID)
	if err != nil {
		return false, fmt.Errorf("Bad: Failed to parse ID (id: %s): %+v", vmID, err)
	}

	vmName := rid.Path["virtualMachines"]
	name := "shutdown-computevm-" + vmName // Auto-shutdown schedule must use this naming format for Compute VMs
	resourceGroup := rid.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if resp.Response.IsHTTPStatus(404) {
			return false, nil
		}
		return false, fmt.Errorf("Bad: Get on devTestLabsGlobalSchedules client (id: %s): %+v", vmID, err)
	}

	return true, nil
}

func testAccAzureRMDevTestLabGlobalShutdownSchedule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctnic-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_B2s"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDevTestLabGlobalShutdownSchedule_autoShutdownBasic(data acceptance.TestData) string {
	template := testAccAzureRMDevTestLabGlobalShutdownSchedule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_global_shutdown_schedule" "test" {
  location           = "${azurerm_resource_group.test.location}"
  target_resource_id = "${azurerm_virtual_machine.test.id}"

  daily_recurrence {
    time = "0100"
  }

  time_zone_id = "Pacific Standard Time"

  notification_settings {
    status = "Disabled"
  }

  tags = {
    environment = "Production"
  }
}
`, template)
}

func testAccAzureRMDevTestLabGlobalShutdownSchedule_autoShutdownBasicUpdate(data acceptance.TestData) string {
	template := testAccAzureRMDevTestLabGlobalShutdownSchedule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_global_shutdown_schedule" "test" {
  location           = "${azurerm_resource_group.test.location}"
  target_resource_id = "${azurerm_virtual_machine.test.id}"
  status             = "Disabled"

  daily_recurrence {
    time = "1100"
  }

  time_zone_id = "Central Standard Time"

  notification_settings {
    time_in_minutes = 30
    webhook_url     = "https://www.bing.com/2/4"
    status          = "Enabled"
  }

  tags = {
    environment = "Production"
  }
}

`, template)
}
