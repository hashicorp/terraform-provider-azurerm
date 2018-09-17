package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMApplicationInsights_basicWeb(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "web")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists("azurerm_application_insights.test"),
					resource.TestCheckResourceAttr("azurerm_application_insights.test", "application_type", "web"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationInsights_basicJava(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "java")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists("azurerm_application_insights.test"),
					resource.TestCheckResourceAttr("azurerm_application_insights.test", "application_type", "java"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationInsights_basicMobileCenter(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "MobileCenter")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists("azurerm_application_insights.test"),
					resource.TestCheckResourceAttr("azurerm_application_insights.test", "application_type", "MobileCenter"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationInsights_basicOther(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "other")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists("azurerm_application_insights.test"),
					resource.TestCheckResourceAttr("azurerm_application_insights.test", "application_type", "other"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationInsights_basicPhone(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "phone")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists("azurerm_application_insights.test"),
					resource.TestCheckResourceAttr("azurerm_application_insights.test", "application_type", "phone"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationInsights_basicStore(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "store")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists("azurerm_application_insights.test"),
					resource.TestCheckResourceAttr("azurerm_application_insights.test", "application_type", "store"),
				),
			},
		},
	})
}

func TestAccAzureRMApplicationInsights_basiciOS(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMApplicationInsights_basic(ri, testLocation(), "ios")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists("azurerm_application_insights.test"),
					resource.TestCheckResourceAttr("azurerm_application_insights.test", "application_type", "ios"),
				),
			},
		},
	})
}

func testCheckAzureRMApplicationInsightsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).appInsightsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_insights" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Application Insights still exists:\n%#v", resp.ApplicationInsightsComponentProperties)
		}
	}

	return nil
}

func testCheckAzureRMApplicationInsightsExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Insights: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).appInsightsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on appInsightsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Application Insights '%q' (resource group: '%q') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMApplicationInsights_basic(rInt int, location string, applicationType string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "%s"
}
`, rInt, location, rInt, applicationType)
}
