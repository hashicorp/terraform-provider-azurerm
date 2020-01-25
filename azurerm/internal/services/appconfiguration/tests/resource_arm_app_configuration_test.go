package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppConfigurationName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "four",
			ErrCount: 1,
		},
		{
			Value:    "5five",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloWorld",
			ErrCount: 0,
		},
		{
			Value:    "helloworld12",
			ErrCount: 0,
		},
		{
			Value:    "hello@world",
			ErrCount: 1,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd3324120",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd332412020",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqfjjfewsqwcdw21ddwqwd33241201",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := appconfiguration.ValidateAppConfigurationName(tc.Value, "azurerm_app_configuration")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure App Configuration Name to trigger a validation error: %v", tc)
		}
	}
}

func TestAccAzureAppConfiguration_free(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_free(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureAppConfiguration_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_standard(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureAppConfiguration_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_free(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureAppConfiguration_requiresImport),
		},
	})
}

func TestAccAzureAppConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureAppConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureAppConfiguration_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureAppConfiguration_completeUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureAppConfigurationExists(data.ResourceName),
				),
			},
		},
	})
}

func testCheckAzureAppConfigurationDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).AppConfiguration.AppConfigurationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_configuration" {
			continue
		}

		id, err := parse.AppConfigurationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testCheckAzureAppConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).AppConfiguration.AppConfigurationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.AppConfigurationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on appConfigurationsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: App Configuration %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testAccAzureAppConfiguration_free(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testacc-appconf%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "free"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureAppConfiguration_standard(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureAppConfiguration_requiresImport(data acceptance.TestData) string {
	template := testAccAzureAppConfiguration_free(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration" "import" {
  name                = "${azurerm_app_configuration.test.name}"
  resource_group_name = "${azurerm_app_configuration.test.resource_group_name}"
  location            = "${azurerm_app_configuration.test.location}"
  sku                 = "${azurerm_app_configuration.test.sku}"

}
`, template)
}

func testAccAzureAppConfiguration_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "free"

  tags = {
    environment = "development"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureAppConfiguration_completeUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "free"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
