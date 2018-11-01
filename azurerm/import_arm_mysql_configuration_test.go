package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMMySQLConfiguration_importCharacterSetServer(t *testing.T) {
	resourceName := "azurerm_mysql_configuration.test"

	ri := acctest.RandInt()
	config := testAccAzureRMMySQLConfiguration_characterSetServer(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPostgreSQLConfiguration_importInteractiveTimeout(t *testing.T) {
	resourceName := "azurerm_mysql_configuration.test"

	ri := acctest.RandInt()
	config := testAccAzureRMMySQLConfiguration_interactiveTimeout(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPostgreSQLConfiguration_importLogSlowAdminStatements(t *testing.T) {
	resourceName := "azurerm_mysql_configuration.test"

	ri := acctest.RandInt()
	config := testAccAzureRMMySQLConfiguration_logSlowAdminStatements(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMySQLConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
