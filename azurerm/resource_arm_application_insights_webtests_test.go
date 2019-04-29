package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMApplicationInsightsWebTests_basicWeb(t *testing.T) {
	resourceName := "azurerm_application_insights_webtest.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMApplicationInsightsWebTests_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWebTestsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWebTestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "provisioning_state", "Succeeded"),
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

func TestAccAzureRMApplicationInsightsWebTests_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_application_insights_webtest.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMApplicationInsightsWebTests_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWebTestsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWebTestExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApplicationInsightsWebTests_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_application_insights_webtest"),
			},
		},
	})
}

func testCheckAzureRMApplicationInsightsWebTestsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).appInsightsWebTestsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_insights_webtest" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		id, err := parseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		resGroup := id.ResourceGroup
		//appInsightsName := id.Path["components"]
		//webTestName := id.Path["webtests"]

		resp, err := conn.Get(ctx, resGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Application Insights WebTest still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMApplicationInsightsWebTestExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up a WebTest
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		id, err := parseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		resGroup := id.ResourceGroup
		//appInsightsName := id.Path["components"]
		//webTestName := id.Path["webtests"]

		conn := testAccProvider.Meta().(*ArmClient).appInsightsWebTestsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on appInsightsWebTestClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Application Insights WebTest '%q' (resource group: '%q') does not exist", name, resGroup)
		}

		return nil
	}
}

func testAccAzureRMApplicationInsightsWebTests_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name		= "acctestappinsights-%d"
  location	    = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_application_insights_webtest" "test" {
  name		    = "acctestappinsightswebtests-%d"
  location		= "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  application_insights_id = "${azurerm_application_insights.test.id}"
  kind		    = "Ping"
  frequency	       = 300
  timeout		 = 120
  enabled		 = true
  geo_locations	   = ["us-tx-sn1-azr"]

  test_configuration = "${file("testdata/application_insights_webtest_test.xml")}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMApplicationInsightsWebTests_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApplicationInsightsWebTests_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_webtest" "import" {
  name		    = "${azurerm_application_insights_webtest.test.name}"
  application_insights_id = "${azurerm_application_insights_webtest.test.application_insights_id}"
}
`, template)
}
