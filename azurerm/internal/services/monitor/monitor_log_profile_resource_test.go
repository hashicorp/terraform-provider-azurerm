package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MonitorLogProfileResource struct {
}

// NOTE: this is a combined test rather than separate split out tests due to
// Azure only being happy about provisioning one per subscription at once
// (which our test suite can't easily workaround)

// this occasionally fails due to the rapid provisioning and deprovisioning,
// running the exact same test afterwards always results in a pass.

func TestAccMonitorLogProfile(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":          testAccMonitorLogProfile_basic,
			"requiresImport": testAccMonitorLogProfile_requiresImport,
			"servicebus":     testAccMonitorLogProfile_servicebus,
			"complete":       testAccMonitorLogProfile_complete,
			"disappears":     testAccMonitorLogProfile_disappears,
			"update":         testAccMonitorLogProfile_update,
		},
		"datasource": {
			"eventhub":       testAccDataSourceMonitorLogProfile_eventhub,
			"storageaccount": testAccDataSourceMonitorLogProfile_storageaccount,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccMonitorLogProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")
	r := MonitorLogProfileResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccMonitorLogProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")
	r := MonitorLogProfileResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_monitor_log_profile"),
		},
	})
}

func testAccMonitorLogProfile_servicebus(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")
	r := MonitorLogProfileResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.servicebusConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func testAccMonitorLogProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")
	r := MonitorLogProfileResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccMonitorLogProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")
	r := MonitorLogProfileResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccMonitorLogProfile_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_log_profile", "test")
	r := MonitorLogProfileResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicConfig,
			TestResource: r,
		}),
	})
}

func (t MonitorLogProfileResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name, err := monitor.ParseLogProfileNameFromID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing log profile name from ID %s: %s", state.ID, err)
	}

	resp, err := clients.Monitor.LogProfilesClient.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("reading log profile %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (t MonitorLogProfileResource) Destroy(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name, err := monitor.ParseLogProfileNameFromID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing log profile name from ID %s: %s", state.ID, err)
	}

	if _, err := clients.Monitor.LogProfilesClient.Delete(ctx, name); err != nil {
		return nil, fmt.Errorf("deleting log profile %q: %+v", state.ID, err)
	}

	return utils.Bool(true), nil
}

func (MonitorLogProfileResource) basicConfig(data acceptance.TestData) string {
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
  account_replication_type = "GRS"
}

resource "azurerm_monitor_log_profile" "test" {
  name = "acctestlp-%d"

  categories = [
    "Action",
  ]

  locations = [
    "%s",
  ]

  storage_account_id = azurerm_storage_account.test.id

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.Locations.Primary)
}

func (r MonitorLogProfileResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_log_profile" "import" {
  name               = azurerm_monitor_log_profile.test.name
  categories         = azurerm_monitor_log_profile.test.categories
  locations          = azurerm_monitor_log_profile.test.locations
  storage_account_id = azurerm_monitor_log_profile.test.storage_account_id

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.basicConfig(data))
}

func (MonitorLogProfileResource) servicebusConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestsbns-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_namespace_authorization_rule" "test" {
  name                = "acctestsbrule-%s"
  namespace_name      = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = true
  send   = true
  manage = true
}

resource "azurerm_monitor_log_profile" "test" {
  name = "acctestlp-%d"

  categories = [
    "Action",
  ]

  locations = [
    "%s",
  ]

  servicebus_rule_id = azurerm_servicebus_namespace_authorization_rule.test.id

  retention_policy {
    enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, data.Locations.Primary)
}

func (MonitorLogProfileResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctestehns-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 2
}

resource "azurerm_monitor_log_profile" "test" {
  name     = "acctestlp-%[1]d"
  location = azurerm_resource_group.test.location
  categories = [
    "Action",
    "Delete",
    "Write",
  ]

  locations = [
    "%[2]s",
    "%[4]s",
  ]

  # RootManageSharedAccessKey is created by default with listen, send, manage permissions
  servicebus_rule_id = "${azurerm_eventhub_namespace.test.id}/authorizationrules/RootManageSharedAccessKey"
  storage_account_id = azurerm_storage_account.test.id

  retention_policy {
    enabled = true
    days    = 7
  }

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Secondary)
}
