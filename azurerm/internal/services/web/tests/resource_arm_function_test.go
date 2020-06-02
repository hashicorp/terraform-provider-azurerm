package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFunction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunction_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "app_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "config"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "test_data"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "script_href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "script_root_href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "config_href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secrets_file_href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "invoke_url_template"),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "Microsoft.Web/sites/functions"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunction_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunction_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "app_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "config"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "test_data"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "script_href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "script_root_href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "config_href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secrets_file_href"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "invoke_url_template"),
					resource.TestCheckResourceAttr(data.ResourceName, "type", "Microsoft.Web/sites/functions"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFunction_requiresImport),
		},
	})
}

func testCheckAzureRMFunctionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_function" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		functionAppName := rs.Primary.Attributes["app_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetFunction(ctx, resourceGroup, functionAppName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMFunctionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		functionAppName := rs.Primary.Attributes["app_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]

		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function: %s", name)
		}

		resp, err := client.GetFunction(ctx, resourceGroup, functionAppName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Function %q (Resource Group: %q, Function App: %q) does not exist", name, resourceGroup, functionAppName)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %s", err)
		}

		return nil
	}
}

func testAccAzureRMFunction_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Free"
    size = "F1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func-app"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
}

resource "azurerm_function" "test" {
  name                = "acctest-%[1]d-func"
  app_name            = azurerm_function_app.test.name
  resource_group_name = azurerm_resource_group.test.name
  config              = "{\"bindings\": [{\"type\": \"http\",\"direction\": \"out\",\"name\": \"res\"},{\"name\": \"req\",\"route\": \"route/{route}\",\"authLevel\": \"anonymous\",\"methods\": [\"post\"],\"direction\": \"in\",\"type\": \"httpTrigger\"}]}"
  files = {
    "index.js" = "module.exports = function (context, req) {\n  context.res = {\n    body: {\n      params: req.params,\n      headers: req.headers\n    }\n  };\n  context.done();\n};"
  }
  test_data = "{\"method\":\"post\",\"queryStringParams\":[{\"name\":\"route\",\"value\":\"route1\"}],\"headers\":[{\"name\":\"content-type\",\"value\":\"application/json\"}],\"body\":\"{\\n  \\\"key\\\":\\\"value\\\"\\n}\"}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func testAccAzureRMFunction_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFunction_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_function" "import" {
  name                = azurerm_function.test.name
  app_name            = azurerm_function.test.app_name
  resource_group_name = azurerm_resource_group.test.name
  config              = azurerm_function.test.config
  files               = azurerm_function.test.files
  test_data           = azurerm_function.test.test_data
}
`, template)
}
