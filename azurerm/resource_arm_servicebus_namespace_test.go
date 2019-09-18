package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMServiceBusNamespace_basic(t *testing.T) {
	resourceName := "azurerm_servicebus_namespace.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceBusNamespace_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccAzureRMServiceBusNamespace_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_servicebus_namespace.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusNamespace_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMServiceBusNamespace_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_servicebus_namespace"),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_readDefaultKeys(t *testing.T) {
	resourceName := "azurerm_servicebus_namespace.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceBusNamespace_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceBusNamespaceNonStandardCasing(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_premium(t *testing.T) {
	resourceName := "azurerm_servicebus_namespace.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceBusNamespace_premium(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_basicCapacity(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceBusNamespace_basicCapacity(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Service Bus SKU \"Basic\" only supports `capacity` of 0"),
			},
		},
	})
}

func TestAccAzureRMServiceBusNamespace_premiumCapacity(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMServiceBusNamespace_premiumCapacity(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Service Bus SKU \"Premium\" only supports `capacity` of 1, 2 or 4"),
			},
		},
	})
}

func testCheckAzureRMServiceBusNamespaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).ServiceBus.NamespacesClient
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

func testCheckAzureRMServiceBusNamespaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		namespaceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Service Bus Namespace: %s", namespaceName)
		}

		client := testAccProvider.Meta().(*ArmClient).ServiceBus.NamespacesClient
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
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "basic"
}
`, rInt, location, rInt)
}

func testAccAzureRMServiceBusNamespace_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace" "import" {
  name                = "${azurerm_servicebus_namespace.test.name}"
  location            = "${azurerm_servicebus_namespace.test.location}"
  resource_group_name = "${azurerm_servicebus_namespace.test.resource_group_name}"
  sku                 = "${azurerm_servicebus_namespace.test.sku}"
}
`, testAccAzureRMServiceBusNamespace_basic(rInt, location))
}

func testAccAzureRMServiceBusNamespaceNonStandardCasing(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
}
`, rInt, location, rInt)
}

func testAccAzureRMServiceBusNamespace_premium(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Premium"
  capacity            = 1
}
`, rInt, location, rInt)
}

func testAccAzureRMServiceBusNamespace_basicCapacity(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
  capacity            = 1
}
`, rInt, location, rInt)
}

func testAccAzureRMServiceBusNamespace_premiumCapacity(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Premium"
  capacity            = 0
}
`, rInt, location, rInt)
}
