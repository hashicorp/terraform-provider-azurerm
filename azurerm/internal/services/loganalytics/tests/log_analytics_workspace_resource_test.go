package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
)

func TestAccAzureRmLogAnalyticsWorkspaceName_validation(t *testing.T) {
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

func TestAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForIngestion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForIngestionEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForIngestionDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForQuery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForQueryEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForQueryDisabled(data),
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

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)

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

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on Log Analytics Workspace Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Log Analytics Workspace '%s' (resource group: '%s') does not exist", id.Name, id.ResourceGroup)
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
  retention_in_days   = 7
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

func testAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForIngestionEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                                        = "acctestLAW-%d"
  location                                    = azurerm_resource_group.test.location
  resource_group_name                         = azurerm_resource_group.test.name
  public_network_access_for_ingestion_enabled = true
  sku                                         = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForIngestionDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                                        = "acctestLAW-%d"
  location                                    = azurerm_resource_group.test.location
  resource_group_name                         = azurerm_resource_group.test.name
  public_network_access_for_ingestion_enabled = false
  sku                                         = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForQueryEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                                    = "acctestLAW-%d"
  location                                = azurerm_resource_group.test.location
  resource_group_name                     = azurerm_resource_group.test.name
  public_network_access_for_query_enabled = true
  sku                                     = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_withPublicNetworkAccessForQueryDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                                    = "acctestLAW-%d"
  location                                = azurerm_resource_group.test.location
  resource_group_name                     = azurerm_resource_group.test.name
  public_network_access_for_query_enabled = false
  sku                                     = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
