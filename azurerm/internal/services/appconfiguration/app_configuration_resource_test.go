package appconfiguration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAppConfigurationResource_free(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAppConfigurationResource_free(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAppConfigurationResource_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAppConfigurationResource_standard(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAppConfigurationResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAppConfigurationResource_free(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAppConfigurationResource_requiresImport),
		},
	})
}

func TestAccAppConfigurationResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAppConfigurationResource_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAppConfigurationResource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAppConfigurationResource_identity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAppConfigurationResource_identityUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAppConfigurationResource_standard(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAppConfigurationResource_identity(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
			data.ImportStep(),
			{
				Config: testAppConfigurationResource_standard(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAppConfigurationResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAppConfigurationResource_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
				),
			},
			{
				Config: testAppConfigurationResource_completeUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
				),
			},
		},
	})
}

func testCheckAppConfigurationDestroy(s *terraform.State) error {
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

func testCheckAppConfigurationExists(resourceName string) resource.TestCheckFunc {
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

func testAppConfigurationResource_free(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testacc-appconf%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "free"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAppConfigurationResource_standard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAppConfigurationResource_requiresImport(data acceptance.TestData) string {
	template := testAppConfigurationResource_free(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_configuration" "import" {
  name                = azurerm_app_configuration.test.name
  resource_group_name = azurerm_app_configuration.test.resource_group_name
  location            = azurerm_app_configuration.test.location
  sku                 = azurerm_app_configuration.test.sku
}
`, template)
}

func testAppConfigurationResource_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "free"

  tags = {
    environment = "development"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAppConfigurationResource_identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "development"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAppConfigurationResource_completeUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_app_configuration" "test" {
  name                = "testaccappconf%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "free"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
