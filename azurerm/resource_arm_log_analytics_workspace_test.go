package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
		_, errors := validateAzureRmLogAnalyticsWorkspaceName(tc.Value, "azurerm_log_analytics")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the AzureRM Log Analytics Workspace Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccAzureRMLogAnalyticsWorkspace_basic(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_log_analytics_workspace.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMLogAnalyticsWorkspace_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_log_analytics_workspace"),
			},
		},
	})
}

func TestAccAzureRMLogAnalyticsWorkspace_complete(t *testing.T) {
	resourceName := "azurerm_log_analytics_workspace.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLogAnalyticsWorkspace_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLogAnalyticsWorkspaceExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMLogAnalyticsWorkspaceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).logAnalytics.WorkspacesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

		conn := testAccProvider.Meta().(*ArmClient).logAnalytics.WorkspacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
func testAccAzureRMLogAnalyticsWorkspace_basic(rInt int, location string) string {
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
`, rInt, location, rInt)
}

func testAccAzureRMLogAnalyticsWorkspace_requiresImport(rInt int, location string) string {
	template := testAccAzureRMLogAnalyticsWorkspace_basic(rInt, location)
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

func testAccAzureRMLogAnalyticsWorkspace_complete(rInt int, location string) string {
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
`, rInt, location, rInt)
}
