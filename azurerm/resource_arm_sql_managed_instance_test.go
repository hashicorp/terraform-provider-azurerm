package azurerm

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func init() {
	resource.AddTestSweepers("azurerm_sql_managed_instance", &resource.Sweeper{
		Name: "azurerm_sql_managed_instance",
		F:    testSweepSQLMiServer,
	})
}

func testSweepSQLMiServer(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return err
	}

	client := (*armClient).sqlMiServersClient
	ctx := (*armClient).StopContext

	log.Printf("Retrieving the SQL Managed Instance..")
	results, err := client.List(ctx)
	if err != nil {
		return fmt.Errorf("Error Listing on SQL Managed Instance: %+v", err)
	}

	for _, server := range results.Values() {
		if !shouldSweepAcceptanceTestResource(*server.Name, *server.Location, region) {
			continue
		}

		resourceId, err := parseAzureResourceID(*server.ID)
		if err != nil {
			return err
		}

		resourceGroup := resourceId.ResourceGroup
		name := resourceId.Path["servers"]

		log.Printf("Deleting SQL Managed Instance '%s' in Resource Group '%s'", name, resourceGroup)
		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return err
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return err
		}
	}

	return nil
}

func TestAccAzureRMSqlMiServer_basic(t *testing.T) {
	resourceName := "azurerm_sql_managed_instance.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMSqlMiServer_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlMiServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlMiServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"administrator_login_password"},
			},
		},
	})
}

func TestAccAzureRMSqlMiServer_disappears(t *testing.T) {
	resourceName := "azurerm_sql_managed_instance.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMSqlMiServer_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlMiServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlMiServerExists(resourceName),
					testCheckAzureRMSqlMiServerDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMSqlMiServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		sqlServerName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for SQL Managed Instance: %s", sqlServerName)
		}

		conn := testAccProvider.Meta().(*ArmClient).sqlMiServersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, sqlServerName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: SQL Managed Instance %s (resource group: %s) does not exist", sqlServerName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get SQL Managed Instance: %v", err)
		}

		return nil
	}
}

func testCheckAzureRMSqlMiServerDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).sqlMiServersClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_managed_instance" {
			continue
		}

		sqlServerName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, sqlServerName)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get SQL Managed Instance: %+v", err)
		}

		return fmt.Errorf("SQL Managed Instance %s still exists", sqlServerName)

	}

	return nil
}

func testCheckAzureRMSqlMiServerDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlMiServersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, serverName)
		if err != nil {
			return err
		}

		return future.WaitForCompletionRef(ctx, client.Client)
	}
}

func testAccAzureRMSqlMiServer_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}
  
resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.0.0/24"  
  route_table_id       = "${azurerm_route_table.test.id}"
}

resource "azurerm_route_table" "test" {
	name                          = "routetable-%d"
	location                      = "${azurerm_resource_group.test.location}"
	resource_group_name           = "${azurerm_resource_group.test.name}"
	disable_bgp_route_propagation = false
  
	route {
	  name           = "RouteToAzureSqlMiMngSvc"
	  address_prefix = "0.0.0.0/0"
	  next_hop_type  = "Internet"
	}
}

resource "azurerm_subnet_route_table_association" "test" {
	subnet_id      = "${azurerm_subnet.test.id}"
	route_table_id = "${azurerm_route_table.test.id}"
}
 
resource "azurerm_sql_managed_instance" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
  license_type				   = "BasePrice"
  subnet_id					   = "${azurerm_subnet.test.id}"

  tags {
	environment = "staging"
	database    = "test"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
