package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMHDInsight_importBasic(t *testing.T) {
	resourceName := "azurerm_hdinsight_cluster.test"
	rInt := acctest.RandInt()
	rString := acctest.RandString(5)
	location := testLocation()
	config := testAccAzureRMHDInsightCluster_basicConfig(rInt, rString, location, "hadoop", 3)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
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

func TestAccAzureRMHDInsight_importVirtualNetwork(t *testing.T) {
	resourceName := "azurerm_hdinsight_cluster.test"
	rInt := acctest.RandInt()
	rString := acctest.RandString(5)
	location := testLocation()
	config := testAccAzureRMHDInsightCluster_basicNetworkConfig(rInt, rString, location, "hadoop")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
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
