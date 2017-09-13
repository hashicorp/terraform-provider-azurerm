package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRmOperationalInsightWorkspaceName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 1,
		},
		{
			Value:    "Ab-c",
			ErrCount: 0,
		},
		{
			Value:    "-ab",
			ErrCount: 1,
		},
		{
			Value:    "ab-",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateAzureRmOperationalInsightWorkspaceName(tc.Value, "azurerm_operational_insight_workspace")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the AzureRM Operational Insight Workspace Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccAzureRMOperationalInsightWorkspace_requiredOnly(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMOperationalInsightWorkspace_requiredOnly(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMOperationalInsightWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMOperationalInsightWorkspaceExists("azurerm_operational_insight_workspace.test"),
				),
			},
		},
	})
}
func TestAccAzureRMOperationalInsightWorkspace_optional(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMOperationalInsightWorkspace_optional(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMOperationalInsightWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMOperationalInsightWorkspaceExists("azurerm_operational_insight_workspace.test"),
				),
			},
		},
	})
}

func testCheckAzureRMOperationalInsightWorkspaceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).workspacesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_operational_insight_workspace" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("OperationalInsight Workspace still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMOperationalInsightWorkspaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for OperationalInsight Workspace: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).workspacesClient

		resp, err := conn.Get(resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on OperationalInsight Workspace Client: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: OperationalInsight Workspace '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMOperationalInsightWorkspace_requiredOnly(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_operational_insight_workspace" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Free"
}
`, rInt, location, rInt)
}

func testAccAzureRMOperationalInsightWorkspace_optional(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_operational_insight_workspace" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  retention_in_days   = 5
}
`, rInt, location, rInt)
}
