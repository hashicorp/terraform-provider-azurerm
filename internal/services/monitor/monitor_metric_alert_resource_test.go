// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorMetricAlertResource struct{}

func TestAccMonitorMetricAlert_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")
	r := MonitorMetricAlertResource{}

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

func TestAccMonitorMetricAlert_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")
	r := MonitorMetricAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_monitor_metric_alert"),
		},
	})
}

func TestAccMonitorMetricAlert_multiCriteria(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")
	r := MonitorMetricAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiCriteria(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorMetricAlert_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")
	r := MonitorMetricAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorMetricAlert_dynamicCriteria(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")
	r := MonitorMetricAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dynamicCriteria(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorMetricAlert_basicAndCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")
	r := MonitorMetricAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccMonitorMetricAlert_multiScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")
	r := MonitorMetricAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiScope(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiScope(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiScope(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorMetricAlert_applicationInsightsWebTest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")
	r := MonitorMetricAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.applicationInsightsWebTest(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t MonitorMetricAlertResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := metricalerts.ParseMetricAlertID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Monitor.MetricAlertsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading (%s): %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (MonitorMetricAlertResource) basic(data acceptance.TestData) string {
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

  tags = {
    test      = "123"
    Example   = "Example123"
    terraform = "Coolllll"
    CUSTOMER  = "CUSTOMERx"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r MonitorMetricAlertResource) requiresImport(data acceptance.TestData) string {
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
`, r.basic(data))
}

func (MonitorMetricAlertResource) dynamicCriteria(data acceptance.TestData) string {
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

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_storage_account.test.id]

  dynamic_criteria {
    metric_namespace  = "Microsoft.Storage/storageAccounts"
    metric_name       = "Availability"
    aggregation       = "Minimum"
    operator          = "GreaterThan"
    alert_sensitivity = "High"

    dimension {
      name     = "ApiName"
      operator = "Include"
      values   = ["*"]
    }

    evaluation_total_count   = 4
    evaluation_failure_count = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (MonitorMetricAlertResource) multiCriteria(data acceptance.TestData) string {
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

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_storage_account.test.id]
  enabled             = true
  auto_mitigate       = false
  severity            = 4
  description         = "This is a complete metric alert acceptance."
  frequency           = "PT30M"
  window_size         = "PT12H"

  criteria {
    metric_namespace       = "Microsoft.Storage/storageAccounts"
    metric_name            = "Transactions"
    aggregation            = "Total"
    operator               = "GreaterThan"
    threshold              = 99
    skip_metric_validation = true

    dimension {
      name     = "GeoType"
      operator = "Include"
      values   = ["Primary"]
    }
  }

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "UsedCapacity"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = 55.5
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (MonitorMetricAlertResource) complete(data acceptance.TestData) string {
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
  description         = "This is a complete metric alert acceptance."
  frequency           = "PT30M"
  window_size         = "PT12H"

  criteria {
    metric_namespace       = "Microsoft.Storage/storageAccounts"
    metric_name            = "Transactions"
    aggregation            = "Total"
    operator               = "GreaterThan"
    threshold              = 99
    skip_metric_validation = true

    dimension {
      name     = "GeoType"
      operator = "Include"
      values   = ["*"]
    }
  }

  action {
    action_group_id = azurerm_monitor_action_group.test1.id
    webhook_properties = {
      from = "terraform"
    }
  }

  action {
    action_group_id = azurerm_monitor_action_group.test2.id
  }

  tags = {
    test          = "456"
    Example       = "Example456"
    Terraform     = "Coolllll"
    tfazurerm     = "Awesome"
    CUSTOMER      = "CUSTOMERx"
    "EXAMPLE.TAG" = "sample"
    "Foo.Bar"     = "Test tag"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (MonitorMetricAlertResource) multiVMTemplate(data acceptance.TestData, count int) string {
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
  address_prefixes     = ["10.0.${count.index}.0/24"]
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
  count               = %[3]d
  name                = "acctestVM-${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test[count.index].id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, data.RandomInteger, data.Locations.Primary, count)
}

func (r MonitorMetricAlertResource) multiScope(data acceptance.TestData, count int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = azurerm_linux_virtual_machine.test.*.id
  dynamic_criteria {
    metric_namespace = "Microsoft.Compute/virtualMachines"
    metric_name      = "CPU Credits Consumed"
    aggregation      = "Average"

    operator                 = "GreaterOrLessThan"
    alert_sensitivity        = "Medium"
    skip_metric_validation   = true
    evaluation_failure_count = 4
    evaluation_total_count   = 5
    ignore_data_before       = "2022-03-02T15:04:05Z"
  }
  window_size              = "PT5M"
  frequency                = "PT5M"
  target_resource_type     = "Microsoft.Compute/virtualMachines"
  target_resource_location = "%s"
}
`, r.multiVMTemplate(data, count), data.RandomInteger, data.Locations.Primary)
}

func (MonitorMetricAlertResource) applicationInsightsWebTestTemplate(data acceptance.TestData) string {
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

func (r MonitorMetricAlertResource) applicationInsightsWebTest(data acceptance.TestData) string {
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
`, r.applicationInsightsWebTestTemplate(data), data.RandomInteger)
}
