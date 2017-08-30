package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMSqlDatabase_importBasic(t *testing.T) {
	resourceName := "azurerm_sql_database.test"

	ri := acctest.RandInt()
	config := testAccAzureRMSqlDatabase_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_mode"},
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_importDataWarehouse(t *testing.T) {
	resourceName := "azurerm_sql_database.test"

	ri := acctest.RandInt()
	config := testAccAzureRMSqlDatabase_dataWarehouse(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_mode"},
			},
		},
	})
}

func TestAccAzureRMSqlDatabase_importElasticPool(t *testing.T) {
	resourceName := "azurerm_sql_database.test"

	ri := acctest.RandInt()
	config := testAccAzureRMSqlDatabase_elasticPool(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_mode"},
			},
		},
	})
}
