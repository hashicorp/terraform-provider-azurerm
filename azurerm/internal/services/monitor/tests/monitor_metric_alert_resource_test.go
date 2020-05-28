package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMMonitorMetricAlert_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorMetricAlert_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorMetricAlert_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_monitor_metric_alert"),
			},
		},
	})
}

func TestAccAzureRMMonitorMetricAlert_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorMetricAlert_basicAndCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMMonitorMetricAlert_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMMonitorMetricAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorMetricAlert_multiScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_multiScope(data, "azurerm_linux_virtual_machine.test1.id", "azurerm_linux_virtual_machine.test2.id"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorMetricAlert_multiScope(data, "azurerm_linux_virtual_machine.test1.id"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorMetricAlert_multiScope(data, "azurerm_linux_virtual_machine.test1.id", "azurerm_linux_virtual_machine.test2.id"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMonitorMetricAlertDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.MetricAlertsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_metric_alert" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Metric alert still exists:\n%#v", resp)
		}
	}
	return nil
}

func testCheckAzureRMMonitorMetricAlertExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.MetricAlertsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Metric Alert Instance: %s", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on monitorMetricAlertsClient: %+v", err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Metric Alert Instance %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMMonitorMetricAlert_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_storage_account.test.id]

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "UsedCapacity"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = 55.5
  }

  window_size = "PT1H"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMMonitorMetricAlert_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMonitorMetricAlert_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_metric_alert" "import" {
  name                = azurerm_monitor_metric_alert.test.name
  resource_group_name = azurerm_monitor_metric_alert.test.resource_group_name
  scopes              = azurerm_monitor_metric_alert.test.scopes

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "UsedCapacity"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = 55.5
  }
  window_size = "PT1H"
}
`, template)
}

func testAccAzureRMMonitorMetricAlert_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa1%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_action_group" "test1" {
  name                = "acctestActionGroup1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag1"
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag2"
}

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_storage_account.test.id]
  enabled             = true
  auto_mitigate       = false
  severity            = 4
  description         = "This is a complete metric alert resource."
  frequency           = "PT30M"
  window_size         = "PT12H"

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "Transactions"
    aggregation      = "Total"
    operator         = "GreaterThan"
    threshold        = 99

    dimension {
      name     = "GeoType"
      operator = "Include"
      values   = ["*"]
    }
  }

  action {
    action_group_id = azurerm_monitor_action_group.test1.id
  }

  action {
    action_group_id = azurerm_monitor_action_group.test2.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMMonitorMetricAlert_vmTemplate(data acceptance.TestData, n int) string {
	return fmt.Sprintf(`
resource "azurerm_subnet" "test%[2]d" {
  name                 = "internal%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.%[2]d.0/24"
}

resource "azurerm_network_interface" "test%[2]d" {
  name                = "acctestnic%[2]d-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test%d.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test%[2]d" {
  name                            = "acctestVM%[2]d-%[1]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@$$w0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test%d.id,
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
`, data.RandomInteger, n, n, n)
}

func testAccAzureRMMonitorMetricAlert_multiVMTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

%s

%s
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger,
		testAccAzureRMMonitorMetricAlert_vmTemplate(data, 1),
		testAccAzureRMMonitorMetricAlert_vmTemplate(data, 2),
	)
}

func testAccAzureRMMonitorMetricAlert_multiScope(data acceptance.TestData, scopes ...string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes = [
	%s
  ]
  dynamic_criteria {
    metric_namespace = "Microsoft.Compute/virtualMachines"
    metric_name      = "Network In"
    aggregation      = "Total"

    operator          = "GreaterOrLessThan"
    alert_sensitivity = "Medium"
  }
  window_size              = "PT1H"
  frequency                = "PT5M"
  target_resource_type     = "Microsoft.Compute/virtualMachines"
  target_resource_location = "%s"
}
`, testAccAzureRMMonitorMetricAlert_multiVMTemplate(data), data.RandomInteger, strings.Join(scopes, ",\n"), data.Locations.Primary)
}

func testAccAzureRMMonitorMetricAlert_vmScopeDynamicCriteria(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

%s

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_linux_virtual_machine.test1.id]
  dynamic_criteria {
    metric_namespace = "Microsoft.Compute/virtualMachines"
    metric_name      = "Network In"
    aggregation      = "Total"

    operator          = "GreaterOrLessThan"
    alert_sensitivity = "Medium"
  }
  window_size              = "PT1H"
  frequency                = "PT5M"
  target_resource_type     = "Microsoft.Compute/virtualMachines"
  target_resource_location = "eastus"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, testAccAzureRMMonitorMetricAlert_vmTemplate(data, 1), data.RandomInteger)
}
