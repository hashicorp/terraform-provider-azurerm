package azurerm

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMCosmosDBAccount_importEventualConsistency(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.Eventual), "", "")

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
	config := testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.Session), "", "")

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
	config := testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.Strong), "", "")

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

func TestAccAzureRMCosmosDBAccount_importBoundedStaleness(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.BoundedStaleness), "", "")

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
	config := testAccAzureRMCosmosDBAccount_boundedStaleness_complete(ri, testLocation())

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

func TestAccAzureRMCosmosDBAccount_importMongoDB(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_mongoDB(ri, testLocation())

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
	config := testAccAzureRMCosmosDBAccount_geoReplicated_customId(ri, testLocation(), testAltLocation())

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
