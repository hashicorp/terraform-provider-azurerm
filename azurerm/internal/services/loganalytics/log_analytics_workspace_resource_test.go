package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogAnalyticsWorkspaceResource struct {
}

func TestAccLogAnalyticsWorkspaceName_validation(t *testing.T) {
	str := acctest.RandString(63)
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "abc",
			ErrCount: 1,
		},
		{
			Value:    "Ab-c",
			ErrCount: 0,
		},
		{
			Value:    "-abc",
			ErrCount: 1,
		},
		{
			Value:    "abc-",
			ErrCount: 1,
		},
		{
			Value:    str,
			ErrCount: 0,
		},
		{
			Value:    str + "a",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := loganalytics.ValidateAzureRmLogAnalyticsWorkspaceName(tc.Value, "azurerm_log_analytics")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the AzureRM Log Analytics Workspace Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccLogAnalyticsWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_workspace"),
		},
	})
}

func TestAccLogAnalyticsWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_freeTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.freeTier(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withDefaultSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withDefaultSku(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withVolumeCap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withVolumeCap(data, 4.5),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_removeVolumeCap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withVolumeCap(data, 5.5),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeVolumeCap(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("daily_quota_gb").HasValue("-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withInternetIngestionEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withInternetIngestionEnabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetIngestionEnabledUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withInternetQueryEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withInternetQueryEnabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetQueryEnabledUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (LogAnalyticsWorkspaceResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LogAnalyticsStorageInsightsID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.WorkspacesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Log Analytics Workspace %q (resource group: %q): %v", id.Name, id.ResourceGroup)
	}

	return utils.Bool(resp.WorkspaceProperties != nil), nil
}

func (LogAnalyticsWorkspaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogAnalyticsWorkspaceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "import" {
  name                = azurerm_log_analytics_workspace.test.name
  location            = azurerm_log_analytics_workspace.test.location
  resource_group_name = azurerm_log_analytics_workspace.test.resource_group_name
  sku                 = "PerGB2018"
}
`, r.basic(data))
}

func (LogAnalyticsWorkspaceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) freeTier(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Free"
  retention_in_days   = 7
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withDefaultSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withVolumeCap(data acceptance.TestData, volumeCapGb float64) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
  daily_quota_gb      = %f

  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, volumeCapGb)
}

func (LogAnalyticsWorkspaceResource) removeVolumeCap(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withInternetIngestionEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                       = "acctestLAW-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  internet_ingestion_enabled = true
  sku                        = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withInternetIngestionEnabledUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                       = "acctestLAW-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  internet_ingestion_enabled = false
  sku                        = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withInternetQueryEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                   = "acctestLAW-%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  internet_query_enabled = true
  sku                    = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withInternetQueryEnabledUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                   = "acctestLAW-%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  internet_query_enabled = false
  sku                    = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
