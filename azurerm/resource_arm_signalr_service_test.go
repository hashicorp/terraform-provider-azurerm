package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMSignalRService_basic(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
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
func TestAccAzureRMSignalRService_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_signalr_service.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
			{
				Config:      testAccAzureRMSignalRService_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_signalr_service"),
			},
		},
	})
}

func TestAccAzureRMSignalRService_standard(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_standardWithCapacity(ri, testLocation(), 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
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

func TestAccAzureRMSignalRService_standardWithCap2(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_standardWithCapacity(ri, testLocation(), 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
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

func TestAccAzureRMSignalRService_skuUpdate(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	freeConfig := testAccAzureRMSignalRService_basic(ri, location)
	standardConfig := testAccAzureRMSignalRService_standardWithCapacity(ri, location, 1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: freeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
			{
				Config: standardConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
			{
				Config: freeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func TestAccAzureRMSignalRService_capacityUpdate(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	standardConfig := testAccAzureRMSignalRService_standardWithCapacity(ri, location, 1)
	standardCap5Config := testAccAzureRMSignalRService_standardWithCapacity(ri, location, 5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: standardConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
			{
				Config: standardCap5Config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
			{
				Config: standardConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func TestAccAzureRMSignalRService_skuAndCapacityUpdate(t *testing.T) {
	resourceName := "azurerm_signalr_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	freeConfig := testAccAzureRMSignalRService_basic(ri, location)
	standardConfig := testAccAzureRMSignalRService_standardWithCapacity(ri, location, 2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: freeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
			{
				Config: standardConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
			{
				Config: freeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "public_port"),
					resource.TestCheckResourceAttrSet(resourceName, "server_port"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func testAccAzureRMSignalRService_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Free_F1"
    capacity = 1
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMSignalRService_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_signalr_service" "import" {
  name                = "${azurerm_signalr_service.test.name}"
  location            = "${azurerm_signalr_service.test.location}"
  resource_group_name = "${azurerm_signalr_service.test.resource_group_name}"

  sku {
    name     = "Free_F1"
    capacity = 1
  }
}
`, testAccAzureRMSignalRService_basic(rInt, location))
}

func testAccAzureRMSignalRService_standardWithCapacity(rInt int, location string, capacity int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "Standard_S1"
    capacity = %d
  }
}
`, rInt, location, rInt, capacity)
}

func testCheckAzureRMSignalRServiceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).signalr.Client
	ctx := testAccProvider.Meta().(*ArmClient).StopContext
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_signalr_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("SignalR service still exists:\n%#v", resp)
		}
	}
	return nil
}

func testCheckAzureRMSignalRServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for SignalR service: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).signalr.Client
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on signalRClient: %+v", err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: SignalR service %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}
