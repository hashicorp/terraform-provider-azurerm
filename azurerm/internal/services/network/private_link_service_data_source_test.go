package network_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PrivateLinkServiceDataSource struct {
}

func TestAccDataSourcePrivateLinkService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_link_service", "test")
	r := PrivateLinkServiceDataSource{}
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("2"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.private_ip_address").HasValue("10.5.1.40"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.private_ip_address_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("nat_ip_configuration.1.private_ip_address").HasValue("10.5.1.41"),
				check.That(data.ResourceName).Key("nat_ip_configuration.1.private_ip_address_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("auto_approval_subscription_ids.0").HasValue(subscriptionId),
				check.That(data.ResourceName).Key("visibility_subscription_ids.0").HasValue(subscriptionId),
				check.That(data.ResourceName).Key("load_balancer_frontend_ip_configuration_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("load_balancer_frontend_ip_configuration_ids.0").Exists(),
			),
		},
	})
}

func (PrivateLinkServiceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_link_service" "test" {
  name                = azurerm_private_link_service.test.name
  resource_group_name = azurerm_private_link_service.test.resource_group_name
}
`, PrivateLinkServiceResource{}.complete(data))
}
