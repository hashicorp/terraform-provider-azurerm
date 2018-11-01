package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMPostgreSQLConfiguration_importBackslashQuote(t *testing.T) {
	resourceName := "azurerm_postgresql_configuration.test"

	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLConfiguration_backslashQuote(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationDestroy,
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

func TestAccAzureRMPostgreSQLConfiguration_importClientMinMessages(t *testing.T) {
	resourceName := "azurerm_postgresql_configuration.test"

	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLConfiguration_clientMinMessages(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationDestroy,
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

func TestAccAzureRMPostgreSQLConfiguration_importDeadlockTimeout(t *testing.T) {
	resourceName := "azurerm_postgresql_configuration.test"

	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLConfiguration_deadlockTimeout(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLConfigurationDestroy,
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
