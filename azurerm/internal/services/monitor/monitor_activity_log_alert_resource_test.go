package monitor_test

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

type MonitorActivityLogAlertResource struct {
}

func TestAccMonitorActivityLogAlert_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")
	r := MonitorActivityLogAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("scopes.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.0.category").HasValue("Recommendation"),
				check.That(data.ResourceName).Key("action.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActivityLogAlert_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")
	r := MonitorActivityLogAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_monitor_activity_log_alert"),
		},
	})
}

func TestAccMonitorActivityLogAlert_singleResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")
	r := MonitorActivityLogAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleResource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("scopes.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.0.operation_name").HasValue("Microsoft.Storage/storageAccounts/write"),
				check.That(data.ResourceName).Key("criteria.0.category").HasValue("Recommendation"),
				check.That(data.ResourceName).Key("criteria.0.resource_id").Exists(),
				check.That(data.ResourceName).Key("action.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActivityLogAlert_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")
	r := MonitorActivityLogAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue("This is just a test acceptance."),
				check.That(data.ResourceName).Key("scopes.#").HasValue("2"),
				check.That(data.ResourceName).Key("criteria.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.0.operation_name").HasValue("Microsoft.Storage/storageAccounts/write"),
				check.That(data.ResourceName).Key("criteria.0.category").HasValue("Recommendation"),
				check.That(data.ResourceName).Key("criteria.0.resource_provider").HasValue("Microsoft.Storage"),
				check.That(data.ResourceName).Key("criteria.0.resource_type").HasValue("Microsoft.Storage/storageAccounts"),
				check.That(data.ResourceName).Key("criteria.0.resource_group").Exists(),
				check.That(data.ResourceName).Key("criteria.0.resource_id").Exists(),
				check.That(data.ResourceName).Key("action.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActivityLogAlert_basicAndCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")
	r := MonitorActivityLogAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("scopes.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.0.category").HasValue("Recommendation"),
				check.That(data.ResourceName).Key("criteria.0.resource_id").HasValue(""),
				check.That(data.ResourceName).Key("criteria.0.caller").HasValue(""),
				check.That(data.ResourceName).Key("criteria.0.level").HasValue(""),
				check.That(data.ResourceName).Key("criteria.0.status").HasValue(""),
				check.That(data.ResourceName).Key("action.#").HasValue("0"),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue("This is just a test acceptance."),
				check.That(data.ResourceName).Key("scopes.#").HasValue("2"),
				check.That(data.ResourceName).Key("criteria.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.0.operation_name").HasValue("Microsoft.Storage/storageAccounts/write"),
				check.That(data.ResourceName).Key("criteria.0.category").HasValue("Recommendation"),
				check.That(data.ResourceName).Key("criteria.0.resource_provider").HasValue("Microsoft.Storage"),
				check.That(data.ResourceName).Key("criteria.0.resource_type").HasValue("Microsoft.Storage/storageAccounts"),
				check.That(data.ResourceName).Key("criteria.0.resource_group").Exists(),
				check.That(data.ResourceName).Key("criteria.0.resource_id").Exists(),
				check.That(data.ResourceName).Key("action.#").HasValue("2"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("scopes.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.#").HasValue("1"),
				check.That(data.ResourceName).Key("criteria.0.category").HasValue("Recommendation"),
				check.That(data.ResourceName).Key("criteria.0.resource_id").HasValue(""),
				check.That(data.ResourceName).Key("criteria.0.caller").HasValue(""),
				check.That(data.ResourceName).Key("criteria.0.level").HasValue(""),
				check.That(data.ResourceName).Key("criteria.0.status").HasValue(""),
				check.That(data.ResourceName).Key("action.#").HasValue("0"),
			),
		},
	})
}

func TestAccMonitorActivityLogAlert_ServiceHealth_basicAndUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")
	r := MonitorActivityLogAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceHealth_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceHealth_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceHealth_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActivityLogAlert_ServiceHealth_basicAndDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_activity_log_alert", "test")
	r := MonitorActivityLogAlertResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceHealth_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceHealth_delete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MonitorActivityLogAlertResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestActivityLogAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_resource_group.test.id]

  criteria {
    category = "Recommendation"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r MonitorActivityLogAlertResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_activity_log_alert" "import" {
  name                = azurerm_monitor_activity_log_alert.test.name
  resource_group_name = azurerm_monitor_activity_log_alert.test.resource_group_name
  scopes              = [azurerm_resource_group.test.id]

  criteria {
    category = "Recommendation"
  }
}
`, r.basic(data))
}

func (MonitorActivityLogAlertResource) singleResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestActivityLogAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_resource_group.test.id]

  criteria {
    operation_name = "Microsoft.Storage/storageAccounts/write"
    category       = "Recommendation"
    resource_id    = azurerm_storage_account.test.id
  }

  action {
    action_group_id = azurerm_monitor_action_group.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger)
}

func (MonitorActivityLogAlertResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test1" {
  name                = "acctestActionGroup1-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag1"
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag2"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestActivityLogAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  enabled             = true
  description         = "This is just a test acceptance."

  scopes = [
    azurerm_resource_group.test.id,
    azurerm_storage_account.test.id,
  ]

  criteria {
    operation_name          = "Microsoft.Storage/storageAccounts/write"
    category                = "Recommendation"
    resource_provider       = "Microsoft.Storage"
    resource_type           = "Microsoft.Storage/storageAccounts"
    resource_group          = azurerm_resource_group.test.name
    resource_id             = azurerm_storage_account.test.id
    recommendation_category = "OperationalExcellence"
    recommendation_impact   = "High"
  }

  action {
    action_group_id = azurerm_monitor_action_group.test1.id
  }

  action {
    action_group_id = azurerm_monitor_action_group.test2.id

    webhook_properties = {
      from = "terraform test"
      to   = "microsoft azure"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger)
}

func (MonitorActivityLogAlertResource) serviceHealth_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test1" {
  name                = "acctestActionGroup1-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag1"
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag2"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestActivityLogAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  enabled             = true
  description         = "This is just a test acceptance."

  scopes = [
    azurerm_resource_group.test.id,
    azurerm_storage_account.test.id,
  ]

  criteria {
    category = "ServiceHealth"
    service_health {
      events    = ["Incident", "Maintenance"]
      services  = ["Action Groups", "Activity Logs & Alerts"]
      locations = ["Global", "West Europe"]
    }
  }

  action {
    action_group_id = azurerm_monitor_action_group.test1.id
  }

  action {
    action_group_id = azurerm_monitor_action_group.test2.id

    webhook_properties = {
      from = "terraform test"
      to   = "microsoft azure"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger)
}

func (MonitorActivityLogAlertResource) serviceHealth_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test1" {
  name                = "acctestActionGroup1-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag1"
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag2"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestActivityLogAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  enabled             = true
  description         = "This is just a test acceptance."

  scopes = [
    azurerm_resource_group.test.id,
    azurerm_storage_account.test.id,
  ]

  criteria {
    category = "ServiceHealth"
    service_health {
      events    = ["Incident", "Maintenance", "ActionRequired", "Security"]
      services  = ["Action Groups"]
      locations = ["Global", "West Europe", "East US"]
    }
  }

  action {
    action_group_id = azurerm_monitor_action_group.test1.id
  }

  action {
    action_group_id = azurerm_monitor_action_group.test2.id

    webhook_properties = {
      from = "terraform test"
      to   = "microsoft azure"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger)
}

func (MonitorActivityLogAlertResource) serviceHealth_delete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test1" {
  name                = "acctestActionGroup1-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag1"
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag2"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestActivityLogAlert-%d"
  resource_group_name = azurerm_resource_group.test.name
  enabled             = true
  description         = "This is just a test acceptance."

  scopes = [
    azurerm_resource_group.test.id,
    azurerm_storage_account.test.id,
  ]

  criteria {
    category = "Recommendation"
  }

  action {
    action_group_id = azurerm_monitor_action_group.test1.id
  }

  action {
    action_group_id = azurerm_monitor_action_group.test2.id

    webhook_properties = {
      from = "terraform test"
      to   = "microsoft azure"
    }
  }
}
	`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger)
}

func (t MonitorActivityLogAlertResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["activityLogAlerts"]

	resp, err := clients.Monitor.ActivityLogAlertsClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading activity log alert (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
