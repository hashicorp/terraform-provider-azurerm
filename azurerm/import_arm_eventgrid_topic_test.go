package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMEventGridTopic_importBasic(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventGridTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridTopic_basic(ri),
			},

			{
				ResourceName:      "azurerm_eventgrid_topic.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMEventGridTopic_importBasicWithTags(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventGridTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridTopic_basicWithTags(ri),
			},

			{
				ResourceName:      "azurerm_eventgrid_topic.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
