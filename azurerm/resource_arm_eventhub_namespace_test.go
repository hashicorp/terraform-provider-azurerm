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

func TestAccAzureRMEventHubNamespace_BasicWithTagsUpdate(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMEventHubNamespace_basic(ri, testLocation())
	postConfig := testAccAzureRMEventHubNamespace_basicWithTagsUpdate(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMEventHubNamespace_BasicWithCapacity(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace.test"
	ri := acctest.RandInt()
	config := testAccAzureRMEventHubNamespace_capacity(ri, testLocation(), 20)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "capacity", "20"),
				),
			},
		},
	})
}

func TestAccAzureRMEventHubNamespace_BasicWithCapacityUpdate(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMEventHubNamespace_capacity(ri, testLocation(), 20)
	postConfig := testAccAzureRMEventHubNamespace_capacity(ri, testLocation(), 2)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "capacity", "20"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMEventHubNamespace_BasicWithSkuUpdate(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMEventHubNamespace_basic(ri, testLocation())
	postConfig := testAccAzureRMEventHubNamespace_standard(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Basic"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMEventHubNamespace_maximumThroughputUnitsUpdate(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMEventHubNamespace_maximumThroughputUnits(ri, testLocation())
	postConfig := testAccAzureRMEventHubNamespace_maximumThroughputUnitsUpdate(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "maximum_throughput_units", "20"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "maximum_throughput_units", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMEventHubNamespaceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).eventHubNamespacesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_namespace" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

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
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, namespaceName)
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

func testAccAzureRMEventHubNamespace_basicWithTagsUpdate(rInt int, location string) string {
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
  tags {
	environment = "Production"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMEventHubNamespace_capacity(rInt int, location string, capacity int) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                     = "acctesteventhubnamespace-%d"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  sku                      = "Basic"
  capacity                 = %d
}
`, rInt, location, rInt, capacity)
}

func testAccAzureRMEventHubNamespace_maximumThroughputUnitsUpdate(rInt int, location string) string {
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
  capacity                 = 1
  auto_inflate_enabled     = true
  maximum_throughput_units = 1
}
`, rInt, location, rInt)
}
