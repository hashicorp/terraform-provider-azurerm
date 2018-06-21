package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMPostgreSQLServer_importBasicNinePointFive(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"

	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_basicNinePointFive(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_importBasicNinePointSix(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"

	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_basicNinePointSix(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_importGeneralPurpose(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"

	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_generalPurpose(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}

func TestAccAzureRMPostgreSQLServer_importMemoryOptimized(t *testing.T) {
	resourceName := "azurerm_postgresql_server.test"

	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLServer_memoryOptimized(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"administrator_login_password", // not returned as sensitive
				},
			},
		},
	})
}
