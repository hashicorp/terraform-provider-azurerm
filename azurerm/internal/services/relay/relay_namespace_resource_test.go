package relay_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMRelayNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayNamespace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayNamespaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "metric_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRelayNamespace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayNamespace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayNamespaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "metric_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMRelayNamespace_requiresImport),
		},
	})
}

func TestAccAzureRMRelayNamespace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayNamespace_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayNamespaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "metric_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMRelayNamespaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Relay.NamespacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on relayNamespacesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Relay Namespace %q (Resource Group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMRelayNamespaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Relay.NamespacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_relay_namespace" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Relay Namespace still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMRelayNamespace_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctestrn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMRelayNamespace_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMRelayNamespace_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_relay_namespace" "import" {
  name                = azurerm_relay_namespace.test.name
  location            = azurerm_relay_namespace.test.location
  resource_group_name = azurerm_relay_namespace.test.resource_group_name

  sku_name = "Standard"
}
`, template)
}

func testAccAzureRMRelayNamespace_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctestrn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Standard"

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
