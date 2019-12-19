package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSubscriptions_basic(t *testing.T) {
	resourceName := "data.azurerm_subscriptions.current"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "azurerm_subscriptions" "current" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "subscriptions.0.subscription_id"),
					resource.TestCheckResourceAttrSet(resourceName, "subscriptions.0.display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "subscriptions.0.tenant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "subscriptions.0.state"),
				),
			},
		},
	})
}
