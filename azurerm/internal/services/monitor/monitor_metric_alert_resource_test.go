package monitor_test

import (
	"fmt"
	"net/http"
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
				Config: testAccAzureRMMonitorMetricAlert_multiScope(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorMetricAlert_multiScope(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorMetricAlert_multiScope(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorMetricAlert_applicationInsightsWebTest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_applicationInsightsWebTest(data),
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

func testAccAzureRMMonitorMetricAlert_multiVMTemplate(data acceptance.TestData, count int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  count                = %[3]d
  name                 = "internal-${count.index}"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.${count.index}.0/24"
}

resource "azurerm_network_interface" "test" {
  count               = %[3]d
  name                = "acctestnic-${count.index}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test[count.index].id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  count                           = %[3]d
  name                            = "acctestVM-${count.index}"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@$$w0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test[count.index].id,
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
`, data.RandomInteger, data.Locations.Primary, count)
}

func testAccAzureRMMonitorMetricAlert_multiScope(data acceptance.TestData, count int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = azurerm_linux_virtual_machine.test.*.id
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
`, testAccAzureRMMonitorMetricAlert_multiVMTemplate(data, count), data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMonitorMetricAlert_applicationInsightsWebTestTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsight-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_web_test" "test" {
  name                    = "acctestAppInsight-webtest-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  kind                    = "ping"
  frequency               = 300
  timeout                 = 60
  enabled                 = true
  geo_locations           = ["us-tx-sn1-azr", "us-il-ch1-azr"]

  configuration = <<XML
<WebTest Name="WebTest1" Id="ABD48585-0831-40CB-9069-682EA6BB3583" Enabled="True" CssProjectStructure="" CssIteration="" Timeout="0" WorkItemIds="" xmlns="http://microsoft.com/schemas/VisualStudio/TeamTest/2010" Description="" CredentialUserName="" CredentialPassword="" PreAuthenticate="True" Proxy="default" StopOnError="False" RecordedResultFile="" ResultsLocale="">
  <Items>
    <Request Method="GET" Guid="a5f10126-e4cd-570d-961c-cea43999a200" Version="1.1" Url="http://microsoft.com" ThinkTime="0" Timeout="300" ParseDependentRequests="True" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0" Encoding="utf-8" ExpectedHttpStatusCode="200" ExpectedResponseUrl="" ReportingName="" IgnoreHttpStatusCode="False" />
  </Items>
</WebTest>
XML
  lifecycle {
    ignore_changes = [tags]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMonitorMetricAlert_applicationInsightsWebTest(data acceptance.TestData) string {
	template := testAccAzureRMMonitorMetricAlert_applicationInsightsWebTestTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes = [
    azurerm_application_insights.test.id,
    azurerm_application_insights_web_test.test.id,
  ]
  application_insights_web_test_location_availability_criteria {
    web_test_id           = azurerm_application_insights_web_test.test.id
    component_id          = azurerm_application_insights.test.id
    failed_location_count = 2
  }
  window_size = "PT15M"
  frequency   = "PT1M"
}
`, template, data.RandomInteger)
}
