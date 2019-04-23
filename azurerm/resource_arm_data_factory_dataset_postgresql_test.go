package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAzureRMDataFactoryDatasetPostgreSQL_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryDatasetPostgreSQL_basic(ri, testLocation())
	resourceName := "azurerm_data_factory_dataset_postgresql.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDatasetPostgreSQLDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryDatasetPostgreSQLExists(resourceName),
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

func TestAccAzureRMDataFactoryDatasetPostgreSQL_update(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryDatasetPostgreSQL_update1(ri, testLocation())
	config2 := testAccAzureRMDataFactoryDatasetPostgreSQL_update2(ri, testLocation())
	resourceName := "azurerm_data_factory_dataset_postgresql.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDatasetPostgreSQLDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryDatasetPostgreSQLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameters.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "annotations.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "schema_column.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "additional_properties.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryDatasetPostgreSQLExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameters.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "annotations.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "schema_column.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "additional_properties.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description 2"),
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

func testCheckAzureRMDataFactoryDatasetPostgreSQLExists(name string) resource.TestCheckFunc {
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

		client := testAccProvider.Meta().(*ArmClient).dataFactoryDatasetClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on dataFactoryDatasetClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Data Factory Dataset PostgreSQL %q (data factory name: %q / resource group: %q) does not exist", name, dataFactoryName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryDatasetPostgreSQLDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).dataFactoryDatasetClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_dataset_postgredql" {
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
			return fmt.Errorf("Data Factory Dataset PostgreSQL still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMDataFactoryDatasetPostgreSQL_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_linked_service_postgresql" "test" {
  name                = "acctestlssql%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  connection_string   = "Host=example;Port=5432;Database=example;UID=example;EncryptionMethod=0;Password=example"
}

resource "azurerm_data_factory_dataset_postgresql" "test" {
  name                = "acctestds%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  linked_service_name = "${azurerm_data_factory_linked_service_postgresql.test.name}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMDataFactoryDatasetPostgreSQL_update1(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_linked_service_postgresql" "test" {
  name                = "acctestlssql%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  connection_string   = "Host=example;Port=5432;Database=example;UID=example;EncryptionMethod=0;Password=example"
}

resource "azurerm_data_factory_dataset_postgresql" "test" {
  name                = "acctestds%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  linked_service_name = "${azurerm_data_factory_linked_service_postgresql.test.name}"
 
  description = "test description"
  annotations = ["test1", "test2", "test3"]
  table_name  = "testTable"
  folder      = "testFolder"

  parameters = {
    foo = "test1"
    bar = "test2"
  }
 
  additional_properties = {
    foo = "test1"
    bar = "test2"
  }
	
  schema_column {
    name        = "test1"
    type        = "Byte"
    description = "description"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMDataFactoryDatasetPostgreSQL_update2(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_linked_service_postgresql" "test" {
  name                = "acctestlssql%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  connection_string   = "Host=example;Port=5432;Database=example;UID=example;EncryptionMethod=0;Password=example"
}

resource "azurerm_data_factory_dataset_postgresql" "test" {
  name                = "acctestds%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  linked_service_name = "${azurerm_data_factory_linked_service_postgresql.test.name}"
 
  description = "test description 2"
  annotations = ["test1", "test2"]
  table_name  = "testTable"
  folder      = "testFolder"

  parameters = {
    foo  = "test1"
    bar  = "test2"
    buzz = "test3"
  }
 
  additional_properties = {
    foo = "test1"
  }
	
  schema_column {
    name        = "test1"
    type        = "Byte"
    description = "description"
  }
	
  schema_column {
    name        = "test2"
    type        = "Byte"
    description = "description"
  }
}
`, rInt, location, rInt, rInt, rInt)
}
