package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMCosmosDBAccount_importBoundedStaleness(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_boundedStaleness(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
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

func TestAccAzureRMCosmosDBAccount_importBoundedStalenessComplete(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_boundedStalenessComplete(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
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

func TestAccAzureRMCosmosDBAccount_importEventualConsistency(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_eventualConsistency(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
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

func TestAccAzureRMCosmosDBAccount_importSession(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_session(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
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

func TestAccAzureRMCosmosDBAccount_importStrong(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_strong(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
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

func TestAccAzureRMCosmosDBAccount_importGeoReplicated(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_geoReplicated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
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
