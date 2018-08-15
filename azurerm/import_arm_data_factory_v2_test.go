package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMDataFactoryV2_importWithTags(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryV2_tags(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_data_factory_v2.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
