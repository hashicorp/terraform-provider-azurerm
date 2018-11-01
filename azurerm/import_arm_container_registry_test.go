package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMContainerRegistry_importBasicClassic(t *testing.T) {
	resourceName := "azurerm_container_registry.test"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMContainerRegistry_basicUnmanaged(ri, rs, testLocation(), "Classic")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_account"},
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_importBasicBasic(t *testing.T) {
	resourceName := "azurerm_container_registry.test"

	ri := acctest.RandInt()
	config := testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Basic")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_account"},
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_importBasicManagedStandard(t *testing.T) {
	resourceName := "azurerm_container_registry.test"

	ri := acctest.RandInt()
	config := testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Standard")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_account"},
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_importBasicManagedPremium(t *testing.T) {
	resourceName := "azurerm_container_registry.test"

	ri := acctest.RandInt()
	config := testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Premium")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_account"},
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_importComplete(t *testing.T) {
	resourceName := "azurerm_container_registry.test"

	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	config := testAccAzureRMContainerRegistry_complete(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"storage_account"},
			},
		},
	})
}
