package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMApplicationInsights_basicWeb(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsights_basic(data, "web"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "web"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationInsights_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsights_basic(data, "web"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "web"),
				),
			},
			{
				Config:      testAccAzureRMApplicationInsights_requiresImport(data, "web"),
				ExpectError: acceptance.RequiresImportError("azurerm_application_insights"),
			},
		},
	})
}

func TestAccAzureRMApplicationInsights_basicJava(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsights_basic(data, "java"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "java"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationInsights_basicMobileCenter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsights_basic(data, "MobileCenter"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "MobileCenter"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationInsights_basicOther(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsights_basic(data, "other"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "other"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationInsights_basicPhone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsights_basic(data, "phone"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "phone"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationInsights_basicStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsights_basic(data, "store"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "store"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationInsights_basiciOS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsights_basic(data, "ios"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "ios"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMApplicationInsightsDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).AppInsights.ComponentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMApplicationInsightsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).AppInsights.ComponentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Insights: %s", name)
		}

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

func TestAccAzureRMApplicationInsights_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsights_complete(data, "web"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "application_type", "web"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_in_days", "120"),
					resource.TestCheckResourceAttr(data.ResourceName, "sampling_percentage", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_data_cap_in_gb", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_data_cap_notifications_disabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMApplicationInsights_basic(data acceptance.TestData, applicationType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, applicationType)
}

func testAccAzureRMApplicationInsights_requiresImport(data acceptance.TestData, applicationType string) string {
	template := testAccAzureRMApplicationInsights_basic(data, applicationType)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights" "import" {
  name                = azurerm_application_insights.test.name
  location            = azurerm_application_insights.test.location
  resource_group_name = azurerm_application_insights.test.resource_group_name
  application_type    = azurerm_application_insights.test.application_type
}
`, template)
}

func testAccAzureRMApplicationInsights_complete(data acceptance.TestData, applicationType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                                  = "acctestappinsights-%d"
  location                              = azurerm_resource_group.test.location
  resource_group_name                   = azurerm_resource_group.test.name
  application_type                      = "%s"
  retention_in_days                     = 120
  sampling_percentage                   = 50
  daily_data_cap_in_gb                  = 50
  daily_data_cap_notifications_disabled = true
  disable_ip_masking                    = true

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, applicationType)
}
