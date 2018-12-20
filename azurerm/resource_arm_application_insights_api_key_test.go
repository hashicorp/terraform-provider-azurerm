package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMApplicationInsightsAPIKey_no_permission(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsightsAPIKey_basic(ri, testLocation(), "[]", "[]")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("The API Key needs to have a Role"),
			},
		},
	})
}

func TestAccAzureRMApplicationInsightsAPIKey_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_application_insights_api_key.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsAPIKey_basic(ri, testLocation(), "[]", `["annotations"]`),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "read_permissions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "write_permissions.#", "1"),
				),
			},
			{
				Config:      testAccAzureRMApplicationInsightsAPIKey_requiresImport(ri, location, "[]", `["annotations"]`),
				ExpectError: testRequiresImportError("azurerm_application_insights_api_key"),
			},
		},
	})
}

func TestAccAzureRMApplicationInsightsAPIKey_read_telemetry_permissions(t *testing.T) {
	resourceName := "azurerm_application_insights_api_key.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsightsAPIKey_basic(ri, testLocation(), `["aggregate", "api", "draft", "extendqueries", "search"]`, "[]")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "read_permissions.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "write_permissions.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"api_key", // not returned from API, sensitive
				},
			},
		},
	})
}

func TestAccAzureRMApplicationInsightsAPIKey_write_annotations_permission(t *testing.T) {
	resourceName := "azurerm_application_insights_api_key.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsightsAPIKey_basic(ri, testLocation(), "[]", `["annotations"]`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "read_permissions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "write_permissions.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"api_key", // not returned from API, sensitive
				},
			},
		},
	})
}

func TestAccAzureRMApplicationInsightsAPIKey_authenticate_permission(t *testing.T) {
	resourceName := "azurerm_application_insights_api_key.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsightsAPIKey_basic(ri, testLocation(), `["agentconfig"]`, "[]")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "read_permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "write_permissions.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"api_key", // not returned from API, sensitive
				},
			},
		},
	})

}

func TestAccAzureRMApplicationInsightsAPIKey_full_permissions(t *testing.T) {
	resourceName := "azurerm_application_insights_api_key.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsightsAPIKey_basic(ri, testLocation(), `["agentconfig", "aggregate", "api", "draft", "extendqueries", "search"]`, `["annotations"]`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsAPIKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "read_permissions.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "write_permissions.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"api_key", // not returned from API, sensitive
				},
			},
		},
	})

}

func testCheckAzureRMApplicationInsightsAPIKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).appInsightsAPIKeyClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_insights_api_key" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		appInsightsName := rs.Primary.Attributes["application_insights_name"]

		resp, err := conn.Get(ctx, resourceGroup, appInsightsName, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Application Insights API Key still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMApplicationInsightsAPIKeyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id, err := parseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		keyID := id.Path["APIKeys"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Insights API Key: %s", name)
		}
		appInsightsName, hasAppInsightsName := rs.Primary.Attributes["application_insights_name"]
		if !hasAppInsightsName {
			return fmt.Errorf("Bad: no Application Insights resource found in state for App Insights API Key: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).appInsightsAPIKeyClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, appInsightsName, keyID)
		if err != nil {
			return fmt.Errorf("Bad: Get on appInsightsAPIKeyClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Application Insights API Key '%q' (resource group: '%q') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMApplicationInsightsAPIKey_basic(rInt int, location, readPerms, writePerms string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "test" {
  name                      = "acctestappinsightsapikey-%d"
  application_insights_name = "${azurerm_application_insights.test.name}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  read_permissions          = %s
  write_permissions         = %s
}
`, rInt, location, rInt, rInt, readPerms, writePerms)
}

func testAccAzureRMApplicationInsightsAPIKey_requiresImport(rInt int, location, readPerms, writePerms string) string {
	template := testAccAzureRMApplicationInsightsAPIKey_basic(rInt, location, readPerms, writePerms)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_api_key" "import" {
  name                      = "${azurerm_application_insights_api_key.test.name}"
  application_insights_name = "${azurerm_application_insights_api_key.test.application_insights_name}"
  resource_group_name       = "${azurerm_application_insights_api_key.test.resource_group_name}"
  read_permissions          = "${azurerm_application_insights_api_key.test.read_permissions}"
  write_permissions         = "${azurerm_application_insights_api_key.test.write_permissions}"
}
`, template)
}
