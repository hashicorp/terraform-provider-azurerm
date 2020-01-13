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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics"
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
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

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

func testCheckAzureRMLogAnalyticsWorkspaceDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).LogAnalytics.WorkspacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_log_analytics_workspace" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

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

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Log Analytics Workspace: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on Log Analytics Workspace Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Log Analytics Workspace '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		return nil
	}
}
func testAccAzureRMLogAnalyticsWorkspace_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMLogAnalyticsWorkspace_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLogAnalyticsWorkspace_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "import" {
  name                = "${azurerm_log_analytics_workspace.test.name}"
  location            = "${azurerm_log_analytics_workspace.test.location}"
  resource_group_name = "${azurerm_log_analytics_workspace.test.resource_group_name}"
  sku                 = "PerGB2018"
}
`, template)
}

func testAccAzureRMLogAnalyticsWorkspace_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
  retention_in_days   = 30

  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
