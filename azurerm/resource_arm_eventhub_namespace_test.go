package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventHubNamespaceCapacity_validation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{
			Value:    21,
			ErrCount: 1,
		},
		{
			Value:    1,
			ErrCount: 0,
		},
		{
			Value:    2,
			ErrCount: 0,
		},
		{
			Value:    3,
			ErrCount: 0,
		},
		{
			Value:    4,
			ErrCount: 0,
		},
		{
			Value:    0,
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateEventHubNamespaceCapacity(tc.Value, "azurerm_eventhub_namespace")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM EventHub Namespace Capacity '%d' to trigger a validation error", tc.Value)
		}
	}
}

func TestAccAzureRMEventHubNamespaceSku_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "Basic",
			ErrCount: 0,
		},
		{
			Value:    "Standard",
			ErrCount: 0,
		},
		{
			Value:    "Premium",
			ErrCount: 1,
		},
		{
			Value:    "Random",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateEventHubNamespaceSku(tc.Value, "azurerm_eventhub_namespace")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM EventHub Namespace Sku '%s' to trigger a validation error", tc.Value)
		}
	}
}

func TestAccAzureRMEventHubNamespace_basic(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMEventHubNamespace_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists("azurerm_eventhub_namespace.test"),
				),
			},
		},
	})
}

func TestAccAzureRMEventHubNamespace_standard(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMEventHubNamespace_standard(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists("azurerm_eventhub_namespace.test"),
				),
			},
		},
	})
}

func TestAccAzureRMEventHubNamespace_readDefaultKeys(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace.test"
	ri := acctest.RandInt()
	config := testAccAzureRMEventHubNamespace_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "default_primary_connection_string", regexp.MustCompile("Endpoint=.+")),
					resource.TestMatchResourceAttr(resourceName, "default_secondary_connection_string", regexp.MustCompile("Endpoint=.+")),
					resource.TestMatchResourceAttr(resourceName, "default_primary_key", regexp.MustCompile(".+")),
					resource.TestMatchResourceAttr(resourceName, "default_secondary_key", regexp.MustCompile(".+")),
				),
			},
		},
	})
}

func TestAccAzureRMEventHubNamespace_maximumThroughputUnits(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace.test"
	ri := acctest.RandInt()
	config := testAccAzureRMEventHubNamespace_maximumThroughputUnits(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMEventHubNamespace_NonStandardCasing(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMEventHubNamespaceNonStandardCasing(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists("azurerm_eventhub_namespace.test"),
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

func testCheckAzureRMEventHubNamespaceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).eventHubNamespacesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_namespace" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMEventHubNamespaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		namespaceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Event Hub Namespace: %s", namespaceName)
		}

		conn := testAccProvider.Meta().(*ArmClient).eventHubNamespacesClient

		resp, err := conn.Get(resourceGroup, namespaceName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Event Hub Namespace %q (resource group: %q) does not exist", namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on eventHubNamespacesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMEventHubNamespace_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
}
`, rInt, location, rInt)
}

func testAccAzureRMEventHubNamespace_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  capacity            = "2"
}

`, rInt, location, rInt)
}

func testAccAzureRMEventHubNamespaceNonStandardCasing(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "basic"
}
`, rInt, location, rInt)
}

func testAccAzureRMEventHubNamespace_maximumThroughputUnits(rInt int, location string) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                     = "acctesteventhubnamespace-%d"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  sku                      = "Standard"
  capacity                 = "2"
  auto_inflate_enabled     = true
  maximum_throughput_units = 20
}
`, rInt, location, rInt)
}
