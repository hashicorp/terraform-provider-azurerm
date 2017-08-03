package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMPostgreSQLConfiguration_importBackslashQuote(t *testing.T) {
	resourceName := "azurerm_postgresql_configuration.test"

	ri := acctest.RandInt()
	config := testAccAzureRMPostgreSQLConfiguration_backslashQuote(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testCheckAzureRMPostgreSQLConfigurationReset(s, "safe_encoding")
		},
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
