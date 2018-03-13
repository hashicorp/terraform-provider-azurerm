package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMSubscriptions_basic(t *testing.T) {
	resourceName := "data.azurerm_subscriptions.current"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "azurerm_subscriptions" "current" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "subscriptions.0.subscription_id"),
					resource.TestCheckResourceAttrSet(resourceName, "subscriptions.0.display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "subscriptions.0.state"),
				),
			},
		},
	})
}
