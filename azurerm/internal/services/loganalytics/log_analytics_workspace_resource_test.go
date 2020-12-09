package loganalytics_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
)

func TestAccAzureRMLogAnalyticsWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLogAnalyticsWorkspace_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_workspace"),
			},
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_freeTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_freeTier(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_withDefaultSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withDefaultSku(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_withVolumeCap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withVolumeCap(data, 4.5),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_removeVolumeCap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withVolumeCap(data, 5.5),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_removeVolumeCap(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_quota_gb", "-1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_withInternetIngestionEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withInternetIngestionEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withInternetIngestionEnabledUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_withInternetQueryEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withInternetQueryEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withInternetQueryEnabledUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMLogAnalyticsWorkspaceDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.WorkspacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_workspace" {
			continue
		}

		id, err := parse.LogAnalyticsWorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.WorkspaceName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Log Analytics Workspace still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMLogAnalyticsWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.WorkspacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.LogAnalyticsWorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.WorkspaceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on Log Analytics Workspace Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Log Analytics Workspace '%s' (resource group: '%s') does not exist", id.WorkspaceName, id.ResourceGroup)
		}

		return nil
	}
}

func testAccAzureRMLogAnalyticsWorkspace_basic(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsWorkspace_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "import" {
  name                = azurerm_log_analytics_workspace.test.name
  location            = azurerm_log_analytics_workspace.test.location
  resource_group_name = azurerm_log_analytics_workspace.test.resource_group_name
  sku                 = "PerGB2018"
}
`, template)
}

func testAccAzureRMLogAnalyticsWorkspace_complete(data acceptance.TestData) string {
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

func testAccAzureRMLogAnalyticsWorkspace_freeTier(data acceptance.TestData) string {
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
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_withDefaultSku(data acceptance.TestData) string {
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
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_withVolumeCap(data acceptance.TestData, volumeCapGb float64) string {
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

func testAccAzureRMLogAnalyticsWorkspace_removeVolumeCap(data acceptance.TestData) string {
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

func testAccAzureRMLogAnalyticsWorkspace_withInternetIngestionEnabled(data acceptance.TestData) string {
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
  retention_in_days          = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_withInternetIngestionEnabledUpdate(data acceptance.TestData) string {
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
  retention_in_days          = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_withInternetQueryEnabled(data acceptance.TestData) string {
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
  retention_in_days      = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_withInternetQueryEnabledUpdate(data acceptance.TestData) string {
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
  retention_in_days      = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
