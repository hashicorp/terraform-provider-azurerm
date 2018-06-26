package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("azurerm_servicebus_namespace", &resource.Sweeper{
		Name: "azurerm_servicebus_namespace",
		F:    testSweepServiceBusNamespace,
	})
}

func testSweepServiceBusNamespace(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return err
	}

	client := (*armClient).serviceBusNamespacesClient
	ctx := armClient.StopContext

	log.Printf("Retrieving the Servicebus Namespaces..")
	results, err := client.List(ctx)
	if err != nil {
		return fmt.Errorf("Error Listing on Servicebus Namespaces: %+v", err)
	}

	for _, profile := range results.Values() {
		if !shouldSweepAcceptanceTestResource(*profile.Name, *profile.Location, region) {
			continue
		}

		resourceId, err := parseAzureResourceID(*profile.ID)
		if err != nil {
			return err
		}

		resourceGroup := resourceId.ResourceGroup
		name := resourceId.Path["namespaces"]

		log.Printf("Deleting Servicebus Namespace %q in Resource Group %q", name, resourceGroup)
		deleteFuture, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return err
		}

		err = deleteFuture.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccAzureRMServiceBusNamespace_basic(t *testing.T) {
	resourceName := "azurerm_servicebus_namespace.test"
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusNamespace_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_readDefaultKeys(t *testing.T) {
	resourceName := "azurerm_servicebus_namespace.test"
	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusNamespace_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(resourceName),
					resource.TestMatchResourceAttr(
						resourceName, "default_primary_connection_string", regexp.MustCompile("Endpoint=.+")),
					resource.TestMatchResourceAttr(
						resourceName, "default_secondary_connection_string", regexp.MustCompile("Endpoint=.+")),
					resource.TestMatchResourceAttr(
						resourceName, "default_primary_key", regexp.MustCompile(".+")),
					resource.TestMatchResourceAttr(
						resourceName, "default_secondary_key", regexp.MustCompile(".+")),
				),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_NonStandardCasing(t *testing.T) {
	resourceName := "azurerm_servicebus_namespace.test"

	ri := acctest.RandInt()
	config := testAccAzureRMServiceBusNamespaceNonStandardCasing(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(resourceName),
				),
			},
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testCheckAzureRMServiceBusNamespaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).serviceBusNamespacesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_namespace" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("ServiceBus Namespace still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMServiceBusNamespaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		namespaceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Service Bus Namespace: %s", namespaceName)
		}

		client := testAccProvider.Meta().(*ArmClient).serviceBusNamespacesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, namespaceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on serviceBusNamespacesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Service Bus Namespace %q (resource group: %q) does not exist", namespaceName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMServiceBusNamespace_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}
resource "azurerm_servicebus_namespace" "test" {
    name = "acctestservicebusnamespace-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    sku = "basic"
}
`, rInt, location, rInt)
}

func testAccAzureRMServiceBusNamespaceNonStandardCasing(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}
resource "azurerm_servicebus_namespace" "test" {
    name = "acctestservicebusnamespace-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    sku = "Basic"
}
`, rInt, location, rInt)
}
