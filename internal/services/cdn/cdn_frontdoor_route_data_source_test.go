package cdn_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontDoorRouteDataSource struct{}

func TestAccCdnFrontDoorRouteDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_route", "test")
	d := CdnFrontDoorRouteDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("forwarding_protocol").HasValue("HttpsOnly"),
			),
		},
	})
}

func (CdnFrontDoorRouteDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_route" "test" {
  name                      = azurerm_cdn_frontdoor_route.test.name
  cdn_frontdoor_endpoint_id = azurerm_cdn_frontdoor_endpoint.test.id
}
`, CdnFrontDoorRouteResource{}.complete(data))
}
