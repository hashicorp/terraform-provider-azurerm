package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMHybridConnection_basic(t *testing.T) {
	resourceName := "azurerm_relay_hybrid_connection.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHybridConnection_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHybridConnectionExists(resourceName),
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

func TestAccAzureRMHybridConnection_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_relay_hybrid_connection.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHybridConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHybridConnection_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHybridConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "metric_id"),
				),
			},
			{
				Config:      testAccAzureRMHybridConnection_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_relay_hybrid_connection"),
			},
		},
	})
}

func testAccAzureRMHybridConnection_basic(rInt int, location string) string {
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
	relay_namespace_name = "acctestrn-%d"
  }
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMHybridConnection_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_relay_namespace" "import" {
	name                 = "acctestrnhc-%d"
	resource_group_name  = "${azurerm_resource_group.test.name}"
	relay_namespace_name = "acctestrn-%d"
}
`, testAccAzureRMHybridConnection_basic(rInt, location))
}

func testCheckAzureRMHybridConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		relayNamespace := rs.Primary.Attributes["relay_namespace_name"]
		name := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API
		client := testAccProvider.Meta().(*ArmClient).Relay.HybridConnectionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testCheckAzureRMHybridConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Relay.HybridConnectionsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_relay_hybrid_connection" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		relayNamespace := rs.Primary.Attributes["relay_namespace_name"]
		name := rs.Primary.Attributes["name"]

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
