package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMApplicationInsightsWebTests_basic(t *testing.T) {
	resourceName := "azurerm_application_insights_web_test.test"
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApplicationInsightsWebTests_complete(t *testing.T) {
	resourceName := "azurerm_application_insights_web_test.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWebTestsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsWebTests_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWebTestExists(resourceName),
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

func TestAccAzureRMApplicationInsightsWebTests_update(t *testing.T) {
	resourceName := "azurerm_application_insights_web_test.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWebTestsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsWebTests_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWebTestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "geo_locations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "300"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "30"),
				),
			},
			{
				Config: testAccAzureRMApplicationInsightsWebTests_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWebTestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "geo_locations.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "900"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "120"),
				),
			},
			{
				Config: testAccAzureRMApplicationInsightsWebTests_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWebTestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "geo_locations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "frequency", "300"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "30"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationInsightsWebTests_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_application_insights_web_test.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMApplicationInsightsWebTests_basic(ri, location)

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
				Config:      testAccAzureRMApplicationInsightsWebTests_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_application_insights_web_test"),
			},
		},
	})
}

func testCheckAzureRMApplicationInsightsWebTestsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).appInsights.WebTestsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_insights_web_test" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]

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
		resGroup := rs.Primary.Attributes["resource_group_name"]
		conn := testAccProvider.Meta().(*ArmClient).appInsights.WebTestsClient
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
  name                = "acctestappinsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_application_insights_web_test" "test" {
  name                    = "acctestappinsightswebtests-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  application_insights_id = "${azurerm_application_insights.test.id}"
  kind                    = "ping"
  geo_locations           = ["us-tx-sn1-azr"]

  configuration = <<XML
<WebTest Name="WebTest1" Id="ABD48585-0831-40CB-9069-682EA6BB3583" Enabled="True" CssProjectStructure="" CssIteration="" Timeout="0" WorkItemIds="" xmlns="http://microsoft.com/schemas/VisualStudio/TeamTest/2010" Description="" CredentialUserName="" CredentialPassword="" PreAuthenticate="True" Proxy="default" StopOnError="False" RecordedResultFile="" ResultsLocale="">
  <Items>
    <Request Method="GET" Guid="a5f10126-e4cd-570d-961c-cea43999a200" Version="1.1" Url="http://microsoft.com" ThinkTime="0" Timeout="300" ParseDependentRequests="True" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0" Encoding="utf-8" ExpectedHttpStatusCode="200" ExpectedResponseUrl="" ReportingName="" IgnoreHttpStatusCode="False" />
  </Items>
</WebTest>
XML
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMApplicationInsightsWebTests_complete(rInt int, location string) string {
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

resource "azurerm_application_insights_web_test" "test" {
  name                    = "acctestappinsightswebtests-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  application_insights_id = "${azurerm_application_insights.test.id}"
  kind                    = "ping"
  frequency               = 900
  timeout                 = 120
  enabled                 = true
  geo_locations           = ["us-tx-sn1-azr", "us-il-ch1-azr"]

  configuration = <<XML
<WebTest Name="WebTest1" Id="ABD48585-0831-40CB-9069-682EA6BB3583" Enabled="True" CssProjectStructure="" CssIteration="" Timeout="0" WorkItemIds="" xmlns="http://microsoft.com/schemas/VisualStudio/TeamTest/2010" Description="" CredentialUserName="" CredentialPassword="" PreAuthenticate="True" Proxy="default" StopOnError="False" RecordedResultFile="" ResultsLocale="">
  <Items>
    <Request Method="GET" Guid="a5f10126-e4cd-570d-961c-cea43999a200" Version="1.1" Url="http://microsoft.com" ThinkTime="0" Timeout="300" ParseDependentRequests="True" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0" Encoding="utf-8" ExpectedHttpStatusCode="200" ExpectedResponseUrl="" ReportingName="" IgnoreHttpStatusCode="False" />
  </Items>
</WebTest>
XML
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMApplicationInsightsWebTests_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApplicationInsightsWebTests_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_web_test" "import" {
  name                    = "${azurerm_application_insights_web_test.test.name}"
  location                = "${azurerm_application_insights_web_test.test.location}"
  resource_group_name     = "${azurerm_application_insights_web_test.test.resource_group_name}"
  application_insights_id = "${azurerm_application_insights_web_test.test.application_insights_id}"
  kind                    = "${azurerm_application_insights_web_test.test.kind}"
  configuration           = "${azurerm_application_insights_web_test.test.configuration}"
}
`, template)
}
