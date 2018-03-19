package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMHDInsightCluster_importBasic(t *testing.T) {
	resourceName := "azurerm_hdinsight_cluster.test"

	ri := acctest.RandInt()
	rStr := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	config := testAccAzureRMHDInsightCluster_basic(ri, rStr, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy,
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
