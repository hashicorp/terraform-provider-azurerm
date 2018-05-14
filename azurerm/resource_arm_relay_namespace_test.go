package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMRelayNamespace_basic(t *testing.T) {
	resourceName := "azurerm_relay_namespace.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMRelayNamespace_basic(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRelayNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayNamespaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "metric_id"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
				),
			},
		},
	})
}

func TestAccAzureRMRelayNamespace_complete(t *testing.T) {
	resourceName := "azurerm_relay_namespace.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccAzureRMRelayNamespace_complete(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRelayNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRelayNamespaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "metric_id"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
				),
			},
		},
	})
}

func testCheckAzureRMRelayNamespaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API
		client := testAccProvider.Meta().(*ArmClient).relayNamespacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).relayNamespacesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMRelayNamespace_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name =  "acctestrn-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Standard"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMRelayNamespace_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_relay_namespace" "test" {
  name =  "acctestrn-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Standard"
  }

  tags {
    "Hello" = "World"
  }
}
`, rInt, location, rInt)
}
