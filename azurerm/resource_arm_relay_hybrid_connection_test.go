package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMRelayHybridConnection_basic(t *testing.T) {
	resourceName := "azurerm_relay_hybrid_connection.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayHybridConnection_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayHybridConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "requires_client_authorization"),
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

func TestAccAzureRMRelayHybridConnection_full(t *testing.T) {
	resourceName := "azurerm_relay_hybrid_connection.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayHybridConnection_full(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayHybridConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "requires_client_authorization"),
					resource.TestCheckResourceAttr(resourceName, "user_metadata", "metadatatest"),
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

func TestAccAzureRMRelayHybridConnection_update(t *testing.T) {
	resourceName := "azurerm_relay_hybrid_connection.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayHybridConnection_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayHybridConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "requires_client_authorization"),
				),
			},
			{
				Config: testAccAzureRMRelayHybridConnection_update(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "requires_client_authorization", "false"),
					resource.TestCheckResourceAttr(resourceName, "user_metadata", "metadataupdated"),
				),
			},
		},
	})
}

func TestAccAzureRMRelayHybridConnection_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_relay_hybrid_connection.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRelayHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRelayHybridConnection_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayHybridConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "requires_client_authorization"),
				),
			},
			{
				Config:      testAccAzureRMRelayHybridConnection_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_relay_hybrid_connection"),
			},
		},
	})
}

func testAccAzureRMRelayHybridConnection_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctestrn-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Standard"
}

resource "azurerm_relay_hybrid_connection" "test" {
	name                 = "acctestrnhc-%d"
	resource_group_name  = "${azurerm_resource_group.test.name}"
	relay_namespace_name = "${azurerm_relay_namespace.test.name}"
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRelayHybridConnection_full(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctestrn-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Standard"
}

resource "azurerm_relay_hybrid_connection" "test" {
	name                 = "acctestrnhc-%d"
	resource_group_name  = "${azurerm_resource_group.test.name}"
	relay_namespace_name = "${azurerm_relay_namespace.test.name}"
	user_metadata        = "metadatatest"
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRelayHybridConnection_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name                = "acctestrn-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name = "Standard"
}

resource "azurerm_relay_hybrid_connection" "test" {
	name                          = "acctestrnhc-%d"
	resource_group_name           = "${azurerm_resource_group.test.name}"
	relay_namespace_name          = "${azurerm_relay_namespace.test.name}"
	requires_client_authorization = false
	user_metadata                 = "metadataupdated"
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMRelayHybridConnection_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_relay_namespace" "import" {
	name                 = "acctestrnhc-%d"
	resource_group_name  = "${azurerm_resource_group.test.name}"
	relay_namespace_name = "${azurerm_relay_namespace.test.name}"
}
`, testAccAzureRMRelayHybridConnection_basic(rInt, location), rInt)
}

func testCheckAzureRMRelayHybridConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		relayNamespace := rs.Primary.Attributes["relay_namespace_name"]

		// Ensure resource group exists in API
		client := acceptance.AzureProvider.Meta().(*clients.Client).Relay.HybridConnectionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
