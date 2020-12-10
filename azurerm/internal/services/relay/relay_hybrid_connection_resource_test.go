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

func TestAccAzureRMRelayHybridConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayHybridConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayHybridConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "requires_client_authorization"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRelayHybridConnection_full(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayHybridConnection_full(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayHybridConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "requires_client_authorization"),
					resource.TestCheckResourceAttr(data.ResourceName, "user_metadata", "metadatatest"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRelayHybridConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayHybridConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayHybridConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "requires_client_authorization"),
				),
			},
			{
				Config: testAccAzureRMRelayHybridConnection_update(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "requires_client_authorization", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "user_metadata", "metadataupdated"),
				),
			},
		},
	})
}

func TestAccAzureRMRelayHybridConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_relay_hybrid_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayHybridConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayHybridConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "requires_client_authorization"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMRelayHybridConnection_requiresImport),
		},
	})
}

func testAccAzureRMRelayHybridConnection_basic(data acceptance.TestData) string {
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

resource "azurerm_relay_hybrid_connection" "test" {
  name                 = "acctestrnhc-%d"
  resource_group_name  = azurerm_resource_group.test.name
  relay_namespace_name = azurerm_relay_namespace.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMRelayHybridConnection_full(data acceptance.TestData) string {
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

resource "azurerm_relay_hybrid_connection" "test" {
  name                 = "acctestrnhc-%d"
  resource_group_name  = azurerm_resource_group.test.name
  relay_namespace_name = azurerm_relay_namespace.test.name
  user_metadata        = "metadatatest"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMRelayHybridConnection_update(data acceptance.TestData) string {
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

resource "azurerm_relay_hybrid_connection" "test" {
  name                          = "acctestrnhc-%d"
  resource_group_name           = azurerm_resource_group.test.name
  relay_namespace_name          = azurerm_relay_namespace.test.name
  requires_client_authorization = false
  user_metadata                 = "metadataupdated"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMRelayHybridConnection_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMRelayHybridConnection_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_relay_hybrid_connection" "import" {
  name                 = azurerm_relay_hybrid_connection.test.name
  resource_group_name  = azurerm_relay_hybrid_connection.test.resource_group_name
  relay_namespace_name = azurerm_relay_hybrid_connection.test.relay_namespace_name
}
`, template)
}

func testCheckAzureRMRelayHybridConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Relay.HybridConnectionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		relayNamespace := rs.Primary.Attributes["relay_namespace_name"]

		// Ensure resource group exists in API

		resp, err := client.Get(ctx, resourceGroup, relayNamespace, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on relayHybridConnectionsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Relay Hybrid Connection %q in Namespace %q (Resource Group: %q) does not exist", name, relayNamespace, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMRelayHybridConnectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Relay.HybridConnectionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_relay_hybrid_connection" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		relayNamespace := rs.Primary.Attributes["relay_namespace_name"]

		resp, err := client.Get(ctx, resourceGroup, relayNamespace, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Relay Hybrid Connection still exists:\n%#v", resp)
		}
	}

	return nil
}
