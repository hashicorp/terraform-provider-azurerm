package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMCosmosDB_importBoundedStaleness(t *testing.T) {
	resourceName := "azurerm_cosmos_db.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDB_boundedStaleness(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBDestroy,
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

func TestAccAzureRMCosmosDB_importEventualConsistency(t *testing.T) {
	resourceName := "azurerm_cosmos_db.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDB_eventualConsistency(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBDestroy,
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

func TestAccAzureRMCosmosDB_importSession(t *testing.T) {
	resourceName := "azurerm_cosmos_db.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDB_session(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBDestroy,
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

func TestAccAzureRMCosmosDB_importStrong(t *testing.T) {
	resourceName := "azurerm_cosmos_db.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDB_strong(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBDestroy,
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

func TestAccAzureRMCosmosDB_importGeoReplicated(t *testing.T) {
	resourceName := "azurerm_cosmos_db.test"

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDB_geoReplicated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBDestroy,
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
