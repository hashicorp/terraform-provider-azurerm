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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataFactoryLinkedServiceMySQL_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMDataFactoryLinkedServiceMySQL_basic(ri, acceptance.Location())
	resourceName := "azurerm_data_factory_linked_service_mysql.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryLinkedServiceMySQLDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceMySQLExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// returned as ***** from the api
					"connection_string",
				},
			},
		},
	})
}

func TestAccAzureRMDataFactoryLinkedServiceMySQL_update(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMDataFactoryLinkedServiceMySQL_update1(ri, acceptance.Location())
	config2 := testAccAzureRMDataFactoryLinkedServiceMySQL_update2(ri, acceptance.Location())
	resourceName := "azurerm_data_factory_linked_service_mysql.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryLinkedServiceMySQLDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceMySQLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameters.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "annotations.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "additional_properties.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryLinkedServiceMySQLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameters.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "annotations.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "additional_properties.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description 2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// returned as ***** from the api
					"connection_string",
				},
			},
		},
	})
}

func testCheckAzureRMDataFactoryLinkedServiceMySQLExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Data Factory: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.LinkedServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on dataFactoryLinkedServiceClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Data Factory Linked Service MySQL %q (data factory name: %q / resource group: %q) does not exist", name, dataFactoryName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryLinkedServiceMySQLDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.LinkedServiceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_linked_service_mysql" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Data Factory Linked Service MySQL still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMDataFactoryLinkedServiceMySQL_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_linked_service_mysql" "test" {
  name                = "acctestlssql%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  connection_string   = "Server=test;Port=3306;Database=test;User=test;SSLMode=1;UseSystemTrustStore=0;Password=test"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDataFactoryLinkedServiceMySQL_update1(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_linked_service_mysql" "test" {
  name                = "acctestlssql%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  connection_string   = "Server=test;Port=3306;Database=test;User=test;SSLMode=1;UseSystemTrustStore=0;Password=test"
  annotations         = ["test1", "test2", "test3"]
  description         = "test description"

  parameters = {
    foo = "test1"
    bar = "test2"
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDataFactoryLinkedServiceMySQL_update2(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_linked_service_mysql" "test" {
  name                = "acctestlssql%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  connection_string   = "Server=test;Port=3306;Database=test;User=test;SSLMode=1;UseSystemTrustStore=0;Password=test"
  annotations         = ["test1", "test2"]
  description         = "test description 2"

  parameters = {
    foo  = "test1"
    bar  = "test2"
    buzz = "test3"
  }

  additional_properties = {
    foo = "test1"
  }
}
`, rInt, location, rInt, rInt)
}
