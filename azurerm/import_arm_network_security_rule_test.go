package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMNetworkSecurityRule_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "azurerm_network_security_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_basic(rInt, testLocation()),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
