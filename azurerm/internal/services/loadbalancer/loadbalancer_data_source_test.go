package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccAzureRMDataSourceLoadBalancer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb", "test")
	d := LoadBalancer{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.dataSourceBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("production"),
				check.That(data.ResourceName).Key("tags.Purpose").HasValue("AcceptanceTests"),
			),
		},
	})
}

func (r LoadBalancer) dataSourceBasic(data acceptance.TestData) string {
	resource := r.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_lb" "test" {
  name                = azurerm_lb.test.name
  resource_group_name = azurerm_lb.test.resource_group_name
}
`, resource)
}
